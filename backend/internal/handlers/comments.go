package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"doodle-clone/internal/database"
	"doodle-clone/internal/middleware"
	"doodle-clone/internal/models"
)

type CommentHandler struct {
	db *pgxpool.Pool
}

func NewCommentHandler(db *pgxpool.Pool) *CommentHandler {
	return &CommentHandler{db: db}
}

// GetComments returns all comments for a poll
// @Summary      Lister les commentaires
// @Description  Retourne tous les commentaires d'un sondage
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "UUID du sondage"
// @Success      200  {object}  map[string]interface{}  "comments, count"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /polls/{id}/comments [get]
func (h *CommentHandler) GetComments(c *gin.Context) {
	pollID := c.Param("id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID is required"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Check if poll exists
	var exists bool
	err := h.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM polls WHERE id = $1)", pollID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	// Get comments with user info
	rows, err := h.db.Query(ctx, `
		SELECT c.id, c.poll_id, c.user_id, c.content, c.created_at,
		       u.id, u.name, u.avatar, u.email
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.poll_id = $1
		ORDER BY c.created_at DESC
		LIMIT 100
	`, pollID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	comments := []models.CommentWithUser{}
	for rows.Next() {
		var comment models.CommentWithUser
		err := rows.Scan(
			&comment.ID, &comment.PollID, &comment.UserID, &comment.Content, &comment.CreatedAt,
			&comment.User.ID, &comment.User.Name, &comment.User.Avatar, &comment.User.Email,
		)
		if err != nil {
			continue
		}
		comments = append(comments, comment)
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"count":    len(comments),
	})
}

// CreateComment creates a new comment
// @Summary      Créer un commentaire
// @Description  Ajoute un commentaire à un sondage
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string                        true  "UUID du sondage"
// @Param        request body      models.CreateCommentRequest  true  "Contenu du commentaire"
// @Success      201  {object}  models.CommentWithUser
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /polls/{id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
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

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Check if poll exists
	var exists bool
	err := h.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM polls WHERE id = $1)", pollID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	// Create comment
	commentID := uuid.New()
	_, err = h.db.Exec(ctx, `
		INSERT INTO comments (id, poll_id, user_id, content)
		VALUES ($1, $2, $3, $4)
	`, commentID, pollID, *userID, req.Content)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Get created comment with user info
	var comment models.CommentWithUser
	err = h.db.QueryRow(ctx, `
		SELECT c.id, c.poll_id, c.user_id, c.content, c.created_at,
		       u.id, u.name, u.avatar
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = $1
	`, commentID).Scan(
		&comment.ID, &comment.PollID, &comment.UserID, &comment.Content, &comment.CreatedAt,
		&comment.User.ID, &comment.User.Name, &comment.User.Avatar,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// UpdateComment updates an existing comment
// @Summary      Mettre à jour un commentaire
// @Description  Modifie un commentaire existant
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id        path      string                        true  "UUID du sondage"
// @Param        commentId path      string                        true  "UUID du commentaire"
// @Param        request body      models.UpdateCommentRequest  true  "Nouveau contenu"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /polls/{id}/comments/{commentId} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	pollID := c.Param("id")
	commentID := c.Param("commentId")

	if pollID == "" || commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID and Comment ID are required"})
		return
	}

	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Get the comment and verify ownership
	var commentUserID uuid.UUID
	err := h.db.QueryRow(ctx, `
		SELECT user_id FROM comments WHERE id = $1 AND poll_id = $2
	`, commentID, pollID).Scan(&commentUserID)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check ownership
	if commentUserID != *userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own comments"})
		return
	}

	// Update comment
	_, err = h.db.Exec(ctx, `
		UPDATE comments SET content = $1 WHERE id = $2
	`, req.Content, commentID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
}

// DeleteComment deletes a comment
// @Summary      Supprimer un commentaire
// @Description  Supprime un commentaire existant
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id        path      string  true  "UUID du sondage"
// @Param        commentId path      string  true  "UUID du commentaire"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /polls/{id}/comments/{commentId} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	pollID := c.Param("id")
	commentID := c.Param("commentId")

	if pollID == "" || commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID and Comment ID are required"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Get the comment and poll creator
	var commentUserID uuid.UUID
	var pollCreatorID uuid.UUID
	err := h.db.QueryRow(ctx, `
		SELECT c.user_id, p.creator_id
		FROM comments c
		JOIN polls p ON c.poll_id = p.id
		WHERE c.id = $1 AND c.poll_id = $2
	`, commentID, pollID).Scan(&commentUserID, &pollCreatorID)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check ownership (comment author or poll creator)
	if commentUserID != *userID && pollCreatorID != *userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own comments"})
		return
	}

	// Delete comment
	_, err = h.db.Exec(ctx, "DELETE FROM comments WHERE id = $1", commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
