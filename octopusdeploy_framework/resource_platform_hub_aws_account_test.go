package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/platformhubaccounts"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccPlatformHubAwsAccountCreate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_aws_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accessKey := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubAwsAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubAwsAccountBasic(localName, name, description, accessKey),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubAwsAccountExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "access_key", accessKey),
				),
			},
		},
	})
}

func TestAccPlatformHubAwsAccountUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_aws_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accessKey := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	updatedName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedAccessKey := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubAwsAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubAwsAccountBasic(localName, name, description, accessKey),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubAwsAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "access_key", accessKey),
				),
			},
			{
				Config: testPlatformHubAwsAccountBasic(localName, updatedName, updatedDescription, updatedAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubAwsAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "access_key", updatedAccessKey),
				),
			},
		},
	})
}

func TestAccPlatformHubAwsAccountImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_aws_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accessKey := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubAwsAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubAwsAccountBasic(localName, name, description, accessKey),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubAwsAccountExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
		},
	})
}

func testPlatformHubAwsAccountBasic(localName string, name string, description string, accessKey string) string {
	return fmt.Sprintf(`resource "octopusdeploy_platform_hub_aws_account" "%s" {
		name        = "%s"
		description = "%s"
		access_key  = "%s"
		secret_key  = "test_secret_key_123456789"
	}`, localName, name, description, accessKey)
}

func testPlatformHubAwsAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Platform Hub AWS Account ID is set")
		}

		client := octoClient
		account, err := platformhubaccounts.GetByID(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error retrieving Platform Hub AWS Account: %s", err)
		}

		if account.GetID() != rs.Primary.ID {
			return fmt.Errorf("Platform Hub AWS Account not found")
		}

		return nil
	}
}

func testPlatformHubAwsAccountCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_platform_hub_aws_account" {
			continue
		}

		client := octoClient
		account, err := platformhubaccounts.GetByID(client, rs.Primary.ID)
		if err == nil && account != nil {
			return fmt.Errorf("Platform Hub AWS Account (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
