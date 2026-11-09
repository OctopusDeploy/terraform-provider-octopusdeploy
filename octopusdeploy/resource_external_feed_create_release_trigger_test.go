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
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

const (
	providerSpaceID = "Spaces-1"
	triggerSpaceID  = "Spaces-2"
	triggerID       = "ProjectTriggers-1"
)

// recordingTransport serves the canned root document that client construction
// needs and records every other request so a test can assert which space the
// provider actually addressed.
type recordingTransport struct {
	requests      []string
	triggerBody   string
	triggerStatus int
}

func (t *recordingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	path := request.URL.Path

	if !strings.Contains(path, "projecttriggers") {
		return jsonResponse(http.StatusOK, constants.ApiReplyRoot), nil
	}

	t.requests = append(t.requests, request.Method+" "+path)

	return jsonResponse(t.triggerStatus, t.triggerBody), nil
}

func jsonResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		// The SDK skips decoding entirely when ContentLength is zero.
		ContentLength: int64(len(body)),
	}
}

// newTestClientForSpace builds a client bound to spaceID, backed by transport.
func newTestClientForSpace(t *testing.T, transport *recordingTransport, spaceID string) *client.Client {
	t.Helper()

	apiURL, err := url.Parse("http://localhost:8080")
	require.NoError(t, err)

	octopusClient, err := client.NewClientForTool(
		&http.Client{Transport: transport},
		apiURL,
		"API-TESTTESTTESTTEST0",
		spaceID,
		"terraform-provider-octopusdeploy-test",
	)
	require.NoError(t, err)

	return octopusClient
}

func newTestTriggerBody(t *testing.T) string {
	t.Helper()

	project := &projects.Project{SpaceID: triggerSpaceID}
	project.ID = "Projects-1"

	trigger := triggers.NewProjectTrigger(
		"test-trigger",
		"",
		false,
		project,
		actions.NewCreateReleaseAction("Channels-1"),
		filters.NewFeedTriggerFilter([]packages.DeploymentActionSlugPackage{
			{DeploymentActionSlug: "test-action"},
		}),
	)
	trigger.ID = triggerID

	body, err := json.Marshal(trigger)
	require.NoError(t, err)

	return string(body)
}

func newTestTriggerResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()

	resourceData := schema.TestResourceDataRaw(t, getExternalFeedCreateReleaseTriggerSchema(), map[string]interface{}{
		"name":       "test-trigger",
		"space_id":   triggerSpaceID,
		"project_id": "Projects-1",
		"channel_id": "Channels-1",
	})
	resourceData.SetId(triggerID)

	return resourceData
}

// The provider is configured for one space while the resource declares another.
// Read has to address the resource's space, not the provider's, or refreshing a
// trigger outside the provider's space fails with "Resource is not found or it
// doesn't exist in the current space context".
func TestExternalFeedCreateReleaseTriggerReadUsesResourceSpace(t *testing.T) {
	transport := &recordingTransport{
		triggerBody:   newTestTriggerBody(t),
		triggerStatus: http.StatusOK,
	}

	octopusClient := newTestClientForSpace(t, transport, providerSpaceID)
	resourceData := newTestTriggerResourceData(t)

	diags := resourceExternalFeedCreateReleaseTriggerRead(context.Background(), resourceData, octopusClient)

	require.False(t, diags.HasError(), "read returned diagnostics: %v", diags)
	require.Equal(t, []string{"GET /api/" + triggerSpaceID + "/projecttriggers/" + triggerID}, transport.requests)
	require.Equal(t, "test-trigger", resourceData.Get("name"))
	require.Equal(t, triggerSpaceID, resourceData.Get("space_id"))
}

func TestExternalFeedCreateReleaseTriggerDeleteUsesResourceSpace(t *testing.T) {
	transport := &recordingTransport{
		triggerBody:   "",
		triggerStatus: http.StatusOK,
	}

	octopusClient := newTestClientForSpace(t, transport, providerSpaceID)
	resourceData := newTestTriggerResourceData(t)

	diags := resourceExternalFeedCreateReleaseTriggerDelete(context.Background(), resourceData, octopusClient)

	require.False(t, diags.HasError(), "delete returned diagnostics: %v", diags)
	require.Equal(t, []string{"DELETE /api/" + triggerSpaceID + "/projecttriggers/" + triggerID}, transport.requests)
	require.Empty(t, resourceData.Id())
}

// With no space_id set the resource keeps falling back to the space the
// provider was configured for.
func TestExternalFeedCreateReleaseTriggerReadFallsBackToProviderSpace(t *testing.T) {
	transport := &recordingTransport{
		triggerBody:   newTestTriggerBody(t),
		triggerStatus: http.StatusOK,
	}

	octopusClient := newTestClientForSpace(t, transport, providerSpaceID)

	resourceData := schema.TestResourceDataRaw(t, getExternalFeedCreateReleaseTriggerSchema(), map[string]interface{}{
		"name":       "test-trigger",
		"project_id": "Projects-1",
		"channel_id": "Channels-1",
	})
	resourceData.SetId(triggerID)

	diags := resourceExternalFeedCreateReleaseTriggerRead(context.Background(), resourceData, octopusClient)

	require.False(t, diags.HasError(), "read returned diagnostics: %v", diags)
	require.Equal(t, []string{"GET /api/" + providerSpaceID + "/projecttriggers/" + triggerID}, transport.requests)
}
