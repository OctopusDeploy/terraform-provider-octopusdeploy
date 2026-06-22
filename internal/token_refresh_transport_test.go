package internal

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

// roundTripFunc adapts a function to http.RoundTripper for testing.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestPassthroughOn200(t *testing.T) {
	transport := NewTokenRefreshTransport("token1", "", "")
	transport.base = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return newResponse(200, "ok"), nil
	})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer token1")

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestNonBearerAuthPassthrough(t *testing.T) {
	var callCount int
	transport := NewTokenRefreshTransport("token1", "", "")
	transport.base = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		callCount++
		return newResponse(401, "unauthorized"), nil
	})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Octopus-ApiKey", "API-KEY123")

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
	if callCount != 1 {
		t.Fatalf("expected 1 call (no retry for non-Bearer auth), got %d", callCount)
	}
}

func TestRefreshFromFile(t *testing.T) {
	dir := t.TempDir()
	tokenFile := filepath.Join(dir, "token")
	os.WriteFile(tokenFile, []byte("new-token-from-file"), 0600)

	var callCount int
	transport := NewTokenRefreshTransport("old-token", tokenFile, "")
	transport.base = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		callCount++
		auth := req.Header.Get("Authorization")
		if auth == "Bearer new-token-from-file" {
			return newResponse(200, "ok"), nil
		}
		return newResponse(401, "unauthorized"), nil
	})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer old-token")

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 after refresh, got %d", resp.StatusCode)
	}
	if callCount != 2 {
		t.Fatalf("expected 2 calls (original + retry), got %d", callCount)
	}
}

func TestRefreshFallbackToEnvVar(t *testing.T) {
	const envVar = "TEST_OCTOPUS_TOKEN_REFRESH"
	t.Setenv(envVar, "new-token-from-env")

	var callCount int
	transport := NewTokenRefreshTransport("old-token", "", envVar)
	transport.base = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		callCount++
		auth := req.Header.Get("Authorization")
		if auth == "Bearer new-token-from-env" {
			return newResponse(200, "ok"), nil
		}
		return newResponse(401, "unauthorized"), nil
	})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer old-token")

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 after env refresh, got %d", resp.StatusCode)
	}
	if callCount != 2 {
		t.Fatalf("expected 2 calls, got %d", callCount)
	}
}

func TestFilePriorityOverEnvVar(t *testing.T) {
	dir := t.TempDir()
	tokenFile := filepath.Join(dir, "token")
	os.WriteFile(tokenFile, []byte("file-token"), 0600)

	const envVar = "TEST_OCTOPUS_TOKEN_PRIORITY"
	t.Setenv(envVar, "env-token")

	transport := NewTokenRefreshTransport("old-token", tokenFile, envVar)
	transport.base = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		auth := req.Header.Get("Authorization")
		if auth == "Bearer file-token" {
			return newResponse(200, "ok"), nil
		}
		return newResponse(401, "unauthorized"), nil
	})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer old-token")

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected file token to take priority, got %d", resp.StatusCode)
	}
}

func TestBodyReplayedOnRetry(t *testing.T) {
	dir := t.TempDir()
	tokenFile := filepath.Join(dir, "token")
	os.WriteFile(tokenFile, []byte("new-token"), 0600)

	var retryBody string
	transport := NewTokenRefreshTransport("old-token", tokenFile, "")
	transport.base = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		auth := req.Header.Get("Authorization")
		if auth == "Bearer new-token" {
			if req.Body != nil {
				data, _ := io.ReadAll(req.Body)
				retryBody = string(data)
			}
			return newResponse(200, "ok"), nil
		}
		return newResponse(401, "unauthorized"), nil
	})

	body := `{"key":"value"}`
	req, _ := http.NewRequest("POST", "http://example.com", bytes.NewReader([]byte(body)))
	req.Header.Set("Authorization", "Bearer old-token")
	req.Header.Set("Content-Type", "application/json")

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if retryBody != body {
		t.Fatalf("expected body %q on retry, got %q", body, retryBody)
	}
}

