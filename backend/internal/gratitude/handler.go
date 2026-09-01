package gratitude

import (
	"database/sql"
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

type CreateInput struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
}

type UpdateInput struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
}

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	today := time.Now().Format("2006-01-02")

	var count int
	err := h.db.QueryRow(
		`SELECT COUNT(*) FROM gratitudes WHERE user_id = $1 AND gratitude_date = $2`,
		userID, today,
	).Scan(&count)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if count >= 3 {
		c.JSON(http.StatusConflict, gin.H{"error": "daily limit of 3 gratitudes reached"})
		return
	}

	var id string
	err = h.db.QueryRow(
		`INSERT INTO gratitudes (user_id, content, gratitude_date) VALUES ($1, $2, $3) RETURNING id`,
		userID, input.Content, today,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create gratitude"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":             id,
		"content":        input.Content,
		"gratitude_date": today,
		"created_at":     time.Now().Format(time.RFC3339),
	})
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	rows, err := h.db.Query(
		`SELECT id, content, gratitude_date, created_at, updated_at 
		 FROM gratitudes WHERE user_id = $1 
		 ORDER BY gratitude_date DESC, created_at DESC`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	var gratitudes []gin.H
	for rows.Next() {
		var id, content, gratitudeDate, createdAt, updatedAt string
		if err := rows.Scan(&id, &content, &gratitudeDate, &createdAt, &updatedAt); err != nil {
			continue
		}
		gratitudes = append(gratitudes, gin.H{
			"id":             id,
			"content":        content,
			"gratitude_date": gratitudeDate,
			"created_at":     createdAt,
			"updated_at":     updatedAt,
		})
	}

	if gratitudes == nil {
		gratitudes = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{"gratitudes": gratitudes})
}

func (h *Handler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	gratitudeID := c.Param("id")

	var id, content, gratitudeDate, createdAt, updatedAt, ownerID string
	err := h.db.QueryRow(
		`SELECT id, user_id, content, gratitude_date, created_at, updated_at 
		 FROM gratitudes WHERE id = $1`,
		gratitudeID,
	).Scan(&id, &ownerID, &content, &gratitudeDate, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "gratitude not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             id,
		"content":        content,
		"gratitude_date": gratitudeDate,
		"created_at":     createdAt,
		"updated_at":     updatedAt,
	})
}

func (h *Handler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	gratitudeID := c.Param("id")

	var input UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ownerID string
	err := h.db.QueryRow(
		`SELECT user_id FROM gratitudes WHERE id = $1`,
		gratitudeID,
	).Scan(&ownerID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "gratitude not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	_, err = h.db.Exec(
		`UPDATE gratitudes SET content = $1, updated_at = NOW() WHERE id = $2`,
		input.Content, gratitudeID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	gratitudeID := c.Param("id")

	var ownerID string
	err := h.db.QueryRow(
		`SELECT user_id FROM gratitudes WHERE id = $1`,
		gratitudeID,
	).Scan(&ownerID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "gratitude not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	_, err = h.db.Exec(`DELETE FROM gratitudes WHERE id = $1`, gratitudeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
