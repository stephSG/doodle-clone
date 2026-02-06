package models

import (
	"time"

	"github.com/google/uuid"
)

// Comment represents a comment on a poll
type Comment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	PollID    uuid.UUID `json:"poll_id" db:"poll_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TableName returns the table name for Comment
func (Comment) TableName() string {
	return "comments"
}

// CommentWithUser includes user information for display
type CommentWithUser struct {
	Comment
	User User `json:"user"`
}

// CreateCommentRequest is the request payload for creating a comment
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// UpdateCommentRequest is the request payload for updating a comment
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}
