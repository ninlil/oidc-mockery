package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig `yaml:"server"`
	Issuer   string       `yaml:"issuer"`
	JWKS     JWKSConfig   `yaml:"jwks"`
	Clients  []Client     `yaml:"clients"`
	Personas []Persona    `yaml:"personas"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port    int    `yaml:"port"`
	BaseURL string `yaml:"base_url"`
}

// JWKSConfig holds JWKS-specific configuration
type JWKSConfig struct {
	KeyType     string `yaml:"key_type"`
	KeyUse      string `yaml:"key_use"`
	KeyID       string `yaml:"key_id"`
	RSAModulus  string `yaml:"rsa_modulus"`
	RSAExponent string `yaml:"rsa_exponent"`
}

// Client represents an OIDC client configuration
type Client struct {
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	RedirectURIs []string `yaml:"redirect_uris"`
	Scopes       []string `yaml:"scopes"`
}

// Persona represents a user persona for testing
type Persona struct {
	ID     string                 `yaml:"id"`
	Name   string                 `yaml:"name"`
	Email  string                 `yaml:"email"`
	Claims map[string]interface{} `yaml:"claims"`
}

// Load reads and parses the configuration file
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// GetClient returns a client by ID
func (c *Config) GetClient(clientID string) *Client {
	for _, client := range c.Clients {
		if client.ClientID == clientID {
			return &client
		}
	}
	return nil
}

// GetPersona returns a persona by ID
func (c *Config) GetPersona(personaID string) *Persona {
	for _, persona := range c.Personas {
		if persona.ID == personaID {
			return &persona
		}
	}
	return nil
}
