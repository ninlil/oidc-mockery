package models

import "time"

// Client represents an OIDC client
type Client struct {
	ID           string
	Secret       string
	RedirectURIs []string
	Scopes       []string
}

// Persona represents a user persona for testing
type Persona struct {
	ID     string
	Name   string
	Email  string
	Claims map[string]interface{}
}

// AuthorizationCode represents an authorization code
type AuthorizationCode struct {
	Code        string
	ClientID    string
	PersonaID   string
	RedirectURI string
	Scopes      []string
	ExpiresAt   time.Time
	Used        bool
}

// AccessToken represents an access token
type AccessToken struct {
	Token     string
	ClientID  string
	PersonaID string
	Scopes    []string
	ExpiresAt time.Time
}

// IDToken represents an ID token
type IDToken struct {
	Token     string
	ClientID  string
	PersonaID string
	ExpiresAt time.Time
}
