package octopusdeploy_framework

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/stretchr/testify/require"
)

func runbookWithID(id string) *runbooks.Runbook {
	runbook := runbooks.NewRunbook("runbook-"+id, "Projects-1")
	runbook.Resource = resources.Resource{ID: id}
	return runbook
}

func TestFilterRunbooksByID(t *testing.T) {
	all := []*runbooks.Runbook{runbookWithID("Runbooks-1"), runbookWithID("Runbooks-2"), runbookWithID("Runbooks-3")}

	t.Run("no ids returns everything", func(t *testing.T) {
		require.Len(t, filterRunbooksByID(all, nil), 3)
		require.Len(t, filterRunbooksByID(all, []string{}), 3)
	})

	t.Run("keeps only the requested ids", func(t *testing.T) {
		filtered := filterRunbooksByID(all, []string{"Runbooks-3", "Runbooks-1"})

		require.Len(t, filtered, 2)
		require.Equal(t, "Runbooks-1", filtered[0].GetID())
		require.Equal(t, "Runbooks-3", filtered[1].GetID())
	})

	t.Run("unknown ids match nothing", func(t *testing.T) {
		require.Empty(t, filterRunbooksByID(all, []string{"Runbooks-99"}))
	})
}
