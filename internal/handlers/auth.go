package handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ninlil/oidc-mockery/internal/config"
	"github.com/ninlil/oidc-mockery/internal/utils"
)

// AuthCodeData represents stored authorization code data
type AuthCodeData struct {
	ClientID     string
	PersonaID    string
	PersonaName  string
	PersonaEmail string
	RedirectURI  string
	ExpiresAt    time.Time
}

// In-memory store for authorization codes (for mockery purposes)
var authCodeStore = make(map[string]AuthCodeData)

// AuthArgs represents an OIDC authorization request parameters
type AuthArgs struct {
	ClientID     string `json:"client_id" from:"query" required:"true"`
	ResponseType string `json:"response_type" from:"query" required:"true"`
	Scope        string `json:"scope" from:"query" required:"true"`
	RedirectURI  string `json:"redirect_uri" from:"query" required:"true"`
	State        string `json:"state" from:"query"`
	Nonce        string `json:"nonce" from:"query"`
}

// AuthPostArgs represents POST parameters for persona selection
type AuthPostArgs struct {
	Body string `from:"body"`
}

// handleAuthPost handles POST requests to the authorization endpoint (persona selection)
func handleAuthPost(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// log.Trace().Msgf("Received auth post request %q", string(body))

	// Extract form fields from body map
	bodyValues, err := url.ParseQuery(string(body))
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	personaID := bodyValues.Get("persona_id")
	personaName := bodyValues.Get("name")
	personaEmail := bodyValues.Get("email")
	clientID := bodyValues.Get("client_id")
	redirectURI := bodyValues.Get("redirect_uri")
	state := bodyValues.Get("state")
	scope := bodyValues.Get("scope")

	var missingFields []string
	if personaID == "" {
		missingFields = append(missingFields, "persona_id")
	}
	if personaName == "" {
		missingFields = append(missingFields, "name")
	}
	if personaEmail == "" {
		missingFields = append(missingFields, "email")
	}
	// Validate required fields
	if len(missingFields) > 0 {
		http.Error(w, fmt.Sprintf("Missing required fields: %s", strings.Join(missingFields, ", ")), http.StatusBadRequest)
		return
	}

	// Validate client
	client := cfg.GetClient(clientID)
	if client == nil {
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		return
	}

	// Redirect to consent page with parameters
	consentURL, _ := url.Parse("/consent")
	consentQuery := consentURL.Query()
	consentQuery.Set("client_id", clientID)
	consentQuery.Set("redirect_uri", redirectURI)
	consentQuery.Set("state", state)
	consentQuery.Set("scope", scope)
	consentQuery.Set("persona_id", personaID)
	consentQuery.Set("persona_name", personaName)
	consentQuery.Set("persona_email", personaEmail)
	consentURL.RawQuery = consentQuery.Encode()

	// Perform the redirect
	http.Redirect(w, r, consentURL.String(), http.StatusFound)
}

