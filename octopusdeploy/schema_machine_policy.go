package octopusdeploy

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandMachinePolicy(d *schema.ResourceData) *machinepolicies.MachinePolicy {
	name := d.Get("name").(string)

	machinePolicy := machinepolicies.NewMachinePolicy(name)
	machinePolicy.ID = d.Id()

	if v, ok := d.GetOk("connection_connect_timeout"); ok {
		machinePolicy.ConnectionConnectTimeout = time.Duration(v.(int))
	}

	if v, ok := d.GetOk("connection_retry_count_limit"); ok {
		machinePolicy.ConnectionRetryCountLimit = int32(v.(int))
	}

	if v, ok := d.GetOk("connection_retry_sleep_interval"); ok {
		machinePolicy.ConnectionRetrySleepInterval = time.Duration(v.(int))
	}

	if v, ok := d.GetOk("connection_retry_time_limit"); ok {
		machinePolicy.ConnectionRetryTimeLimit = time.Duration(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		machinePolicy.Description = v.(string)
	}

	if v, ok := d.GetOk("is_default"); ok {
		machinePolicy.IsDefault = v.(bool)
	}

	if v, ok := d.GetOk("machine_cleanup_policy"); ok {
		if len(v.(*schema.Set).List()) > 0 {
			machinePolicy.MachineCleanupPolicy = expandMachineCleanupPolicy(v)
		}
	}

	if v, ok := d.GetOk("machine_connectivity_policy"); ok {
		if len(v.(*schema.Set).List()) > 0 {
			machinePolicy.MachineConnectivityPolicy = expandMachineConnectivityPolicy(v)
		}
	}

	if v, ok := d.GetOk("machine_health_check_policy"); ok {
		if len(v.(*schema.Set).List()) > 0 {
			machinePolicy.MachineHealthCheckPolicy = expandMachineHealthCheckPolicy(v)
		}
	}

	if v, ok := d.GetOk("machine_update_policy"); ok {
		if len(v.(*schema.Set).List()) > 0 {
			machinePolicy.MachineUpdatePolicy = expandMachineUpdatePolicy(v)
		}
	}

	if v, ok := d.GetOk("machine_package_cache_retention_policy"); ok {
		if len(v.(*schema.Set).List()) > 0 {
			machinePolicy.MachinePackageCacheRetentionPolicy = expandMachinePackageCacheRetentionPolicy(v)
		}
	}

	if v, ok := d.GetOk("name"); ok {
		machinePolicy.Name = v.(string)
	}

	if v, ok := d.GetOk("polling_request_queue_timeout"); ok {
		machinePolicy.PollingRequestQueueTimeout = time.Duration(v.(int))
	}

	if v, ok := d.GetOk("space_id"); ok {
		machinePolicy.SpaceID = v.(string)
	}

	return machinePolicy

}

func flattenMachinePolicy(machinePolicy *machinepolicies.MachinePolicy) map[string]interface{} {
	if machinePolicy == nil {
		return nil
	}

	return map[string]interface{}{
		"connection_connect_timeout":             machinePolicy.ConnectionConnectTimeout,
		"connection_retry_count_limit":           machinePolicy.ConnectionRetryCountLimit,
		"connection_retry_sleep_interval":        machinePolicy.ConnectionRetrySleepInterval,
		"connection_retry_time_limit":            machinePolicy.ConnectionRetryTimeLimit,
		"description":                            machinePolicy.Description,
		"id":                                     machinePolicy.GetID(),
		"is_default":                             machinePolicy.IsDefault,
		"machine_cleanup_policy":                 flattenMachineCleanupPolicy(machinePolicy.MachineCleanupPolicy),
		"machine_connectivity_policy":            flattenMachineConnectivityPolicy(machinePolicy.MachineConnectivityPolicy),
		"machine_health_check_policy":            flattenMachineHealthCheckPolicy(machinePolicy.MachineHealthCheckPolicy),
		"machine_update_policy":                  flattenMachineUpdatePolicy(machinePolicy.MachineUpdatePolicy),
		"machine_package_cache_retention_policy": flattenMachinePackageCacheRetentionPolicy(machinePolicy.MachinePackageCacheRetentionPolicy),
		"name":                                   machinePolicy.Name,
		"polling_request_queue_timeout":          machinePolicy.PollingRequestQueueTimeout,
		"space_id":                               machinePolicy.SpaceID,
	}
}

func getMachinePolicyDataSchema() map[string]*schema.Schema {
	dataSchema := getMachinePolicySchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"ids": getQueryIDs(),
		"machine_policies": {
			Computed:    true,
			Description: "A list of machine policies that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    false,
			Type:        schema.TypeList,
		},
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
		"space_id":     getSpaceIDSchema(),
	}
}

func getMachinePolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_connect_timeout": {
			Default:     time.Minute,
			Optional:    true,
			Type:        schema.TypeInt,
			Description: "In nanoseconds. Minimum value: 10000000000 (10 seconds).",
		},
		"connection_retry_count_limit": {
			Default:  5,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"connection_retry_sleep_interval": {
			Default:     time.Second,
			Optional:    true,
			Type:        schema.TypeInt,
			Description: "In nanoseconds.",
		},
		"connection_retry_time_limit": {
			Default:     5 * time.Minute,
			Optional:    true,
			Type:        schema.TypeInt,
			Description: "In nanoseconds.",
		},
		"description": getDescriptionSchema("machine policy"),
		"id":          getIDSchema(),
		"is_default": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"machine_cleanup_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getMachineCleanupPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"machine_connectivity_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getMachineConnectivityPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"machine_health_check_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getMachineHealthCheckPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"machine_update_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getMachineUpdatePolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"machine_package_cache_retention_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getMachinePackageCacheRetentionPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"name": getNameSchema(true),
		"polling_request_queue_timeout": {
			Default:     2 * time.Minute,
			Optional:    true,
			Type:        schema.TypeInt,
			Description: "In nanoseconds.",
		},
		"space_id": getSpaceIDSchema(),
	}
}

func setMachinePolicy(ctx context.Context, d *schema.ResourceData, machinePolicy *machinepolicies.MachinePolicy) error {
	d.Set("connection_connect_timeout", machinePolicy.ConnectionConnectTimeout)
	d.Set("connection_retry_count_limit", machinePolicy.ConnectionRetryCountLimit)
	d.Set("connection_retry_sleep_interval", machinePolicy.ConnectionRetrySleepInterval)
	d.Set("connection_retry_time_limit", machinePolicy.ConnectionRetryTimeLimit)
	d.Set("description", machinePolicy.Description)
	d.Set("id", machinePolicy.GetID())
	d.Set("is_default", machinePolicy.IsDefault)
	d.Set("name", machinePolicy.Name)
	d.Set("polling_request_queue_timeout", machinePolicy.PollingRequestQueueTimeout)
	d.Set("space_id", machinePolicy.SpaceID)

	if err := d.Set("machine_cleanup_policy", flattenMachineCleanupPolicy(machinePolicy.MachineCleanupPolicy)); err != nil {
		return fmt.Errorf("error setting machine_cleanup_policy: %s", err)
	}

	if err := d.Set("machine_connectivity_policy", flattenMachineConnectivityPolicy(machinePolicy.MachineConnectivityPolicy)); err != nil {
		return fmt.Errorf("error setting machine_connectivity_policy: %s", err)
	}

	if err := d.Set("machine_health_check_policy", flattenMachineHealthCheckPolicy(machinePolicy.MachineHealthCheckPolicy)); err != nil {
		return fmt.Errorf("error setting machine_health_check_policy: %s", err)
	}

	if err := d.Set("machine_update_policy", flattenMachineUpdatePolicy(machinePolicy.MachineUpdatePolicy)); err != nil {
		return fmt.Errorf("error setting machine_update_policy: %s", err)
	}

	if err := d.Set("machine_package_cache_retention_policy", flattenMachinePackageCacheRetentionPolicy(machinePolicy.MachinePackageCacheRetentionPolicy)); err != nil {
		return fmt.Errorf("error setting machine_package_cache_retention_policy: %s", err)
	}

	return nil
}

func expandMachinePackageCacheRetentionPolicy(values interface{}) *machinepolicies.MachinePackageCacheRetentionPolicy {
	if values == nil {
		return nil
	}
	flattenedValues := values.(*schema.Set)
	if len(flattenedValues.List()) == 0 {
		return nil
	}

	flattenedMap := flattenedValues.List()[0].(map[string]interface{})

	machinePackageCacheRetentionPolicy := machinepolicies.NewDefaultMachinePackageCacheRetentionPolicy()

	if v, ok := flattenedMap["strategy"]; ok {
		var strategy = v.(string)

		if strategy == "Default" {
			return machinePackageCacheRetentionPolicy
		}

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

	if machineUpdatePolicy.Strategy == "Default" {
		return []interface{}{map[string]interface{}{
			"strategy": machineUpdatePolicy.Strategy,
		}}
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
