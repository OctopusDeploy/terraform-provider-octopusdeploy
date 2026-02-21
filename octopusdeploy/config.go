package octopusdeploy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Config holds Address and the APIKey of the Octopus Deploy server
type Config struct {
	Address     string
	APIKey      string
	AccessToken string
	SpaceID     string
}

// Start of OctoAI patch

type headerRoundTripper struct {
	Transport http.RoundTripper
	Headers   map[string]string
}

func (h *headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}
	return h.Transport.RoundTrip(req)
}

func getHttpClient(ctx context.Context, octopusUrl *url.URL) (*http.Client, *url.URL, error) {
	if !isDirectlyAccessibleOctopusInstance(octopusUrl) {
		tflog.Info(ctx, "[SPACEBUILDER] Enabled Octopus AI Assistant redirection service")
		return createHttpClient(octopusUrl)
	}

	tflog.Info(ctx, "[SPACEBUILDER] Did not enable Octopus AI Assistant redirection service")

	return nil, octopusUrl, nil
}

func getRedirectionBypass() []string {
	hostnames := []string{}
	hostnamesJson := os.Getenv("REDIRECTION_BYPASS")
	if hostnamesJson == "" {
		return []string{} // Default to empty slice if not set
	}

	err := json.Unmarshal([]byte(hostnamesJson), &hostnames)
	if err != nil {
		return []string{}
	}

	return hostnames
}

func getRedirectionForce() bool {
	redirectionForce := os.Getenv("REDIRECTION_FORCE")
	return strings.ToLower(redirectionForce) == "true"
}

// isDirectlyAccessibleOctopusInstance determines if the host should be contacted directly
func isDirectlyAccessibleOctopusInstance(octopusUrl *url.URL) bool {
	serviceEnabled, found := os.LookupEnv("REDIRECTION_SERVICE_ENABLED")

	if !found || serviceEnabled != "true" {
		return true
	}

	bypassList := getRedirectionBypass()

	// Allow bypassing specific domains via environment variable
	if slices.Contains(bypassList, octopusUrl.Hostname()) {
		return true
	}

	// Allow forcing all traffic through the redirection service
	if getRedirectionForce() {
		return false
	}

	return strings.HasSuffix(octopusUrl.Hostname(), ".octopus.app") ||
		strings.HasSuffix(octopusUrl.Hostname(), ".testoctopus.com") ||
		octopusUrl.Hostname() == "localhost" ||
		octopusUrl.Hostname() == "127.0.0.1"
}

func createHttpClient(octopusUrl *url.URL) (*http.Client, *url.URL, error) {

	serviceApiKey, found := os.LookupEnv("REDIRECTION_SERVICE_API_KEY")

	if !found {
		return nil, nil, errors.New("REDIRECTION_SERVICE_API_KEY is required")
	}

	redirectionHost, found := os.LookupEnv("REDIRECTION_HOST")

	if !found {
		return nil, nil, errors.New("REDIRECTION_HOST is required")
	}

	redirectionHostUrl, err := url.Parse("https://" + redirectionHost)

	if err != nil {
		return nil, nil, err
	}

	headers := map[string]string{
		"X_REDIRECTION_UPSTREAM_HOST":   octopusUrl.Hostname(),
		"X_REDIRECTION_SERVICE_API_KEY": serviceApiKey,
	}

	return &http.Client{
		Transport: &headerRoundTripper{
			Transport: http.DefaultTransport,
			Headers:   headers,
		},
	}, redirectionHostUrl, nil
}

// End of OctoAI patch

// Client returns a new Octopus Deploy client
func (c *Config) Client(ctx context.Context) (*client.Client, diag.Diagnostics) {
	octopus, err := getClientForDefaultSpace(ctx, c)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if len(c.SpaceID) > 0 {
		space, err := spaces.GetByID(octopus, c.SpaceID)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		octopus, err = getClientForSpace(ctx, c, space.GetID())
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	return octopus, nil
}

func getClientForDefaultSpace(ctx context.Context, c *Config) (*client.Client, error) {
	return getClientForSpace(ctx, c, "")
}

func getClientForSpace(ctx context.Context, c *Config, spaceID string) (*client.Client, error) {
	apiURL, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}

	credential, err := getApiCredential(c)
	if err != nil {
		return nil, err
	}

	// Start of OctoAI patch
	httpClient, url, err := getHttpClient(ctx, apiURL)
	if err != nil {
		return nil, err
	}

	tflog.Info(ctx, "[SPACEBUILDER] Directing requests from "+apiURL.String())
	tflog.Info(ctx, "[SPACEBUILDER] Directing requests to redirector at "+url.String())

	return client.NewClientWithCredentials(httpClient, url, credential, spaceID, "TerraformProvider")
	// End of OctoAI patch
}

func getApiCredential(c *Config) (client.ICredential, error) {
	if c.APIKey != "" {
		credential, err := client.NewApiKey(c.APIKey)
		if err != nil {
			return nil, err
		}

		return credential, nil
	}

	if c.AccessToken != "" {
		credential, err := client.NewAccessToken(c.AccessToken)
		if err != nil {
			return nil, err
		}

		return credential, nil
	}

	return nil, fmt.Errorf("either an APIKey or an AccessToken is required to connect to the Octopus Server instance")
}
