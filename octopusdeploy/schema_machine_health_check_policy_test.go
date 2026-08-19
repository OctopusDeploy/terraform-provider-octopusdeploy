package octopusdeploy

import (
	"testing"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func healthCheckPolicySet(values map[string]interface{}) *schema.Set {
	resource := &schema.Resource{Schema: getMachineHealthCheckPolicySchema()}

	block := map[string]interface{}{
		"bash_health_check_policy":       []interface{}{},
		"health_check_cron":              "",
		"health_check_cron_timezone":     "UTC",
		"health_check_interval":          0,
		"health_check_schedule_type":     "",
		"health_check_type":              "RunScript",
		"powershell_health_check_policy": []interface{}{},
	}
	for k, v := range values {
		block[k] = v
	}

	return schema.NewSet(schema.HashResource(resource), []interface{}{block})
}

func TestExpandDefaultsToIntervalWhenNothingSpecified(t *testing.T) {
	policy := expandMachineHealthCheckPolicy(healthCheckPolicySet(nil))

	require.NotNil(t, policy)
	assert.Equal(t, 24*time.Hour, policy.HealthCheckInterval, "omitting the interval must keep the historic 24h default")
	assert.Empty(t, policy.HealthCheckCron)
}

func TestExpandHonoursExplicitInterval(t *testing.T) {
	policy := expandMachineHealthCheckPolicy(healthCheckPolicySet(map[string]interface{}{
		"health_check_interval": int(6 * time.Hour),
	}))

	require.NotNil(t, policy)
	assert.Equal(t, 6*time.Hour, policy.HealthCheckInterval)
	assert.Empty(t, policy.HealthCheckCron)
}

func TestExpandCronOmitsInterval(t *testing.T) {
	policy := expandMachineHealthCheckPolicy(healthCheckPolicySet(map[string]interface{}{
		"health_check_cron": "0 0 0 1 1 * 2099",
	}))

	require.NotNil(t, policy)
	assert.Equal(t, "0 0 0 1 1 * 2099", policy.HealthCheckCron)
	assert.Equal(t, time.Duration(0), policy.HealthCheckInterval, "a cron schedule must not also send an interval")
}

func TestExpandCronWinsOverExplicitInterval(t *testing.T) {
	policy := expandMachineHealthCheckPolicy(healthCheckPolicySet(map[string]interface{}{
		"health_check_cron":     "0 0 0 1 1 * 2099",
		"health_check_interval": int(24 * time.Hour),
	}))

	require.NotNil(t, policy)
	assert.Equal(t, "0 0 0 1 1 * 2099", policy.HealthCheckCron)
	assert.Equal(t, time.Duration(0), policy.HealthCheckInterval)
}

func TestExpandNeverOmitsBoth(t *testing.T) {
	policy := expandMachineHealthCheckPolicy(healthCheckPolicySet(map[string]interface{}{
		"health_check_schedule_type": healthCheckScheduleTypeNever,
	}))

	require.NotNil(t, policy)
	assert.Empty(t, policy.HealthCheckCron)
	assert.Equal(t, time.Duration(0), policy.HealthCheckInterval)
}

func TestExpandNeverIgnoresLeftoverIntervalAndCron(t *testing.T) {
	policy := expandMachineHealthCheckPolicy(healthCheckPolicySet(map[string]interface{}{
		"health_check_schedule_type": healthCheckScheduleTypeNever,
		"health_check_cron":          "0 0 0 1 1 * 2099",
		"health_check_interval":      int(24 * time.Hour),
	}))

	require.NotNil(t, policy)
	assert.Empty(t, policy.HealthCheckCron)
	assert.Equal(t, time.Duration(0), policy.HealthCheckInterval)
}

func TestExpandExplicitIntervalTypeIgnoresCron(t *testing.T) {
	policy := expandMachineHealthCheckPolicy(healthCheckPolicySet(map[string]interface{}{
		"health_check_schedule_type": healthCheckScheduleTypeInterval,
		"health_check_cron":          "0 0 0 1 1 * 2099",
		"health_check_interval":      int(time.Hour),
	}))

	require.NotNil(t, policy)
	assert.Empty(t, policy.HealthCheckCron)
	assert.Equal(t, time.Hour, policy.HealthCheckInterval)
}

func TestFlattenDerivesScheduleType(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		interval time.Duration
		cron     string
		expected string
	}{
		{"interval", 24 * time.Hour, "", healthCheckScheduleTypeInterval},
		{"cron", 0, "0 0 0 1 1 * 2099", healthCheckScheduleTypeCron},
		{"never", 0, "", healthCheckScheduleTypeNever},
		{"interval wins when the server somehow returns both", 24 * time.Hour, "0 0 * * *", healthCheckScheduleTypeInterval},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			policy := machinepolicies.NewMachineHealthCheckPolicy()
			policy.HealthCheckInterval = testCase.interval
			policy.HealthCheckCron = testCase.cron

			flattened := flattenMachineHealthCheckPolicy(policy)

			require.Len(t, flattened, 1)
			assert.Equal(t, testCase.expected, flattened[0].(map[string]interface{})["health_check_schedule_type"])
		})
	}
}

func TestResolveHealthCheckScheduleType(t *testing.T) {
	assert.Equal(t, healthCheckScheduleTypeInterval, resolveHealthCheckScheduleType("", ""))
	assert.Equal(t, healthCheckScheduleTypeCron, resolveHealthCheckScheduleType("", "0 0 * * *"))
	assert.Equal(t, healthCheckScheduleTypeNever, resolveHealthCheckScheduleType(healthCheckScheduleTypeNever, "0 0 * * *"))
	assert.Equal(t, healthCheckScheduleTypeInterval, resolveHealthCheckScheduleType(healthCheckScheduleTypeInterval, "0 0 * * *"))
}
