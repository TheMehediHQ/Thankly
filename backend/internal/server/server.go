package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thankly/backend/internal/database"
	"github.com/thankly/backend/internal/middleware"
	"github.com/thankly/backend/internal/auth"
	"github.com/thankly/backend/internal/gratitude"
	"github.com/thankly/backend/internal/journal"
	"github.com/thankly/backend/internal/subscription"
	"github.com/thankly/backend/internal/user"
)

type Server struct {
	router *gin.Engine
	db     *database.DB
}

func New() (*Server, error) {
	db, err := database.New()
	if err != nil {
		return nil, err
	}

	router := gin.Default()

	router.Use(middleware.CORS())
	router.Use(middleware.SecureHeaders())

	s := &Server{
		router: router,
		db:     db,
	}

	s.setupRoutes()

	return s, nil
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api")

	authHandler := auth.NewHandler(s.db)
	userHandler := user.NewHandler(s.db)
	gratitudeHandler := gratitude.NewHandler(s.db)
	journalHandler := journal.NewHandler(s.db)
	subscriptionHandler := subscription.NewHandler(s.db)

	// Auth routes
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.POST("/logout", middleware.AuthRequired(), authHandler.Logout)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthRequired())
	{
		// User
		protected.GET("/me", userHandler.GetProfile)
		protected.PATCH("/me", userHandler.UpdateProfile)
		protected.PATCH("/me/password", userHandler.ChangePassword)

		// Gratitudes
		protected.POST("/gratitudes", gratitudeHandler.Create)
		protected.GET("/gratitudes", gratitudeHandler.List)
		protected.GET("/gratitudes/:id", gratitudeHandler.Get)
		protected.PATCH("/gratitudes/:id", gratitudeHandler.Update)
		protected.DELETE("/gratitudes/:id", gratitudeHandler.Delete)

		// Journal
		protected.GET("/journal/today", journalHandler.Today)
		protected.GET("/journal/history", journalHandler.History)
		protected.GET("/journal/:date", journalHandler.ByDate)

		// Subscription
		protected.GET("/subscription", subscriptionHandler.Get)
		protected.POST("/subscription/checkout", subscriptionHandler.Checkout)
		protected.POST("/subscription/cancel", subscriptionHandler.Cancel)
	}

	// Webhooks (no auth, verified by signature)
	api.POST("/webhooks/payment", subscriptionHandler.Webhook)

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})
}

func (s *Server) Run(addr string) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.db.Close()
	return nil
}
