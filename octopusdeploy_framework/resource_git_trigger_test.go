package octopusdeploy_framework

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
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

const (
	gitTriggerProviderSpaceID = "Spaces-1"
	gitTriggerSpaceIDValue    = "Spaces-2"
	gitTriggerIDValue         = "ProjectTriggers-1"
)

// gitTriggerRecordingTransport serves the canned root document that client
// construction needs and records every other request, so a test can assert
// which space the provider actually addressed.
type gitTriggerRecordingTransport struct {
	requests []string
	body     string
	// failUpdate makes the trigger PUT return a server error.
	failUpdate bool
}

func (t *gitTriggerRecordingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	switch {
	case strings.Contains(request.URL.Path, "projecttriggers"):
		t.requests = append(t.requests, request.Method+" "+request.URL.Path)

		if t.failUpdate && request.Method == http.MethodPut {
			return gitTriggerErrorResponse(), nil
		}

		return gitTriggerResponse(http.StatusOK, t.body), nil
	case strings.Contains(request.URL.Path, "/projects/"):
		return gitTriggerResponse(http.StatusOK, `{"Id":"Projects-1","SpaceId":"`+gitTriggerSpaceIDValue+`","Name":"test-project","LifecycleId":"Lifecycles-1","ProjectGroupId":"ProjectGroups-1"}`), nil
	default:
		return gitTriggerResponse(http.StatusOK, constants.ApiReplyRoot), nil
	}
}

func gitTriggerResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		// The SDK skips decoding entirely when ContentLength is zero.
		ContentLength: int64(len(body)),
	}
}

func gitTriggerErrorResponse() *http.Response {
	return gitTriggerResponse(http.StatusInternalServerError, `{"ErrorMessage":"something went wrong"}`)
}

func newGitTriggerResource(t *testing.T, transport *gitTriggerRecordingTransport) *gitTriggerResource {
	t.Helper()

	apiURL, err := url.Parse("http://localhost:8080")
	require.NoError(t, err)

	octopusClient, err := client.NewClientForTool(
		&http.Client{Transport: transport},
		apiURL,
		"API-TESTTESTTESTTEST0",
		gitTriggerProviderSpaceID,
		"terraform-provider-octopusdeploy-test",
	)
	require.NoError(t, err)

	return &gitTriggerResource{Config: &Config{Client: octopusClient}}
}

func newGitTriggerBody(t *testing.T) string {
	t.Helper()

	project := &projects.Project{SpaceID: gitTriggerSpaceIDValue}
	project.ID = "Projects-1"

	trigger := triggers.NewProjectTrigger(
		"test-trigger",
		"",
		false,
		project,
		actions.NewCreateReleaseAction("Channels-1"),
		filters.NewGitTriggerFilter([]filters.GitTriggerSource{{
			DeploymentActionSlug: "test-action",
			GitDependencyName:    "",
			IncludeFilePaths:     []string{"include/me"},
			ExcludeFilePaths:     []string{"exclude/me"},
		}}),
	)
	trigger.ID = gitTriggerIDValue

	body, err := json.Marshal(trigger)
	require.NoError(t, err)

	return string(body)
}

// gitTriggerState builds resource state holding a trigger that declares
// spaceID, whatever space the provider itself is pointed at.
func gitTriggerState(t *testing.T, ctx context.Context, spaceID string) tfsdk.State {
	t.Helper()

	schemaResponse := &resource.SchemaResponse{}
	(&gitTriggerResource{}).Schema(ctx, resource.SchemaRequest{}, schemaResponse)

	state := tfsdk.State{
		Schema: schemaResponse.Schema,
		Raw:    tftypes.NewValue(schemaResponse.Schema.Type().TerraformType(ctx), nil),
	}

	sources := types.ListValueMust(types.ObjectType{AttrTypes: sourcesObjectType()}, []attr.Value{})

	diags := state.Set(ctx, &schemas.GitTriggerResourceModel{
		Name:        types.StringValue("test-trigger"),
		Description: types.StringValue(""),
		SpaceId:     types.StringValue(spaceID),
		ProjectId:   types.StringValue("Projects-1"),
		ChannelId:   types.StringValue("Channels-1"),
		Sources:     sources,
		IsDisabled:  types.BoolValue(false),
		ResourceModel: schemas.ResourceModel{
			ID: types.StringValue(gitTriggerIDValue),
		},
	})
	require.False(t, diags.HasError(), "setting state: %v", diags)

	return state
}

// The provider is configured for one space while the resource declares
// another. Read has to address the resource's space, not the provider's.
func TestGitTriggerReadUsesResourceSpace(t *testing.T) {
	ctx := context.Background()

	transport := &gitTriggerRecordingTransport{body: newGitTriggerBody(t)}
	gitTrigger := newGitTriggerResource(t, transport)

	state := gitTriggerState(t, ctx, gitTriggerSpaceIDValue)
	response := &resource.ReadResponse{State: state}

	gitTrigger.Read(ctx, resource.ReadRequest{State: state}, response)

	require.False(t, response.Diagnostics.HasError(), "read returned diagnostics: %v", response.Diagnostics)
	require.Equal(t, []string{"GET /api/" + gitTriggerSpaceIDValue + "/projecttriggers/" + gitTriggerIDValue}, transport.requests)
}

func TestGitTriggerDeleteUsesResourceSpace(t *testing.T) {
	ctx := context.Background()

	transport := &gitTriggerRecordingTransport{}
	gitTrigger := newGitTriggerResource(t, transport)

	state := gitTriggerState(t, ctx, gitTriggerSpaceIDValue)
	response := &resource.DeleteResponse{State: state}

	gitTrigger.Delete(ctx, resource.DeleteRequest{State: state}, response)

	require.False(t, response.Diagnostics.HasError(), "delete returned diagnostics: %v", response.Diagnostics)
	require.Equal(t, []string{"DELETE /api/" + gitTriggerSpaceIDValue + "/projecttriggers/" + gitTriggerIDValue}, transport.requests)
}

// A failed update used to be reported as a success, because the error from the
// update call was assigned and never checked before state was written.
func TestGitTriggerUpdateReportsFailure(t *testing.T) {
	ctx := context.Background()

	transport := &gitTriggerRecordingTransport{body: newGitTriggerBody(t), failUpdate: true}
	gitTrigger := newGitTriggerResource(t, transport)

	state := gitTriggerState(t, ctx, gitTriggerSpaceIDValue)
	response := &resource.UpdateResponse{State: state}

	gitTrigger.Update(ctx, resource.UpdateRequest{Plan: tfsdk.Plan(state), State: state}, response)

	require.True(t, response.Diagnostics.HasError(), "a failed update must surface a diagnostic")
}

// With no space_id set the resource keeps falling back to the space the
// provider was configured for.
func TestGitTriggerReadFallsBackToProviderSpace(t *testing.T) {
	ctx := context.Background()

	transport := &gitTriggerRecordingTransport{body: newGitTriggerBody(t)}
	gitTrigger := newGitTriggerResource(t, transport)

	state := gitTriggerState(t, ctx, "")
	response := &resource.ReadResponse{State: state}

	gitTrigger.Read(ctx, resource.ReadRequest{State: state}, response)

	require.False(t, response.Diagnostics.HasError(), "read returned diagnostics: %v", response.Diagnostics)
	require.Equal(t, []string{"GET /api/" + gitTriggerProviderSpaceID + "/projecttriggers/" + gitTriggerIDValue}, transport.requests)
}
