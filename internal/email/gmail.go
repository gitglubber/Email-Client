package email

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gitglubber/email-client/internal/models"
	"google.golang.org/api/gmail/v1"
)

type GmailService struct {
	service *gmail.Service
}

func NewGmailService(service *gmail.Service) *GmailService {
	return &GmailService{service: service}
}

func (g *GmailService) ListEmails(ctx context.Context, maxResults int64, pageToken string) ([]*models.Email, string, error) {
	call := g.service.Users.Messages.List("me").MaxResults(maxResults)
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}

	resp, err := call.Do()
	if err != nil {
		return nil, "", fmt.Errorf("failed to list emails: %w", err)
	}

	emails := make([]*models.Email, 0, len(resp.Messages))
	for _, msg := range resp.Messages {
		email, err := g.GetEmail(ctx, msg.Id)
		if err != nil {
			continue // Skip errors for individual emails
		}
		emails = append(emails, email)
	}

	return emails, resp.NextPageToken, nil
}

func (g *GmailService) GetEmail(ctx context.Context, id string) (*models.Email, error) {
	msg, err := g.service.Users.Messages.Get("me", id).Format("full").Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get email: %w", err)
	}

	email := &models.Email{
		ID:       msg.Id,
		ThreadID: msg.ThreadId,
		Labels:   msg.LabelIds,
	}

	// Parse headers
	for _, header := range msg.Payload.Headers {
		switch header.Name {
		case "From":
			email.From = header.Value
		case "To":
			email.To = strings.Split(header.Value, ",")
		case "Subject":
			email.Subject = header.Value
		case "Date":
			email.Date, _ = time.Parse(time.RFC1123Z, header.Value)
		}
	}

	// Parse body
	email.Body, email.BodyHTML = g.parseBody(msg.Payload)

	// Check if read
	for _, label := range msg.LabelIds {
		if label == "UNREAD" {
			email.IsRead = false
			break
		} else {
			email.IsRead = true
		}
		if label == "STARRED" {
			email.IsStarred = true
		}
		if label == "SPAM" {
			email.IsSpam = true
		}
	}

	return email, nil
}

func (g *GmailService) parseBody(payload *gmail.MessagePart) (string, string) {
	var plainText, htmlText string

	if payload.MimeType == "text/plain" && payload.Body.Data != "" {
		data, _ := base64.URLEncoding.DecodeString(payload.Body.Data)
		plainText = string(data)
	} else if payload.MimeType == "text/html" && payload.Body.Data != "" {
		data, _ := base64.URLEncoding.DecodeString(payload.Body.Data)
		htmlText = string(data)
	}

	// Check parts for multipart messages
	for _, part := range payload.Parts {
		if part.MimeType == "text/plain" && part.Body.Data != "" {
			data, _ := base64.URLEncoding.DecodeString(part.Body.Data)
			plainText = string(data)
		} else if part.MimeType == "text/html" && part.Body.Data != "" {
			data, _ := base64.URLEncoding.DecodeString(part.Body.Data)
			htmlText = string(data)
		}
	}

	return plainText, htmlText
}

func (g *GmailService) SendEmail(ctx context.Context, req *models.EmailSendRequest) error {
	var message strings.Builder
	message.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(req.To, ",")))
	if len(req.Cc) > 0 {
		message.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(req.Cc, ",")))
	}
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", req.Subject))
	if req.IsHTML {
		message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}
	message.WriteString("\r\n")
	message.WriteString(req.Body)

	encoded := base64.URLEncoding.EncodeToString([]byte(message.String()))

	msg := &gmail.Message{
		Raw: encoded,
	}

	_, err := g.service.Users.Messages.Send("me", msg).Do()
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (g *GmailService) ReplyToEmail(ctx context.Context, emailID string, body string) error {
	// Get original email to construct reply
	original, err := g.service.Users.Messages.Get("me", emailID).Format("full").Do()
	if err != nil {
		return fmt.Errorf("failed to get original email: %w", err)
	}

	var to, subject string
	for _, header := range original.Payload.Headers {
		if header.Name == "From" {
			to = header.Value
		}
		if header.Name == "Subject" {
			subject = header.Value
			if !strings.HasPrefix(strings.ToLower(subject), "re:") {
				subject = "Re: " + subject
			}
		}
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString(fmt.Sprintf("In-Reply-To: %s\r\n", original.Id))
	message.WriteString(fmt.Sprintf("References: %s\r\n", original.Id))
	message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	message.WriteString("\r\n")
	message.WriteString(body)

	encoded := base64.URLEncoding.EncodeToString([]byte(message.String()))

	msg := &gmail.Message{
		Raw:      encoded,
		ThreadId: original.ThreadId,
	}

	_, err = g.service.Users.Messages.Send("me", msg).Do()
	if err != nil {
		return fmt.Errorf("failed to send reply: %w", err)
	}

	return nil
}
