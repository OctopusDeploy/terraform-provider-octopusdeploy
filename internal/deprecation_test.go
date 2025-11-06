package internal

import (
	"os"
	"sync"
	"testing"
)

func TestIsDeprecatedResourceEnabled(t *testing.T) {
	tests := []struct {
		name           string
		envValue       string
		deprecationKey string
		expected       bool
	}{
		{
			name:           "empty environment variable",
			envValue:       "",
			deprecationKey: "Process_v1.0.0",
			expected:       false,
		},
		{
			name:           "single matching key",
			envValue:       "Process_v1.0.0",
			deprecationKey: "Process_v1.0.0",
			expected:       true,
		},
		{
			name:           "multiple keys with match",
			envValue:       "Process_v1.0.0,SomeOther_v2.0.0",
			deprecationKey: "Process_v1.0.0",
			expected:       true,
		},
		{
			name:           "multiple keys no match",
			envValue:       "SomeOther_v1.0.0,AnotherOne_v2.0.0",
			deprecationKey: "Process_v1.0.0",
			expected:       false,
		},
		{
			name:           "whitespace handling",
			envValue:       " Process_v1.0.0 , SomeOther_v2.0.0 ",
			deprecationKey: "Process_v1.0.0",
			expected:       true,
		},
		{
			name:           "empty values in list",
			envValue:       "Process_v1.0.0,,SomeOther_v2.0.0",
			deprecationKey: "Process_v1.0.0",
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deprecationCache = nil
			deprecationCacheOnce = sync.Once{}

			if tt.envValue != "" {
				os.Setenv(DeprecationReversalsEnvVar, tt.envValue)
			} else {
				os.Unsetenv(DeprecationReversalsEnvVar)
			}
			defer os.Unsetenv(DeprecationReversalsEnvVar)

			result := IsDeprecatedResourceEnabled(tt.deprecationKey)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetDeprecatedResourceError(t *testing.T) {
	err := GetDeprecatedResourceError("octopusdeploy_test_resource", "Test_v1.0.0")

	expectedSubstrings := []string{
		"octopusdeploy_test_resource",
		"deprecated and disabled",
		"permanently removed",
		"TF_OCTOPUS_DEPRECATION_REVERSALS=Test_v1.0.0",
	}

	for _, substr := range expectedSubstrings {
		if !contains(err.Error(), substr) {
			t.Errorf("expected error message to contain %q, got: %s", substr, err.Error())
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
