package models

import (
	"time"

	"github.com/google/uuid"
)

// Notification represents a notification sent to users
type Notification struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	PollID       uuid.UUID  `json:"poll_id" db:"poll_id"`
	UserID       *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	Type         string     `json:"type" db:"type"` // event_reminder, new_vote, new_comment, etc.
	Status       string     `json:"status" db:"status"` // pending, sent, failed
	ScheduledAt  time.Time  `json:"scheduled_at" db:"scheduled_at"`
	SentAt       *time.Time `json:"sent_at,omitempty" db:"sent_at"`
	ErrorMessage *string    `json:"error_message,omitempty" db:"error_message"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// NotificationSetting represents admin-configurable notification settings
type NotificationSetting struct {
	Key         string    `json:"key" db:"key"`
	Value       string    `json:"value" db:"value"`
	Description string    `json:"description" db:"description"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Notification types
const (
	NotificationTypeEventReminder = "event_reminder"
	NotificationTypeNewVote       = "new_vote"
	NotificationTypeNewComment    = "new_comment"
	NotificationTypeFinalDate     = "final_date"
)

// Notification statuses
const (
	NotificationStatusPending = "pending"
	NotificationStatusSent    = "sent"
	NotificationStatusFailed  = "failed"
)

// Default notification settings keys
const (
	SettingReminderEnabled    = "reminder_enabled"
	SettingReminderHours      = "reminder_hours"
	SettingNewVoteEnabled     = "new_vote_enabled"
	SettingNewCommentEnabled  = "new_comment_enabled"
	SettingFinalDateEnabled   = "final_date_enabled"
)