// handleAuthTemplate handles GET requests to the authorization endpoint and renders HTML template
func handleAuthTemplate(cfg *config.Config, w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	clientID := query.Get("client_id")
	responseType := query.Get("response_type")
	state := query.Get("state")
	scope := query.Get("scope")
	redirectURI := query.Get("redirect_uri")
	nonce := query.Get("nonce")

	// Validate client
	client := cfg.GetClient(clientID)
	if client == nil {
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		return
	}

	// Validate redirect URI
	if !utils.ValidateRedirectURI(redirectURI, client.RedirectURIs) {
		http.Error(w, "Invalid redirect_uri", http.StatusBadRequest)
		return
	}

	// Parse and render login template
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	templateData := struct {
		Personas []config.Persona
		Request  *AuthArgs
	}{
		Personas: cfg.Personas,
		Request: &AuthArgs{
			ClientID:     clientID,
			ResponseType: responseType,
			Scope:        scope,
			RedirectURI:  redirectURI,
			State:        state,
			Nonce:        nonce},
	}

	// Set content type and render template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, templateData); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

// handleConsentTemplate handles consent form display and submission
func handleConsentTemplate(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Parse query parameters
		query := r.URL.Query()
		clientID := query.Get("client_id")
		redirectURI := query.Get("redirect_uri")
		state := query.Get("state")
		scope := query.Get("scope")
		personaID := query.Get("persona_id")
		personaName := query.Get("persona_name")
		personaEmail := query.Get("persona_email")

		// Validate required persona fields
		if personaID == "" || personaName == "" || personaEmail == "" {
			http.Error(w, "Missing required persona fields", http.StatusBadRequest)
			return
		}

		// Create persona from form data
		persona := &config.Persona{
			ID:    personaID,
			Name:  personaName,
			Email: personaEmail,
		}

		// Parse template
		tmpl, err := template.ParseFiles("templates/consent.html")
		if err != nil {
			http.Error(w, "Template parsing error", http.StatusInternalServerError)
			return
		}

		// Prepare data for template
		scopes := []string{"openid", "profile", "email"}
		if scope != "" {
			scopes = strings.Split(scope, " ")
		}

		data := struct {
			ClientID    string
			RedirectURI string
			State       string
			Scopes      []string
			Persona     *config.Persona
		}{
			ClientID:    clientID,
			RedirectURI: redirectURI,
			State:       state,
			Scopes:      scopes,
			Persona:     persona,
		}

		w.Header().Set("Content-Type", "text/html")
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		// Handle consent form submission
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		action := r.FormValue("action")

		if action == "deny" {
			// Redirect with error
			redirectURI := r.FormValue("redirect_uri")
			state := r.FormValue("state")

			redirectURL, err := url.Parse(redirectURI)
			if err != nil {
				http.Error(w, "Invalid redirect_uri", http.StatusBadRequest)
				return
			}
			redirectQuery := redirectURL.Query()
			redirectQuery.Set("error", "access_denied")
			redirectQuery.Set("state", state)
			redirectURL.RawQuery = redirectQuery.Encode()

			http.Redirect(w, r, redirectURL.String(), http.StatusFound)
			return
		}

		if action == "" {
			http.Error(w, "Missing action parameter", http.StatusBadRequest)
			return
		}

		// Handle allow - generate authorization code
		clientID := r.FormValue("client_id")
		personaID := r.FormValue("persona_id")
		personaName := r.FormValue("persona_name")
		personaEmail := r.FormValue("persona_email")
		redirectURI := r.FormValue("redirect_uri")
		state := r.FormValue("state")

		// Validate required fields
		if personaID == "" || personaName == "" || personaEmail == "" {
			http.Error(w, "Missing required persona fields", http.StatusBadRequest)
			return
		}

		var personaFound bool
		for _, persona := range cfg.Personas {
			if persona.ID == personaID {
				personaFound = true
				break
			}
		}
		if !personaFound {
			cfg.Personas = append(cfg.Personas, config.Persona{
				ID:    personaID,
				Name:  personaName,
				Email: personaEmail,
			})
		}

		// Validate client
		client := cfg.GetClient(clientID)
		if client == nil {
			http.Error(w, "Invalid client_id", http.StatusBadRequest)
			return
		}

		// Generate authorization code
		authCode := fmt.Sprintf("auth_%d", time.Now().Unix())

		// Store the authorization code (in memory for this mockery)
		authCodeStore[authCode] = AuthCodeData{
			ClientID:     clientID,
			PersonaID:    personaID,
			PersonaName:  personaName,
			PersonaEmail: personaEmail,
			RedirectURI:  redirectURI,
			ExpiresAt:    time.Now().Add(10 * time.Minute),
		}

		// Redirect back to client with authorization code
		redirectURL, err := url.Parse(redirectURI)
		if err != nil {
			http.Error(w, "Invalid redirect_uri", http.StatusBadRequest)
			return
		}
		redirectQuery := redirectURL.Query()
		redirectQuery.Set("code", authCode)
		redirectQuery.Set("state", state)
		redirectURL.RawQuery = redirectQuery.Encode()

		// fmt.Printf("Redirecting to: %s\n", redirectURL)
		http.Redirect(w, r, redirectURL.String(), http.StatusFound)
		// fmt.Printf("Redirect sent\n")
	}
}
