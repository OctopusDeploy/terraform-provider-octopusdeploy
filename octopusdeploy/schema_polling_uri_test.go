package octopusdeploy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Octopus Server lowercases the host of a polling URI, while
// octopusdeploy_polling_subscription_id only ever generates uppercase subscription
// IDs. Without diff suppression that mismatch shows up as a permanent in-place
// update. These tests guard the wiring, which is what actually regressed: the
// suppressor existed for kubernetes_agent_deployment_target but was never
// attached to the worker or the polling tentacle target.

func TestPollingURIDiffSuppress(t *testing.T) {
	suppress := getPollingURIDiffSuppressFunc()

	testCases := []struct {
		name     string
		old      string
		new      string
		expected bool
	}{
		{"identical", "poll://abcdef0123456789/", "poll://abcdef0123456789/", true},
		{"server lowercase vs generated uppercase", "poll://abcdef0123456789/", "poll://ABCDEF0123456789/", true},
		{"mixed casing", "poll://AbCdEf0123456789/", "poll://aBcDeF0123456789/", true},
		{"different subscription id", "poll://abcdef0123456789/", "poll://fedcba9876543210/", false},
		{"trailing slash removed", "poll://abcdef0123456789/", "poll://abcdef0123456789", false},
		{"empty vs populated", "", "poll://abcdef0123456789/", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, suppress("uri", testCase.old, testCase.new, nil))
		})
	}
}

func TestPollingURIAttributesSuppressCasingOnResources(t *testing.T) {
	testCases := []struct {
		name      string
		resource  *schema.Resource
		attribute string
	}{
		{"kubernetes agent worker", resourceKubernetesAgentWorker(), "uri"},
		{"kubernetes agent deployment target", resourceKubernetesAgentDeploymentTarget(), "uri"},
		{"polling tentacle deployment target", resourcePollingTentacleDeploymentTarget(), "tentacle_url"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			attribute, ok := testCase.resource.Schema[testCase.attribute]
			require.True(t, ok, "expected a %s attribute", testCase.attribute)
			require.NotNil(t, attribute.DiffSuppressFunc, "expected %s to suppress casing differences", testCase.attribute)

			assert.True(t, attribute.DiffSuppressFunc(testCase.attribute, "poll://abcdef0123456789/", "poll://ABCDEF0123456789/", nil))
			assert.False(t, attribute.DiffSuppressFunc(testCase.attribute, "poll://abcdef0123456789/", "poll://fedcba9876543210/", nil))
			assert.True(t, attribute.DiffSuppressOnRefresh, "expected %s to suppress casing differences on refresh too", testCase.attribute)
		})
	}
}

func TestPollingURIAttributesDoNotSuppressOnDataSources(t *testing.T) {
	testCases := []struct {
		name      string
		schema    map[string]*schema.Schema
		attribute string
	}{
		{"kubernetes agent workers", getKubernetesAgentWorkerSchemaForDataSource(), "uri"},
		{"kubernetes agent deployment targets", getKubernetesAgentDeploymentTargetSchemaForDataSource(), "uri"},
		{"polling tentacle deployment targets", getPollingTentacleDeploymentTargetSchemaForDataSource(), "tentacle_url"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			attribute, ok := testCase.schema[testCase.attribute]
			require.True(t, ok, "expected a %s attribute", testCase.attribute)
			assert.Nil(t, attribute.DiffSuppressFunc, "data source attributes are computed and never diffed")
			assert.False(t, attribute.DiffSuppressOnRefresh)
		})
	}
}
