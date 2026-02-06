package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"log"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"doodle-clone/internal/database"
	"doodle-clone/internal/middleware"
	"doodle-clone/internal/models"
)

type PollHandler struct {
	db                  *pgxpool.Pool
	notificationHandler *NotificationHandler
}

func NewPollHandler(db *pgxpool.Pool) *PollHandler {
	return &PollHandler{db: db}
}

func (h *PollHandler) SetNotificationHandler(nh *NotificationHandler) {
	h.notificationHandler = nh
}

// ListPolls returns a list of public polls
// @Summary      Lister les sondages
// @Description  Retourne la liste des sondages publics
// @Tags         polls
// @Accept       json
// @Produce      json
// @Param        search query string false "Terme de recherche dans le titre"
// @Success      200  {object}  map[string]interface{}  "polls, count"
// @Failure      500  {object}  map[string]string
// @Router       /polls [get]
func (h *PollHandler) ListPolls(c *gin.Context) {
	ctx, cancel := database.GetContext()
	defer cancel()

	// Query params
	search := c.Query("search")
	limit := 20
	offset := 0

	// Build query
	query := `
		SELECT p.id, p.title, p.description, p.location, p.creator_id, p.access_code, p.expires_at,
		       p.allow_multiple, p.allow_maybe, p.anonymous, p.limit_votes, p.max_votes_per_user,
		       p.final_date, p.created_at, p.updated_at,
		       u.id, u.name, u.avatar,
		       COUNT(DISTINCT v.user_id) as participant_count
		FROM polls p
		LEFT JOIN users u ON p.creator_id = u.id
		LEFT JOIN votes v ON p.id = v.poll_id
		WHERE p.expires_at IS NULL OR p.expires_at > CURRENT_TIMESTAMP
	`

	args := []interface{}{}
	argCount := 1

	if search != "" {
		query += " AND p.title ILIKE $" + string(rune('0'+argCount))
		args = append(args, "%"+search+"%")
		argCount++
	}

	query += ` GROUP BY p.id, p.title, p.description, p.location, p.creator_id, p.access_code, p.expires_at,
	                    p.allow_multiple, p.allow_maybe, p.anonymous, p.limit_votes, p.max_votes_per_user,
	                    p.final_date, p.created_at, p.updated_at, u.id, u.name, u.avatar
	            ORDER BY p.created_at DESC LIMIT $` + string(rune('0'+argCount)) + " OFFSET $" + string(rune('0'+argCount+1))
	args = append(args, limit, offset)

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch polls"})
		return
	}
	defer rows.Close()

	var polls []models.Poll
	for rows.Next() {
		var poll models.Poll
		var creator models.User
		var avatar sql.NullString
		var participantCount int

		err := rows.Scan(
			&poll.ID, &poll.Title, &poll.Description, &poll.Location, &poll.CreatorID, &poll.AccessCode, &poll.ExpiresAt,
			&poll.AllowMultiple, &poll.AllowMaybe, &poll.Anonymous, &poll.LimitVotes, &poll.MaxVotesPerUser,
			&poll.FinalDate, &poll.CreatedAt, &poll.UpdatedAt,
			&creator.ID, &creator.Name, &avatar,
			&participantCount,
		)
		if err != nil {
			log.Printf("Error scanning poll: %v", err)
			continue
		}
		if avatar.Valid {
			creator.Avatar = avatar.String
		}

		poll.Creator = &creator
		polls = append(polls, poll)
	}

	c.JSON(http.StatusOK, gin.H{
		"polls": polls,
		"count": len(polls),
	})
}

