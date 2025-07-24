package handlers

import (
	"github.com/ninlil/oidc-mockery/internal/config"
	"github.com/ninlil/oidc-mockery/internal/utils"
)

// JWKSResponse represents the JWKS endpoint response
type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// handleJWKS handles the JWKS endpoint
func handleJWKS(cfg *config.Config) *JWKSResponse {
	// Use configuration values with defaults
	keyType := cfg.JWKS.KeyType
	if keyType == "" {
		keyType = "RSA"
	}

	keyUse := cfg.JWKS.KeyUse
	if keyUse == "" {
		keyUse = "sig"
	}

	keyID := cfg.JWKS.KeyID
	if keyID == "" {
		keyID = "mock-key-id"
	}

	rsaModulus := cfg.JWKS.RSAModulus
	if rsaModulus == "" {
		rsaModulus = utils.GetMockRSAModulus()
	}

	rsaExponent := cfg.JWKS.RSAExponent
	if rsaExponent == "" {
		rsaExponent = "AQAB" // 65537 in base64url
	}

	jwk := JWK{
		Kty: keyType,
		Use: keyUse,
		Kid: keyID,
		N:   rsaModulus,
		E:   rsaExponent,
	}

	return &JWKSResponse{
		Keys: []JWK{jwk},
	}
}
