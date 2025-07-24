# OIDC-Mockery

A developer tool that provides a complete mock OpenID Connect (OIDC) server for testing client applications without needing a real OIDC provider.

> ⚠️ **Development Tool Only** - This is designed for development and testing purposes. Do not use in production environments.

## ✅ **COMPLETED FEATURES**

### 🔐 Core OIDC Endpoints
- **Discovery** (`/.well-known/openid-configuration`) - OIDC metadata
- **Authorization** (`/auth`) - Login with manual and quick persona selection
- **Consent** (`/consent`) - OAuth consent page with Allow/Deny options
- **Token** (`/token`) - Authorization code exchange for JWT tokens
- **UserInfo** (`/userinfo`) - Claims retrieval with access tokens
- **JWKS** (`/jwks`) - Public keys for JWT verification

### 🎯 Additional Features
- **Manual Login** - Input fields for creating new personas on-the-fly
- **Quick Login** - One-click login with existing personas
- **Dynamic Persona Management** - New personas automatically added to configuration
- **OAuth Compliance** - Proper consent denial handling with `access_denied` error responses
- **Debug Interface** (`/debug/login`, `/debug/callback`) - Self-testing with JWT parsing
- **Static Files** (`/static/*`) - CSS/JS assets
- **404 Handling** - Error pages with endpoint guidance
- **Load Testing** - k6 test suite with 100% success validation

### 🛠️ Technical Stack
- **Go + Butler Framework** - HTTP routing and middleware
- **HTML Templates** - Login, consent, debug, and error pages
- **CSS/JS** - Form handling and interactions
- **YAML Configuration** - Flexible client and persona management with dynamic updates
- **In-Memory Storage** - Authorization codes with 10-minute expiration
- **RSA JWT Signing** - Development-compatible token generation

## ⚙️ Configuration

The `config.yaml` file includes complete OIDC server configuration:

```yaml
server:
  port: 8081
  base_url: "http://localhost:8081"

issuer: "http://localhost:8081"

# JWT signing configuration
jwks:
  key_type: "RSA"
  key_use: "sig"
  key_id: "mock-key-id"
  rsa_modulus: "mock-rsa-modulus-base64url-encoded"
  rsa_exponent: "AQAB"

clients:
  - client_id: "test-client"
    client_secret: "test-secret"
    redirect_uris:
      - "http://localhost:8081/debug/callback"
    scopes: ["openid", "profile", "email"]

personas:
  - id: "user1"
    name: "John Doe"
    email: "john.doe@example.com"
    claims:
      given_name: "John"
      family_name: "Doe"
      preferred_username: "johndoe"
  - id: "user2"
    name: "Jane Smith"
    email: "jane.smith@example.com"
    claims:
      given_name: "Jane"
      family_name: "Smith"
      preferred_username: "janesmith"
```

## 📁 Project Structure

```
oidc-mockery/
├── main.go                    # Server entry point
├── config.yaml               # OIDC configuration
├── internal/
│   ├── config/config.go      # Configuration loading
│   ├── handlers/             # HTTP endpoints
│   │   ├── routes.go         # Route definitions
│   │   ├── discovery.go      # OIDC discovery
│   │   ├── auth.go           # Authorization + templates
│   │   ├── token.go          # Token exchange
│   │   ├── userinfo.go       # User claims
│   │   ├── jwks.go           # Key distribution
│   │   ├── static.go         # Static files
│   │   └── debug.go          # Debug interface
│   ├── models/               # Data structures
│   └── utils/                # JWT & crypto utilities
├── templates/                # HTML templates
│   ├── login.html            # Manual + quick persona selection
│   ├── consent.html          # OAuth consent with Allow/Deny buttons
│   ├── debug-*.html          # Debug interface
│   └── 404.html              # Error page
├── static/
│   ├── css/style.css         # Styling & themes
│   └── js/app.js             # Form handling & interactions
└── test/full-flow.k6         # Complete OIDC flow testing
```

## 🚀 Quick Start

```bash
# Start the server (already built)
./bin/oidc-mockery

# Test endpoints
curl http://localhost:8081/.well-known/openid-configuration

# Use debug interface
open http://localhost:8081/debug/login

# Run load tests
k6 run test/full-flow.k6
```

## 🔄 OIDC Flow

Complete Authorization Code flow implementation:
```
GET  /auth     → Login (persona selection)
POST /auth     → Redirect to consent
GET  /consent  → Consent page
POST /consent  → Generate auth code & redirect
POST /token    → Exchange code for JWT tokens
GET  /userinfo → Retrieve user claims
```

## 🔒 Security & 🛠️ Customization

**Security Features (Development Use Only):**
- RSA JWT signing with proper key management
- Redirect URI validation and authorization code expiration
- Client authentication and secure random code generation

⚠️ **Note**: This is a development/testing tool - not intended for production use.

**Add Personas/Clients:**
```yaml
# Edit config.yaml
personas:
  - id: "admin"
    name: "Admin User"
    email: "admin@example.com"
    claims: { role: "administrator" }

clients:
  - client_id: "mobile-app"
    client_secret: "mobile-secret"
    redirect_uris: ["myapp://callback"]
```

**Rebuild:** `go build -o bin/oidc-mockery . && ./bin/oidc-mockery`

## 🎯 Use Cases

Perfect for testing **Web Apps**, **Mobile Apps**, **Development Teams**, **CI/CD Pipelines**, and **Workshops/Demos**.
