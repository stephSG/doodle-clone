package models

import (
	"time"

	"github.com/google/uuid"
)

// Vote represents a user's vote on a date option
type Vote struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	PollID      uuid.UUID  `json:"poll_id" db:"poll_id"`
	DateOptionID uuid.UUID `json:"date_option_id" db:"date_option_id"`
	UserID      *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	UserName    string     `json:"user_name" db:"user_name"` // Display name (user's name or custom for anonymous)
	Response    string     `json:"response" db:"response"`   // "yes", "no", "maybe"
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// TableName returns the table name for Vote
func (Vote) TableName() string {
	return "votes"
}

// CreateVoteRequest is the request payload for creating/updating a vote
type CreateVoteRequest struct {
	Response string     `json:"response" binding:"omitempty,oneof=yes no maybe"` // For single vote with date_option_id query param
	Votes    []VoteItem `json:"votes"` // For multiple votes at once
	UserName string     `json:"user_name"` // Name for anonymous voting
}

// VoteItem represents a single vote item
type VoteItem struct {
	DateOptionID uuid.UUID `json:"date_option_id" binding:"required"`
	Response     string    `json:"response" binding:"required,oneof=yes no maybe"`
}

// UpdateVoteRequest is the request payload for updating a vote
type UpdateVoteRequest struct {
	Response string `json:"response" binding:"required,oneof=yes no maybe"`
}

// VoteWithUser includes user information for display
type VoteWithUser struct {
	Vote
	User     *User `json:"user,omitempty"`
	IsCreator bool  `json:"is_creator"` // Is this the poll creator
}
