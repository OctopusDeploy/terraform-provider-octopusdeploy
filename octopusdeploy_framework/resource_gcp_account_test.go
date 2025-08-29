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

func TestAccOctopusDeployGCPAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_gcp_account." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	jsonKey := `{\"type\":\"service_account\",\"project_id\":\"test-project-` + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) + `\",\"private_key_id\":\"test-key\"}`
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccGCPAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccGCPAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttrSet(prefix, "json_key"),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttrSet(prefix, "space_id"),
					resource.TestCheckResourceAttr(prefix, "tenants.#", "0"),
				),
				Config: testAccGCPAccountBasic(localName, name, description, jsonKey, tenantedDeploymentParticipation),
			},
		},
	})
}

func TestAccOctopusDeployGCPAccountUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_gcp_account." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	jsonKey := `{\"type\":\"service_account\",\"project_id\":\"test-project-` + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) + `\",\"private_key_id\":\"test-key\"}`
	newJsonKey := `{\"type\":\"service_account\",\"project_id\":\"test-project-` + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) + `\",\"private_key_id\":\"new-test-key\"}`
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccGCPAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccGCPAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttrSet(prefix, "json_key"),
				),
				Config: testAccGCPAccountBasic(localName, name, description, jsonKey, tenantedDeploymentParticipation),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccGCPAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
					resource.TestCheckResourceAttrSet(prefix, "json_key"),
				),
				Config: testAccGCPAccountBasic(localName, name, newDescription, newJsonKey, tenantedDeploymentParticipation),
			},
		},
	})
}

func TestAccOctopusDeployGCPAccountMinimal(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_gcp_account." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	jsonKey := `{\"type\":\"service_account\",\"project_id\":\"test-project-` + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) + `\",\"private_key_id\":\"test-key\"}`

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccGCPAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccGCPAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttrSet(prefix, "json_key"),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccGCPAccountMinimal(localName, name, jsonKey),
			},
		},
	})
}

func TestAccOctopusDeployGCPAccountTenanted(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_gcp_account." + localName
	tenantLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	jsonKey := `{\"type\":\"service_account\",\"project_id\":\"test-project-` + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) + `\",\"private_key_id\":\"test-key\"}`
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccGCPAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccGCPAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttr(prefix, "tenants.#", "1"),
				),
				Config: testAccGCPAccountTenanted(localName, tenantLocalName, name, description, tenantName, jsonKey, tenantedDeploymentParticipation),
			},
		},
	})
}

func TestAccOctopusDeployGCPAccountImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_gcp_account." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	jsonKey := `{\"type\":\"service_account\",\"project_id\":\"test-project-` + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) + `\",\"private_key_id\":\"test-key\"}`
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccGCPAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccGCPAccountBasic(localName, name, description, jsonKey, tenantedDeploymentParticipation),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"json_key"}, // sensitive field
				ImportStateIdFunc: testAccGCPAccountImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccGCPAccountBasic(localName string, name string, description string, jsonKey string, tenantedDeploymentParticipation core.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_gcp_account" "%s" {
		name                              = "%s"
		description                       = "%s"
		json_key                          = "%s"
		tenanted_deployment_participation = "%s"
		tenants                           = []
	}`, localName, name, description, jsonKey, tenantedDeploymentParticipation)
}

func testAccGCPAccountMinimal(localName string, name string, jsonKey string) string {
	return fmt.Sprintf(`resource "octopusdeploy_gcp_account" "%s" {
		name     = "%s"
		json_key = "%s"
	}`, localName, name, jsonKey)
}

func testAccGCPAccountTenanted(localName string, tenantLocalName string, name string, description string, tenantName string, jsonKey string, tenantedDeploymentParticipation core.TenantedDeploymentMode) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_tenant" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_gcp_account" "%s" {
		name                              = "%s"
		description                       = "%s"
		json_key                          = "%s"
		tenanted_deployment_participation = "%s"
		tenants                           = [octopusdeploy_tenant.%s.id]
	}`, tenantLocalName, tenantName, localName, name, description, jsonKey, tenantedDeploymentParticipation, tenantLocalName)
}

func testAccGCPAccountExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		accountID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), accountID); err != nil {
			return err
		}

		return nil
	}
}

func testAccGCPAccountCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_gcp_account" {
			continue
		}

		if account, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("GCP account (%s) still exists", account.GetID())
		}
	}

	return nil
}

func testAccGCPAccountImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}