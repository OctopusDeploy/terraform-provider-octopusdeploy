package octopusdeploy_framework

import (
	"testing"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy"
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
// See https://github.com/OctopusDeploy/terraform-provider-octopusdeploy/issues/269.
//
// The affected resources are still implemented with the SDK, so these tests read
// the schemas through the exported SDK provider.

func TestSdkPollingURIDiffSuppress(t *testing.T) {
	attribute := sdkResourceAttribute(t, "octopusdeploy_kubernetes_agent_worker", "uri")
	require.NotNil(t, attribute.DiffSuppressFunc)

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
			assert.Equal(t, testCase.expected, attribute.DiffSuppressFunc("uri", testCase.old, testCase.new, nil))
		})
	}
}

func TestSdkPollingURIAttributesSuppressCasingOnResources(t *testing.T) {
	testCases := []struct {
		resource  string
		attribute string
	}{
		{"octopusdeploy_kubernetes_agent_worker", "uri"},
		{"octopusdeploy_kubernetes_agent_deployment_target", "uri"},
		{"octopusdeploy_polling_tentacle_deployment_target", "tentacle_url"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.resource, func(t *testing.T) {
			attribute := sdkResourceAttribute(t, testCase.resource, testCase.attribute)
			require.NotNil(t, attribute.DiffSuppressFunc, "expected %s to suppress casing differences", testCase.attribute)

			assert.True(t, attribute.DiffSuppressFunc(testCase.attribute, "poll://abcdef0123456789/", "poll://ABCDEF0123456789/", nil))
			assert.False(t, attribute.DiffSuppressFunc(testCase.attribute, "poll://abcdef0123456789/", "poll://fedcba9876543210/", nil))
			assert.True(t, attribute.DiffSuppressOnRefresh, "expected %s to suppress casing differences on refresh too", testCase.attribute)
		})
	}
}

func TestSdkPollingURIAttributesDoNotSuppressOnDataSources(t *testing.T) {
	testCases := []struct {
		dataSource string
		listField  string
		attribute  string
	}{
		{"octopusdeploy_kubernetes_agent_workers", "kubernetes_agent_workers", "uri"},
		{"octopusdeploy_kubernetes_agent_deployment_targets", "kubernetes_agent_deployment_targets", "uri"},
		{"octopusdeploy_polling_tentacle_deployment_targets", "polling_tentacle_deployment_targets", "tentacle_url"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.dataSource, func(t *testing.T) {
			dataSource, ok := octopusdeploy.Provider().DataSourcesMap[testCase.dataSource]
			require.True(t, ok, "expected the SDK provider to register %s", testCase.dataSource)

			list, ok := dataSource.Schema[testCase.listField]
			require.True(t, ok, "expected a %s attribute", testCase.listField)
			elem, ok := list.Elem.(*schema.Resource)
			require.True(t, ok, "expected %s to be a list of nested resources", testCase.listField)

			attribute, ok := elem.Schema[testCase.attribute]
			require.True(t, ok, "expected a %s attribute", testCase.attribute)
			assert.Nil(t, attribute.DiffSuppressFunc, "data source attributes are computed and never diffed")
			assert.False(t, attribute.DiffSuppressOnRefresh)
		})
	}
}

func sdkResourceAttribute(t *testing.T, resourceName string, attributeName string) *schema.Schema {
	t.Helper()

	resource, ok := octopusdeploy.Provider().ResourcesMap[resourceName]
	require.True(t, ok, "expected the SDK provider to register %s", resourceName)

	attribute, ok := resource.Schema[attributeName]
	require.True(t, ok, "expected %s to have a %s attribute", resourceName, attributeName)

	return attribute
}
