package octopusdeploy_framework

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTagTestClient returns a client backed by a fake server whose tag set
// endpoint answers GET with an empty tag set and PUT with the given status.
func newTagTestClient(t *testing.T, putStatus int, putBody string) *client.Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.Contains(r.URL.Path, "/tagsets/") && r.Method == http.MethodPut:
			w.WriteHeader(putStatus)
			_, _ = w.Write([]byte(putBody))
		case strings.Contains(r.URL.Path, "/tagsets/"):
			_, _ = w.Write([]byte(`{"Id":"TagSets-1","SpaceId":"Spaces-1","Name":"set","Tags":[{"Id":"TagSets-1/Existing","Name":"Existing"}]}`))
		default:
			_, _ = w.Write([]byte(`{"Links":{}}`))
		}
	}))
	t.Cleanup(server.Close)

	u, err := url.Parse(server.URL)
	require.NoError(t, err)
	c, err := client.NewClient(server.Client(), u, "API-ABCDEFGHIJKLMNOPQRSTUVWXYZ", "Spaces-1")
	require.NoError(t, err)
	return c
}

func tagCreateModel(name string) *schemas.TagResourceModel {
	return &schemas.TagResourceModel{
		Name:          types.StringValue(name),
		Color:         types.StringValue("#6e6e6e"),
		TagSetId:      types.StringValue("TagSets-1"),
		TagSetSpaceId: types.StringValue("Spaces-1"),
	}
}

func TestTagCreateReportsTagSetUpdateError(t *testing.T) {
	c := newTagTestClient(t, http.StatusBadRequest,
		`{"ErrorMessage":"The following tag names already exist in other Target tag sets: 'Chicago'"}`)

	var diags = tagCreate(context.Background(), tagCreateModel("Chicago"), nil, c)

	require.True(t, diags.HasError(), "server error must be reported")
	assert.Contains(t, diags.Errors()[0].Detail(), "already exist")
}

func TestTagCreateReportsDuplicateName(t *testing.T) {
	c := newTagTestClient(t, http.StatusOK, `{}`)

	diags := tagCreate(context.Background(), tagCreateModel("Existing"), nil, c)

	require.True(t, diags.HasError())
	assert.Contains(t, diags.Errors()[0].Detail(), "'Existing' is already in use")
}
