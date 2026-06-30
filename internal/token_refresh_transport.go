package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// TokenRefreshTransport is an http.RoundTripper that intercepts 401 responses
// for Bearer token auth, refreshes the token from a file or environment variable,
// and retries the request once. This supports OIDC token refresh for long-running
// Terraform operations.
type TokenRefreshTransport struct {
	base             http.RoundTripper
	mu               sync.Mutex
	currentToken     string
	accessTokenFile  string
	envVarName       string
	onTokenRefreshed func(newToken string)
}

// NewTokenRefreshTransport creates a new TokenRefreshTransport.
// initialToken is the token used at provider init time.
// accessTokenFile is the path to a file containing the token (optional, empty to disable).
// envVarName is the environment variable to read for a refreshed token.
func NewTokenRefreshTransport(initialToken string, accessTokenFile string, envVarName string) *TokenRefreshTransport {
	return &TokenRefreshTransport{
		base:            http.DefaultTransport,
		currentToken:    initialToken,
		accessTokenFile: accessTokenFile,
		envVarName:      envVarName,
	}
}

// SetOnTokenRefreshed sets a callback invoked after a successful token refresh.
// Typically used to update HttpSession.DefaultHeaders.
func (t *TokenRefreshTransport) SetOnTokenRefreshed(fn func(newToken string)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.onTokenRefreshed = fn
}

func (t *TokenRefreshTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	authHeader := req.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		// Not Bearer auth (e.g. API key) — pass through untouched
		return t.base.RoundTrip(req)
	}

	// Proactive sync: if Sling/HttpSession set a stale token, override with current
	t.mu.Lock()
	currentToken := t.currentToken
	t.mu.Unlock()

	req = cloneRequest(req)
	req.Header.Set("Authorization", "Bearer "+currentToken)

	// Buffer body for potential retry
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body for token refresh retry: %w", err)
		}
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != http.StatusUnauthorized {
		return resp, nil
	}

	// 401 — attempt token refresh
	newToken, refreshErr := t.refreshToken()
	if refreshErr != nil {
		return resp, nil // Return original 401
	}
	if newToken == "" {
		return resp, nil
	}

	// Close the original response body before retry
	resp.Body.Close()

	// Retry with new token
	retryReq := cloneRequest(req)
	retryReq.Header.Set("Authorization", "Bearer "+newToken)
	if bodyBytes != nil {
		retryReq.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		retryReq.ContentLength = int64(len(bodyBytes))
	}

	return t.base.RoundTrip(retryReq)
}

// refreshToken attempts to read a fresh token from file or env var.
// Returns empty string if no new token available.
func (t *TokenRefreshTransport) refreshToken() (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	oldToken := t.currentToken
	newToken := ""

	// Try file first
	if t.accessTokenFile != "" {
		data, err := os.ReadFile(t.accessTokenFile)
		if err != nil {
			log.Printf("[WARN] Failed to read access token file %q: %v", t.accessTokenFile, err)
		} else {
			candidate := strings.TrimSpace(string(data))
			if candidate != "" {
				newToken = candidate
			}
		}
	}

	// Fall back to env var
	if newToken == "" && t.envVarName != "" {
		candidate := strings.TrimSpace(os.Getenv(t.envVarName))
		if candidate != "" {
			newToken = candidate
		}
	}

	if newToken == "" || newToken == oldToken {
		// No fresh token or same token — don't retry
		return "", nil
	}

	t.currentToken = newToken
	if t.onTokenRefreshed != nil {
		t.onTokenRefreshed(newToken)
	}

	log.Printf("[INFO] OIDC access token refreshed successfully")
	return newToken, nil
}

// cloneRequest creates a shallow copy of the request with a cloned Header map.
func cloneRequest(req *http.Request) *http.Request {
	r := req.Clone(req.Context())
	return r
}
