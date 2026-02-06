package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"doodle-clone/internal/database"
	"doodle-clone/internal/email"
	"doodle-clone/internal/middleware"
	"doodle-clone/internal/models"
)

type NotificationHandler struct {
	db     *pgxpool.Pool
	email  *email.Sender
	ticker *time.Ticker
	stopCh chan struct{}
}

func NewNotificationHandler(db *pgxpool.Pool, emailSender *email.Sender) *NotificationHandler {
	return &NotificationHandler{
		db:     db,
		email:  emailSender,
		stopCh: make(chan struct{}),
	}
}

// StartBackgroundWorker starts the background worker to process notifications
func (h *NotificationHandler) StartBackgroundWorker() {
	// Run every 5 minutes
	h.ticker = time.NewTicker(5 * time.Minute)
	log.Println("Notification worker started")

	// Run immediately on start
	go h.processPendingNotifications()

	go func() {
		for {
			select {
			case <-h.ticker.C:
				h.processPendingNotifications()
			case <-h.stopCh:
				h.ticker.Stop()
				log.Println("Notification worker stopped")
				return
			}
		}
	}()
}

// StopBackgroundWorker stops the background worker
func (h *NotificationHandler) StopBackgroundWorker() {
	close(h.stopCh)
}

// processPendingNotifications processes all pending notifications
func (h *NotificationHandler) processPendingNotifications() {
	ctx, cancel := database.GetContext(30 * time.Second)
	defer cancel()

	// Get pending notifications that are due
	now := time.Now()
	rows, err := h.db.Query(ctx, `
		SELECT id, poll_id, user_id, type, scheduled_at
		FROM notifications
		WHERE status = $1 AND scheduled_at <= $2
		ORDER BY scheduled_at ASC
		LIMIT 100
	`, models.NotificationStatusPending, now)
	if err != nil {
		log.Printf("Error fetching pending notifications: %v", err)
		return
	}
	defer rows.Close()

	var notifications []struct {
		ID     uuid.UUID
		PollID uuid.UUID
		UserID *uuid.UUID
		Type   string
	}

	for rows.Next() {
		var n struct {
			ID     uuid.UUID
			PollID uuid.UUID
			UserID *uuid.UUID
			Type   string
		}
		if err := rows.Scan(&n.ID, &n.PollID, &n.UserID, &n.Type); err != nil {
			log.Printf("Error scanning notification: %v", err)
			continue
		}
		notifications = append(notifications, n)
	}

	log.Printf("Processing %d pending notifications", len(notifications))

	for _, n := range notifications {
		if err := h.sendNotification(ctx, n.ID, n.PollID, n.UserID, n.Type); err != nil {
			log.Printf("Failed to send notification %s: %v", n.ID, err)
			h.updateNotificationStatus(ctx, n.ID, models.NotificationStatusFailed, err.Error())
		} else {
			h.updateNotificationStatus(ctx, n.ID, models.NotificationStatusSent, "")
		}
	}
}