// GetPoll returns a single poll with all details
// Supports both UUID and access_code as the ID parameter
// @Summary      Obtenir un sondage
// @Description  Retourne les détails d'un sondage (par UUID ou code d'accès)
// @Tags         polls
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "UUID du sondage ou code d'accès"
// @Success      200  {object}  map[string]interface{}  "poll, date_options, comments, votes"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /polls/{id} [get]
func (h *PollHandler) GetPoll(c *gin.Context) {
	pollID := c.Param("id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID is required"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Get poll with creator info - try by UUID first, then by access_code
	var poll models.Poll
	var creator models.User
	var creatorAvatar sql.NullString

	// First try to find by UUID
	err := h.db.QueryRow(ctx, `
		SELECT p.id, p.title, p.description, p.location, p.creator_id, p.access_code, p.expires_at,
		       p.allow_multiple, p.allow_maybe, p.anonymous, p.limit_votes, p.max_votes_per_user,
		       p.final_date, p.created_at, p.updated_at,
		       u.id, u.name, u.avatar, u.email
		FROM polls p
		LEFT JOIN users u ON p.creator_id = u.id
		WHERE p.id = $1
	`, pollID).Scan(
		&poll.ID, &poll.Title, &poll.Description, &poll.Location, &poll.CreatorID, &poll.AccessCode, &poll.ExpiresAt,
		&poll.AllowMultiple, &poll.AllowMaybe, &poll.Anonymous, &poll.LimitVotes, &poll.MaxVotesPerUser,
		&poll.FinalDate, &poll.CreatedAt, &poll.UpdatedAt,
		&creator.ID, &creator.Name, &creatorAvatar, &creator.Email,
	)

	// If not found by UUID, try by access_code
	if err != nil {
		err = h.db.QueryRow(ctx, `
			SELECT p.id, p.title, p.description, p.location, p.creator_id, p.access_code, p.expires_at,
			       p.allow_multiple, p.allow_maybe, p.anonymous, p.limit_votes, p.max_votes_per_user,
			       p.final_date, p.created_at, p.updated_at,
			       u.id, u.name, u.avatar, u.email
			FROM polls p
			LEFT JOIN users u ON p.creator_id = u.id
			WHERE p.access_code = $1
		`, pollID).Scan(
			&poll.ID, &poll.Title, &poll.Description, &poll.Location, &poll.CreatorID, &poll.AccessCode, &poll.ExpiresAt,
			&poll.AllowMultiple, &poll.AllowMaybe, &poll.Anonymous, &poll.LimitVotes, &poll.MaxVotesPerUser,
			&poll.FinalDate, &poll.CreatedAt, &poll.UpdatedAt,
			&creator.ID, &creator.Name, &creatorAvatar, &creator.Email,
		)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
			return
		}
		log.Printf("Error fetching poll: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch poll"})
		return
	}

	if creatorAvatar.Valid {
		creator.Avatar = creatorAvatar.String
	}

	poll.Creator = &creator

	// Get date options with stats
	dateOptions, err := h.getDateOptionsWithStats(ctx, poll.ID)
	if err != nil {
		log.Printf("Error fetching date options for poll %s: %v", pollID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch date options"})
		return
	}

	// Get comments
	comments, err := h.getComments(ctx, poll.ID)
	if err != nil {
		comments = []models.CommentWithUser{}
	}

	// Get votes with user info
	votes, err := h.getVotesWithUsers(ctx, poll.ID)
	if err != nil {
		votes = []models.VoteWithUser{}
	}

	c.JSON(http.StatusOK, gin.H{
		"poll":         poll,
		"date_options": dateOptions,
		"comments":     comments,
		"votes":        votes,
	})
}

