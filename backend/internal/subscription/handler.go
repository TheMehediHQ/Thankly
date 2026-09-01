package subscription

import (
	"database/sql"
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

func (h *Handler) Get(c *gin.Context) {
	userID := c.GetString("user_id")

	var id, plan, status string
	err := h.db.QueryRow(
		`SELECT id, plan, status FROM subscriptions 
		 WHERE user_id = $1 
		 ORDER BY created_at DESC LIMIT 1`,
		userID,
	).Scan(&id, &plan, &status)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{
			"plan":   "free",
			"status": "active",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     id,
		"plan":   plan,
		"status": status,
	})
}

func (h *Handler) Checkout(c *gin.Context) {
	// TODO: create Stripe checkout session
	c.JSON(http.StatusOK, gin.H{
		"checkout_url": "https://checkout.stripe.com/...",
	})
}

func (h *Handler) Cancel(c *gin.Context) {
	// TODO: cancel Stripe subscription
	c.JSON(http.StatusOK, gin.H{"message": "subscription cancelled"})
}

func (h *Handler) Webhook(c *gin.Context) {
	// TODO: verify Stripe webhook signature and process event
	c.JSON(http.StatusOK, gin.H{"received": true})
}
