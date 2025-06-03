package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpandDefaultMachinePackageCacheRetentionPolicy(t *testing.T) {
	actual := expandMachinePackageCacheRetentionPolicy(nil)
	require.Nil(t, actual)

	actual = expandMachinePackageCacheRetentionPolicy(&schema.Set{})
	require.Nil(t, actual)

	dataSet := schema.NewSet(schema.HashResource(&schema.Resource{
		Schema: getMachinePackageCacheRetentionPolicySchema(),
	}), []interface{}{
		map[string]interface{}{
			"strategy": "Default",
		},
	})

	expected := machinepolicies.NewDefaultMachinePackageCacheRetentionPolicy()

	actual = expandMachinePackageCacheRetentionPolicy(dataSet)
	require.Equal(t, expected, actual)
}

func TestExpandMachinePackageCacheRetentionPolicy(t *testing.T) {
	actual := expandMachinePackageCacheRetentionPolicy(nil)
	require.Nil(t, actual)

	actual = expandMachinePackageCacheRetentionPolicy(&schema.Set{})
	require.Nil(t, actual)

	dataSet := schema.NewSet(schema.HashResource(&schema.Resource{
		Schema: getMachinePackageCacheRetentionPolicySchema(),
	}), []interface{}{
		map[string]interface{}{
			"strategy":                     "Quantities",
			"quantity_of_packages_to_keep": 4,
			"package_unit":                 "Items",
			"quantity_of_versions_to_keep": 6,
			"version_unit":                 "Items",
		},
	})

	expected := machinepolicies.NewMachinePackageCacheRetentionPolicy("Quantities", 4, "Items", 6, "Items")

	actual = expandMachinePackageCacheRetentionPolicy(dataSet)
	require.Equal(t, expected, actual)
}

func TestFlattenDefaultMachinePackageCacheRetentionPolicy(t *testing.T) {
	actual := flattenMachinePackageCacheRetentionPolicy(nil)
	require.Nil(t, actual)

	expanded := machinepolicies.NewDefaultMachinePackageCacheRetentionPolicy()

	actual = flattenMachinePackageCacheRetentionPolicy(expanded)
	expected := []interface{}{
		map[string]interface{}{
			"strategy": "Default",
		},
	}
	require.Equal(t, expected, actual)
}
