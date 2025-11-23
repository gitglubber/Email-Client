package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gitglubber/email-client/internal/config"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	handler := NewHandler(cfg)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.GET("/google", handler.GoogleAuth)
		auth.GET("/google/callback", handler.GoogleCallback)
		auth.GET("/microsoft", handler.MicrosoftAuth)
		auth.GET("/microsoft/callback", handler.MicrosoftCallback)
		auth.GET("/logout", handler.Logout)
		auth.GET("/me", handler.AuthMiddleware(), handler.GetMe)
	}

	// API routes (protected)
	api := router.Group("/api")
	api.Use(handler.AuthMiddleware())
	{
		// Email endpoints
		emails := api.Group("/emails")
		{
			emails.GET("", handler.ListEmails)
			emails.GET("/:id", handler.GetEmail)
			emails.POST("/send", handler.SendEmail)
			emails.POST("/:id/reply", handler.ReplyToEmail)
			emails.POST("/:id/ai-reply", handler.AIGenerateReply)
			emails.POST("/:id/categorize", handler.AICategorizeEmail)
			emails.POST("/:id/spam-check", handler.AICheckSpam)
		}

		// Calendar endpoints
		cal := api.Group("/calendar")
		{
			cal.GET("/events", handler.ListCalendarEvents)
			cal.POST("/events", handler.CreateCalendarEvent)
			cal.POST("/ai-schedule", handler.AISuggestScheduling)
		}
	}

	return router
}
