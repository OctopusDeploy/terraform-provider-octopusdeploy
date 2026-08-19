package octopusdeploy

import (
	"encoding/json"
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
		"health_check_interval":          int(24 * time.Hour),
		"health_check_type":              "RunScript",
		"powershell_health_check_policy": []interface{}{},
	}
	for k, v := range values {
		block[k] = v
	}

	return schema.NewSet(schema.HashResource(resource), []interface{}{block})
}

func TestExpandKeepsConfiguredInterval(t *testing.T) {
	policy := expandMachineHealthCheckPolicy(healthCheckPolicySet(nil))

	require.NotNil(t, policy)
	assert.Equal(t, 24*time.Hour, policy.HealthCheckInterval)
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

func TestExpandZeroIntervalMeansNoHealthChecks(t *testing.T) {
	policy := expandMachineHealthCheckPolicy(healthCheckPolicySet(map[string]interface{}{
		"health_check_interval": 0,
	}))

	require.NotNil(t, policy)
	assert.Equal(t, time.Duration(0), policy.HealthCheckInterval)
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
		"health_check_interval": int(6 * time.Hour),
	}))

	require.NotNil(t, policy)
	assert.Equal(t, "0 0 0 1 1 * 2099", policy.HealthCheckCron)
	assert.Equal(t, time.Duration(0), policy.HealthCheckInterval)
}

// clientOmitsZeroInterval reports whether the pinned go-octopusdeploy can express "no health
// checks". Before v2.116.0 it stringified a zero interval as "00:00:00", which the server reads
// as an interval of zero rather than as an absent one.
func clientOmitsZeroInterval(t *testing.T) bool {
	probe := machinepolicies.NewMachineHealthCheckPolicy()
	probe.HealthCheckInterval = 0

	data, err := json.Marshal(probe)
	require.NoError(t, err)

	var fields map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &fields))

	_, present := fields["HealthCheckInterval"]
	return !present
}

// The expanded struct is only half the contract: a zero interval is correct solely because the
// client omits it. Sending "00:00:00" instead would health check every target roughly every minute.
func TestExpandedPolicyOmitsIntervalOnTheWire(t *testing.T) {
	if !clientOmitsZeroInterval(t) {
		t.Skip("pinned go-octopusdeploy still serialises a zero interval as \"00:00:00\"; bump the go.mod pin to run this")
	}

	for _, testCase := range []struct {
		name  string
		block map[string]interface{}
	}{
		{"no health checks", map[string]interface{}{"health_check_interval": 0}},
		{"cron schedule", map[string]interface{}{"health_check_cron": "0 0 0 1 1 * 2099"}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			policy := expandMachineHealthCheckPolicy(healthCheckPolicySet(testCase.block))

			data, err := json.Marshal(policy)
			require.NoError(t, err)

			var fields map[string]interface{}
			require.NoError(t, json.Unmarshal(data, &fields))

			_, present := fields["HealthCheckInterval"]
			assert.False(t, present, "HealthCheckInterval must be absent, not \"00:00:00\"")
		})
	}
}

func TestFlattenRoundTripsThroughExpand(t *testing.T) {
	resource := &schema.Resource{Schema: getMachineHealthCheckPolicySchema()}

	for _, testCase := range []struct {
		name             string
		interval         time.Duration
		cron             string
		expectedInterval time.Duration
		expectedCron     string
	}{
		{"no health checks", 0, "", 0, ""},
		{"cron", 0, "0 0 0 1 1 * 2099", 0, "0 0 0 1 1 * 2099"},
		{"interval", 6 * time.Hour, "", 6 * time.Hour, ""},
		{"policy carrying both keeps the cron", 24 * time.Hour, "0 0 0 1 1 * 2099", 0, "0 0 0 1 1 * 2099"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			original := machinepolicies.NewMachineHealthCheckPolicy()
			original.HealthCheckInterval = testCase.interval
			original.HealthCheckCron = testCase.cron

			flattened := flattenMachineHealthCheckPolicy(original)

			// The script policy blocks carry pointers that the set hasher cannot serialize, and
			// they are not part of the schedule being round tripped.
			block := flattened[0].(map[string]interface{})
			block["bash_health_check_policy"] = []interface{}{}
			block["powershell_health_check_policy"] = []interface{}{}
			block["health_check_interval"] = int(original.HealthCheckInterval)

			reExpanded := expandMachineHealthCheckPolicy(schema.NewSet(schema.HashResource(resource), []interface{}{block}))

			require.NotNil(t, reExpanded)
			assert.Equal(t, testCase.expectedInterval, reExpanded.HealthCheckInterval)
			assert.Equal(t, testCase.expectedCron, reExpanded.HealthCheckCron)
		})
	}
}
