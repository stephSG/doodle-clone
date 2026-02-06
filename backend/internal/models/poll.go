package models

import (
	"time"

	"github.com/google/uuid"
)

// Poll represents an event/survey for scheduling
type Poll struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Title           string     `json:"title" db:"title"`
	Description     string     `json:"description" db:"description"`
	Location        string     `json:"location" db:"location"`
	CreatorID       uuid.UUID  `json:"creator_id" db:"creator_id"`
	Creator         *User      `json:"creator,omitempty" db:"-"` // Populated by join
	ExpiresAt       *time.Time `json:"expires_at" db:"expires_at"`
	AllowMultiple   bool       `json:"allow_multiple" db:"allow_multiple"`   // Allow voting for multiple dates
	AllowMaybe      bool       `json:"allow_maybe" db:"allow_maybe"`         // Allow "maybe" response
	Anonymous       bool       `json:"anonymous" db:"anonymous"`             // Allow anonymous voting
	LimitVotes      bool       `json:"limit_votes" db:"limit_votes"`         // Limit number of votes per user
	MaxVotesPerUser int        `json:"max_votes_per_user" db:"max_votes_per_user"`
	FinalDate       *uuid.UUID `json:"final_date,omitempty" db:"final_date"` // ID of the chosen DateOption
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// TableName returns the table name for Poll
func (Poll) TableName() string {
	return "polls"
}

// CreatePollRequest is the request payload for creating a poll
type CreatePollRequest struct {
	Title           string     `json:"title" binding:"required,min=3,max=200"`
	Description     string     `json:"description" binding:"max=2000"`
	Location        string     `json:"location" binding:"max=500"`
	ExpiresAt       *time.Time `json:"expires_at"`
	AllowMultiple   bool       `json:"allow_multiple"`
	AllowMaybe      bool       `json:"allow_maybe"`
	Anonymous       bool       `json:"anonymous"`
	LimitVotes      bool       `json:"limit_votes"`
	MaxVotesPerUser int        `json:"max_votes_per_user"`
	Dates           []DateRequest `json:"dates" binding:"required,min=1"`
}

// DateRequest represents a date option in the create request
type DateRequest struct {
	StartTime time.Time  `json:"start_time" binding:"required"`
	EndTime   *time.Time `json:"end_time"`
}

// UpdatePollRequest is the request payload for updating a poll
type UpdatePollRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Location    *string    `json:"location"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// SetFinalDateRequest is the request payload for setting the final date
type SetFinalDateRequest struct {
	DateOptionID uuid.UUID `json:"date_option_id" binding:"required"`
}

// IsExpired checks if the poll has expired
func (p *Poll) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return p.ExpiresAt.Before(time.Now())
}

// CanVote checks if a user can still vote on this poll
func (p *Poll) CanVote() bool {
	return !p.IsExpired()
}
