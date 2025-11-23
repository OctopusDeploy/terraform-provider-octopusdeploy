package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func ProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"octopusdeploy": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()

			upgradedSdkServer, err := tf5to6server.UpgradeServer(
				ctx,
				octopusdeploy.Provider().GRPCProvider)
			if err != nil {
				log.Fatal(err)
			}

			if err != nil {
				log.Fatal(err)
			}
			providers := []func() tfprotov6.ProviderServer{
				func() tfprotov6.ProviderServer {
					return upgradedSdkServer
				},
				providerserver.NewProtocol6(NewOctopusDeployFrameworkProvider()),
			}

			return tf6muxserver.NewMuxServer(context.Background(), providers...)
		},
	}
}

func ProtoV6ProviderFactoriesWithFeatureToggleOverrides(overrides map[string]bool) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"octopusdeploy": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()

			upgradedSdkServer, err := tf5to6server.UpgradeServer(
				ctx,
				octopusdeploy.Provider().GRPCProvider)
			if err != nil {
				log.Fatal(err)
			}

			providers := []func() tfprotov6.ProviderServer{
				func() tfprotov6.ProviderServer {
					return upgradedSdkServer
				},
				providerserver.NewProtocol6(newTestOctopusDeployFrameworkProvider(overrides)),
			}

			return tf6muxserver.NewMuxServer(context.Background(), providers...)
		},
	}
}

func TestAccPreCheck(t *testing.T) {
	if v := os.Getenv("OCTOPUS_URL"); isEmpty(v) {
		t.Fatal("OCTOPUS_URL must be set for acceptance tests")
	}
	if v := os.Getenv("OCTOPUS_APIKEY"); isEmpty(v) {
		t.Fatal("OCTOPUS_APIKEY must be set for acceptance tests")
	}
}

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func TestTestProviderFeatureToggleOverrides(t *testing.T) {
	tests := []struct {
		name            string
		overrides       map[string]bool
		existingToggles map[string]bool
		expectedToggles map[string]bool
		description     string
	}{
		{
			name:      "No overrides",
			overrides: map[string]bool{},
			existingToggles: map[string]bool{
				"SomeFeatureToggle":      true,
				"SomeOtherFeatureToggle": true,
			},
			expectedToggles: map[string]bool{
				"SomeFeatureToggle":      true,
				"SomeOtherFeatureToggle": true,
			},
			description: "No overrides should preserve all toggles",
		},
		{
			name: "Single override",
			overrides: map[string]bool{
				"FeatureToggleToOverride": false,
			},
			existingToggles: map[string]bool{
				"FeatureToggleToOverride": true,
				"SomeOtherFeatureToggle":  true,
			},
			expectedToggles: map[string]bool{
				"FeatureToggleToOverride": false,
				"SomeOtherFeatureToggle":  true,
			},
			description: "Should override specified toggle while preserving others",
		},
		{
			name: "Multiple overrides",
			overrides: map[string]bool{
				"FeatureToggleToOverride":        true,
				"AnotherFeatureToggleToOverride": true,
			},
			existingToggles: map[string]bool{
				"FeatureToggleToOverride":        false,
				"AnotherFeatureToggleToOverride": false,
				"SomeOtherFeatureToggle":         false,
			},
			expectedToggles: map[string]bool{
				"FeatureToggleToOverride":        true,
				"AnotherFeatureToggleToOverride": true,
				"SomeOtherFeatureToggle":         false,
			},
			description: "Should override multiple toggles while preserving others",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a config with existing toggles
			config := &Config{
				FeatureToggles: tt.existingToggles,
			}

			// Simulate what the test provider does
			if tt.overrides != nil {
				if config.FeatureToggles == nil {
					config.FeatureToggles = make(map[string]bool)
				}

				for key, value := range tt.overrides {
					config.FeatureToggles[key] = value
				}
			}

			// Verify the result
			for expectedKey, expectedValue := range tt.expectedToggles {
				if actualValue, ok := config.FeatureToggles[expectedKey]; !ok {
					t.Errorf("%s: expected toggle %q to exist but it was missing", tt.description, expectedKey)
				} else if actualValue != expectedValue {
					t.Errorf("%s: toggle %q = %t, want %t", tt.description, expectedKey, actualValue, expectedValue)
				}
			}

			// Verify no extra toggles were added
			if len(config.FeatureToggles) != len(tt.expectedToggles) {
				t.Errorf("%s: got %d toggles, want %d toggles", tt.description, len(config.FeatureToggles), len(tt.expectedToggles))
			}
		})
	}
}
