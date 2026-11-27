package octopusdeploy_framework

import (
	"math"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/uritemplates"
	"github.com/stretchr/testify/require"
)

// The step template data source filters server side, which only works when the query struct
// carries `uri` tags matching the variables in the URI template. actiontemplates.Query does;
// actiontemplates.ActionTemplateSearch does not, and silently degrades to fetching the first
// page unfiltered. These tests fail loudly if either of those facts changes in the SDK.
func TestActionTemplatesCollectionUriTemplateSendsFilters(t *testing.T) {
	query := actiontemplates.Query{
		PartialName: "Hello World",
		Take:        math.MaxInt32,
	}

	uri := expandActionTemplatesCollectionUri(t, query, "Spaces-1")

	require.Contains(t, uri, "partialName=Hello%20World")
	require.Contains(t, uri, "take=2147483647")
	require.Contains(t, uri, "/api/Spaces-1/actiontemplates")
}

func TestActionTemplatesCollectionUriTemplateSendsIdsAsCommaSeparatedList(t *testing.T) {
	// The server matches a comma separated ids parameter, which is what the template
	// expansion produces for a slice.
	query := actiontemplates.Query{IDs: []string{"ActionTemplates-1", "ActionTemplates-2"}}

	uri := expandActionTemplatesCollectionUri(t, query, "Spaces-1")

	require.Contains(t, uri, "ids=ActionTemplates-1%2CActionTemplates-2")
}

func TestActionTemplatesCollectionUriTemplateHasNoIdPathSegment(t *testing.T) {
	// A {/id} segment would expand into a single-resource path that cannot be decoded as a
	// collection, so the collection template must not carry one.
	require.NotContains(t, actionTemplatesCollectionUriTemplate, "{/id}")
}

func TestActionTemplateSearchCannotFilterServerSide(t *testing.T) {
	search := actiontemplates.ActionTemplateSearch{
		ID:   "ActionTemplates-1",
		Name: "Hello World",
	}

	uri := expandActionTemplatesCollectionUri(t, search, "Spaces-1")

	require.Equal(t, "/api/Spaces-1/actiontemplates", uri)
}

func expandActionTemplatesCollectionUri(t *testing.T, query any, spaceID string) string {
	t.Helper()

	values, ok := uritemplates.Struct2map(query)
	require.True(t, ok)
	values["spaceId"] = spaceID

	uri, err := uritemplates.NewUriTemplateCache().Expand(actionTemplatesCollectionUriTemplate, values)
	require.NoError(t, err)

	return uri
}
