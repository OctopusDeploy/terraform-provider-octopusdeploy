package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployPollingSubscriptionIdBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_polling_subscription_id." + localName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccPollingSubscriptionIdCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccPollingSubscriptionIdExists(prefix),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttrSet(prefix, "polling_uri"),
					resource.TestCheckResourceAttrWith(prefix, "polling_uri", testAccPollingSubscriptionIdValidateURI),
					resource.TestCheckResourceAttrWith(prefix, "id", testAccPollingSubscriptionIdValidateID),
				),
				Config: testAccPollingSubscriptionIdBasic(localName),
			},
		},
	})
}

func TestAccOctopusDeployPollingSubscriptionIdWithDependencies(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_polling_subscription_id." + localName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccPollingSubscriptionIdCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccPollingSubscriptionIdExists(prefix),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttrSet(prefix, "polling_uri"),
					resource.TestCheckResourceAttr(prefix, "dependencies.project_id", "Projects-1"),
					resource.TestCheckResourceAttr(prefix, "dependencies.environment_id", "Environments-1"),
				),
				Config: testAccPollingSubscriptionIdWithDependencies(localName),
			},
		},
	})
}

func TestAccOctopusDeployPollingSubscriptionIdDependenciesForceNew(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_polling_subscription_id." + localName

	var originalId string

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccPollingSubscriptionIdCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccPollingSubscriptionIdExists(prefix),
					testAccPollingSubscriptionIdStoreId(prefix, &originalId),
					resource.TestCheckResourceAttr(prefix, "dependencies.project_id", "Projects-1"),
				),
				Config: testAccPollingSubscriptionIdWithDependencies(localName),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccPollingSubscriptionIdExists(prefix),
					testAccPollingSubscriptionIdCheckIdChanged(prefix, &originalId),
					resource.TestCheckResourceAttr(prefix, "dependencies.project_id", "Projects-2"),
				),
				Config: testAccPollingSubscriptionIdWithDependenciesChanged(localName),
			},
		},
	})
}

func testAccPollingSubscriptionIdBasic(localName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_polling_subscription_id" "%s" {
	}`, localName)
}

func testAccPollingSubscriptionIdWithDependencies(localName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_polling_subscription_id" "%s" {
		dependencies = {
			project_id = "Projects-1"
			environment_id = "Environments-1"
		}
	}`, localName)
}

func testAccPollingSubscriptionIdWithDependenciesChanged(localName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_polling_subscription_id" "%s" {
		dependencies = {
			project_id = "Projects-2"
			environment_id = "Environments-1"
		}
	}`, localName)
}

func testAccPollingSubscriptionIdExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("Not found: %s", prefix)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No polling subscription ID is set")
		}

		return nil
	}
}

func testAccPollingSubscriptionIdCheckDestroy(s *terraform.State) error {
	// Since this is a local resource that doesn't actually exist in Octopus,
	// and it's properly removed from state on destroy, we don't need to check anything
	return nil
}

func testAccPollingSubscriptionIdValidateURI(value string) error {
	// URI should be in format: poll://[ID]/
	if len(value) < 8 || value[:7] != "poll://" || value[len(value)-1:] != "/" {
		return fmt.Errorf("polling_uri should be in format 'poll://[ID]/', got: %s", value)
	}
	return nil
}

func testAccPollingSubscriptionIdValidateID(value string) error {
	// ID should be a 20-character string
	if len(value) != 20 {
		return fmt.Errorf("id should be 20 characters long, got: %d characters", len(value))
	}
	return nil
}

func testAccPollingSubscriptionIdStoreId(prefix string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("Not found: %s", prefix)
		}
		*id = rs.Primary.ID
		return nil
	}
}

func testAccPollingSubscriptionIdCheckIdChanged(prefix string, originalId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("Not found: %s", prefix)
		}
		
		if rs.Primary.ID == *originalId {
			return fmt.Errorf("polling subscription ID should have changed due to dependencies ForceNew, but it didn't")
		}
		
		return nil
	}
}