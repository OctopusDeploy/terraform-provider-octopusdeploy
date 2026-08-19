package octopusdeploy

import (
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	healthCheckScheduleTypeInterval = "Interval"
	healthCheckScheduleTypeCron     = "Cron"
	healthCheckScheduleTypeNever    = "Never"

	defaultHealthCheckInterval = 24 * time.Hour
)

// resolveHealthCheckScheduleType falls back to the schedule implied by the configuration
// when the type is not stated, so existing configurations keep their current behaviour.
func resolveHealthCheckScheduleType(scheduleType string, cron string) string {
	if scheduleType != "" {
		return scheduleType
	}
	if cron != "" {
		return healthCheckScheduleTypeCron
	}
	return healthCheckScheduleTypeInterval
}

func expandMachineHealthCheckPolicy(values interface{}) *machinepolicies.MachineHealthCheckPolicy {
	if values == nil {
		return nil
	}
	flattenedValues := values.(*schema.Set)
	if len(flattenedValues.List()) == 0 {
		return nil
	}

	flattenedMap := flattenedValues.List()[0].(map[string]interface{})

	machineHealthCheckPolicy := machinepolicies.NewMachineHealthCheckPolicy()

	if v, ok := flattenedMap["bash_health_check_policy"]; ok {
		if len(v.([]interface{})) > 0 {
			machineHealthCheckPolicy.BashHealthCheckPolicy = expandMachineScriptPolicy(v)
		}
	}

	if v, ok := flattenedMap["health_check_cron_timezone"]; ok {
		if s := v.(string); len(s) > 0 {
			machineHealthCheckPolicy.HealthCheckCronTimezone = s
		}
	}

	cron, _ := flattenedMap["health_check_cron"].(string)
	interval, _ := flattenedMap["health_check_interval"].(int)
	scheduleType, _ := flattenedMap["health_check_schedule_type"].(string)

	// The server distinguishes the three schedules by which fields are absent, so only
	// the field belonging to the resolved schedule is populated.
	switch resolveHealthCheckScheduleType(scheduleType, cron) {
	case healthCheckScheduleTypeNever:
		machineHealthCheckPolicy.HealthCheckCron = ""
		machineHealthCheckPolicy.HealthCheckInterval = 0
	case healthCheckScheduleTypeCron:
		machineHealthCheckPolicy.HealthCheckCron = cron
		machineHealthCheckPolicy.HealthCheckInterval = 0
	default:
		machineHealthCheckPolicy.HealthCheckCron = ""
		if interval > 0 {
			machineHealthCheckPolicy.HealthCheckInterval = time.Duration(interval)
		} else {
			machineHealthCheckPolicy.HealthCheckInterval = defaultHealthCheckInterval
		}
	}

	if v, ok := flattenedMap["health_check_type"]; ok {
		machineHealthCheckPolicy.HealthCheckType = v.(string)
	}

	if v, ok := flattenedMap["powershell_health_check_policy"]; ok {
		if len(v.([]interface{})) > 0 {
			machineHealthCheckPolicy.PowerShellHealthCheckPolicy = expandMachineScriptPolicy(v)
		}
	}

	return machineHealthCheckPolicy
}

func flattenMachineHealthCheckPolicy(machineHealthCheckPolicy *machinepolicies.MachineHealthCheckPolicy) []interface{} {
	if machineHealthCheckPolicy == nil {
		return nil
	}

	scheduleType := healthCheckScheduleTypeNever
	switch {
	case machineHealthCheckPolicy.HealthCheckInterval > 0:
		scheduleType = healthCheckScheduleTypeInterval
	case machineHealthCheckPolicy.HealthCheckCron != "":
		scheduleType = healthCheckScheduleTypeCron
	}

	return []interface{}{map[string]interface{}{
		"bash_health_check_policy":       flattenMachineScriptPolicy(machineHealthCheckPolicy.BashHealthCheckPolicy),
		"health_check_cron":              machineHealthCheckPolicy.HealthCheckCron,
		"health_check_cron_timezone":     machineHealthCheckPolicy.HealthCheckCronTimezone,
		"health_check_interval":          machineHealthCheckPolicy.HealthCheckInterval,
		"health_check_schedule_type":     scheduleType,
		"health_check_type":              machineHealthCheckPolicy.HealthCheckType,
		"powershell_health_check_policy": flattenMachineScriptPolicy(machineHealthCheckPolicy.PowerShellHealthCheckPolicy),
	}}
}

func getMachineHealthCheckPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bash_health_check_policy": {
			Elem:     &schema.Resource{Schema: getMachineScriptPolicySchema()},
			MaxItems: 1,
			Required: true,
			Type:     schema.TypeList,
		},
		"health_check_cron": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"health_check_cron_timezone": {
			Default:  "UTC",
			Optional: true,
			Type:     schema.TypeString,
		},
		"health_check_interval": {
			Computed:    true,
			Optional:    true,
			Type:        schema.TypeInt,
			Description: "In nanoseconds. Defaults to 24 hours. Ignored unless health_check_schedule_type is Interval.",
		},
		"health_check_schedule_type": {
			Computed:    true,
			Optional:    true,
			Type:        schema.TypeString,
			Description: "The health check schedule: Interval, Cron, or Never. Defaults to Cron when health_check_cron is set, otherwise Interval. Set to Never to disable automatic health checks.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				healthCheckScheduleTypeInterval,
				healthCheckScheduleTypeCron,
				healthCheckScheduleTypeNever,
			}, false)),
		},
		"health_check_type": {
			Default:  "RunScript",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"OnlyConnectivity",
				"RunScript",
			}, false)),
		},
		"powershell_health_check_policy": {
			Elem:     &schema.Resource{Schema: getMachineScriptPolicySchema()},
			MaxItems: 1,
			Required: true,
			Type:     schema.TypeList,
		},
	}
}
