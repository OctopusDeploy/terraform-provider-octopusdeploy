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

func TestAccOctopusDeployStaticWorkerPoolBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_static_worker_pool." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	isDefault := false
	sortOrder := acctest.RandIntRange(50, 100)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccStaticWorkerPoolCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccStaticWorkerPoolExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "is_default", strconv.FormatBool(isDefault)),
					resource.TestCheckResourceAttr(prefix, "sort_order", strconv.Itoa(sortOrder)),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttrSet(prefix, "space_id"),
				),
				Config: testAccStaticWorkerPoolBasic(localName, name, description, isDefault, sortOrder),
			},
		},
	})
}

func TestAccOctopusDeployStaticWorkerPoolUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_static_worker_pool." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	isDefault := false
	sortOrder := acctest.RandIntRange(50, 100)
	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccStaticWorkerPoolCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccStaticWorkerPoolExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "sort_order", strconv.Itoa(sortOrder)),
				),
				Config: testAccStaticWorkerPoolBasic(localName, name, description, isDefault, sortOrder),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccStaticWorkerPoolExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
					resource.TestCheckResourceAttr(prefix, "sort_order", strconv.Itoa(sortOrder)), // Keep same sort order
				),
				Config: testAccStaticWorkerPoolBasic(localName, name, newDescription, isDefault, sortOrder), // Use original sort order
			},
		},
	})
}

func TestAccOctopusDeployStaticWorkerPoolMinimal(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_static_worker_pool." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccStaticWorkerPoolCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccStaticWorkerPoolExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccStaticWorkerPoolMinimal(localName, name),
			},
		},
	})
}

func TestAccOctopusDeployStaticWorkerPoolImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_static_worker_pool." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	isDefault := false
	sortOrder := acctest.RandIntRange(50, 100)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccStaticWorkerPoolCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccStaticWorkerPoolBasic(localName, name, description, isDefault, sortOrder),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStaticWorkerPoolImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccStaticWorkerPoolBasic(localName string, name string, description string, isDefault bool, sortOrder int) string {
	return fmt.Sprintf(`resource "octopusdeploy_static_worker_pool" "%s" {
		name        = "%s"
		description = "%s"
		is_default  = %v
		sort_order  = %v
	}`, localName, name, description, isDefault, sortOrder)
}

func testAccStaticWorkerPoolMinimal(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_static_worker_pool" "%s" {
		name = "%s"
	}`, localName, name)
}

func testAccStaticWorkerPoolExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		workerPoolID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := workerpools.GetByID(octoClient, octoClient.GetSpaceID(), workerPoolID); err != nil {
			return err
		}

		return nil
	}
}

func testAccStaticWorkerPoolCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_static_worker_pool" {
			continue
		}

		if workerPool, err := workerpools.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("static worker pool (%s) still exists", workerPool.GetID())
		}
	}

	return nil
}

func testAccStaticWorkerPoolImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}