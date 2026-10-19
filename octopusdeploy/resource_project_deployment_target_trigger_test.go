package octopusdeploy

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/constants"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/filters"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

const (
	targetTriggerProviderSpaceID = "Spaces-1"
	targetTriggerSpaceID         = "Spaces-2"
	targetTriggerID              = "ProjectTriggers-1"
)

// targetTriggerTransport serves the canned root document that client
// construction needs and records every other request, so a test can assert
// which space the provider actually addressed.
type targetTriggerTransport struct {
	requests []string
	body     string
}

func (t *targetTriggerTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if !strings.Contains(request.URL.Path, "projecttriggers") {
		return targetTriggerResponse(constants.ApiReplyRoot), nil
	}

	t.requests = append(t.requests, request.Method+" "+request.URL.Path)

	return targetTriggerResponse(t.body), nil
}

func targetTriggerResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		// The SDK skips decoding entirely when ContentLength is zero.
		ContentLength: int64(len(body)),
	}
}

func newTargetTriggerClient(t *testing.T, transport *targetTriggerTransport) *client.Client {
	t.Helper()

	apiURL, err := url.Parse("http://localhost:8080")
	require.NoError(t, err)

	octopusClient, err := client.NewClientForTool(
		&http.Client{Transport: transport},
		apiURL,
		"API-TESTTESTTESTTEST0",
		targetTriggerProviderSpaceID,
		"terraform-provider-octopusdeploy-test",
	)
	require.NoError(t, err)

	return octopusClient
}

func newTargetTriggerBody(t *testing.T) string {
	t.Helper()

	project := &projects.Project{SpaceID: targetTriggerSpaceID}
	project.ID = "Projects-1"

	trigger := triggers.NewProjectTrigger(
		"test-trigger",
		"",
		false,
		project,
		actions.NewAutoDeployAction(false),
		filters.NewDeploymentTargetFilter(
			[]string{"Environments-1"},
			[]string{"MachineHealthy"},
			[]string{"Machine"},
			[]string{"web-server"},
		),
	)
	trigger.ID = targetTriggerID

	body, err := json.Marshal(trigger)
	require.NoError(t, err)

	return string(body)
}

func newTargetTriggerResourceData(t *testing.T, spaceID string) *schema.ResourceData {
	t.Helper()

	raw := map[string]interface{}{
		"name":       "test-trigger",
		"project_id": "Projects-1",
	}
	if spaceID != "" {
		raw["space_id"] = spaceID
	}

	resourceData := schema.TestResourceDataRaw(t, getProjectDeploymentTargetTriggerSchema(), raw)
	resourceData.SetId(targetTriggerID)

	return resourceData
}

// The provider is configured for one space while the resource declares
// another. Read has to address the resource's space, not the provider's.
func TestProjectDeploymentTargetTriggerReadUsesResourceSpace(t *testing.T) {
	transport := &targetTriggerTransport{body: newTargetTriggerBody(t)}
	octopusClient := newTargetTriggerClient(t, transport)
	resourceData := newTargetTriggerResourceData(t, targetTriggerSpaceID)

	diags := resourceProjectDeploymentTargetTriggerRead(context.Background(), resourceData, octopusClient)

	require.False(t, diags.HasError(), "read returned diagnostics: %v", diags)
	require.Equal(t, []string{"GET /api/" + targetTriggerSpaceID + "/projecttriggers/" + targetTriggerID}, transport.requests)
	require.Equal(t, targetTriggerSpaceID, resourceData.Get("space_id"))
	require.Equal(t, "test-trigger", resourceData.Get("name"))
}

func TestProjectDeploymentTargetTriggerDeleteUsesResourceSpace(t *testing.T) {
	transport := &targetTriggerTransport{}
	octopusClient := newTargetTriggerClient(t, transport)
	resourceData := newTargetTriggerResourceData(t, targetTriggerSpaceID)

	diags := resourceProjectDeploymentTargetTriggerDelete(context.Background(), resourceData, octopusClient)

	require.False(t, diags.HasError(), "delete returned diagnostics: %v", diags)
	require.Equal(t, []string{"DELETE /api/" + targetTriggerSpaceID + "/projecttriggers/" + targetTriggerID}, transport.requests)
	require.Empty(t, resourceData.Id())
}

// Configurations written before space_id existed leave it unset, and must keep
// addressing the space the provider was configured for.
func TestProjectDeploymentTargetTriggerReadFallsBackToProviderSpace(t *testing.T) {
	transport := &targetTriggerTransport{body: newTargetTriggerBody(t)}
	octopusClient := newTargetTriggerClient(t, transport)
	resourceData := newTargetTriggerResourceData(t, "")

	diags := resourceProjectDeploymentTargetTriggerRead(context.Background(), resourceData, octopusClient)

	require.False(t, diags.HasError(), "read returned diagnostics: %v", diags)
	require.Equal(t, []string{"GET /api/" + targetTriggerProviderSpaceID + "/projecttriggers/" + targetTriggerID}, transport.requests)
}
