package octopusdeploy

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

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

type HeaderRoundTripper struct {
	Transport http.RoundTripper
	Headers   map[string]string
}

func (h *HeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}
	return h.Transport.RoundTrip(req)
}

func getHttpClient(octopusUrl string) (*http.Client, error) {
	if !isDirectlyAccessibleOctopusInstance(octopusUrl) {
		log.Printf("[INFO] Enabled Octopus AI Assistant redirection service")
		return createHttpClient(octopusUrl)
	}

	log.Printf("[INFO] Did not enable Octopus AI Assistant redirection service")

	return nil, nil
}

// isDirectlyAccessibleOctopusInstance determines if the host should be contacted directly
func isDirectlyAccessibleOctopusInstance(octopusUrl string) bool {
	serviceEnabled, found := os.LookupEnv("REDIRECTION_SERVICE_ENABLED")

	if !found || serviceEnabled != "true" {
		return true
	}

	parsedUrl, err := url.Parse(octopusUrl)

	// Contact the server directly if the URL is invalid
	if err != nil {
		return true
	}

	return strings.HasSuffix(parsedUrl.Hostname(), ".octopus.app") ||
		strings.HasSuffix(parsedUrl.Hostname(), ".testoctopus.com") ||
		parsedUrl.Hostname() == "localhost" ||
		parsedUrl.Hostname() == "127.0.0.1"
}

func createHttpClient(octopusUrl string) (*http.Client, error) {

	serviceApiKey, found := os.LookupEnv("REDIRECTION_SERVICE_API_KEY")

	if !found {
		return nil, errors.New("REDIRECTION_SERVICE_API_KEY is required")
	}

	parsedUrl, err := url.Parse(octopusUrl)

	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"X_REDIRECTION_UPSTREAM_HOST":   parsedUrl.Hostname(),
		"X_REDIRECTION_SERVICE_API_KEY": serviceApiKey,
	}

	return &http.Client{
		Transport: &HeaderRoundTripper{
			Transport: http.DefaultTransport,
			Headers:   headers,
		},
	}, nil
}

// End of OctoAI patch

// Client returns a new Octopus Deploy client
func (c *Config) Client() (*client.Client, diag.Diagnostics) {
	octopus, err := getClientForDefaultSpace(c)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if len(c.SpaceID) > 0 {
		space, err := spaces.GetByID(octopus, c.SpaceID)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		octopus, err = getClientForSpace(c, space.GetID())
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	return octopus, nil
}

func getClientForDefaultSpace(c *Config) (*client.Client, error) {
	return getClientForSpace(c, "")
}

func getClientForSpace(c *Config, spaceID string) (*client.Client, error) {
	apiURL, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}

	credential, err := getApiCredential(c)
	if err != nil {
		return nil, err
	}

	// Start of OctoAI patch
	httpClient, err := getHttpClient(c.Address)
	if err != nil {
		return nil, err
	}

	return client.NewClientWithCredentials(httpClient, apiURL, credential, spaceID, "TerraformProvider")
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
