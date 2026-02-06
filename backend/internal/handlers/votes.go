package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"doodle-clone/internal/database"
	"doodle-clone/internal/middleware"
	"doodle-clone/internal/models"
)

type VoteHandler struct {
	db *pgxpool.Pool
}

func NewVoteHandler(db *pgxpool.Pool) *VoteHandler {
	return &VoteHandler{db: db}
}

// GetVotes returns all votes for a poll
func (h *VoteHandler) GetVotes(c *gin.Context) {
	pollID := c.Param("id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID is required"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Get poll info to check if anonymous is allowed
	var poll models.Poll
	err := h.db.QueryRow(ctx, `
		SELECT id, title, allow_multiple, anonymous, creator_id
		FROM polls WHERE id = $1
	`, pollID).Scan(&poll.ID, &poll.Title, &poll.AllowMultiple, &poll.Anonymous, &poll.CreatorID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	// Get votes grouped by user
	rows, err := h.db.Query(ctx, `
		SELECT v.id, v.poll_id, v.date_option_id, v.user_id, v.user_name, v.response, v.created_at,
		       u.id, u.name, u.avatar
		FROM votes v
		LEFT JOIN users u ON v.user_id = u.id
		WHERE v.poll_id = $1
		ORDER BY v.created_at
	`, pollID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch votes"})
		return
	}
	defer rows.Close()

	votesByUser := make(map[uuid.UUID][]models.VoteWithUser)
	anonymousVotes := []models.VoteWithUser{}

	for rows.Next() {
		var vote models.VoteWithUser
		var userIDPtr *uuid.UUID
		var user models.User

		err := rows.Scan(
			&vote.ID, &vote.PollID, &vote.DateOptionID, &userIDPtr, &vote.UserName, &vote.Response, &vote.CreatedAt,
			&user.ID, &user.Name, &user.Avatar,
		)
		if err != nil {
			continue
		}

		if userIDPtr != nil {
			vote.UserID = userIDPtr
			vote.User = &user
			vote.IsCreator = user.ID == poll.CreatorID

			if _, exists := votesByUser[*userIDPtr]; !exists {
				votesByUser[*userIDPtr] = []models.VoteWithUser{}
			}
			votesByUser[*userIDPtr] = append(votesByUser[*userIDPtr], vote)
		} else {
			anonymousVotes = append(anonymousVotes, vote)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"votes_by_user":   votesByUser,
		"anonymous_votes": anonymousVotes,
		"poll":            poll,
	})
}

// CreateVote creates a new vote
func (h *VoteHandler) CreateVote(c *gin.Context) {
	pollID := c.Param("id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID is required"})
		return
	}

	var req models.CreateVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Vote bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Vote request: response=%q, votes_count=%d", req.Response, len(req.Votes))

	ctx, cancel := database.GetContext()
	defer cancel()

	// Get poll info
	var poll models.Poll
	var creatorID uuid.UUID
	log.Printf("Fetching poll info for: %s", pollID)
	err := h.db.QueryRow(ctx, `
		SELECT id, title, allow_multiple, allow_maybe, anonymous, limit_votes, max_votes_per_user, creator_id, expires_at
		FROM polls WHERE id = $1
	`, pollID).Scan(&poll.ID, &poll.Title, &poll.AllowMultiple, &poll.AllowMaybe, &poll.Anonymous,
		&poll.LimitVotes, &poll.MaxVotesPerUser, &creatorID, &poll.ExpiresAt)
	if err != nil {
		log.Printf("Error fetching poll: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	// Check if poll is expired
	if poll.ExpiresAt != nil && poll.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This poll has expired"})
		return
	}

	// Validate "maybe" response
	if req.Response == "maybe" && !poll.AllowMaybe {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This poll does not allow 'maybe' responses"})
		return
	}

	// Get user info
	userID := middleware.GetCurrentUser(c)
	var userName string
	var userIDPtr *uuid.UUID
	log.Printf("User ID: %v", userID)

	if userID != nil {
		// Authenticated user - get their name
		err = h.db.QueryRow(ctx, "SELECT name FROM users WHERE id = $1", *userID).Scan(&userName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
			return
		}
		userIDPtr = userID
	} else {
		// Anonymous user
		if !poll.Anonymous {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You must be logged in to vote on this poll"})
			return
		}
		// For anonymous, require a name
		if req.Votes == nil || len(req.Votes) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_name is required for anonymous votes"})
			return
		}
	}

	// Handle single vote (legacy format) or multiple votes
	votesToCreate := []models.VoteItem{}
	if len(req.Votes) > 0 {
		votesToCreate = req.Votes
	} else if req.Response != "" {
		// Single vote from date_option_id query param or body
		dateOptionID := c.Query("date_option_id")
		if dateOptionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "date_option_id is required"})
			return
		}
		parsedID, err := uuid.Parse(dateOptionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_option_id"})
			return
		}
		votesToCreate = []models.VoteItem{
			{DateOptionID: parsedID, Response: req.Response},
		}
	}

	// Validate votes
	for _, voteItem := range votesToCreate {
		// Check if date option exists
		var exists bool
		err = h.db.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM date_options WHERE id = $1 AND poll_id = $2)
		`, voteItem.DateOptionID, pollID).Scan(&exists)

		if err != nil || !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date option"})
			return
		}

		// Validate "maybe" for each vote
		if voteItem.Response == "maybe" && !poll.AllowMaybe {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This poll does not allow 'maybe' responses"})
			return
		}
	}

	// Check vote limits
	if poll.LimitVotes && userID != nil {
		// Count existing votes for this user on this poll
		var existingVoteCount int
		err = h.db.QueryRow(ctx, `
			SELECT COUNT(*) FROM votes WHERE poll_id = $1 AND user_id = $2
		`, pollID, *userID).Scan(&existingVoteCount)

		if err == nil && existingVoteCount+len(votesToCreate) > poll.MaxVotesPerUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "You have exceeded the maximum number of votes for this poll",
				"limit": poll.MaxVotesPerUser,
			})
			return
		}
	}

	// Create votes
	createdVotes := []models.Vote{}
	for _, voteItem := range votesToCreate {
		voteID := uuid.New()

		// For authenticated users, use their name
		displayName := userName
		if displayName == "" && userID == nil {
			displayName = "Anonymous"
		}

		log.Printf("Creating vote: pollID=%s, dateOptionID=%s, userID=%v, displayName=%s, response=%s",
			pollID, voteItem.DateOptionID, userIDPtr, displayName, voteItem.Response)

		_, err = h.db.Exec(ctx, `
			INSERT INTO votes (id, poll_id, date_option_id, user_id, user_name, response)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (poll_id, date_option_id, user_id)
			DO UPDATE SET response = $6, user_name = $5
		`, voteID, pollID, voteItem.DateOptionID, userIDPtr, displayName, voteItem.Response)

		if err != nil {
			log.Printf("Error creating vote: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
			return
		}

		createdVotes = append(createdVotes, models.Vote{
			ID:          voteID,
			PollID:      uuid.MustParse(pollID),
			DateOptionID: voteItem.DateOptionID,
			UserID:      userIDPtr,
			UserName:    displayName,
			Response:    voteItem.Response,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"votes": createdVotes,
		"message": "Vote(s) recorded successfully",
	})
}

// UpdateVote updates an existing vote
func (h *VoteHandler) UpdateVote(c *gin.Context) {
	pollID := c.Param("id")
	voteID := c.Param("voteId")

	if pollID == "" || voteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID and Vote ID are required"})
		return
	}

	var req models.UpdateVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Get the vote
	var vote models.Vote
	err := h.db.QueryRow(ctx, `
		SELECT id, poll_id, date_option_id, user_id, response
		FROM votes WHERE id = $1
	`, voteID).Scan(&vote.ID, &vote.PollID, &vote.DateOptionID, &vote.UserID, &vote.Response)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vote not found"})
		return
	}

	// Check ownership
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	if vote.UserID == nil || *vote.UserID != *userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own votes"})
		return
	}

	// Get poll info to check if "maybe" is allowed
	var allowMaybe bool
	err = h.db.QueryRow(ctx, "SELECT allow_maybe FROM polls WHERE id = $1", pollID).Scan(&allowMaybe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get poll info"})
		return
	}

	if req.Response == "maybe" && !allowMaybe {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This poll does not allow 'maybe' responses"})
		return
	}

	// Update vote
	_, err = h.db.Exec(ctx, `
		UPDATE votes SET response = $1 WHERE id = $2
	`, req.Response, voteID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote updated successfully"})
}

// DeleteVote deletes a vote
func (h *VoteHandler) DeleteVote(c *gin.Context) {
	pollID := c.Param("id")
	voteID := c.Param("voteId")

	if pollID == "" || voteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID and Vote ID are required"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Get the vote
	var vote models.Vote
	err := h.db.QueryRow(ctx, `
		SELECT id, user_id, poll_id
		FROM votes WHERE id = $1
	`, voteID).Scan(&vote.ID, &vote.UserID, &vote.PollID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vote not found"})
		return
	}

	// Check ownership
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Allow user to delete their own vote
	if vote.UserID != nil && *vote.UserID != *userID {
		// Also allow poll creator to delete any vote
		var creatorID uuid.UUID
		err = h.db.QueryRow(ctx, "SELECT creator_id FROM polls WHERE id = $1", pollID).Scan(&creatorID)
		if err != nil || creatorID != *userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own votes"})
			return
		}
	}

	// Delete vote
	_, err = h.db.Exec(ctx, "DELETE FROM votes WHERE id = $1", voteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote deleted successfully"})
}

// GetUserVotes returns votes made by the current user
func (h *VoteHandler) GetUserVotes(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	rows, err := h.db.Query(ctx, `
		SELECT v.id, v.poll_id, v.date_option_id, v.response, v.created_at,
		       p.title, p.location, p.expires_at, do.start_time
		FROM votes v
		JOIN polls p ON v.poll_id = p.id
		JOIN date_options do ON v.date_option_id = do.id
		WHERE v.user_id = $1
		ORDER BY v.created_at DESC
	`, *userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch votes"})
		return
	}
	defer rows.Close()

	type VoteDetail struct {
		models.Vote
		PollTitle   string     `json:"poll_title"`
		PollLocation string    `json:"poll_location"`
		PollExpires *time.Time `json:"poll_expires"`
		StartTime   time.Time  `json:"start_time"`
	}

	votes := []VoteDetail{}
	for rows.Next() {
		var v VoteDetail
		err := rows.Scan(
			&v.ID, &v.PollID, &v.DateOptionID, &v.Response, &v.CreatedAt,
			&v.PollTitle, &v.PollLocation, &v.PollExpires, &v.StartTime,
		)
		if err != nil {
			continue
		}
		votes = append(votes, v)
	}

	c.JSON(http.StatusOK, gin.H{
		"votes": votes,
		"count": len(votes),
	})
}

// now returns the current time
func now() time.Time {
	return time.Now()
}
