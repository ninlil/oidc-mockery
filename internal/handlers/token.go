package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ninlil/oidc-mockery/internal/config"
	"github.com/ninlil/oidc-mockery/internal/utils"
)

// TokenArgs represents an OIDC token request parameters
type TokenArgs struct {
	Body string `from:"body"`
}

// TokenResponse represents an OIDC token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token"`
}

// handleToken handles the token endpoint
func handleToken(cfg *config.Config, args *TokenArgs) (*TokenResponse, int, error) {
	// Extract form fields from body map
	bodyValues, err := url.ParseQuery(args.Body)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err)
	}
	grantType := bodyValues.Get("grant_type")
	code := bodyValues.Get("code")                // TODO: Validate authorization code
	redirectURI := bodyValues.Get("redirect_uri") // TODO: Validate redirect URI
	clientID := bodyValues.Get("client_id")
	clientSecret := bodyValues.Get("client_secret")

	_ = code        // Silence unused variable warning
	_ = redirectURI // Silence unused variable warning

	// Validate grant type
	if grantType != "authorization_code" {
		return nil, http.StatusBadRequest, fmt.Errorf("unsupported_grant_type")
	}

	// Validate client credentials
	client := cfg.GetClient(clientID)
	if client == nil || client.ClientSecret != clientSecret {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid_client")
	}

	// Validate authorization code
	authData, exists := authCodeStore[code]
	if !exists {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid_grant")
	}

	// Check if code has expired
	if time.Now().After(authData.ExpiresAt) {
		delete(authCodeStore, code) // Clean up expired code
		return nil, http.StatusBadRequest, fmt.Errorf("invalid_grant")
	}

	// Validate client ID matches
	if authData.ClientID != clientID {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid_grant")
	}

	// Validate redirect URI matches
	if authData.RedirectURI != redirectURI {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid_grant")
	}

	// Get the persona from the stored auth code
	persona := cfg.GetPersona(authData.PersonaID)
	if persona == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("server_error: persona not found")
	}

	// Remove used authorization code
	delete(authCodeStore, code)

	// Generate tokens
	accessToken := utils.GenerateAccessToken()
	idToken, err := utils.GenerateIDToken(cfg, persona, clientID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("server_error: %v", err)
	}

	response := &TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		IDToken:     idToken,
	}

	return response, http.StatusOK, nil
}
