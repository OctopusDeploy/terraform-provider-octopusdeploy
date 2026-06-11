package octopusdeploy_framework

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projectgroups"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// TestSpace is an isolated Octopus space created for the lifetime of a single
// acceptance test. Tests are not well isolated when they all share the default
// space (Spaces-1): they create entities with overlapping names, mutate
// space-wide settings, and contend on the same database rows, which makes the
// suite flaky under parallelism. Giving each test its own space removes that
// shared state entirely.
//
// Usage is intentionally a single line; teardown is registered with
// t.Cleanup, so tests never call delete themselves:
//
//	space := NewTestSpace(t)
//	// use space.ID in HCL and assertions, space.Client for direct API checks
type TestSpace struct {
	// ID is the Octopus ID of the created space (e.g. "Spaces-123"). Use this
	// in HCL `space_id` attributes and TestCheckResourceAttr assertions in
	// place of the hard-coded "Spaces-1".
	ID string
	// Name is the generated space name.
	Name string
	// Client is an Octopus client scoped to this space. Use it for direct API
	// lookups in CheckDestroy / Check functions instead of the package-level
	// octoClient, which is scoped to the default space.
	Client *client.Client
	// LifecycleID is the ID of the space's auto-created default lifecycle.
	LifecycleID string
	// ProjectGroupID is the ID of the space's auto-created default project
	// group.
	ProjectGroupID string
}

// ProviderConfig returns an HCL `provider "octopusdeploy"` block pinned to this
// space. Prepend it to a test's configuration so that the provider performs ALL
// operations (create, read, import, destroy) against the isolated space. This
// is required for operations that do not carry a space_id on the resource
// itself — notably `terraform import`, where the import ID is just the resource
// ID and the provider must already know which space to look in.
func (s *TestSpace) ProviderConfig() string {
	return providerSpaceConfig(s.ID)
}

// providerSpaceConfig returns an HCL `provider "octopusdeploy"` block pinned to
// the given space ID. Config builders that only have a space ID (not the full
// TestSpace) prepend this so every config they emit targets the test space —
// including the empty-config import step, which reuses the provider configured
// by the preceding step.
func providerSpaceConfig(spaceID string) string {
	return fmt.Sprintf(`
provider "octopusdeploy" {
	space_id = "%s"
}
`, spaceID)
}

// NewTestSpace creates a fresh, isolated space with its task queue stopped and
// registers cleanup so the space is removed when the test finishes. It fails
// the test on any error. Acceptance tests only manipulate configuration and
// never need the task queue running, so stopping it avoids the space spending
// resources processing tasks.
//
// The space is created with the calling API key's user as a space manager so
// that the same credentials retain full access to the new space.
func NewTestSpace(t *testing.T) *TestSpace {
	t.Helper()

	if octoClient == nil {
		t.Fatal("octoClient is nil; acceptance tests must run with a shared container (-createSharedContainer=true) or TF_ACC_LOCAL set")
	}

	// Space names are limited to 50 chars; a 20-char alpha suffix keeps the
	// name unique while staying well under the limit.
	name := "TF-" + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	me, err := octoClient.Users.GetMe()
	if err != nil {
		t.Fatalf("failed to get current user while creating test space: %s", err)
	}

	newSpace := spaces.NewSpace(name)
	newSpace.SpaceManagersTeamMembers = []string{me.GetID()}

	createdSpace, err := octoClient.Spaces.Add(newSpace)
	if err != nil {
		t.Fatalf("failed to create test space %q: %s", name, err)
	}

	// A space can't be created with a stopped task queue via the create call;
	// it needs a subsequent update (mirroring the space resource's Create).
	createdSpace.TaskQueueStopped = true
	if createdSpace, err = spaces.Update(octoClient, createdSpace); err != nil {
		t.Fatalf("failed to stop task queue for test space %q: %s", createdSpace.GetID(), err)
	}

	t.Cleanup(func() {
		space, err := spaces.GetByID(octoClient, createdSpace.GetID())
		if err != nil {
			t.Logf("failed to read test space %q for cleanup: %s", createdSpace.GetID(), err)
			return
		}

		if !space.TaskQueueStopped {
			space.TaskQueueStopped = true
			if _, err := spaces.Update(octoClient, space); err != nil {
				t.Logf("failed to stop task queue for test space %q: %s", space.GetID(), err)
				return
			}
		}

		if err := octoClient.Spaces.DeleteByID(space.GetID()); err != nil {
			t.Logf("failed to delete test space %q: %s", space.GetID(), err)
		}
	})

	spaceClient, err := octoclient.CreateClient(os.Getenv("OCTOPUS_URL"), createdSpace.GetID(), os.Getenv("OCTOPUS_APIKEY"))
	if err != nil {
		t.Fatalf("failed to create space-scoped client for %q: %s", createdSpace.GetID(), err)
	}

	// A new space's built-in entities (default lifecycle, project group) are
	// seeded asynchronously by the server. Wait for them to appear, both to
	// expose their space-specific IDs and to ensure the space is ready to
	// accept writes before a test starts applying Terraform. Until they exist,
	// the API returns "Resource is not found ... in the current space context".
	lifecycleID, projectGroupID := waitForSpaceDefaults(t, spaceClient, createdSpace.GetID())

	return &TestSpace{
		ID:             createdSpace.GetID(),
		Name:           createdSpace.GetName(),
		Client:         spaceClient,
		LifecycleID:    lifecycleID,
		ProjectGroupID: projectGroupID,
	}
}

func waitForSpaceDefaults(t *testing.T, spaceClient *client.Client, spaceID string) (lifecycleID string, projectGroupID string) {
	t.Helper()

	deadline := time.Now().Add(1 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		lifecycleID, projectGroupID, lastErr = getSpaceDefaults(spaceClient, spaceID)
		if lastErr == nil {
			return lifecycleID, projectGroupID
		}
		time.Sleep(2 * time.Second)
	}

	t.Fatalf("timed out waiting for default lifecycle and project group in space %q: %s", spaceID, lastErr)
	return "", ""
}

func getSpaceDefaults(spaceClient *client.Client, spaceID string) (lifecycleID string, projectGroupID string, err error) {
	allLifecycles, err := lifecycles.GetAll(spaceClient, spaceID)
	if err != nil {
		return "", "", err
	}
	if len(allLifecycles) == 0 {
		return "", "", fmt.Errorf("no lifecycles found yet")
	}

	allProjectGroups, err := projectgroups.GetAll(spaceClient, spaceID)
	if err != nil {
		return "", "", err
	}
	if len(allProjectGroups) == 0 {
		return "", "", fmt.Errorf("no project groups found yet")
	}

	return allLifecycles[0].GetID(), allProjectGroups[0].GetID(), nil
}
