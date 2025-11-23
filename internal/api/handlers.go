package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitglubber/email-client/internal/ai"
	"github.com/gitglubber/email-client/internal/auth"
	"github.com/gitglubber/email-client/internal/calendar"
	"github.com/gitglubber/email-client/internal/config"
	"github.com/gitglubber/email-client/internal/email"
	"github.com/gitglubber/email-client/internal/models"
	"github.com/gitglubber/email-client/internal/store"
	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type Handler struct {
	cfg          *config.Config
	sessionStore *store.SessionStore
	googleAuth   *auth.GoogleAuth
	msAuth       *auth.MicrosoftAuth
	aiService    *ai.Service
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg:          cfg,
		sessionStore: store.NewSessionStore(),
		googleAuth:   auth.NewGoogleAuth(cfg),
		msAuth:       auth.NewMicrosoftAuth(cfg),
		aiService:    ai.NewService(cfg),
	}
}

// Auth middleware
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		user, exists := h.sessionStore.Get(sessionID)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session not found"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// Google Auth
func (h *Handler) GoogleAuth(c *gin.Context) {
	state := uuid.New().String()
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	url := h.googleAuth.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	savedState, _ := c.Cookie("oauth_state")

	if state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	code := c.Query("code")
	user, err := h.googleAuth.HandleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sessionID := uuid.New().String()
	h.sessionStore.Set(sessionID, user)

	c.SetCookie("session_id", sessionID, 86400*7, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, h.cfg.FrontendURL)
}

// Microsoft Auth
func (h *Handler) MicrosoftAuth(c *gin.Context) {
	state := uuid.New().String()
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	url := h.msAuth.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) MicrosoftCallback(c *gin.Context) {
	state := c.Query("state")
	savedState, _ := c.Cookie("oauth_state")

	if state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	code := c.Query("code")
	user, err := h.msAuth.HandleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sessionID := uuid.New().String()
	h.sessionStore.Set(sessionID, user)

	c.SetCookie("session_id", sessionID, 86400*7, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, h.cfg.FrontendURL)
}

