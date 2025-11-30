package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ninlil/oidc-mockery/internal/config"
)

// DebugLoginHandler handles GET requests to the debug login page
func handleDebugLogin(cfg *config.Config, w http.ResponseWriter, r *http.Request) {

	// Get the first client for the demo
	if len(cfg.Clients) == 0 {
		http.Error(w, "No clients configured", http.StatusInternalServerError)
		return
	}

	client := &cfg.Clients[0]

	// Construct the authorization URL
	baseURL := fmt.Sprintf("http://%s", r.Host)
	authURL := fmt.Sprintf("%s/auth/authorize", baseURL)

	// Build query parameters
	params := url.Values{}
	params.Set("client_id", client.ClientID)
	params.Set("response_type", "code")
	params.Set("scope", "openid profile email")
	params.Set("redirect_uri", fmt.Sprintf("%s/debug/callback", baseURL))
	params.Set("state", fmt.Sprintf("debug_state_%d", time.Now().Unix()))

	fullAuthURL := fmt.Sprintf("%s?%s", authURL, params.Encode())

	// Load template from file
	tmpl, err := template.ParseFiles("templates/debug-login.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := struct {
		ClientID    string
		AuthURL     string
		RedirectURI string
	}{
		ClientID:    client.ClientID,
		AuthURL:     fullAuthURL,
		RedirectURI: fmt.Sprintf("%s/debug/callback", baseURL),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

// DebugCallbackHandler handles the OAuth callback for debug purposes
func handleDebugCallback(cfg *config.Config, w http.ResponseWriter, r *http.Request) {

	// Get the first client for the demo
	if len(cfg.Clients) == 0 {
		http.Error(w, "No clients configured", http.StatusInternalServerError)
		return
	}

	client := &cfg.Clients[0]

	// Extract authorization code and state from query parameters
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")

	baseURL := fmt.Sprintf("http://%s", r.Host)
	var tokenResponse *TokenResponse
	var tokenError string

	// If we have an authorization code, exchange it for tokens
	if code != "" && errorParam == "" {
		// Create token request
		tokenData := url.Values{}
		tokenData.Set("grant_type", "authorization_code")
		tokenData.Set("code", code)
		tokenData.Set("redirect_uri", fmt.Sprintf("%s/debug/callback", baseURL))
		tokenData.Set("client_id", client.ClientID)
		tokenData.Set("client_secret", client.ClientSecret)

		// Make token request
		tokenReq := &TokenArgs{Body: tokenData.Encode()}
		response, statusCode, err := handleToken(cfg, tokenReq)
		if err != nil || statusCode != 200 {
			tokenError = fmt.Sprintf("Token error (status %d): %v", statusCode, err)
		} else {
			tokenResponse = response
		}
	}

	// Parse JWT if we have an ID token
	var parsedJWT map[string]interface{}
	var jwtHeader map[string]interface{}
	var jwtError string

	if tokenResponse != nil && tokenResponse.IDToken != "" {
		parts := strings.Split(tokenResponse.IDToken, ".")
		if len(parts) == 3 {
			// Parse header
			headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
			if err != nil {
				jwtError = fmt.Sprintf("Failed to decode JWT header: %v", err)
			} else {
				if err := json.Unmarshal(headerBytes, &jwtHeader); err != nil {
					jwtError = fmt.Sprintf("Failed to parse JWT header: %v", err)
				}
			}

			// Parse payload
			payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
			if err != nil {
				jwtError = fmt.Sprintf("Failed to decode JWT payload: %v", err)
			} else {
				if err := json.Unmarshal(payloadBytes, &parsedJWT); err != nil {
					jwtError = fmt.Sprintf("Failed to parse JWT payload: %v", err)
				}
			}
		} else {
			jwtError = "Invalid JWT format"
		}
	}

	// HTML template for debug callback page
	tmpl, err := template.ParseFiles("templates/debug-callback.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Convert timestamps to human readable format
	var issuedAtTime, expiresAtTime string
	if parsedJWT != nil {
		if iat, ok := parsedJWT["iat"].(float64); ok {
			issuedAtTime = time.Unix(int64(iat), 0).Format("2006-01-02 15:04:05 UTC")
		}
		if exp, ok := parsedJWT["exp"].(float64); ok {
			expiresAtTime = time.Unix(int64(exp), 0).Format("2006-01-02 15:04:05 UTC")
		}
	}

	// Convert JSON to pretty printed strings
	var jwtHeaderJSON, jwtPayloadJSON string
	if jwtHeader != nil {
		if headerBytes, err := json.MarshalIndent(jwtHeader, "", "  "); err == nil {
			jwtHeaderJSON = string(headerBytes)
		}
	}
	if parsedJWT != nil {
		if payloadBytes, err := json.MarshalIndent(parsedJWT, "", "  "); err == nil {
			jwtPayloadJSON = string(payloadBytes)
		}
	}

	data := struct {
		Code             string
		State            string
		Error            string
		ErrorDescription string
		TokenResponse    *TokenResponse
		TokenError       string
		ParsedJWT        map[string]interface{}
		JWTError         string
		JWTHeaderJSON    string
		JWTPayloadJSON   string
		IssuedAtTime     string
		ExpiresAtTime    string
	}{
		Code:             code,
		State:            state,
		Error:            errorParam,
		ErrorDescription: r.URL.Query().Get("error_description"),
		TokenResponse:    tokenResponse,
		TokenError:       tokenError,
		ParsedJWT:        parsedJWT,
		JWTError:         jwtError,
		JWTHeaderJSON:    jwtHeaderJSON,
		JWTPayloadJSON:   jwtPayloadJSON,
		IssuedAtTime:     issuedAtTime,
		ExpiresAtTime:    expiresAtTime,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

// Handle404 handles all unmatched routes with a helpful 404 page
func handle404(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	// Load template from file
	tmpl, err := template.ParseFiles("templates/404.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Method    string
		Path      string
		Host      string
		UserAgent string
	}{
		Method:    r.Method,
		Path:      r.URL.Path,
		Host:      r.Host,
		UserAgent: r.Header.Get("User-Agent"),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}