// CreatePoll creates a new poll
// @Summary      Créer un sondage
// @Description  Crée un nouveau sondage
// @Tags         polls
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreatePollRequest true "Données du sondage"
// @Success      201  {object}  models.Poll
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /polls [post]
func (h *PollHandler) CreatePoll(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var req models.CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Generate unique access code
	accessCode := generateAccessCode()
	for {
		var exists bool
		err := h.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM polls WHERE access_code = $1)", accessCode).Scan(&exists)
		if err != nil || !exists {
			break
		}
		accessCode = generateAccessCode()
	}

	// Create poll
	pollID := uuid.New()
	_, err := h.db.Exec(ctx, `
		INSERT INTO polls (id, title, description, location, creator_id, access_code, expires_at,
		                  allow_multiple, allow_maybe, anonymous, limit_votes, max_votes_per_user)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, pollID, req.Title, req.Description, req.Location, *userID, accessCode, req.ExpiresAt,
		req.AllowMultiple, req.AllowMaybe, req.Anonymous, req.LimitVotes, req.MaxVotesPerUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create poll"})
		return
	}

	// Create date options
	for _, date := range req.Dates {
		dateOptionID := uuid.New()
		_, err = h.db.Exec(ctx, `
			INSERT INTO date_options (id, poll_id, start_time, end_time)
			VALUES ($1, $2, $3, $4)
		`, dateOptionID, pollID, date.StartTime, date.EndTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create date options"})
			return
		}
	}

	// Get created poll
	var poll models.Poll
	err = h.db.QueryRow(ctx, `
		SELECT id, title, description, location, creator_id, access_code, expires_at,
		       allow_multiple, allow_maybe, anonymous, limit_votes, max_votes_per_user,
		       final_date, created_at, updated_at
		FROM polls WHERE id = $1
	`, pollID).Scan(&poll.ID, &poll.Title, &poll.Description, &poll.Location, &poll.CreatorID, &poll.AccessCode, &poll.ExpiresAt,
		&poll.AllowMultiple, &poll.AllowMaybe, &poll.Anonymous, &poll.LimitVotes, &poll.MaxVotesPerUser,
		&poll.FinalDate, &poll.CreatedAt, &poll.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created poll"})
		return
	}

	c.JSON(http.StatusCreated, poll)
}

// generateAccessCode generates a random 8-character access code
func generateAccessCode() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // No confusing chars like I, O, 0, 1
	b := make([]byte, 8)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// UpdatePoll updates an existing poll
// @Summary      Mettre à jour un sondage
// @Description  Met à jour un sondage existant (réservé au créateur)
// @Tags         polls
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string                      true  "UUID du sondage"
// @Param        request body      models.UpdatePollRequest  true  "Champs à mettre à jour"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /polls/{id} [put]
func (h *PollHandler) UpdatePoll(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	pollID := c.Param("id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID is required"})
		return
	}

	var req models.UpdatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Check if user is the creator
	var creatorID uuid.UUID
	err := h.db.QueryRow(ctx, "SELECT creator_id FROM polls WHERE id = $1", pollID).Scan(&creatorID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if creatorID != *userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own polls"})
		return
	}

	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Title != nil {
		updates = append(updates, "title = $"+string(rune('0'+argCount)))
		args = append(args, *req.Title)
		argCount++
	}
	if req.Description != nil {
		updates = append(updates, "description = $"+string(rune('0'+argCount)))
		args = append(args, *req.Description)
		argCount++
	}
	if req.Location != nil {
		updates = append(updates, "location = $"+string(rune('0'+argCount)))
		args = append(args, *req.Location)
		argCount++
	}
	if req.ExpiresAt != nil {
		updates = append(updates, "expires_at = $"+string(rune('0'+argCount)))
		args = append(args, *req.ExpiresAt)
		argCount++
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, pollID)

	query := "UPDATE polls SET " + string(updates[0])
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += " WHERE id = $" + string(rune('0'+argCount))

	_, err = h.db.Exec(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update poll"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Poll updated successfully"})
}

// DeletePoll deletes a poll
// @Summary      Supprimer un sondage
// @Description  Supprime un sondage (réservé au créateur)
// @Tags         polls
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "UUID du sondage"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /polls/{id} [delete]
func (h *PollHandler) DeletePoll(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	pollID := c.Param("id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID is required"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Check if user is the creator
	var creatorID uuid.UUID
	err := h.db.QueryRow(ctx, "SELECT creator_id FROM polls WHERE id = $1", pollID).Scan(&creatorID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if creatorID != *userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own polls"})
		return
	}

	// Delete poll (cascade will delete related records)
	_, err = h.db.Exec(ctx, "DELETE FROM polls WHERE id = $1", pollID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete poll"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Poll deleted successfully"})
}

// SetFinalDate sets the final date for a poll
// @Summary      Fixer la date finale
// @Description  Fixe la date finale d'un sondage (réservé au créateur)
// @Tags         polls
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string                        true  "UUID du sondage"
// @Param        request body      models.SetFinalDateRequest  true  "ID de l'option de date"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /polls/{id}/final [post]
func (h *PollHandler) SetFinalDate(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	pollID := c.Param("id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID is required"})
		return
	}

	var req models.SetFinalDateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Check if user is the creator
	var creatorID uuid.UUID
	err := h.db.QueryRow(ctx, "SELECT creator_id FROM polls WHERE id = $1", pollID).Scan(&creatorID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if creatorID != *userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the creator can set the final date"})
		return
	}

	// Verify the date option belongs to this poll
	var exists bool
	err = h.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM date_options WHERE id = $1 AND poll_id = $2)
	`, req.DateOptionID, pollID).Scan(&exists)

	if err != nil || !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date option"})
		return
	}

	// Set final date
	_, err = h.db.Exec(ctx, "UPDATE polls SET final_date = $1 WHERE id = $2", req.DateOptionID, pollID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set final date"})
		return
	}

	// Schedule reminder notifications
	if h.notificationHandler != nil {
		pollUUID, _ := uuid.Parse(pollID)
		go h.notificationHandler.ScheduleReminderForPoll(pollUUID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Final date set successfully"})
}

// GetUserPolls returns polls created by the current user
// @Summary      Mes sondages
// @Description  Retourne la liste des sondages créés par l'utilisateur
// @Tags         polls
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "polls, count"
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/polls [get]
func (h *PollHandler) GetUserPolls(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	rows, err := h.db.Query(ctx, `
		SELECT p.id, p.title, p.description, p.location, p.expires_at,
		       p.final_date, p.created_at, p.updated_at,
		       COUNT(DISTINCT v.user_id) as participant_count
		FROM polls p
		LEFT JOIN votes v ON p.id = v.poll_id
		WHERE p.creator_id = $1
		GROUP BY p.id, p.title, p.description, p.location, p.expires_at,
		         p.final_date, p.created_at, p.updated_at
		ORDER BY p.created_at DESC
	`, *userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch polls"})
		return
	}
	defer rows.Close()

	var polls []models.Poll
	for rows.Next() {
		var poll models.Poll
		var participantCount int

		err := rows.Scan(
			&poll.ID, &poll.Title, &poll.Description, &poll.Location, &poll.ExpiresAt,
			&poll.FinalDate, &poll.CreatedAt, &poll.UpdatedAt,
			&participantCount,
		)
		if err != nil {
			continue
		}

		polls = append(polls, poll)
	}

	c.JSON(http.StatusOK, gin.H{
		"polls": polls,
		"count": len(polls),
	})
}

// AddDateOption adds a new date option to a poll
// @Summary      Ajouter une option de date
// @Description  Ajoute une nouvelle option de date à un sondage
// @Tags         polls
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string                        true  "UUID du sondage"
// @Param        request body      models.AddDateOptionRequest  true  "Date et heure"
// @Success      201  {object}  models.DateOption
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /polls/{id}/dates [post]
func (h *PollHandler) AddDateOption(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	pollID := c.Param("id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID is required"})
		return
	}

	var req models.AddDateOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Check if user is the creator
	var creatorID uuid.UUID
	err := h.db.QueryRow(ctx, "SELECT creator_id FROM polls WHERE id = $1", pollID).Scan(&creatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if creatorID != *userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the creator can add date options"})
		return
	}

	// Create date option
	dateOptionID := uuid.New()
	_, err = h.db.Exec(ctx, `
		INSERT INTO date_options (id, poll_id, start_time, end_time)
		VALUES ($1, $2, $3, $4)
	`, dateOptionID, pollID, req.StartTime, req.EndTime)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create date option"})
		return
	}

	// Get created date option
	var dateOption models.DateOption
	err = h.db.QueryRow(ctx, `
		SELECT id, poll_id, start_time, end_time, created_at
		FROM date_options WHERE id = $1
	`, dateOptionID).Scan(&dateOption.ID, &dateOption.PollID, &dateOption.StartTime, &dateOption.EndTime, &dateOption.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created date option"})
		return
	}

	c.JSON(http.StatusCreated, dateOption)
}

