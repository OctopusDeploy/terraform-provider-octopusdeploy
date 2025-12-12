package octopusdeploy_framework

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workerpools"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployDynamicWorkerPoolBasic(t *testing.T) {
	t.Skip("Skipping dynamic worker pool test - worker images do not exist in CI environment")
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_dynamic_worker_pool." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	workerType := "WindowsDefault"
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	isDefault := false
	sortOrder := acctest.RandIntRange(50, 100)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccDynamicWorkerPoolCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccDynamicWorkerPoolExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "worker_type", workerType),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "is_default", strconv.FormatBool(isDefault)),
					resource.TestCheckResourceAttr(prefix, "sort_order", strconv.Itoa(sortOrder)),
					resource.TestCheckResourceAttr(prefix, "space_id", "Spaces-1"),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttr(prefix, "can_add_workers", "true"), // This should be computed
				),
				Config: testAccDynamicWorkerPoolBasic(localName, name, workerType, description, isDefault, sortOrder),
			},
		},
	})
}

func TestAccOctopusDeployDynamicWorkerPoolUpdate(t *testing.T) {
	t.Skip("Skipping dynamic worker pool test - worker images do not exist in CI environment")
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_dynamic_worker_pool." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	workerType := "WindowsDefault"
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	isDefault := false
	sortOrder := acctest.RandIntRange(50, 100)
	newSortOrder := acctest.RandIntRange(101, 200)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccDynamicWorkerPoolCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccDynamicWorkerPoolExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "sort_order", strconv.Itoa(sortOrder)),
				),
				Config: testAccDynamicWorkerPoolBasic(localName, name, workerType, description, isDefault, sortOrder),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccDynamicWorkerPoolExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
					resource.TestCheckResourceAttr(prefix, "sort_order", strconv.Itoa(newSortOrder)),
				),
				Config: testAccDynamicWorkerPoolBasic(localName, newName, workerType, newDescription, isDefault, newSortOrder),
			},
		},
	})
}

func TestAccOctopusDeployDynamicWorkerPoolMinimal(t *testing.T) {
	t.Skip("Skipping dynamic worker pool test - worker images do not exist in CI environment")
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_dynamic_worker_pool." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	workerType := "WindowsDefault"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccDynamicWorkerPoolCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccDynamicWorkerPoolExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "worker_type", workerType),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccDynamicWorkerPoolMinimal(localName, name, workerType),
			},
		},
	})
}

func TestAccOctopusDeployDynamicWorkerPoolImport(t *testing.T) {
	t.Skip("Skipping dynamic worker pool test - worker images do not exist in CI environment")
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_dynamic_worker_pool." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	workerType := "WindowsDefault"
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	isDefault := false
	sortOrder := acctest.RandIntRange(50, 100)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccDynamicWorkerPoolCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDynamicWorkerPoolBasic(localName, name, workerType, description, isDefault, sortOrder),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDynamicWorkerPoolImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccDynamicWorkerPoolBasic(localName, name, workerType, description string, isDefault bool, sortOrder int) string {
	return fmt.Sprintf(`resource "octopusdeploy_dynamic_worker_pool" "%s" {
		name        = "%s"
		worker_type = "%s"
		description = "%s"
		is_default  = %v
		sort_order  = %v
		space_id    = "Spaces-1"
	}`, localName, name, workerType, description, isDefault, sortOrder)
}

func testAccDynamicWorkerPoolMinimal(localName, name, workerType string) string {
	return fmt.Sprintf(`resource "octopusdeploy_dynamic_worker_pool" "%s" {
		name        = "%s"
		worker_type = "%s"
		space_id    = "Spaces-1"
	}`, localName, name, workerType)
}

func testAccDynamicWorkerPoolExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		workerPoolID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := workerpools.GetByID(octoClient, octoClient.GetSpaceID(), workerPoolID); err != nil {
			return err
		}

		return nil
	}
}

func testAccDynamicWorkerPoolCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_dynamic_worker_pool" {
			continue
		}

		if workerPool, err := workerpools.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("dynamic worker pool (%s) still exists", workerPool.GetID())
		}
	}

	return nil
}

func testAccDynamicWorkerPoolImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}