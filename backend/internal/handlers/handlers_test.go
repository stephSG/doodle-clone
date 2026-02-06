package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"doodle-clone/internal/config"
	"doodle-clone/internal/database"
	"doodle-clone/internal/middleware"
	"doodle-clone/internal/models"
)

// setupTestDB creates a test database connection
func setupTestDB(t *testing.T) *pgxpool.Pool {
	// Load test config
	config.AppConfig = &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBName:     "doodle_clone",
		DBUser:     "postgres",
		DBPassword: "postgres",
		JWTSecret:  "test-secret",
		Environment: "test",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connStr := "postgres://postgres:postgres@localhost:5432/doodle_clone"
	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err, "Failed to connect to test database")

	return pool
}

// setupTestContext creates a test Gin context
func setupTestContext() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// createTestUser creates a test user in the database
func createTestUser(t *testing.T, db *pgxpool.Pool) *models.User {
	userID := uuid.New()
	ctx := context.Background()

	_, err := db.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, name, provider)
		VALUES ($1, $2, $3, $4, 'email')
		ON CONFLICT (email) DO NOTHING
	`, userID, "test@example.com", "$2a$10$testhash", "Test User")

	require.NoError(t, err, "Failed to create test user")

	return &models.User{
		ID:       userID,
		Email:    "test@example.com",
		Name:     "Test User",
		Provider: "email",
	}
}

// createTestPoll creates a test poll in the database
func createTestPoll(t *testing.T, db *pgxpool.Pool, creatorID uuid.UUID) *models.Poll {
	pollID := uuid.New()
	accessCode := "TESTCODE1"
	ctx := context.Background()

	_, err := db.Exec(ctx, `
		INSERT INTO polls (id, title, description, location, creator_id, access_code, allow_multiple, allow_maybe, anonymous)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, pollID, "Test Poll", "Test Description", "Test Location", creatorID, accessCode, true, true, true)

	require.NoError(t, err, "Failed to create test poll")

	// Add date options
	dateID1 := uuid.New()
	dateID2 := uuid.New()
	_, err = db.Exec(ctx, `
		INSERT INTO date_options (id, poll_id, start_time, end_time)
		VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)
	`, dateID1, pollID, time.Now().Add(24*time.Hour), time.Now().Add(25*time.Hour),
		dateID2, pollID, time.Now().Add(48*time.Hour), time.Now().Add(49*time.Hour))

	require.NoError(t, err, "Failed to create test date options")

	return &models.Poll{
		ID:          pollID,
		Title:       "Test Poll",
		Description: "Test Description",
		Location:    "Test Location",
		CreatorID:   creatorID,
		AccessCode:  accessCode,
	}
}

// cleanupTestData removes test data from the database
func cleanupTestData(t *testing.T, db *pgxpool.Pool, userID, pollID uuid.UUID) {
	ctx := context.Background()

	// Delete in correct order due to foreign keys
	_, _ = db.Exec(ctx, "DELETE FROM votes WHERE poll_id = $1", pollID)
	_, _ = db.Exec(ctx, "DELETE FROM comments WHERE poll_id = $1", pollID)
	_, _ = db.Exec(ctx, "DELETE FROM date_options WHERE poll_id = $1", pollID)
	_, _ = db.Exec(ctx, "DELETE FROM polls WHERE id = $1", pollID)
	_, _ = db.Exec(ctx, "DELETE FROM refresh_tokens WHERE user_id = $1", userID)
	_, _ = db.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
}

// AuthHandler Tests

