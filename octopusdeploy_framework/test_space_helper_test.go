package octopusdeploy_framework

import (
	"os"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
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
}

// NewTestSpace creates a fresh, isolated space with its task queue stopped and
// registers cleanup so the space is removed when the test finishes. Most
// acceptance tests only manipulate configuration and never need the task queue
// running, so stopping it avoids the space spending resources processing tasks.
// Use NewTestSpaceWithTaskQueue if a test needs the queue running (e.g. it
// triggers deployments or runbook runs).
func NewTestSpace(t *testing.T) *TestSpace {
	return NewTestSpaceWithTaskQueue(t, false)
}

// NewTestSpaceWithTaskQueue creates a fresh, isolated space and registers
// cleanup so the space is removed when the test finishes. It fails the test on
// any error. When taskQueueEnabled is false the space is created with its task
// queue stopped.
//
// The space is created with the calling API key's user as a space manager so
// that the same credentials retain full access to the new space.
func NewTestSpaceWithTaskQueue(t *testing.T, taskQueueEnabled bool) *TestSpace {
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
	if !taskQueueEnabled {
		createdSpace.TaskQueueStopped = true
		if createdSpace, err = spaces.Update(octoClient, createdSpace); err != nil {
			t.Fatalf("failed to stop task queue for test space %q: %s", createdSpace.GetID(), err)
		}
	}

	t.Cleanup(func() {
		// A space can only be deleted while its task queue is stopped. It is
		// stopped at creation by default, but a test may have started it (or
		// requested it enabled), so ensure it is stopped before deleting.
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

	// Build a client scoped to the new space for direct API checks. The URL
	// and API key are set on the environment by TestMain for both the shared
	// container and TF_ACC_LOCAL paths.
	spaceClient, err := octoclient.CreateClient(os.Getenv("OCTOPUS_URL"), createdSpace.GetID(), os.Getenv("OCTOPUS_APIKEY"))
	if err != nil {
		t.Fatalf("failed to create space-scoped client for %q: %s", createdSpace.GetID(), err)
	}

	return &TestSpace{
		ID:     createdSpace.GetID(),
		Name:   createdSpace.GetName(),
		Client: spaceClient,
	}
}
