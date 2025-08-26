package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployAWSOpenIDConnectAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_aws_openid_connect_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	roleArn := "arn:aws:iam::123456789012:role/test-role"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAWSOpenIDConnectAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAWSOpenIDConnectAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "role_arn", roleArn),
					resource.TestCheckResourceAttr(prefix, "space_id", "Spaces-1"),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(core.TenantedDeploymentModeUntenanted)),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccAWSOpenIDConnectAccountBasic(localName, name, description, roleArn),
			},
		},
	})
}

func TestAccOctopusDeployAWSOpenIDConnectAccountUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_aws_openid_connect_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	roleArn := "arn:aws:iam::123456789012:role/test-role"
	newRoleArn := "arn:aws:iam::123456789012:role/updated-test-role"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAWSOpenIDConnectAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAWSOpenIDConnectAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "role_arn", roleArn),
				),
				Config: testAccAWSOpenIDConnectAccountBasic(localName, name, description, roleArn),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAWSOpenIDConnectAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
					resource.TestCheckResourceAttr(prefix, "role_arn", newRoleArn),
				),
				Config: testAccAWSOpenIDConnectAccountBasic(localName, newName, newDescription, newRoleArn),
			},
		},
	})
}

func TestAccOctopusDeployAWSOpenIDConnectAccountWithSubjectKeys(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_aws_openid_connect_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	roleArn := "arn:aws:iam::123456789012:role/test-role"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAWSOpenIDConnectAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAWSOpenIDConnectAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "execution_subject_keys.#", "1"),
					resource.TestCheckResourceAttr(prefix, "execution_subject_keys.0", "space"),
					resource.TestCheckResourceAttr(prefix, "health_subject_keys.#", "1"),
					resource.TestCheckResourceAttr(prefix, "health_subject_keys.0", "target"),
					resource.TestCheckResourceAttr(prefix, "account_test_subject_keys.#", "1"),
					resource.TestCheckResourceAttr(prefix, "account_test_subject_keys.0", "type"),
				),
				Config: testAccAWSOpenIDConnectAccountWithSubjectKeys(localName, name, description, roleArn),
			},
		},
	})
}

func TestAccOctopusDeployAWSOpenIDConnectAccountImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_aws_openid_connect_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	roleArn := "arn:aws:iam::123456789012:role/test-role"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAWSOpenIDConnectAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAWSOpenIDConnectAccountBasic(localName, name, description, roleArn),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAWSOpenIDConnectAccountImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccAWSOpenIDConnectAccountBasic(localName, name, description, roleArn string) string {
	return fmt.Sprintf(`resource "octopusdeploy_aws_openid_connect_account" "%s" {
		name         = "%s"
		description  = "%s"
		role_arn     = "%s"
		space_id     = "Spaces-1"
	}`, localName, name, description, roleArn)
}

func testAccAWSOpenIDConnectAccountWithSubjectKeys(localName, name, description, roleArn string) string {
	return fmt.Sprintf(`resource "octopusdeploy_aws_openid_connect_account" "%s" {
		name                     = "%s"
		description              = "%s"
		role_arn                 = "%s"
		space_id                 = "Spaces-1"
		execution_subject_keys   = ["space"]
		health_subject_keys      = ["target"]
		account_test_subject_keys = ["type"]
	}`, localName, name, description, roleArn)
}

func testAccAWSOpenIDConnectAccountExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		accountID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), accountID); err != nil {
			return err
		}

		return nil
	}
}

func testAccAWSOpenIDConnectAccountCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_aws_openid_connect_account" {
			continue
		}

		if account, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("AWS OIDC account (%s) still exists", account.GetID())
		}
	}

	return nil
}

func testAccAWSOpenIDConnectAccountImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}