package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/platformhubaccounts"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccPlatformHubAwsOpenIDConnectAccountCreate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_aws_openid_connect_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	roleArn := "arn:aws:iam::123456789012:role/TestRole"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubAwsOpenIDConnectAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubAwsOpenIDConnectAccountBasic(localName, name, description, roleArn),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubAwsOpenIDConnectAccountExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "role_arn", roleArn),
					resource.TestCheckResourceAttr(resourceName, "session_duration", "3600"),
				),
			},
		},
	})
}

func TestAccPlatformHubAwsOpenIDConnectAccountUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_aws_openid_connect_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	roleArn := "arn:aws:iam::123456789012:role/TestRole"

	updatedName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedRoleArn := "arn:aws:iam::123456789012:role/UpdatedRole"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubAwsOpenIDConnectAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubAwsOpenIDConnectAccountBasic(localName, name, description, roleArn),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubAwsOpenIDConnectAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "role_arn", roleArn),
				),
			},
			{
				Config: testPlatformHubAwsOpenIDConnectAccountBasic(localName, updatedName, updatedDescription, updatedRoleArn),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubAwsOpenIDConnectAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "role_arn", updatedRoleArn),
				),
			},
		},
	})
}

func TestAccPlatformHubAwsOpenIDConnectAccountImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_aws_openid_connect_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	roleArn := "arn:aws:iam::123456789012:role/TestRole"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubAwsOpenIDConnectAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubAwsOpenIDConnectAccountBasic(localName, name, description, roleArn),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubAwsOpenIDConnectAccountExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testPlatformHubAwsOpenIDConnectAccountBasic(localName string, name string, description string, roleArn string) string {
	return fmt.Sprintf(`resource "octopusdeploy_platform_hub_aws_openid_connect_account" "%s" {
		name                     = "%s"
		description              = "%s"
		role_arn                 = "%s"
		session_duration         = "3600"
		execution_subject_keys   = ["space", "environment"]
		health_subject_keys      = ["space"]
		account_test_subject_keys = ["space"]
	}`, localName, name, description, roleArn)
}

func testPlatformHubAwsOpenIDConnectAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Platform Hub AWS OpenID Connect Account ID is set")
		}

		client := octoClient
		account, err := platformhubaccounts.GetByID(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error retrieving Platform Hub AWS OpenID Connect Account: %s", err)
		}

		if account.GetID() != rs.Primary.ID {
			return fmt.Errorf("Platform Hub AWS OpenID Connect Account not found")
		}

		return nil
	}
}

func testPlatformHubAwsOpenIDConnectAccountCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_platform_hub_aws_openid_connect_account" {
			continue
		}

		client := octoClient
		account, err := platformhubaccounts.GetByID(client, rs.Primary.ID)
		if err == nil && account != nil {
			return fmt.Errorf("Platform Hub AWS OpenID Connect Account (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
