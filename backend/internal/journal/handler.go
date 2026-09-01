package journal

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thankly/backend/internal/database"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Today(c *gin.Context) {
	userID := c.GetString("user_id")
	today := time.Now().Format("2006-01-02")

	rows, err := h.db.Query(
		`SELECT id, content, gratitude_date, created_at 
		 FROM gratitudes 
		 WHERE user_id = $1 AND gratitude_date = $2 
		 ORDER BY created_at ASC`,
		userID, today,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	var entries []gin.H
	for rows.Next() {
		var id, content, gratitudeDate, createdAt string
		if err := rows.Scan(&id, &content, &gratitudeDate, &createdAt); err != nil {
			continue
		}
		entries = append(entries, gin.H{
			"id":             id,
			"content":        content,
			"gratitude_date": gratitudeDate,
			"created_at":     createdAt,
		})
	}

	if entries == nil {
		entries = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"date":    today,
		"entries": entries,
		"count":   len(entries),
		"limit":   3,
		"remaining": 3 - len(entries),
	})
}

func (h *Handler) History(c *gin.Context) {
	userID := c.GetString("user_id")

	rows, err := h.db.Query(
		`SELECT id, content, gratitude_date, created_at 
		 FROM gratitudes 
		 WHERE user_id = $1 
		 ORDER BY gratitude_date DESC, created_at ASC`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	var entries []gin.H
	for rows.Next() {
		var id, content, gratitudeDate, createdAt string
		if err := rows.Scan(&id, &content, &gratitudeDate, &createdAt); err != nil {
			continue
		}
		entries = append(entries, gin.H{
			"id":             id,
			"content":        content,
			"gratitude_date": gratitudeDate,
			"created_at":     createdAt,
		})
	}

	if entries == nil {
		entries = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{"history": entries})
}

func (h *Handler) ByDate(c *gin.Context) {
	userID := c.GetString("user_id")
	date := c.Param("date")

	if _, err := time.Parse("2006-01-02", date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	rows, err := h.db.Query(
		`SELECT id, content, gratitude_date, created_at 
		 FROM gratitudes 
		 WHERE user_id = $1 AND gratitude_date = $2 
		 ORDER BY created_at ASC`,
		userID, date,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	var entries []gin.H
	for rows.Next() {
		var id, content, gratitudeDate, createdAt string
		if err := rows.Scan(&id, &content, &gratitudeDate, &createdAt); err != nil {
			continue
		}
		entries = append(entries, gin.H{
			"id":             id,
			"content":        content,
			"gratitude_date": gratitudeDate,
			"created_at":     createdAt,
		})
	}

	if entries == nil {
		entries = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"date":    date,
		"entries": entries,
	})
}