// Logout
func (h *Handler) Logout(c *gin.Context) {
	sessionID, _ := c.Cookie("session_id")
	h.sessionStore.Delete(sessionID)
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// Get current user
func (h *Handler) GetMe(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	c.JSON(http.StatusOK, user)
}

// List emails
func (h *Handler) ListEmails(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	maxResults, _ := strconv.ParseInt(c.DefaultQuery("max", "50"), 10, 64)
	pageToken := c.Query("page_token")

	var emails []*models.Email
	var nextToken string
	var err error

	if user.Provider == "google" {
		gmailSrv, err := h.googleAuth.GetGmailService(c.Request.Context(), user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewGmailService(gmailSrv)
		emails, nextToken, err = emailSrv.ListEmails(c.Request.Context(), maxResults, pageToken)
	} else if user.Provider == "microsoft" {
		graphClient, err := h.msAuth.GetGraphClient(user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewMicrosoftEmailService(graphClient)
		emails, nextToken, err = emailSrv.ListEmails(c.Request.Context(), int(maxResults), pageToken)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"emails":     emails,
		"next_token": nextToken,
	})
}

// Get single email
func (h *Handler) GetEmail(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	emailID := c.Param("id")

	var emailData *models.Email
	var err error

	if user.Provider == "google" {
		gmailSrv, err := h.googleAuth.GetGmailService(c.Request.Context(), user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewGmailService(gmailSrv)
		emailData, err = emailSrv.GetEmail(c.Request.Context(), emailID)
	} else if user.Provider == "microsoft" {
		graphClient, err := h.msAuth.GetGraphClient(user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewMicrosoftEmailService(graphClient)
		emailData, err = emailSrv.GetEmail(c.Request.Context(), emailID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, emailData)
}

// Send email
func (h *Handler) SendEmail(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var req models.EmailSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error

	if user.Provider == "google" {
		gmailSrv, err := h.googleAuth.GetGmailService(c.Request.Context(), user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewGmailService(gmailSrv)
		err = emailSrv.SendEmail(c.Request.Context(), &req)
	} else if user.Provider == "microsoft" {
		graphClient, err := h.msAuth.GetGraphClient(user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewMicrosoftEmailService(graphClient)
		err = emailSrv.SendEmail(c.Request.Context(), &req)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email sent"})
}

// Reply to email
func (h *Handler) ReplyToEmail(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	emailID := c.Param("id")

	var req struct {
		Body string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error

	if user.Provider == "google" {
		gmailSrv, err := h.googleAuth.GetGmailService(c.Request.Context(), user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewGmailService(gmailSrv)
		err = emailSrv.ReplyToEmail(c.Request.Context(), emailID, req.Body)
	} else if user.Provider == "microsoft" {
		graphClient, err := h.msAuth.GetGraphClient(user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewMicrosoftEmailService(graphClient)
		err = emailSrv.ReplyToEmail(c.Request.Context(), emailID, req.Body)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reply sent"})
}

// AI: Generate reply
func (h *Handler) AIGenerateReply(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	emailID := c.Param("id")

	var req struct {
		Context string `json:"context"`
		Tone    string `json:"tone"`
	}
	c.ShouldBindJSON(&req)

	// Get the email first
	var emailData *models.Email
	var err error

	if user.Provider == "google" {
		gmailSrv, err := h.googleAuth.GetGmailService(c.Request.Context(), user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewGmailService(gmailSrv)
		emailData, err = emailSrv.GetEmail(c.Request.Context(), emailID)
	} else if user.Provider == "microsoft" {
		graphClient, err := h.msAuth.GetGraphClient(user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewMicrosoftEmailService(graphClient)
		emailData, err = emailSrv.GetEmail(c.Request.Context(), emailID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	reply, err := h.aiService.GenerateReply(c.Request.Context(), emailData, req.Context, req.Tone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reply)
}

// AI: Categorize email
func (h *Handler) AICategorizeEmail(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	emailID := c.Param("id")

	var emailData *models.Email
	var err error

	if user.Provider == "google" {
		gmailSrv, err := h.googleAuth.GetGmailService(c.Request.Context(), user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewGmailService(gmailSrv)
		emailData, err = emailSrv.GetEmail(c.Request.Context(), emailID)
	} else if user.Provider == "microsoft" {
		graphClient, err := h.msAuth.GetGraphClient(user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewMicrosoftEmailService(graphClient)
		emailData, err = emailSrv.GetEmail(c.Request.Context(), emailID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	category, err := h.aiService.CategorizeEmail(c.Request.Context(), emailData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// AI: Check spam
func (h *Handler) AICheckSpam(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	emailID := c.Param("id")

	var emailData *models.Email
	var err error

	if user.Provider == "google" {
		gmailSrv, err := h.googleAuth.GetGmailService(c.Request.Context(), user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewGmailService(gmailSrv)
		emailData, err = emailSrv.GetEmail(c.Request.Context(), emailID)
	} else if user.Provider == "microsoft" {
		graphClient, err := h.msAuth.GetGraphClient(user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailSrv := email.NewMicrosoftEmailService(graphClient)
		emailData, err = emailSrv.GetEmail(c.Request.Context(), emailID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	spamCheck, err := h.aiService.CheckSpam(c.Request.Context(), emailData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, spamCheck)
}

// List calendar events
func (h *Handler) ListCalendarEvents(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	startDate := time.Now()
	endDate := time.Now().AddDate(0, 1, 0) // Next month

	if start := c.Query("start"); start != "" {
		startDate, _ = time.Parse(time.RFC3339, start)
	}
	if end := c.Query("end"); end != "" {
		endDate, _ = time.Parse(time.RFC3339, end)
	}

	var events []*models.CalendarEvent
	var err error

	if user.Provider == "google" {
		gmailSrv, err := h.googleAuth.GetGmailService(c.Request.Context(), user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create calendar service from same client
		calSrv, err := calendar.NewService(c.Request.Context(), option.WithHTTPClient(gmailSrv.Client()))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		calendarSrv := calendar.NewGoogleCalendarService(calSrv)
		events, err = calendarSrv.ListEvents(c.Request.Context(), startDate, endDate)
	} else if user.Provider == "microsoft" {
		graphClient, err := h.msAuth.GetGraphClient(user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		calendarSrv := calendar.NewMicrosoftCalendarService(graphClient)
		events, err = calendarSrv.ListEvents(c.Request.Context(), startDate, endDate)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// Create calendar event
func (h *Handler) CreateCalendarEvent(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var req models.CalendarEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var event *models.CalendarEvent
	var err error

	if user.Provider == "google" {
		gmailSrv, err := h.googleAuth.GetGmailService(c.Request.Context(), user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		calSrv, err := calendar.NewService(c.Request.Context(), option.WithHTTPClient(gmailSrv.Client()))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		calendarSrv := calendar.NewGoogleCalendarService(calSrv)
		event, err = calendarSrv.CreateEvent(c.Request.Context(), &req)
	} else if user.Provider == "microsoft" {
		graphClient, err := h.msAuth.GetGraphClient(user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		calendarSrv := calendar.NewMicrosoftCalendarService(graphClient)
		event, err = calendarSrv.CreateEvent(c.Request.Context(), &req)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// AI: Suggest scheduling
func (h *Handler) AISuggestScheduling(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var req models.AIScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing events for context
	startDate := time.Now()
	endDate := time.Now().AddDate(0, 0, 7) // Next week

	var events []*models.CalendarEvent
	var err error

	if user.Provider == "google" {
		gmailSrv, err := h.googleAuth.GetGmailService(c.Request.Context(), user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		calSrv, err := calendar.NewService(c.Request.Context(), option.WithHTTPClient(gmailSrv.Client()))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		calendarSrv := calendar.NewGoogleCalendarService(calSrv)
		events, err = calendarSrv.ListEvents(c.Request.Context(), startDate, endDate)
	} else if user.Provider == "microsoft" {
		graphClient, err := h.msAuth.GetGraphClient(user.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		calendarSrv := calendar.NewMicrosoftCalendarService(graphClient)
		events, err = calendarSrv.ListEvents(c.Request.Context(), startDate, endDate)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	suggestion, err := h.aiService.SuggestScheduling(c.Request.Context(), &req, events)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, suggestion)
}
