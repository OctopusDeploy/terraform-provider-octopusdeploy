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

func TestAccOctopusDeployAmazonWebServicesAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_aws_account." + localName

	accessKey := acctest.RandString(acctest.RandIntRange(20, 255))
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secretKey := acctest.RandString(acctest.RandIntRange(20, 255))
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAmazonWebServicesAccountCheckDestroy,
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAmazonWebServicesAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "access_key", accessKey),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "secret_key", secretKey),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
				),
				Config: testAmazonWebServicesAccountBasic(localName, name, description, accessKey, secretKey, tenantedDeploymentParticipation),
			},
		},
	})
}

func TestAccOctopusDeployAmazonWebServicesAccountUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_aws_account." + localName

	accessKey := acctest.RandString(acctest.RandIntRange(20, 255))
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secretKey := acctest.RandString(acctest.RandIntRange(20, 255))
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	newAccessKey := acctest.RandString(acctest.RandIntRange(20, 255))
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAmazonWebServicesAccountCheckDestroy,
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAmazonWebServicesAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "access_key", accessKey),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testAmazonWebServicesAccountBasic(localName, name, description, accessKey, secretKey, tenantedDeploymentParticipation),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAmazonWebServicesAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "access_key", newAccessKey),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testAmazonWebServicesAccountBasic(localName, name, newDescription, newAccessKey, secretKey, tenantedDeploymentParticipation),
			},
		},
	})
}

func TestAccOctopusDeployAmazonWebServicesAccountMinimal(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_aws_account." + localName

	accessKey := acctest.RandString(acctest.RandIntRange(20, 255))
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secretKey := acctest.RandString(acctest.RandIntRange(20, 255))

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAmazonWebServicesAccountCheckDestroy,
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAmazonWebServicesAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "access_key", accessKey),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "secret_key", secretKey),
				),
				Config: testAmazonWebServicesAccountMinimal(localName, name, accessKey, secretKey),
			},
		},
	})
}

func TestAccOctopusDeployAmazonWebServicesAccountImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_aws_account." + localName

	accessKey := acctest.RandString(acctest.RandIntRange(20, 255))
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secretKey := acctest.RandString(acctest.RandIntRange(20, 255))
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAmazonWebServicesAccountCheckDestroy,
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAmazonWebServicesAccountBasic(localName, name, description, accessKey, secretKey, tenantedDeploymentParticipation),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"secret_key",
				},
			},
		},
	})
}

func testAmazonWebServicesAccountBasic(localName string, name string, description string, accessKey string, secretKey string, tenantedDeploymentParticipation core.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_aws_account" "%s" {
		access_key = "%s"
		description = "%s"
		name = "%s"
		secret_key = "%s"
		tenanted_deployment_participation = "%s"
	}`, localName, accessKey, description, name, secretKey, tenantedDeploymentParticipation)
}

func testAmazonWebServicesAccountMinimal(localName string, name string, accessKey string, secretKey string) string {
	return fmt.Sprintf(`resource "octopusdeploy_aws_account" "%s" {
		access_key = "%s"
		name = "%s"
		secret_key = "%s"
	}`, localName, accessKey, name, secretKey)
}

func testAmazonWebServicesAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if _, err := accounts.GetByID(octoClient, rs.Primary.Attributes["space_id"], rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAmazonWebServicesAccountCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_aws_account" {
			continue
		}

		if account, err := accounts.GetByID(octoClient, rs.Primary.Attributes["space_id"], rs.Primary.ID); err == nil {
			return fmt.Errorf("Amazon Web Services account (%s) still exists", account.GetID())
		}
	}

	return nil
}