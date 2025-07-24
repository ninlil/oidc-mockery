package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/ninlil/oidc-mockery/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// IDTokenClaims represents the claims in an ID token
type IDTokenClaims struct {
	jwt.RegisteredClaims
	Name              string `json:"name,omitempty"`
	Email             string `json:"email,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
}

// GenerateIDToken creates a JWT ID token
func GenerateIDToken(cfg *config.Config, persona *config.Persona, clientID string) (string, error) {
	// TODO: Use proper RSA private key for signing
	// For MVP, we'll use a symmetric key
	signingKey := []byte("mock-signing-key-change-in-production")

	claims := IDTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Subject:   persona.ID,
			Audience:  []string{clientID},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Name:  persona.Name,
		Email: persona.Email,
	}

	// Add additional claims from persona
	if givenName, ok := persona.Claims["given_name"].(string); ok {
		claims.GivenName = givenName
	}
	if familyName, ok := persona.Claims["family_name"].(string); ok {
		claims.FamilyName = familyName
	}
	if username, ok := persona.Claims["preferred_username"].(string); ok {
		claims.PreferredUsername = username
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

// GenerateAccessToken creates a random access token
func GenerateAccessToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

// GenerateAuthCode creates a random authorization code
func GenerateAuthCode() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

// ValidateRedirectURI validates if the redirect URI is allowed for the client
func ValidateRedirectURI(redirectURI string, allowedURIs []string) bool {
	for _, allowed := range allowedURIs {
		if redirectURI == allowed {
			return true
		}
	}
	return false
}

// GetMockRSAModulus returns a placeholder RSA modulus for JWKS
func GetMockRSAModulus() string {
	// This is a placeholder - in production, use actual RSA public key modulus
	return "mock-rsa-modulus-base64url-encoded"
}
