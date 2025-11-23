package email

import (
	"context"
	"fmt"

	"github.com/gitglubber/email-client/internal/models"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

type MicrosoftEmailService struct {
	client *msgraphsdk.GraphServiceClient
}

func NewMicrosoftEmailService(client *msgraphsdk.GraphServiceClient) *MicrosoftEmailService {
	return &MicrosoftEmailService{client: client}
}

func (m *MicrosoftEmailService) ListEmails(ctx context.Context, maxResults int, skipToken string) ([]*models.Email, string, error) {
	top := int32(maxResults)
	config := &users.ItemMessagesRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.ItemMessagesRequestBuilderGetQueryParameters{
			Top: &top,
		},
	}

	if skipToken != "" {
		config.QueryParameters.Skip = &skipToken
	}

	messages, err := m.client.Me().Messages().Get(ctx, config)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list emails: %w", err)
	}

	emails := make([]*models.Email, 0)
	for _, msg := range messages.GetValue() {
		email := m.convertToEmail(msg)
		emails = append(emails, email)
	}

	var nextToken string
	if messages.GetOdataNextLink() != nil {
		nextToken = *messages.GetOdataNextLink()
	}

	return emails, nextToken, nil
}

func (m *MicrosoftEmailService) GetEmail(ctx context.Context, id string) (*models.Email, error) {
	msg, err := m.client.Me().Messages().ByMessageId(id).Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get email: %w", err)
	}

	return m.convertToEmail(msg), nil
}

func (m *MicrosoftEmailService) convertToEmail(msg graphmodels.Messageable) *models.Email {
	email := &models.Email{
		ID:      *msg.GetId(),
		Subject: *msg.GetSubject(),
		IsRead:  *msg.GetIsRead(),
	}

	if msg.GetReceivedDateTime() != nil {
		email.Date = *msg.GetReceivedDateTime()
	}

	if msg.GetFrom() != nil && msg.GetFrom().GetEmailAddress() != nil {
		email.From = *msg.GetFrom().GetEmailAddress().GetAddress()
	}

	toRecipients := msg.GetToRecipients()
	email.To = make([]string, 0, len(toRecipients))
	for _, recipient := range toRecipients {
		if recipient.GetEmailAddress() != nil {
			email.To = append(email.To, *recipient.GetEmailAddress().GetAddress())
		}
	}

	if msg.GetBody() != nil {
		bodyContent := msg.GetBody().GetContent()
		if bodyContent != nil {
			if *msg.GetBody().GetContentType() == graphmodels.HTML_BODYTYPE {
				email.BodyHTML = *bodyContent
			} else {
				email.Body = *bodyContent
			}
		}
	}

	if msg.GetFlag() != nil && msg.GetFlag().GetFlagStatus() != nil {
		flagged := *msg.GetFlag().GetFlagStatus() == graphmodels.FLAGGED_FOLLOWUPFLAGSTATUS
		email.IsStarred = flagged
	}

	return email
}

func (m *MicrosoftEmailService) SendEmail(ctx context.Context, req *models.EmailSendRequest) error {
	message := graphmodels.NewMessage()
	subject := req.Subject
	message.SetSubject(&subject)

	// Set body
	body := graphmodels.NewItemBody()
	bodyContent := req.Body
	body.SetContent(&bodyContent)
	if req.IsHTML {
		contentType := graphmodels.HTML_BODYTYPE
		body.SetContentType(&contentType)
	} else {
		contentType := graphmodels.TEXT_BODYTYPE
		body.SetContentType(&contentType)
	}
	message.SetBody(body)

	// Set recipients
	toRecipients := make([]graphmodels.Recipientable, 0, len(req.To))
	for _, addr := range req.To {
		recipient := graphmodels.NewRecipient()
		emailAddr := graphmodels.NewEmailAddress()
		address := addr
		emailAddr.SetAddress(&address)
		recipient.SetEmailAddress(emailAddr)
		toRecipients = append(toRecipients, recipient)
	}
	message.SetToRecipients(toRecipients)

	// Send
	err := m.client.Me().SendMail().Post(ctx, &users.ItemSendMailPostRequestBody{
		Message: message,
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (m *MicrosoftEmailService) ReplyToEmail(ctx context.Context, emailID string, body string) error {
	comment := body
	replyBody := &users.ItemMessagesItemReplyPostRequestBody{}
	replyBody.SetComment(&comment)

	err := m.client.Me().Messages().ByMessageId(emailID).Reply().Post(ctx, replyBody, nil)
	if err != nil {
		return fmt.Errorf("failed to reply to email: %w", err)
	}

	return nil
}
