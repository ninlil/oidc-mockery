package handlers

import (
	"net/http"

	"github.com/ninlil/oidc-mockery/internal/config"

	"github.com/ninlil/butler/router"
)

// GetRoutes returns all OIDC endpoints as router.Route slice
func GetRoutes(cfg *config.Config) []router.Route {
	// Create handlers with config closure
	discoveryHandler := func() *DiscoveryResponse {
		return handleDiscovery(cfg)
	}

	// Template-based auth handler that serves HTML
	authHandler := func(w http.ResponseWriter, r *http.Request) {
		handleAuthTemplate(cfg, w, r)
	}

	authPostHandler := func(w http.ResponseWriter, r *http.Request) {
		handleAuthPost(cfg, w, r)
	}

	// Template-based consent handler that serves HTML
	consentHandler := func(w http.ResponseWriter, r *http.Request) {
		handleConsentTemplate(cfg, w, r)
	}

	// Static file handler
	staticHandler := func(w http.ResponseWriter, r *http.Request) {
		HandleStatic(w, r)
	}

	// Debug handlers for testing
	debugLoginHandler := func(w http.ResponseWriter, r *http.Request) {
		handleDebugLogin(cfg, w, r)
	}

	debugCallbackHandler := func(w http.ResponseWriter, r *http.Request) {
		handleDebugCallback(cfg, w, r)
	}

	// 404 fallback handler for all unmatched routes
	notFoundHandler := func(w http.ResponseWriter, r *http.Request) {
		handle404(cfg, w, r)
	}

	tokenHandler := func(args *TokenArgs) (*TokenResponse, int, error) {
		return handleToken(cfg, args)
	}

	userInfoHandler := func(args *UserInfoArgs) (*UserInfoResponse, int, error) {
		return handleUserInfo(cfg, args)
	}

	jwksHandler := func() *JWKSResponse {
		return handleJWKS(cfg)
	}

	return []router.Route{
		{Name: "discovery", Method: "GET", Path: "/.well-known/openid-configuration", Handler: discoveryHandler},
		{Name: "static", Method: "GET", Path: "/static/*", Handler: staticHandler},
		{Name: "auth", Method: "GET", Path: "/auth", Handler: authHandler},
		{Name: "authPost", Method: "POST", Path: "/auth", Handler: authPostHandler},
		{Name: "consent", Method: "GET", Path: "/consent", Handler: consentHandler},
		{Name: "consentPost", Method: "POST", Path: "/consent", Handler: consentHandler},
		{Name: "token", Method: "POST", Path: "/token", Handler: tokenHandler},
		{Name: "userinfo", Method: "GET", Path: "/userinfo", Handler: userInfoHandler},
		{Name: "userinfoPost", Method: "POST", Path: "/userinfo", Handler: userInfoHandler},
		{Name: "jwks", Method: "GET", Path: "/jwks", Handler: jwksHandler},
		{Name: "debugLogin", Method: "GET", Path: "/debug/login", Handler: debugLoginHandler},
		{Name: "debugCallback", Method: "GET", Path: "/debug/callback", Handler: debugCallbackHandler},
		{Name: "notFound", Method: "*", Path: "/*", Handler: notFoundHandler},
	}
}