// Helper functions

func (h *PollHandler) getDateOptionsWithStats(ctx context.Context, pollID uuid.UUID) ([]models.DateOptionWithStats, error) {
	log.Printf("Fetching date options for poll: %s", pollID)
	rows, err := h.db.Query(ctx, `
		SELECT d.id, d.poll_id, d.start_time, d.end_time, d.created_at,
		       COALESCE(SUM(CASE WHEN v.response = 'yes' THEN 1 ELSE 0 END), 0) as yes_count,
		       COALESCE(SUM(CASE WHEN v.response = 'no' THEN 1 ELSE 0 END), 0) as no_count,
		       COALESCE(SUM(CASE WHEN v.response = 'maybe' THEN 1 ELSE 0 END), 0) as maybe_count,
		       COALESCE(COUNT(v.id), 0) as total_votes
		FROM date_options d
		LEFT JOIN votes v ON d.id = v.date_option_id
		WHERE d.poll_id = $1
		GROUP BY d.id
		ORDER BY d.start_time
	`, pollID)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []models.DateOptionWithStats
	for rows.Next() {
		var opt models.DateOptionWithStats
		err := rows.Scan(
			&opt.ID, &opt.PollID, &opt.StartTime, &opt.EndTime, &opt.CreatedAt,
			&opt.YesCount, &opt.NoCount, &opt.MaybeCount, &opt.TotalVotes,
		)
		if err != nil {
			continue
		}
		options = append(options, opt)
	}

	return options, nil
}

