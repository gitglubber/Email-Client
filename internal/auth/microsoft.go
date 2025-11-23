package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/gitglubber/email-client/internal/config"
	"github.com/gitglubber/email-client/internal/models"
	"github.com/google/uuid"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type MicrosoftAuth struct {
	config *oauth2.Config
}

func NewMicrosoftAuth(cfg *config.Config) *MicrosoftAuth {
	return &MicrosoftAuth{
		config: &oauth2.Config{
			ClientID:     cfg.MicrosoftClientID,
			ClientSecret: cfg.MicrosoftClientSecret,
			RedirectURL:  cfg.MicrosoftRedirectURL,
			Scopes: []string{
				"https://graph.microsoft.com/User.Read",
				"https://graph.microsoft.com/Mail.Read",
				"https://graph.microsoft.com/Mail.Send",
				"https://graph.microsoft.com/Mail.ReadWrite",
				"https://graph.microsoft.com/Calendars.ReadWrite",
				"offline_access",
			},
			Endpoint: microsoft.AzureADEndpoint("common"),
		},
	}
}

func (m *MicrosoftAuth) GetAuthURL(state string) string {
	return m.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (m *MicrosoftAuth) HandleCallback(ctx context.Context, code string) (*models.User, error) {
	token, err := m.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get user info using the token
	client := m.config.Client(ctx, token)
	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(
		&tokenProvider{token: token.AccessToken}, []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to create graph client: %w", err)
	}

	userInfo, err := graphClient.Me().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Email:        *userInfo.GetUserPrincipalName(),
		Name:         *userInfo.GetDisplayName(),
		Provider:     "microsoft",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenExpiry:  token.Expiry,
		CreatedAt:    time.Now(),
	}

	return user, nil
}

func (m *MicrosoftAuth) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := m.config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken, nil
}

func (m *MicrosoftAuth) GetGraphClient(accessToken string) (*msgraphsdk.GraphServiceClient, error) {
	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(
		&tokenProvider{token: accessToken}, []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to create graph client: %w", err)
	}

	return client, nil
}

// Simple token provider for Microsoft Graph SDK
type tokenProvider struct {
	token string
}

func (t *tokenProvider) GetAuthorizationToken(ctx context.Context, url string, additionalAuthenticationContext map[string]interface{}) (string, error) {
	return t.token, nil
}

func (t *tokenProvider) GetAllowedHostsValidator() *authentication.AllowedHostsValidator {
	return &authentication.AllowedHostsValidator{
		AllowedHosts: []string{"graph.microsoft.com"},
	}
}
