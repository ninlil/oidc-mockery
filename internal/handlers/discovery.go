package handlers

import (
	"github.com/ninlil/oidc-mockery/internal/config"
)

// DiscoveryResponse represents the OIDC discovery document
type DiscoveryResponse struct {
	Issuer                  string   `json:"issuer"`
	AuthorizationEndpoint   string   `json:"authorization_endpoint"`
	TokenEndpoint           string   `json:"token_endpoint"`
	UserInfoEndpoint        string   `json:"userinfo_endpoint"`
	JWKSUri                 string   `json:"jwks_uri"`
	ResponseTypesSupported  []string `json:"response_types_supported"`
	SubjectTypesSupported   []string `json:"subject_types_supported"`
	IDTokenSigningAlgValues []string `json:"id_token_signing_alg_values_supported"`
	ScopesSupported         []string `json:"scopes_supported"`
}

// handleDiscovery handles the OIDC discovery endpoint
func handleDiscovery(cfg *config.Config) *DiscoveryResponse {
	return &DiscoveryResponse{
		Issuer:                  cfg.Issuer,
		AuthorizationEndpoint:   cfg.Server.BaseURL + "/auth",
		TokenEndpoint:           cfg.Server.BaseURL + "/token",
		UserInfoEndpoint:        cfg.Server.BaseURL + "/userinfo",
		JWKSUri:                 cfg.Server.BaseURL + "/jwks",
		ResponseTypesSupported:  []string{"code"},
		SubjectTypesSupported:   []string{"public"},
		IDTokenSigningAlgValues: []string{"RS256"},
		ScopesSupported:         []string{"openid", "profile", "email"},
	}
}
