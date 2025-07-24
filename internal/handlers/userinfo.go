package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ninlil/oidc-mockery/internal/config"
)

// UserInfoArgs represents the userinfo endpoint parameters
type UserInfoArgs struct {
	Authorization string `json:"Authorization" from:"header" required:"true"`
}

// UserInfoResponse represents the userinfo endpoint response
type UserInfoResponse struct {
	Sub               string `json:"sub"`
	Name              string `json:"name,omitempty"`
	Email             string `json:"email,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
}

// handleUserInfo handles the userinfo endpoint
func handleUserInfo(cfg *config.Config, args *UserInfoArgs) (*UserInfoResponse, int, error) {
	// Extract bearer token from Authorization header
	if args.Authorization == "" || !strings.HasPrefix(args.Authorization, "Bearer ") {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid_token")
	}

	token := strings.TrimPrefix(args.Authorization, "Bearer ")
	_ = token // TODO: Validate access token

	// TODO: Validate access token and retrieve associated persona
	// For MVP, we'll return the first persona's claims
	if len(cfg.Personas) == 0 {
		return nil, http.StatusInternalServerError, fmt.Errorf("no_personas_configured")
	}

	persona := cfg.Personas[0]

	response := &UserInfoResponse{
		Sub:   persona.ID,
		Name:  persona.Name,
		Email: persona.Email,
	}

	// Add additional claims from persona
	if givenName, ok := persona.Claims["given_name"].(string); ok {
		response.GivenName = givenName
	}
	if familyName, ok := persona.Claims["family_name"].(string); ok {
		response.FamilyName = familyName
	}
	if username, ok := persona.Claims["preferred_username"].(string); ok {
		response.PreferredUsername = username
	}

	return response, http.StatusOK, nil
}
