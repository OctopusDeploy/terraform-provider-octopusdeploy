package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandMachinePackageCacheRetentionPolicy(values interface{}) *machinepolicies.MachinePackageCacheRetentionPolicy {
	if values == nil {
		return nil
	}
	flattenedValues := values.(*schema.Set)
	if len(flattenedValues.List()) == 0 {
		return nil
	}

	flattenedMap := flattenedValues.List()[0].(map[string]interface{})

	machinePackageCacheRetentionPolicy := machinepolicies.NewMachinePackageCacheRetentionPolicy()

	if v, ok := flattenedMap["strategy"]; ok {
		machinePackageCacheRetentionPolicy.Strategy = v.(string)
	}

	if v, ok := flattenedMap["quantity_of_packages_to_keep"]; ok {
		int32Val := int32(v.(int))
		machinePackageCacheRetentionPolicy.QuantityOfPackagesToKeep = int32Val
	}

	if v, ok := flattenedMap["package_unit"]; ok {
		var stringPackageUnit = v.(string)
		machinePackageCacheRetentionPolicy.PackageUnit = stringPackageUnit
	}

	if v, ok := flattenedMap["quantity_of_versions_to_keep"]; ok {
		int32Val := int32(v.(int))
		machinePackageCacheRetentionPolicy.QuantityOfVersionsToKeep = int32Val
	}

	if v, ok := flattenedMap["version_unit"]; ok {
		var stringVersionUnit = v.(string)
		machinePackageCacheRetentionPolicy.VersionUnit = stringVersionUnit
	}

	return machinePackageCacheRetentionPolicy
}

func flattenMachinePackageCacheRetentionPolicy(machineUpdatePolicy *machinepolicies.MachinePackageCacheRetentionPolicy) []interface{} {
	if machineUpdatePolicy == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"strategy":                     machineUpdatePolicy.Strategy,
		"quantity_of_packages_to_keep": machineUpdatePolicy.QuantityOfPackagesToKeep,
		"package_unit":                 machineUpdatePolicy.PackageUnit,
		"quantity_of_versions_to_keep": machineUpdatePolicy.QuantityOfVersionsToKeep,
		"version_unit":                 machineUpdatePolicy.VersionUnit,
	}}
}

func getMachinePackageCacheRetentionPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"strategy": {
			Required: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"Default",
				"Quantities",
			}, false)),
			Description: "The behaviour of the cache retention policy. Valid values are `Default` (let Octopus decide), `Quantities` (keep by a specified number of packages and versions).",
		},
		"quantity_of_packages_to_keep": {
			Optional:    true,
			Type:        schema.TypeInt,
			Description: "The number of packages to keep.",
		},
		"package_unit": {
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"Items",
			}, false)),
			Description: "The method of counting packages when applying the Quantities strategy.",
		},
		"quantity_of_versions_to_keep": {
			Optional:    true,
			Type:        schema.TypeInt,
			Description: "The number of package versions to keep.",
		},
		"version_unit": {
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"Items",
			}, false)),
			Description: "The method of counting package versions when applying the Quantities strategy.",
		},
	}
}
