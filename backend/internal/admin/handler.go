package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thankly/backend/internal/database"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) ListUsers(c *gin.Context) {
	rows, err := h.db.Query(
		`SELECT id, name, email, role, created_at FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	var users []gin.H
	for rows.Next() {
		var id, name, email, role, createdAt string
		if err := rows.Scan(&id, &name, &email, &role, &createdAt); err != nil {
			continue
		}
		users = append(users, gin.H{
			"id":         id,
			"name":       name,
			"email":      email,
			"role":       role,
			"created_at": createdAt,
		})
	}

	if users == nil {
		users = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *Handler) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")

	var input struct {
		Role string `json:"role" binding:"required,oneof=user admin super_admin"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role must be: user, admin, or super_admin"})
		return
	}

	callerRole := c.GetString("role")
	if input.Role == "super_admin" && callerRole != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only super_admin can assign super_admin role"})
		return
	}

	result, err := h.db.Exec(
		`UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`,
		input.Role, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role updated", "role": input.Role})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	callerID := c.GetString("user_id")
	if callerID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete yourself"})
		return
	}

	result, err := h.db.Exec(`DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func (h *Handler) Stats(c *gin.Context) {
	var totalUsers, totalGratitudes, premiumUsers int

	h.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&totalUsers)
	h.db.QueryRow(`SELECT COUNT(*) FROM gratitudes`).Scan(&totalGratitudes)
	h.db.QueryRow(`SELECT COUNT(*) FROM subscriptions WHERE plan = 'premium'`).Scan(&premiumUsers)

	c.JSON(http.StatusOK, gin.H{
		"total_users":      totalUsers,
		"total_gratitudes": totalGratitudes,
		"premium_users":    premiumUsers,
	})
}
