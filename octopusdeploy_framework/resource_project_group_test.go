package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projectgroups"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployProjectGroupBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project_group." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:                 testProjectGroupCheckDestroy,
		ProtoV6ProviderFactories:     ProtoV6ProviderFactories(),
		PreCheck:                     func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployProjectGroupExists(prefix),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
				),
				Config: testProjectGroupBasic(localName, name, description),
			},
		},
	})
}

func TestAccOctopusDeployProjectGroupUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project_group." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:                 testProjectGroupCheckDestroy,
		ProtoV6ProviderFactories:     ProtoV6ProviderFactories(),
		PreCheck:                     func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployProjectGroupExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
				),
				Config: testProjectGroupBasic(localName, name, description),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployProjectGroupExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
				),
				Config: testProjectGroupBasic(localName, newName, newDescription),
			},
		},
	})
}

func TestAccOctopusDeployProjectGroupMinimal(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project_group." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:                 testProjectGroupCheckDestroy,
		ProtoV6ProviderFactories:     ProtoV6ProviderFactories(),
		PreCheck:                     func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployProjectGroupExists(prefix),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", ""),
				),
				Config: testProjectGroupMinimal(localName, name),
			},
		},
	})
}

func TestAccOctopusDeployProjectGroupImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_project_group." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:                 testProjectGroupCheckDestroy,
		ProtoV6ProviderFactories:     ProtoV6ProviderFactories(),
		PreCheck:                     func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testProjectGroupBasic(localName, name, description),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"space_id",
				},
			},
		},
	})
}

func testProjectGroupBasic(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project_group" "%s" {
		name = "%s"
		description = "%s"
	}`, localName, name, description)
}

func testProjectGroupMinimal(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project_group" "%s" {
		name = "%s"
	}`, localName, name)
}

func testOctopusDeployProjectGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if _, err := projectgroups.GetByID(octoClient, rs.Primary.Attributes["space_id"], rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testProjectGroupCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_project_group" {
			continue
		}

		if projectGroup, err := projectgroups.GetByID(octoClient, rs.Primary.Attributes["space_id"], rs.Primary.ID); err == nil {
			return fmt.Errorf("project group (%s) still exists", projectGroup.GetID())
		}
	}

	return nil
}