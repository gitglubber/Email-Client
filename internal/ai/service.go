package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gitglubber/email-client/internal/config"
	"github.com/gitglubber/email-client/internal/models"
	openai "github.com/sashabaranov/go-openai"
)

type Service struct {
	client *openai.Client
	model  string
}

func NewService(cfg *config.Config) *Service {
	config := openai.DefaultConfig(cfg.OpenAIAPIKey)
	if cfg.OpenAIBaseURL != "" {
		config.BaseURL = cfg.OpenAIBaseURL
	}

	client := openai.NewClientWithConfig(config)

	return &Service{
		client: client,
		model:  cfg.OpenAIModel,
	}
}

// CategorizeEmail uses AI to categorize an email into appropriate folders/labels
func (s *Service) CategorizeEmail(ctx context.Context, email *models.Email) (*models.AICategoryResponse, error) {
	prompt := fmt.Sprintf(`Analyze the following email and categorize it. Provide a primary category, confidence score (0-1), reasoning, and suggested labels.

Categories: Work, Personal, Finance, Shopping, Social, Promotions, Updates, Important, Other

Email Details:
From: %s
Subject: %s
Body: %s

Respond in JSON format:
{
  "category": "category name",
  "confidence": 0.0-1.0,
  "reasoning": "brief explanation",
  "suggested_labels": ["label1", "label2"]
}`, email.From, email.Subject, truncateText(email.Body, 500))

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an AI assistant that helps categorize emails. Always respond with valid JSON only, no additional text.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to categorize email: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	var result models.AICategoryResponse
	content := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Remove markdown code blocks if present
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &result, nil
}

// GenerateReply generates an AI-powered email reply
func (s *Service) GenerateReply(ctx context.Context, email *models.Email, userContext string, tone string) (*models.AIReplyResponse, error) {
	if tone == "" {
		tone = "professional"
	}

	prompt := fmt.Sprintf(`Generate a %s reply to the following email.

Original Email:
From: %s
Subject: %s
Body: %s

Additional Context: %s

Generate:
1. A complete email reply (2-4 paragraphs)
2. Three alternative shorter suggestions

Respond in JSON format:
{
  "generated_reply": "full email reply text",
  "suggestions": ["suggestion 1", "suggestion 2", "suggestion 3"]
}`, tone, email.From, email.Subject, truncateText(email.Body, 800), userContext)

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an AI email assistant that helps write professional, clear, and effective email replies. Always respond with valid JSON only.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.7,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate reply: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	var result models.AIReplyResponse
	content := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Remove markdown code blocks if present
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &result, nil
}

// CheckSpam uses AI to determine if an email is spam
func (s *Service) CheckSpam(ctx context.Context, email *models.Email) (*models.AISpamCheckResponse, error) {
	prompt := fmt.Sprintf(`Analyze the following email and determine if it's spam. Consider:
- Suspicious sender addresses
- Phishing attempts
- Unrealistic offers
- Poor grammar/spelling
- Urgency tactics
- Suspicious links

Email Details:
From: %s
Subject: %s
Body: %s

Respond in JSON format:
{
  "is_spam": true/false,
  "confidence": 0.0-1.0,
  "reasoning": "brief explanation of why this is or isn't spam"
}`, email.From, email.Subject, truncateText(email.Body, 600))

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an AI spam detection system. Analyze emails for spam indicators. Always respond with valid JSON only.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.2,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to check spam: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	var result models.AISpamCheckResponse
	content := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Remove markdown code blocks if present
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &result, nil
}

// SuggestScheduling uses AI to suggest optimal meeting times
func (s *Service) SuggestScheduling(ctx context.Context, request *models.AIScheduleRequest, existingEvents []*models.CalendarEvent) (*models.AIScheduleResponse, error) {
	// Build context of existing events
	eventsContext := "Existing calendar events:\n"
	for _, event := range existingEvents {
		eventsContext += fmt.Sprintf("- %s: %s to %s\n",
			event.Title,
			event.Start.Format("2006-01-02 15:04"),
			event.End.Format("2006-01-02 15:04"))
	}

	duration := request.Duration
	if duration == 0 {
		duration = 60 // default 1 hour
	}

	prompt := fmt.Sprintf(`Based on the user's calendar and the meeting request, suggest 3-5 optimal meeting times for the next 7 days.

%s

Meeting Request:
Description: %s
Duration: %d minutes
Attendees: %s

Consider:
- Working hours (9 AM - 6 PM)
- Avoid conflicts with existing events
- Leave buffer time between meetings
- Prefer mid-morning or mid-afternoon slots

Respond in JSON format with suggested times:
{
  "suggested_times": ["2024-01-15T10:00:00Z", "2024-01-15T14:00:00Z", "2024-01-16T11:00:00Z"],
  "reasoning": "explanation of why these times were chosen"
}`, eventsContext, request.Description, duration, strings.Join(request.Attendees, ", "))

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an AI calendar assistant that helps schedule meetings intelligently. Always respond with valid JSON only.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.5,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to suggest scheduling: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	var result struct {
		SuggestedTimes []string `json:"suggested_times"`
		Reasoning      string   `json:"reasoning"`
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Remove markdown code blocks if present
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Convert string times to time.Time
	times := make([]time.Time, 0, len(result.SuggestedTimes))
	for _, timeStr := range result.SuggestedTimes {
		t, err := time.Parse(time.RFC3339, timeStr)
		if err == nil {
			times = append(times, t)
		}
	}

	return &models.AIScheduleResponse{
		SuggestedTimes: times,
		Reasoning:      result.Reasoning,
	}, nil
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}