func TestConcurrent401sOnlyOneRefresh(t *testing.T) {
	dir := t.TempDir()
	tokenFile := filepath.Join(dir, "token")
	os.WriteFile(tokenFile, []byte("refreshed-token"), 0600)

	var refreshCount atomic.Int32
	transport := NewTokenRefreshTransport("old-token", tokenFile, "")

	// Track refresh via callback
	transport.SetOnTokenRefreshed(func(newToken string) {
		refreshCount.Add(1)
	})

	transport.base = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		auth := req.Header.Get("Authorization")
		if auth == "Bearer refreshed-token" {
			return newResponse(200, "ok"), nil
		}
		return newResponse(401, "unauthorized"), nil
	})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			req.Header.Set("Authorization", "Bearer old-token")
			resp, err := transport.RoundTrip(req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if resp.StatusCode != 200 {
				t.Errorf("expected 200, got %d", resp.StatusCode)
			}
		}()
	}
	wg.Wait()

	// Only one refresh should happen — subsequent goroutines find token already updated
	count := refreshCount.Load()
	if count != 1 {
		t.Fatalf("expected exactly 1 refresh, got %d", count)
	}
}

func TestNoFreshTokenReturnsOriginal401(t *testing.T) {
	// No file, no env var
	transport := NewTokenRefreshTransport("old-token", "", "")
	transport.base = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return newResponse(401, "unauthorized"), nil
	})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer old-token")

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401 when no fresh token, got %d", resp.StatusCode)
	}
}

func TestTokenUnchangedAfterReadNoRetry(t *testing.T) {
	dir := t.TempDir()
	tokenFile := filepath.Join(dir, "token")
	// File contains same token as current
	os.WriteFile(tokenFile, []byte("old-token"), 0600)

	var callCount int
	transport := NewTokenRefreshTransport("old-token", tokenFile, "")
	transport.base = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		callCount++
		return newResponse(401, "unauthorized"), nil
	})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer old-token")

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
	if callCount != 1 {
		t.Fatalf("expected 1 call (no retry with same token), got %d", callCount)
	}
}

func TestOnTokenRefreshedCallbackInvoked(t *testing.T) {
	dir := t.TempDir()
	tokenFile := filepath.Join(dir, "token")
	os.WriteFile(tokenFile, []byte("callback-token"), 0600)

	var callbackToken string
	transport := NewTokenRefreshTransport("old-token", tokenFile, "")
	transport.SetOnTokenRefreshed(func(newToken string) {
		callbackToken = newToken
	})
	transport.base = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		auth := req.Header.Get("Authorization")
		if auth == "Bearer callback-token" {
			return newResponse(200, "ok"), nil
		}
		return newResponse(401, "unauthorized"), nil
	})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer old-token")

	transport.RoundTrip(req)

	if callbackToken != "callback-token" {
		t.Fatalf("expected callback with %q, got %q", "callback-token", callbackToken)
	}
}

func TestProactiveSyncOverridesStaleHeader(t *testing.T) {
	transport := NewTokenRefreshTransport("current-token", "", "")

	var receivedAuth string
	transport.base = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		receivedAuth = req.Header.Get("Authorization")
		return newResponse(200, "ok"), nil
	})

	// Request has stale token (set by Sling/HttpSession)
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer stale-token")

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if receivedAuth != "Bearer current-token" {
		t.Fatalf("expected proactive sync to current token, got %q", receivedAuth)
	}
}

func TestFileMissingFallsToEnvVar(t *testing.T) {
	const envVar = "TEST_OCTOPUS_FILE_MISSING"
	t.Setenv(envVar, "env-fallback-token")

	transport := NewTokenRefreshTransport("old-token", "/nonexistent/path/token", envVar)
	transport.base = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		auth := req.Header.Get("Authorization")
		if auth == "Bearer env-fallback-token" {
			return newResponse(200, "ok"), nil
		}
		return newResponse(401, "unauthorized"), nil
	})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer old-token")

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected env var fallback, got %d", resp.StatusCode)
	}
}
