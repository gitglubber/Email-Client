package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Provider     string    `json:"provider"` // "google" or "microsoft"
	AccessToken  string    `json:"-"`
	RefreshToken string    `json:"-"`
	TokenExpiry  time.Time `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Email struct {
	ID          string    `json:"id"`
	ThreadID    string    `json:"thread_id"`
	From        string    `json:"from"`
	To          []string  `json:"to"`
	Subject     string    `json:"subject"`
	Body        string    `json:"body"`
	BodyHTML    string    `json:"body_html"`
	Date        time.Time `json:"date"`
	IsRead      bool      `json:"is_read"`
	IsStarred   bool      `json:"is_starred"`
	Labels      []string  `json:"labels"`
	Attachments []string  `json:"attachments"`
	AICategory  string    `json:"ai_category,omitempty"`
	IsSpam      bool      `json:"is_spam"`
}

type EmailSendRequest struct {
	To      []string `json:"to" binding:"required"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
	Subject string   `json:"subject" binding:"required"`
	Body    string   `json:"body" binding:"required"`
	IsHTML  bool     `json:"is_html"`
}

type AIReplyRequest struct {
	EmailID string `json:"email_id" binding:"required"`
	Context string `json:"context"`
	Tone    string `json:"tone"` // professional, casual, friendly, etc.
}

type AIReplyResponse struct {
	GeneratedReply string `json:"generated_reply"`
	Suggestions    []string `json:"suggestions"`
}

type AICategoryResponse struct {
	Category    string  `json:"category"`
	Confidence  float64 `json:"confidence"`
	Reasoning   string  `json:"reasoning"`
	SuggestedLabels []string `json:"suggested_labels"`
}

type AISpamCheckResponse struct {
	IsSpam     bool    `json:"is_spam"`
	Confidence float64 `json:"confidence"`
	Reasoning  string  `json:"reasoning"`
}

type CalendarEvent struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Location    string    `json:"location"`
	Attendees   []string  `json:"attendees"`
	Organizer   string    `json:"organizer"`
	IsAllDay    bool      `json:"is_all_day"`
}

type CalendarEventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Start       time.Time `json:"start" binding:"required"`
	End         time.Time `json:"end" binding:"required"`
	Location    string    `json:"location"`
	Attendees   []string  `json:"attendees"`
	IsAllDay    bool      `json:"is_all_day"`
}

type AIScheduleRequest struct {
	Description string   `json:"description" binding:"required"`
	Attendees   []string `json:"attendees"`
	Duration    int      `json:"duration"` // in minutes
}

type AIScheduleResponse struct {
	SuggestedTimes []time.Time `json:"suggested_times"`
	Event          *CalendarEvent `json:"event,omitempty"`
	Reasoning      string `json:"reasoning"`
}