// sendNotification sends a notification
func (h *NotificationHandler) sendNotification(ctx context.Context, notificationID, pollID uuid.UUID, userID *uuid.UUID, notificationType string) error {
	// Get poll details
	var poll struct {
		Title      string
		AccessCode string
		CreatorID  uuid.UUID
	}
	err := h.db.QueryRow(ctx, `
		SELECT title, access_code, creator_id FROM polls WHERE id = $1
	`, pollID).Scan(&poll.Title, &poll.AccessCode, &poll.CreatorID)
	if err != nil {
		return fmt.Errorf("failed to get poll: %w", err)
	}

	// Get the final date if set
	var dateOption struct {
		StartTime time.Time
	}
	var hasFinalDate bool
	_ = h.db.QueryRow(ctx, `
		SELECT start_time FROM date_options
		WHERE poll_id = $1 AND id = (SELECT final_date FROM polls WHERE id = $1)
	`, pollID).Scan(&dateOption.StartTime)
	if err == nil {
		hasFinalDate = true
	}

	var recipientEmail string
	var recipientName string

	if userID != nil {
		// Get user details
		err = h.db.QueryRow(ctx, `
			SELECT email, name FROM users WHERE id = $1
		`, *userID).Scan(&recipientEmail, &recipientName)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
	} else {
		// Send to poll creator
		err = h.db.QueryRow(ctx, `
			SELECT email, name FROM users WHERE id = $1
		`, poll.CreatorID).Scan(&recipientEmail, &recipientName)
		if err != nil {
			return fmt.Errorf("failed to get creator: %w", err)
		}
	}

	// Build email based on type
	var subject, body string
	baseURL := "http://localhost:5173" // TODO: from config
	pollURL := fmt.Sprintf("%s/poll/%s", baseURL, poll.AccessCode)

	switch notificationType {
	case models.NotificationTypeEventReminder:
		subject = fmt.Sprintf("Rappel: %s", poll.Title)
		body = fmt.Sprintf(`
			<h2>Rappel pour l'événement : %s</h2>
			<p>Bonjour %s,</p>
			<p>Ceci est un rappel pour l'événement <strong>%s</strong>.</p>
			%s
			<p><a href="%s" style="padding: 10px 20px; background: #4F46E5; color: white; text-decoration: none; border-radius: 5px;">Voir le sondage</a></p>
		`, poll.Title, recipientName, poll.Title,
			func() string {
				if hasFinalDate {
					return fmt.Sprintf("<p>Date et heure: <strong>%s</strong></p>", dateOption.StartTime.Format("02/01/2006 à 15:04"))
				}
				return ""
			}(), pollURL)

	case models.NotificationTypeNewVote:
		subject = fmt.Sprintf("Nouveau vote pour: %s", poll.Title)
		body = fmt.Sprintf(`
			<h2>Nouveau vote enregistré</h2>
			<p>Bonjour %s,</p>
			<p>Un nouveau vote a été enregistré pour le sondage <strong>%s</strong>.</p>
			<p><a href="%s" style="padding: 10px 20px; background: #4F46E5; color: white; text-decoration: none; border-radius: 5px;">Voir le sondage</a></p>
		`, recipientName, poll.Title, pollURL)

	case models.NotificationTypeNewComment:
		subject = fmt.Sprintf("Nouveau commentaire pour: %s", poll.Title)
		body = fmt.Sprintf(`
			<h2>Nouveau commentaire</h2>
			<p>Bonjour %s,</p>
			<p>Un nouveau commentaire a été ajouté au sondage <strong>%s</strong>.</p>
			<p><a href="%s" style="padding: 10px 20px; background: #4F46E5; color: white; text-decoration: none; border-radius: 5px;">Voir le sondage</a></p>
		`, recipientName, poll.Title, pollURL)

	case models.NotificationTypeFinalDate:
		subject = fmt.Sprintf("Date fixée pour: %s", poll.Title)
		body = fmt.Sprintf(`
			<h2>Date finale fixée</h2>
			<p>Bonjour %s,</p>
			<p>La date finale a été fixée pour le sondage <strong>%s</strong>.</p>
			%s
			<p><a href="%s" style="padding: 10px 20px; background: #4F46E5; color: white; text-decoration: none; border-radius: 5px;">Voir le sondage</a></p>
		`, recipientName, poll.Title,
			func() string {
				if hasFinalDate {
					return fmt.Sprintf("<p>Date retenue: <strong>%s</strong></p>", dateOption.StartTime.Format("02/01/2006 à 15:04"))
				}
				return ""
			}(), pollURL)
	}

	// Send email
	if err := h.email.Send([]string{recipientEmail}, subject, body); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// updateNotificationStatus updates the status of a notification
func (h *NotificationHandler) updateNotificationStatus(ctx context.Context, notificationID uuid.UUID, status, errorMessage string) {
	now := time.Now()
	var errMsg *string
	if errorMessage != "" {
		errMsg = &errorMessage
	}

	_, err := h.db.Exec(ctx, `
		UPDATE notifications SET status = $1, sent_at = $2, error_message = $3
		WHERE id = $4
	`, status, now, errMsg, notificationID)
	if err != nil {
		log.Printf("Failed to update notification status: %v", err)
	}
}

// ScheduleReminderForPoll schedules reminder notifications for all participants of a poll
func (h *NotificationHandler) ScheduleReminderForPoll(pollID uuid.UUID) error {
	ctx, cancel := database.GetContext()
	defer cancel()

	// Check if reminders are enabled
	var enabled bool
	var hoursBefore int
	err := h.db.QueryRow(ctx, `
		SELECT CASE WHEN value = 'true' THEN true ELSE false END
		FROM notification_settings WHERE key = $1
	`, models.SettingReminderEnabled).Scan(&enabled)
	if err != nil || !enabled {
		return nil // Reminders disabled
	}

	err = h.db.QueryRow(ctx, `
		SELECT value::int FROM notification_settings WHERE key = $1
	`, models.SettingReminderHours).Scan(&hoursBefore)
	if err != nil {
		hoursBefore = 1 // Default to 1 hour
	}

	// Get poll final date
	var finalDateID *uuid.UUID
	err = h.db.QueryRow(ctx, `SELECT final_date FROM polls WHERE id = $1`, pollID).Scan(&finalDateID)
	if err != nil || finalDateID == nil {
		return nil // No final date set
	}

	var eventTime time.Time
	err = h.db.QueryRow(ctx, `
		SELECT start_time FROM date_options WHERE id = $1
	`, *finalDateID).Scan(&eventTime)
	if err != nil {
		return nil
	}

	// Calculate reminder time
	reminderTime := eventTime.Add(-time.Duration(hoursBefore) * time.Hour)

	if reminderTime.Before(time.Now()) {
		return nil // Event already passed or reminder time passed
	}

	// Get all participants (users who voted)
	rows, err := h.db.Query(ctx, `
		SELECT DISTINCT user_id FROM votes WHERE poll_id = $1 AND user_id IS NOT NULL
	`, pollID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			continue
		}

		// Check if notification already scheduled
		var exists bool
		h.db.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM notifications
			WHERE poll_id = $1 AND user_id = $2 AND type = $3 AND status = 'pending')
		`, pollID, userID, models.NotificationTypeEventReminder).Scan(&exists)

		if !exists {
			h.db.Exec(ctx, `
				INSERT INTO notifications (id, poll_id, user_id, type, scheduled_at)
				VALUES ($1, $2, $3, $4, $5)
			`, uuid.New(), pollID, userID, models.NotificationTypeEventReminder, reminderTime)
		}
	}

	return nil
}

// GetNotificationSettings returns all notification settings
// @Summary      Obtenir les paramètres de notification
// @Description  Retourne tous les paramètres de notification (admin uniquement)
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "settings"
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notifications/settings [get]
func (h *NotificationHandler) GetNotificationSettings(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Check if user is admin
	var isAdmin bool
	h.db.QueryRow(ctx, `SELECT email IN ($1, $2, $3) FROM users WHERE id = $4`, "steph.leminhnhut@gmail.com", "stephane.le@gmail.com", "stephane.consulting.ai@gmail.com", *userID).Scan(&isAdmin)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	rows, err := h.db.Query(ctx, `SELECT key, value, description, updated_at FROM notification_settings ORDER BY key`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	}
	defer rows.Close()

	var settings []models.NotificationSetting
	for rows.Next() {
		var s models.NotificationSetting
		if err := rows.Scan(&s.Key, &s.Value, &s.Description, &s.UpdatedAt); err != nil {
			continue
		}
		settings = append(settings, s)
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// UpdateNotificationSetting updates a notification setting
// @Summary      Mettre à jour un paramètre de notification
// @Description  Met à jour un paramètre de notification (admin uniquement)
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body      object  true  "key et value"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notifications/settings [put]
func (h *NotificationHandler) UpdateNotificationSetting(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Check if user is admin
	var isAdmin bool
	h.db.QueryRow(ctx, `SELECT email IN ($1, $2, $3) FROM users WHERE id = $4`, "steph.leminhnhut@gmail.com", "stephane.le@gmail.com", "stephane.consulting.ai@gmail.com", *userID).Scan(&isAdmin)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	_, err := h.db.Exec(ctx, `
		INSERT INTO notification_settings (key, value, description, updated_at)
		VALUES ($1, $2, '', CURRENT_TIMESTAMP)
		ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = CURRENT_TIMESTAMP
	`, req.Key, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update setting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting updated"})
}
