package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"doodle-clone/internal/config"
	"doodle-clone/internal/database"
	"doodle-clone/internal/middleware"
	"doodle-clone/internal/models"
)

type AuthHandler struct {
	db          *pgxpool.Pool
	oauthConfig *oauth2.Config
}

func NewAuthHandler(db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{
		db: db,
		oauthConfig: &oauth2.Config{
			RedirectURL:  config.AppConfig.GoogleRedirectURL,
			ClientID:     config.AppConfig.GoogleClientID,
			ClientSecret: config.AppConfig.GoogleClientSecret,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// generateToken creates a JWT token for a user
func (h *AuthHandler) generateToken(userID uuid.UUID, email string) (string, error) {
	claims := middleware.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AppConfig.JWTExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

// generateRefreshToken creates a refresh token
func (h *AuthHandler) generateRefreshToken(userID uuid.UUID) (string, error) {
	// Generate random token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)

	// Store in database
	ctx, cancel := database.GetContext()
	defer cancel()

	_, err := h.db.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`, userID, token, time.Now().Add(config.AppConfig.RefreshExpiry))

	return token, err
}

// validateRefreshToken checks if a refresh token is valid
func (h *AuthHandler) validateRefreshToken(token string) (uuid.UUID, error) {
	ctx, cancel := database.GetContext()
	defer cancel()

	var userID uuid.UUID
	var expiresAt time.Time

	err := h.db.QueryRow(ctx, `
		SELECT user_id, expires_at
		FROM refresh_tokens
		WHERE token = $1
	`, token).Scan(&userID, &expiresAt)

	if err != nil {
		return uuid.Nil, errors.New("invalid refresh token")
	}

	if time.Now().After(expiresAt) {
		// Delete expired token
		h.deleteRefreshToken(token)
		return uuid.Nil, errors.New("refresh token expired")
	}

	return userID, nil
}

func (h *AuthHandler) deleteRefreshToken(token string) {
	ctx, cancel := database.GetContext()
	defer cancel()
	h.db.Exec(ctx, "DELETE FROM refresh_tokens WHERE token = $1", token)
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Check if user already exists
	var exists bool
	err := h.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	userID := uuid.New()
	_, err = h.db.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, name, provider)
		VALUES ($1, $2, $3, $4, 'email')
	`, userID, req.Email, string(hashedPassword), req.Name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate tokens
	token, err := h.generateToken(userID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := h.generateRefreshToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Set httpOnly cookie for refresh token
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(config.AppConfig.RefreshExpiry.Seconds()),
		"/",
		"",
		config.IsProduction(),
		true, // httpOnly
	)

	user := models.User{
		ID:        userID,
		Email:     req.Email,
		Name:      req.Name,
		Provider:  "email",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Get user
	var user models.User
	var avatar sql.NullString
	err := h.db.QueryRow(ctx, `
		SELECT id, email, password_hash, name, avatar, provider, created_at, updated_at
		FROM users WHERE email = $1
	`, req.Email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &avatar, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
	if avatar.Valid {
		user.Avatar = avatar.String
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		log.Printf("Login DB error for email %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check provider
	if user.Provider != "email" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please use Google OAuth to login"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	token, err := h.generateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := h.generateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Set httpOnly cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(config.AppConfig.RefreshExpiry.Seconds()),
		"/",
		"",
		config.IsProduction(),
		true,
	)

	user.PasswordHash = ""
	c.JSON(http.StatusOK, models.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	})
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token"})
		return
	}

	// Validate refresh token
	userID, err := h.validateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get user
	ctx, cancel := database.GetContext()
	defer cancel()

	var user models.User
	err = h.db.QueryRow(ctx, `
		SELECT id, email, name, avatar, provider, created_at, updated_at
		FROM users WHERE id = $1
	`, userID).Scan(&user.ID, &user.Email, &user.Name, &user.Avatar, &user.Provider, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Generate new tokens
	token, err := h.generateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	newRefreshToken, err := h.generateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Delete old token
	h.deleteRefreshToken(refreshToken)

	// Set new cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		int(config.AppConfig.RefreshExpiry.Seconds()),
		"/",
		"",
		config.IsProduction(),
		true,
	)

	c.JSON(http.StatusOK, models.AuthResponse{
		Token:        token,
		RefreshToken: newRefreshToken,
		User:         user,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		h.deleteRefreshToken(refreshToken)
	}

	// Clear cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		config.IsProduction(),
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetMe returns the current user
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	var user models.User
	err := h.db.QueryRow(ctx, `
		SELECT id, email, name, avatar, provider, created_at, updated_at
		FROM users WHERE id = $1
	`, *userID).Scan(&user.ID, &user.Email, &user.Name, &user.Avatar, &user.Provider, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GoogleLogin initiates Google OAuth flow
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	// Generate state token for CSRF protection
	state := generateStateToken()

	// Set cookie - more permissive for development
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", state, 600, "/", "", false, false)

	// Build OAuth URL
	url := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles Google OAuth callback
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// Verify state
	state := c.Query("state")
	storedState, err := c.Cookie("oauth_state")
	if err != nil || state != storedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}
	c.SetCookie("oauth_state", "", -1, "/", "", config.IsProduction(), true)

	// Exchange code for token
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in request"})
		return
	}

	token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Failed to exchange token: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to exchange token: %v", err)})
		return
	}

	// Get user info
	client := h.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info"})
		return
	}

	// Parse user info
	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse user info"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Check if user exists
	var user models.User
	err = h.db.QueryRow(ctx, `
		SELECT id, email, name, avatar, provider, created_at, updated_at
		FROM users WHERE email = $1
	`, googleUser.Email).Scan(&user.ID, &user.Email, &user.Name, &user.Avatar, &user.Provider, &user.CreatedAt, &user.UpdatedAt)

	var userID uuid.UUID
	userName := googleUser.Name
	if userName == "" {
		userName = googleUser.Email
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Create new user
			userID = uuid.New()
			_, err = h.db.Exec(ctx, `
				INSERT INTO users (id, email, name, avatar, provider)
				VALUES ($1, $2, $3, $4, 'google')
			`, userID, googleUser.Email, userName, googleUser.Picture)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	} else {
		userID = user.ID
	}

	// Generate tokens
	email := googleUser.Email
	if email == "" && user.ID != uuid.Nil {
		email = user.Email
	}

	tokenStr, err := h.generateToken(userID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	_, err = h.generateRefreshToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Redirect to frontend with token
	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", config.AppConfig.FrontendURL, tokenStr)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// generateStateToken creates a random state token for OAuth
func generateStateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// UpdateProfile updates user profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Check if email is already taken by another user
	var exists bool
	err := h.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)
	`, req.Email, *userID).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already taken"})
		return
	}

	// Update user
	_, err = h.db.Exec(ctx, `
		UPDATE users SET name = $1, email = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`, req.Name, req.Email, *userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// ChangePassword changes user password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Get current password hash
	var currentHash string
	err := h.db.QueryRow(ctx, "SELECT password_hash FROM users WHERE id = $1", *userID).Scan(&currentHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.OldPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid current password"})
		return
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update password
	_, err = h.db.Exec(ctx, `
		UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, string(newHash), *userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
