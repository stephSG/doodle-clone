package models

import (
	"time"

	"github.com/google/uuid"
)

// DateOption represents a date/time option for a poll
type DateOption struct {
	ID        uuid.UUID `json:"id" db:"id"`
	PollID    uuid.UUID `json:"poll_id" db:"poll_id"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty" db:"end_time"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// TableName returns the table name for DateOption
func (DateOption) TableName() string {
	return "date_options"
}

// DateOptionWithStats includes vote statistics
type DateOptionWithStats struct {
	DateOption
	YesCount   int `json:"yes_count"`
	NoCount    int `json:"no_count"`
	MaybeCount int `json:"maybe_count"`
	TotalVotes int `json:"total_votes"`
}

// AddDateOptionRequest is the request payload for adding a date option
type AddDateOptionRequest struct {
	StartTime time.Time  `json:"start_time" binding:"required"`
	EndTime   *time.Time `json:"end_time"`
}

// IsFinalDate checks if this date option is the final selected date
func (d *DateOption) IsFinalDate(poll *Poll) bool {
	if poll.FinalDate == nil {
		return false
	}
	return d.ID == *poll.FinalDate
}
