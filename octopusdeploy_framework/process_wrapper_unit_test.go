package octopusdeploy_framework

import (
	"strings"
	"testing"
)

func TestSplitImportSpaceID(t *testing.T) {
	cases := []struct {
		name          string
		importID      string
		expectedSpace string
		expectedRest  []string
	}{
		{
			name:         "a bare process keeps working",
			importID:     "deploymentprocess-Projects-123",
			expectedRest: []string{"deploymentprocess-Projects-123"},
		},
		{
			name:         "a process and step keep working",
			importID:     "deploymentprocess-Projects-123:00000000-0000-0000-0000-000000000001",
			expectedRest: []string{"deploymentprocess-Projects-123", "00000000-0000-0000-0000-000000000001"},
		},
		{
			name:          "a leading space is taken off a process",
			importID:      "Spaces-2:deploymentprocess-Projects-123",
			expectedSpace: "Spaces-2",
			expectedRest:  []string{"deploymentprocess-Projects-123"},
		},
		{
			name:          "a leading space is taken off a process and step",
			importID:      "Spaces-2:deploymentprocess-Projects-123:00000000-0000-0000-0000-000000000001",
			expectedSpace: "Spaces-2",
			expectedRest:  []string{"deploymentprocess-Projects-123", "00000000-0000-0000-0000-000000000001"},
		},
		{
			name:          "a leading space is taken off a process, parent and child",
			importID:      "Spaces-12:deploymentprocess-Projects-123:00000000-0000-0000-0000-000000000010:00000000-0000-0000-0000-000000000012",
			expectedSpace: "Spaces-12",
			expectedRest:  []string{"deploymentprocess-Projects-123", "00000000-0000-0000-0000-000000000010", "00000000-0000-0000-0000-000000000012"},
		},
		{
			name:         "a space on its own is not treated as a prefix",
			importID:     "Spaces-2",
			expectedRest: []string{"Spaces-2"},
		},
		{
			name:         "a runbook process is not mistaken for a space",
			importID:     "runbookprocess-Projects-123",
			expectedRest: []string{"runbookprocess-Projects-123"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			space, rest := splitImportSpaceID(tc.importID)

			if space != tc.expectedSpace {
				t.Errorf("expected space %q, got %q", tc.expectedSpace, space)
			}

			if strings.Join(rest, ":") != strings.Join(tc.expectedRest, ":") {
				t.Errorf("expected remainder %v, got %v", tc.expectedRest, rest)
			}
		})
	}
}