func (h *PollHandler) getComments(ctx context.Context, pollID uuid.UUID) ([]models.CommentWithUser, error) {
	rows, err := h.db.Query(ctx, `
		SELECT c.id, c.poll_id, c.user_id, c.content, c.created_at,
		       u.id, u.name, u.avatar
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.poll_id = $1
		ORDER BY c.created_at DESC
	`, pollID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.CommentWithUser
	for rows.Next() {
		var comment models.CommentWithUser
		err := rows.Scan(
			&comment.ID, &comment.PollID, &comment.UserID, &comment.Content, &comment.CreatedAt,
			&comment.User.ID, &comment.User.Name, &comment.User.Avatar,
		)
		if err != nil {
			continue
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (h *PollHandler) getVotesWithUsers(ctx context.Context, pollID uuid.UUID) ([]models.VoteWithUser, error) {
	rows, err := h.db.Query(ctx, `
		SELECT v.id, v.poll_id, v.date_option_id, v.user_id, v.user_name, v.response, v.created_at,
		       u.id, u.name, u.avatar
		FROM votes v
		LEFT JOIN users u ON v.user_id = u.id
		WHERE v.poll_id = $1
		ORDER BY v.created_at DESC
	`, pollID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes []models.VoteWithUser
	for rows.Next() {
		var vote models.VoteWithUser
		var userID sql.NullString
		var userIDPtr *uuid.UUID
		var avatar sql.NullString
		var userName sql.NullString
		var userID2 sql.NullString
		var userName2 sql.NullString

		err := rows.Scan(
			&vote.ID, &vote.PollID, &vote.DateOptionID, &userID, &userName, &vote.Response, &vote.CreatedAt,
			&userID2, &userName2, &avatar,
		)
		if err != nil {
			log.Printf("Error scanning vote: %v", err)
			continue
		}

		// Handle nullable user_id for anonymous votes
		if userID.Valid {
			uid, err := uuid.Parse(userID.String)
			if err == nil {
				userIDPtr = &uid
			}
		}

		vote.UserID = userIDPtr
		vote.UserName = userName.String

		// Only populate user fields if we have a user
		if userIDPtr != nil && userID2.Valid {
			uid, err := uuid.Parse(userID2.String)
			if err == nil {
				vote.User.ID = uid
			}
			vote.User.Name = userName2.String
			if avatar.Valid {
				vote.User.Avatar = avatar.String
			}
		} else {
			// Anonymous vote - clear user fields
			vote.User.ID = uuid.Nil
			vote.User.Name = ""
			vote.User.Avatar = ""
		}

		votes = append(votes, vote)
	}

	return votes, nil
}
