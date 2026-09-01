package user

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thankly/backend/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

type UpdateProfileInput struct {
	Name string `json:"name" binding:"omitempty,min=2,max=100"`
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var id, name, email, role, createdAt string
	err := h.db.QueryRow(
		`SELECT id, name, email, role, created_at FROM users WHERE id = $1`,
		userID,
	).Scan(&id, &name, &email, &role, &createdAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         id,
		"name":       name,
		"email":      email,
		"role":       role,
		"created_at": createdAt,
	})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Exec(
		`UPDATE users SET name = $1, updated_at = NOW() WHERE id = $2`,
		input.Name, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")

	var input ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var currentHash string
	err := h.db.QueryRow(
		`SELECT password_hash FROM users WHERE id = $1`,
		userID,
	).Scan(&currentHash)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(input.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current password is incorrect"})
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	_, err = h.db.Exec(
		`UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`,
		string(newHash), userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}
