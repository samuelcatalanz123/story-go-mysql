// Package oauth implements "Sign in with Google" using OpenID Connect.
//
// IMPORTANT: this needs real Google credentials to work. Create an OAuth 2.0
// "Web application" client at https://console.cloud.google.com and set
// GOOGLE_CLIENT_ID / GOOGLE_CLIENT_SECRET / GOOGLE_REDIRECT_URI. Without them
// the server still starts, but the Google login endpoint reports it is not
// configured. (See .env.example.)
package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

const (
	googleIssuer       = "https://accounts.google.com"
	googleTokenURL     = "https://oauth2.googleapis.com/token"
	googleProviderName = "google"
)

// GoogleAuthenticator exchanges an OAuth authorization code for a verified
// Google identity. The id_token's signature is checked against Google's public
// keys (JWKS), so we never trust unverified client input.
type GoogleAuthenticator struct {
	clientID     string
	clientSecret string
	redirectURI  string
	verifier     *oidc.IDTokenVerifier
	httpClient   *http.Client
}

// NewGoogleAuthenticator sets up the authenticator, fetching Google's OpenID
// discovery document. Returns an error if Google is unreachable.
func NewGoogleAuthenticator(ctx context.Context, clientID, clientSecret, redirectURI string) (*GoogleAuthenticator, error) {
	provider, err := oidc.NewProvider(ctx, googleIssuer)
	if err != nil {
		return nil, fmt.Errorf("google oidc discovery: %w", err)
	}
	return &GoogleAuthenticator{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		verifier:     provider.Verifier(&oidc.Config{ClientID: clientID}),
		httpClient:   http.DefaultClient,
	}, nil
}

// ProviderName identifies this provider in the oauth_accounts table.
func (g *GoogleAuthenticator) ProviderName() string { return googleProviderName }

// Exchange swaps the authorization code (plus the PKCE code_verifier) for a
// Google id_token, verifies it and returns the user's immutable subject and
// email.
func (g *GoogleAuthenticator) Exchange(ctx context.Context, code, codeVerifier string) (subject, email string, err error) {
	idToken, err := g.fetchIDToken(ctx, code, codeVerifier)
	if err != nil {
		return "", "", err
	}

	verified, err := g.verifier.Verify(ctx, idToken)
	if err != nil {
		return "", "", fmt.Errorf("verify id_token: %w", err)
	}

	var claims struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
	}
	if err := verified.Claims(&claims); err != nil {
		return "", "", fmt.Errorf("read id_token claims: %w", err)
	}
	if claims.Sub == "" {
		return "", "", fmt.Errorf("id_token without subject")
	}
	return claims.Sub, claims.Email, nil
}

// fetchIDToken calls Google's token endpoint to exchange the code for tokens
// and returns the raw id_token string.
func (g *GoogleAuthenticator) fetchIDToken(ctx context.Context, code, codeVerifier string) (string, error) {
	form := url.Values{
		"code":          {code},
		"client_id":     {g.clientID},
		"client_secret": {g.clientSecret},
		"redirect_uri":  {g.redirectURI},
		"grant_type":    {"authorization_code"},
		"code_verifier": {codeVerifier},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, googleTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("google token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("google token endpoint returned %d", resp.StatusCode)
	}

	var body struct {
		IDToken string `json:"id_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("decode google token response: %w", err)
	}
	if body.IDToken == "" {
		return "", fmt.Errorf("google response without id_token")
	}
	return body.IDToken, nil
}