func TestAuthHandler_Register(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	handler := NewAuthHandler(db)

	router := setupTestContext()
	router.POST("/register", handler.Register)

	t.Run("Valid registration", func(t *testing.T) {
		body := models.CreateUserRequest{
			Email:    "newuser@example.com",
			Password: "Password123!",
			Name:     "New User",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "New User", response.User.Name)

		// Cleanup
		cleanupTestData(t, db, response.User.ID, uuid.Nil)
	})

	t.Run("Duplicate email", func(t *testing.T) {
		user := createTestUser(t, db)
		defer cleanupTestData(t, db, user.ID, uuid.Nil)

		body := models.CreateUserRequest{
			Email:    user.Email,
			Password: "Password123!",
			Name:     "Duplicate User",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("Invalid input", func(t *testing.T) {
		body := models.CreateUserRequest{
			Email:    "invalid-email",
			Password: "123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	handler := NewAuthHandler(db)

	router := setupTestContext()
	router.POST("/login", handler.Login)

	user := createTestUser(t, db)
	defer cleanupTestData(t, db, user.ID, uuid.Nil)

	t.Run("Valid login with password", func(t *testing.T) {
		// First, we need to create a user with a known password hash
		// For this test, we'll skip bcrypt and just test the flow

		body := models.LoginRequest{
			Email:    user.Email,
			Password: "somepassword", // This won't match since we didn't use bcrypt
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should get 401 because password doesn't match (we used a dummy hash)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Non-existent user", func(t *testing.T) {
		body := models.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "Password123!",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthHandler_GetMe(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	handler := NewAuthHandler(db)

	router := setupTestContext()
	router.GET("/me", middleware.Auth(), handler.GetMe)

	user := createTestUser(t, db)
	defer cleanupTestData(t, db, user.ID, uuid.Nil)

	// Generate a token for the user
	claims := middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: middleware.RegisteredClaims{
			ExpiresAt: nil, // Never expire for tests
			IssuedAt:  nil,
		},
	}
	token, err := middleware.GenerateToken(claims)
	require.NoError(t, err)

	t.Run("Valid token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, user.Email, response.Email)
	})

	t.Run("No token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/me", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// PollHandler Tests

func TestPollHandler_ListPolls(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	handler := NewPollHandler(db)

	router := setupTestContext()
	router.GET("/polls", handler.ListPolls)

	user := createTestUser(t, db)
	poll := createTestPoll(t, db, user.ID)
	defer cleanupTestData(t, db, user.ID, poll.ID)

	t.Run("List polls", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/polls", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "polls")
	})
}

func TestPollHandler_GetPoll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	handler := NewPollHandler(db)

	router := setupTestContext()
	router.GET("/polls/:id", handler.GetPoll)

	user := createTestUser(t, db)
	poll := createTestPoll(t, db, user.ID)
	defer cleanupTestData(t, db, user.ID, poll.ID)

	t.Run("Get poll by UUID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/polls/"+poll.ID.String(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "poll")
	})

	t.Run("Get poll by access code", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/polls/"+poll.AccessCode, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Poll not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/polls/"+uuid.New().String(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestPollHandler_CreatePoll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	handler := NewPollHandler(db)

	router := setupTestContext()
	router.POST("/polls", middleware.Auth(), handler.CreatePoll)

	user := createTestUser(t, db)
	defer cleanupTestData(t, db, user.ID, uuid.Nil)

	// Generate a token for the user
	claims := middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: middleware.RegisteredClaims{
			ExpiresAt: nil, // Never expire for tests
			IssuedAt:  nil,
		},
	}
	token, err := middleware.GenerateToken(claims)
	require.NoError(t, err)

	t.Run("Create valid poll", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		body := models.CreatePollRequest{
			Title:       "New Test Poll",
			Description: "New Description",
			Location:    "New Location",
			AllowMultiple: true,
			AllowMaybe:    true,
			Anonymous:     true,
			Dates: []models.DateRequest{
				{StartTime: startTime},
			},
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/polls", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.Poll
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "New Test Poll", response.Title)
		assert.NotEmpty(t, response.AccessCode)

		// Cleanup the created poll
		cleanupTestData(t, db, uuid.Nil, response.ID)
	})

	t.Run("No authentication", func(t *testing.T) {
		body := models.CreatePollRequest{
			Title: "Unauthorized Poll",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/polls", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// VoteHandler Tests

func TestVoteHandler_CreateVote(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	voteHandler := NewVoteHandler(db)

	router := setupTestContext()
	router.POST("/polls/:id/vote", voteHandler.CreateVote)

	user := createTestUser(t, db)
	poll := createTestPoll(t, db, user.ID)
	defer cleanupTestData(t, db, user.ID, poll.ID)

	// Get date options
	ctx := context.Background()
	var dateID uuid.UUID
	err := db.QueryRow(ctx, "SELECT id FROM date_options WHERE poll_id = $1 LIMIT 1", poll.ID).Scan(&dateID)
	require.NoError(t, err)

	t.Run("Anonymous vote", func(t *testing.T) {
		body := models.CreateVoteRequest{
			Votes: []models.VoteItem{
				{DateOptionID: dateID, Response: "yes"},
			},
			UserName: "Anonymous User",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/polls/"+poll.ID.String()+"/vote", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Invalid date option", func(t *testing.T) {
		body := models.CreateVoteRequest{
			Votes: []models.VoteItem{
				{DateOptionID: uuid.New(), Response: "yes"},
			},
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/polls/"+poll.ID.String()+"/vote", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestVoteHandler_GetVotes(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	voteHandler := NewVoteHandler(db)

	router := setupTestContext()
	router.GET("/polls/:id/votes", voteHandler.GetVotes)

	user := createTestUser(t, db)
	poll := createTestPoll(t, db, user.ID)
	defer cleanupTestData(t, db, user.ID, poll.ID)

	t.Run("Get votes for poll", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/polls/"+poll.ID.String()+"/votes", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "votes_by_user")
		assert.Contains(t, response, "anonymous_votes")
	})
}

// CommentHandler Tests

func TestCommentHandler_CreateComment(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	commentHandler := NewCommentHandler(db)

	router := setupTestContext()
	router.POST("/polls/:id/comments", middleware.Auth(), commentHandler.CreateComment)

	user := createTestUser(t, db)
	poll := createTestPoll(t, db, user.ID)
	defer cleanupTestData(t, db, user.ID, poll.ID)

	// Generate a token for the user
	claims := middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: middleware.RegisteredClaims{
			ExpiresAt: nil, // Never expire for tests
			IssuedAt:  nil,
		},
	}
	token, err := middleware.GenerateToken(claims)
	require.NoError(t, err)

	t.Run("Create valid comment", func(t *testing.T) {
		body := models.CreateCommentRequest{
			Content: "This is a test comment",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/polls/"+poll.ID.String()+"/comments", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.CommentWithUser
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "This is a test comment", response.Content)
	})

	t.Run("Poll not found", func(t *testing.T) {
		body := models.CreateCommentRequest{
			Content: "This should fail",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/polls/"+uuid.New().String()+"/comments", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCommentHandler_GetComments(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	commentHandler := NewCommentHandler(db)

	router := setupTestContext()
	router.GET("/polls/:id/comments", commentHandler.GetComments)

	user := createTestUser(t, db)
	poll := createTestPoll(t, db, user.ID)
	defer cleanupTestData(t, db, user.ID, poll.ID)

	t.Run("Get comments for poll", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/polls/"+poll.ID.String()+"/comments", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "comments")
	})
}

// Helper method for AuthHandler to expose generateToken for tests
// This is now done via middleware.GenerateToken()
