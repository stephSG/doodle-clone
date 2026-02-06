package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jung-kurt/gofpdf"
	"doodle-clone/internal/database"
	"doodle-clone/internal/models"
)

type ExportHandler struct {
	db *pgxpool.Pool
}

func NewExportHandler(db *pgxpool.Pool) *ExportHandler {
	return &ExportHandler{db: db}
}

// ExportPDF generates a PDF export of a poll
// @Summary      Exporter en PDF
// @Description  Génère un fichier PDF du sondage avec les résultats
// @Tags         exports
// @Accept       json
// @Produce      application/pdf
// @Param        id   path      string  true  "UUID du sondage"
// @Success      200  {file}  file
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /polls/{id}/export/pdf [get]
func (h *ExportHandler) ExportPDF(c *gin.Context) {
	pollID := c.Param("id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID is required"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Get poll data
	poll, err := h.getPollWithDetails(ctx, uuid.MustParse(pollID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch poll"})
		return
	}

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Title
	pdf.Cell(40, 10, "Poll: "+poll.Title)
	pdf.Ln(12)

	// Description
	if poll.Description != "" {
		pdf.SetFont("Arial", "", 12)
		pdf.Cell(40, 10, "Description: "+poll.Description)
		pdf.Ln(10)
	}

	// Location
	if poll.Location != "" {
		pdf.Cell(40, 10, "Location: "+poll.Location)
		pdf.Ln(10)
	}

	// Creator
	pdf.Cell(40, 10, "Created by: "+poll.Creator.Name)
	pdf.Ln(10)

	pdf.Ln(5)

	// Date options with votes
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Voting Results")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)

	for _, do := range poll.DateOptions {
		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(40, 8, formatDate(do.StartTime))
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 10)
		pdf.Cell(10, 6, "Yes: "+fmt.Sprint(do.YesCount))
		pdf.Cell(20, 6, "No: "+fmt.Sprint(do.NoCount))
		pdf.Cell(20, 6, "Maybe: "+fmt.Sprint(do.MaybeCount))
		pdf.Ln(8)

		// List voters
		for _, vote := range poll.Votes {
			if vote.DateOptionID == do.ID {
				prefix := "  "
				if vote.Response == "yes" {
					prefix += "[✓] "
				} else if vote.Response == "no" {
					prefix += "[✗] "
				} else {
					prefix += "[?] "
				}
				pdf.Cell(40, 6, prefix+vote.UserName)
				pdf.Ln(6)
			}
		}
		pdf.Ln(3)
	}

	// Final date
	if poll.FinalDate != nil {
		pdf.SetFont("Arial", "B", 12)
		pdf.Ln(5)
		pdf.Cell(40, 10, "Final Date Selected:")
		pdf.Ln(7)
		pdf.SetFont("Arial", "", 11)
		for _, do := range poll.DateOptions {
			if do.ID == *poll.FinalDate {
				pdf.Cell(40, 8, formatDate(do.StartTime))
				pdf.Ln(8)
				break
			}
		}
	}

	// Set headers
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=poll_%s.pdf", pollID))

	// Write to buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}

// ExportICS generates an ICS calendar file for a poll
// @Summary      Exporter en ICS
// @Description  Génère un fichier calendrier (Google Calendar, Outlook)
// @Tags         exports
// @Accept       json
// @Produce      text/calendar
// @Param        id   path      string  true  "UUID du sondage"
// @Success      200  {file}  file
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /polls/{id}/export/ics [get]
func (h *ExportHandler) ExportICS(c *gin.Context) {
	pollID := c.Param("id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID is required"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Get poll data
	poll, err := h.getPollWithDetails(ctx, uuid.MustParse(pollID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch poll"})
		return
	}

	// If final date is set, export only that
	// Otherwise export all dates
	var datesToExport []models.DateOptionWithStats
	if poll.FinalDate != nil {
		for _, do := range poll.DateOptions {
			if do.ID == *poll.FinalDate {
				datesToExport = []models.DateOptionWithStats{do}
				break
			}
		}
	} else {
		datesToExport = poll.DateOptions
	}

	// Generate ICS content
	ics := "BEGIN:VCALENDAR\r\n"
	ics += "VERSION:2.0\r\n"
	ics += "PRODID:-//Doodle Clone//EN\r\n"
	ics += "CALSCALE:GREGORIAN\r\n"
	ics += "METHOD:PUBLISH\r\n"

	for _, do := range datesToExport {
		ics += "BEGIN:VEVENT\r\n"
		ics += fmt.Sprintf("DTSTART:%s\r\n", formatICSDate(do.StartTime))
		if do.EndTime != nil {
			ics += fmt.Sprintf("DTEND:%s\r\n", formatICSDate(*do.EndTime))
		}
		ics += fmt.Sprintf("SUMMARY:%s\r\n", poll.Title)
		if poll.Location != "" {
			ics += fmt.Sprintf("LOCATION:%s\r\n", poll.Location)
		}
		if poll.Description != "" {
			ics += fmt.Sprintf("DESCRIPTION:%s\\r\\\\nVote count: %d yes, %d no, %d maybe\r\n",
				poll.Description, do.YesCount, do.NoCount, do.MaybeCount)
		}
		ics += fmt.Sprintf("UID:%s@doodleclone\r\n", do.ID.String())
		ics += "END:VEVENT\r\n"
	}

	ics += "END:VCALENDAR\r\n"

	// Set headers
	c.Header("Content-Type", "text/calendar")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=poll_%s.ics", pollID))

	c.String(http.StatusOK, ics)
}

// ExportCSV generates a CSV export of votes
// @Summary      Exporter en CSV
// @Description  Génère un fichier CSV des votes pour analyse
// @Tags         exports
// @Accept       json
// @Produce      text/csv
// @Param        id   path      string  true  "UUID du sondage"
// @Success      200  {file}  file
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /polls/{id}/export/csv [get]
func (h *ExportHandler) ExportCSV(c *gin.Context) {
	pollID := c.Param("id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID is required"})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	// Get poll data
	poll, err := h.getPollWithDetails(ctx, uuid.MustParse(pollID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch poll"})
		return
	}

	// Generate CSV content
	csv := "Participant"

	// Header row with all dates
	for _, do := range poll.DateOptions {
		csv += ";" + formatDateTime(do.StartTime)
	}
	csv += "\r\n"

	// Group votes by user
	votesByUser := make(map[string]map[uuid.UUID]string)
	for _, vote := range poll.Votes {
		if _, exists := votesByUser[vote.UserName]; !exists {
			votesByUser[vote.UserName] = make(map[uuid.UUID]string)
		}
		votesByUser[vote.UserName][vote.DateOptionID] = vote.Response
	}

	// Data rows
	for userName, votes := range votesByUser {
		csv += userName
		for _, do := range poll.DateOptions {
			if response, exists := votes[do.ID]; exists {
				csv += ";" + response
			} else {
				csv += ";"
			}
		}
		csv += "\r\n"
	}

	// Summary row
	csv += "\r\n;Summary"
	for _, do := range poll.DateOptions {
		csv += fmt.Sprintf(";Yes:%d No:%d Maybe:%d", do.YesCount, do.NoCount, do.MaybeCount)
	}
	csv += "\r\n"

	// Set headers
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=poll_%s.csv", pollID))

	c.String(http.StatusOK, csv)
}

// Helper structs and functions

type PollExport struct {
	models.Poll
	Creator      models.User
	DateOptions  []models.DateOptionWithStats
	Votes        []models.Vote
	Comments     []models.CommentWithUser
}

func (h *ExportHandler) getPollWithDetails(ctx context.Context, pollID uuid.UUID) (*PollExport, error) {
	// Get poll with creator
	var poll PollExport
	err := h.db.QueryRow(ctx, `
		SELECT p.id, p.title, p.description, p.location, p.creator_id, p.expires_at,
		       p.allow_multiple, p.allow_maybe, p.anonymous, p.limit_votes, p.max_votes_per_user,
		       p.final_date, p.created_at, p.updated_at,
		       u.id, u.name, u.avatar, u.email
		FROM polls p
		LEFT JOIN users u ON p.creator_id = u.id
		WHERE p.id = $1
	`, pollID).Scan(
		&poll.ID, &poll.Title, &poll.Description, &poll.Location, &poll.CreatorID, &poll.ExpiresAt,
		&poll.AllowMultiple, &poll.AllowMaybe, &poll.Anonymous, &poll.LimitVotes, &poll.MaxVotesPerUser,
		&poll.FinalDate, &poll.CreatedAt, &poll.UpdatedAt,
		&poll.Creator.ID, &poll.Creator.Name, &poll.Creator.Avatar, &poll.Creator.Email,
	)

	if err != nil {
		return nil, err
	}

	// Get date options with stats
	rows, err := h.db.Query(ctx, `
		SELECT do.id, do.poll_id, do.start_time, do.end_time, do.created_at,
		       COALESCE(SUM(CASE WHEN v.response = 'yes' THEN 1 ELSE 0 END), 0) as yes_count,
		       COALESCE(SUM(CASE WHEN v.response = 'no' THEN 1 ELSE 0 END), 0) as no_count,
		       COALESCE(SUM(CASE WHEN v.response = 'maybe' THEN 1 ELSE 0 END), 0) as maybe_count
		FROM date_options do
		LEFT JOIN votes v ON do.id = v.date_option_id
		WHERE do.poll_id = $1
		GROUP BY do.id
		ORDER BY do.start_time
	`, pollID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var do models.DateOptionWithStats
		err := rows.Scan(
			&do.ID, &do.PollID, &do.StartTime, &do.EndTime, &do.CreatedAt,
			&do.YesCount, &do.NoCount, &do.MaybeCount,
		)
		if err != nil {
			continue
		}
		poll.DateOptions = append(poll.DateOptions, do)
	}

	// Get votes
	voteRows, err := h.db.Query(ctx, `
		SELECT id, poll_id, date_option_id, user_id, user_name, response, created_at
		FROM votes WHERE poll_id = $1
		ORDER BY created_at
	`, pollID)

	if err == nil {
		defer voteRows.Close()
		for voteRows.Next() {
			var vote models.Vote
			err := voteRows.Scan(
				&vote.ID, &vote.PollID, &vote.DateOptionID, &vote.UserID, &vote.UserName, &vote.Response, &vote.CreatedAt,
			)
			if err != nil {
				continue
			}
			poll.Votes = append(poll.Votes, vote)
		}
	}

	return &poll, nil
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func formatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func formatICSDate(t time.Time) string {
	return t.Format("20060102T150405")
}
