package email

import (
	"fmt"
	"strings"

	"doodle-clone/internal/config"
	"gopkg.in/gomail.v2"
)

// Sender handles email sending
type Sender struct {
	dialer *gomail.Dialer
	from   string
}

// NewSender creates a new email sender
func NewSender() *Sender {
	return &Sender{
		dialer: gomail.NewDialer(
			config.AppConfig.SMTPHost,
			toInt(config.AppConfig.SMTPPort),
			config.AppConfig.SMTPUser,
			config.AppConfig.SMTPPassword,
		),
		from: config.AppConfig.SMTPFrom,
	}
}

// Send sends an email
func (s *Sender) Send(to []string, subject, body string) error {
	if s.dialer == nil || len(to) == 0 {
		return fmt.Errorf("email not configured or no recipients")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to[0])
	if len(to) > 1 {
		m.SetHeader("Bcc", to[1:]...)
	}
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}

// SendPollNotification sends a notification about a new poll
func (s *Sender) SendPollNotification(to []string, pollTitle, pollURL, creatorName string) error {
	subject := fmt.Sprintf("New Poll: %s", pollTitle)
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.button { display: inline-block; padding: 12px 24px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 4px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>You've been invited to a poll!</h2>
				<p>%s has created a new poll titled <strong>%s</strong>.</p>
				<p>Please click the button below to vote:</p>
				<a href="%s" class="button">Vote Now</a>
				<p>Or copy this link to your browser:<br>%s</p>
				<hr>
				<p><small>This is an automated message. Please do not reply.</small></p>
			</div>
		</body>
		</html>
	`, creatorName, pollTitle, pollURL, pollURL)

	return s.Send(to, subject, body)
}

// SendVoteNotification sends a notification about a new vote
func (s *Sender) SendVoteNotification(to string, pollTitle, pollURL, voterName string) error {
	subject := fmt.Sprintf("New Vote on: %s", pollTitle)
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.button { display: inline-block; padding: 12px 24px; background-color: #2196F3; color: white; text-decoration: none; border-radius: 4px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>New vote received!</h2>
				<p><strong>%s</strong> has voted on your poll <strong>%s</strong>.</p>
				<a href="%s" class="button">View Poll</a>
				<hr>
				<p><small>To unsubscribe from these notifications, please visit your settings page.</small></p>
			</div>
		</body>
		</html>
	`, voterName, pollTitle, pollURL)

	return s.Send([]string{to}, subject, body)
}

// SendCommentNotification sends a notification about a new comment
func (s *Sender) SendCommentNotification(to string, pollTitle, pollURL, commenterName, comment string) error {
	subject := fmt.Sprintf("New Comment on: %s", pollTitle)
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.button { display: inline-block; padding: 12px 24px; background-color: #FF9800; color: white; text-decoration: none; border-radius: 4px; }
				.comment { background-color: #f5f5f5; padding: 10px; border-left: 3px solid #FF9800; margin: 10px 0; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>New comment!</h2>
				<p><strong>%s</strong> commented on <strong>%s</strong>:</p>
				<div class="comment">%s</div>
				<a href="%s" class="button">View Poll</a>
				<hr>
				<p><small>To unsubscribe from these notifications, please visit your settings page.</small></p>
			</div>
		</body>
		</html>
	`, commenterName, pollTitle, comment, pollURL)

	return s.Send([]string{to}, subject, body)
}

// SendFinalDateNotification sends a notification when final date is set
func (s *Sender) SendFinalDateNotification(to []string, pollTitle, pollURL, finalDate string, creatorName string) error {
	subject := fmt.Sprintf("Final Date Selected: %s", pollTitle)
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.button { display: inline-block; padding: 12px 24px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 4px; }
				.date-box { background-color: #e8f5e9; padding: 15px; border-radius: 4px; margin: 15px 0; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>A date has been selected!</h2>
				<p><strong>%s</strong> has selected the final date for <strong>%s</strong>:</p>
				<div class="date-box">
					<h3>%s</h3>
				</div>
				<a href="%s" class="button">View Poll</a>
				<hr>
				<p><small>Please mark your calendar!</small></p>
			</div>
		</body>
		</html>
	`, creatorName, pollTitle, finalDate, pollURL)

	return s.Send(to, subject, body)
}

// SendExpirationReminder sends a reminder before poll expires
func (s *Sender) SendExpirationReminder(to []string, pollTitle, pollURL, expiresAt string) error {
	subject := fmt.Sprintf("Reminder: %s expires soon!", pollTitle)
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.button { display: inline-block; padding: 12px 24px; background-color: #f44336; color: white; text-decoration: none; border-radius: 4px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Don't forget to vote!</h2>
				<p>The poll <strong>%s</strong> will expire on:</p>
				<h3>%s</h3>
				<p>Make sure to cast your vote before it's too late!</p>
				<a href="%s" class="button">Vote Now</a>
				<hr>
				<p><small>This is an automated reminder.</small></p>
			</div>
		</body>
		</html>
	`, pollTitle, expiresAt, pollURL)

	return s.Send(to, subject, body)
}

// IsValidEmail checks if the email format is valid
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if parts[0] == "" || parts[1] == "" {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}

// IsConfigured returns true if email is properly configured
func IsConfigured() bool {
	return config.AppConfig.SMTPHost != "" &&
		config.AppConfig.SMTPPort != "" &&
		config.AppConfig.SMTPUser != "" &&
		config.AppConfig.SMTPPassword != ""
}

func toInt(s string) int {
	var i int
	if _, err := fmt.Sscanf(s, "%d", &i); err == nil {
		return i
	}
	return 587 // Default SMTP port
}
