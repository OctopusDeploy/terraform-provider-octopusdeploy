package octopusdeploy_framework

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/version"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/configuration"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Config struct {
	Address        string
	ApiKey         string
	AccessToken    string
	SpaceID        string
	Client         *client.Client
	OctopusVersion string
	// Can be nil when server doesn't support feature toggles API endpoint
	FeatureToggles map[string]bool
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

func (c *Config) SetOctopus(ctx context.Context) diag.Diagnostics {
	tflog.Debug(ctx, "SetOctopus")

	diags := diag.Diagnostics{}

	if clientError := c.GetClient(ctx); clientError != nil {
		diags.AddError("failed to load client", clientError.Error())
		return diags
	}

	if versionError := c.SetOctopusVersion(ctx); versionError != nil {
		diags.AddError("failed to load Octopus Server version", versionError.Error())
		return diags
	}

	c.SetFeatureToggles(ctx)

	tflog.Debug(ctx, "SetOctopus completed")
	return diags
}

func (c *Config) GetClient(ctx context.Context) error {
	tflog.Debug(ctx, "GetClient")

	octopus, err := getClientForDefaultSpace(c, ctx)
	if err != nil {
		return err
	}

	if len(c.SpaceID) > 0 {
		space, err := spaces.GetByID(octopus, c.SpaceID)
		if err != nil {
			return err
		}

		octopus, err = getClientForSpace(c, ctx, space.GetID())
		if err != nil {
			return err
		}
	}

	c.Client = octopus

	createdClient := octopus != nil
	tflog.Debug(ctx, fmt.Sprintf("GetClient completed: %t", createdClient))
	return nil
}

func (c *Config) SetFeatureToggles(ctx context.Context) {
	tflog.Debug(ctx, "SetFeatureToggles")

	response, err := configuration.Get(c.Client, &configuration.FeatureToggleConfigurationQuery{})
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Unable to load feature toggles: %q", err.Error()))
		c.FeatureToggles = nil
		return
	}

	features := make(map[string]bool, len(response.FeatureToggles))
	for _, feature := range response.FeatureToggles {
		features[feature.Name] = feature.IsEnabled
	}

	c.FeatureToggles = features

	tflog.Debug(ctx, fmt.Sprintf("SetFeatureToggles completed with %d features", len(c.FeatureToggles)))
}

func (c *Config) SetOctopusVersion(ctx context.Context) error {
	tflog.Debug(ctx, "SetOctopusVersion")

	root, err := client.GetServerRoot(c.Client)
	if err != nil {
		return err
	}

	c.OctopusVersion = root.Version
	tflog.Debug(ctx, fmt.Sprintf("SetOctopusVersion completed with %s", c.OctopusVersion))

	return nil
}

func getClientForDefaultSpace(c *Config, ctx context.Context) (*client.Client, error) {
	return getClientForSpace(c, ctx, "")
}

func getClientForSpace(c *Config, ctx context.Context, spaceID string) (*client.Client, error) {
	apiURL, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}

	credential, err := getApiCredential(c, ctx)
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

func getApiCredential(c *Config, ctx context.Context) (client.ICredential, error) {
	tflog.Debug(ctx, "GetClient: Trying the following auth methods in order of priority - APIKey, AccessToken")

	if c.ApiKey != "" {
		tflog.Debug(ctx, "GetClient: Attempting to authenticate with API Key")
		credential, err := client.NewApiKey(c.ApiKey)
		if err != nil {
			return nil, err
		}

		return credential, nil
	} else {
		tflog.Debug(ctx, "GetClient: No API Key found")
	}

	if c.AccessToken != "" {
		tflog.Debug(ctx, "GetClient: Attempting to authenticate with Access Token")
		credential, err := client.NewAccessToken(c.AccessToken)
		if err != nil {
			return nil, err
		}

		return credential, nil
	} else {
		tflog.Debug(ctx, "GetClient: No Access Token found")
	}

	return nil, fmt.Errorf("either an APIKey or an AccessToken is required to connect to the Octopus Server instance")
}

func DataSourceConfiguration(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) *Config {
	if req.ProviderData == nil {
		return nil
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return config
}

func ResourceConfiguration(req resource.ConfigureRequest, resp *resource.ConfigureResponse) *Config {
	if req.ProviderData == nil {
		return nil
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return config
}

// FeatureToggleEnabled Reports whether feature toggle enabled on connected Octopus Server instance.
func (c *Config) FeatureToggleEnabled(toggle string) bool {
	if c.FeatureToggles == nil {
		return true
	}

	if enabled, ok := c.FeatureToggles[toggle]; ok {
		return enabled
	}

	return false
}

// EnsureResourceCompatibilityByFeature Reports whether resource is compatible with current instance of Octopus Server by .
// Returns diagnostics with error when resource is incompatible and empty diagnostics for compatible resources
func (c *Config) EnsureResourceCompatibilityByFeature(resourceName string, toggle string) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if c.FeatureToggleEnabled(toggle) {
		return diags
	}

	summary := fmt.Sprintf("The '%s' resource is not supported by the connected Octopus Deploy instance", resourceName)
	detail := fmt.Sprintf("This resource requires feature toggle '%s' to be enabled.", toggle)
	diags.AddError(summary, detail)

	return diags
}

// EnsureResourceCompatibilityByVersion Reports whether resource is compatible with current version of Octopus Server.
// Returns diagnostics with error when resource is incompatible and empty diagnostics for compatible resources
//
// Example: '2025.1' - first version where resource can be used
func (c *Config) EnsureResourceCompatibilityByVersion(resourceName string, version string) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if c.IsVersionSameOrGreaterThan(version) {
		return diags
	}

	summary := fmt.Sprintf("The '%s' resource is not supported by the current Octopus Deploy server version", resourceName)
	detail := fmt.Sprintf("This resource requires Octopus Deploy server version %s or later. The connected server is running version %s, which is incompatible with this resource.", version, c.OctopusVersion)
	diags.AddError(summary, detail)

	return diags
}

func (c *Config) IsVersionSameOrGreaterThan(minVersion string) bool {
	if c.OctopusVersion == "0.0.0-local" {
		return true // Always true for local instance
	}

	diff := version.Compare(fmt.Sprintf("go%s", c.OctopusVersion), fmt.Sprintf("go%s", minVersion))

	return diff == 1 || diff == 0
}
