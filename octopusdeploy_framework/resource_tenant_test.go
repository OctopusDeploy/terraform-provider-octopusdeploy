package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
)

func TestAccTenantBasic(t *testing.T) {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_tenant." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTenantCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "is_disabled", strconv.FormatBool(false)),
				),
				Config: testAccTenantBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, localName, name, description, false),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", newDescription),
					resource.TestCheckResourceAttr(resourceName, "is_disabled", strconv.FormatBool(false)),
				),
				Config: testAccTenantBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, localName, name, newDescription, false),
			},
		},
	})
}

func testAccTenantBasic(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, projectDescription string, environmentLocalName string, environmentName string, localName string, name string, description string, isDisabled bool) string {
	allowDynamicInfrastructure := false
	environmentDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, 2)+"\n"+
		testAccEnvironment(environmentLocalName, environmentName, environmentDescription, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+`
	resource "octopusdeploy_tenant" "%s" {
		description = "%s"
		name        = "%s"
		is_disabled = %v
	}

	resource "octopusdeploy_tenant_project" "project_environment" {
		tenant_id = octopusdeploy_tenant.%s.id
		project_id   = "${octopusdeploy_project.%s.id}"
		environment_ids = ["${octopusdeploy_environment.%s.id}"]
	}`, localName, description, name, isDisabled, localName, projectLocalName, environmentLocalName)
}

// Covers https://github.com/OctopusDeploy/terraform-provider-octopusdeploy/issues/236:
// tags written in config have to reach the server, tags changed on the server have
// to show up as drift, and an empty list has to clear them.
func TestAccTenantTags(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_tenant." + localName
	firstTag := localName + "/first"
	secondTag := localName + "/second"

	var tenantID string

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTenantCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTenantWithTags(localName, "[octopusdeploy_tag.first.canonical_tag_name]"),
				Check: resource.ComposeTestCheckFunc(
					testTenantExists(resourceName),
					captureTenantID(resourceName, &tenantID),
					resource.TestCheckResourceAttr(resourceName, "tenant_tags.#", "1"),
					testTenantTagsOnServer(resourceName, firstTag),
				),
			},
			{
				Config: testAccTenantWithTags(localName, "[octopusdeploy_tag.second.canonical_tag_name]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tenant_tags.#", "1"),
					testTenantTagsOnServer(resourceName, secondTag),
				),
			},
			{
				// Change the tags behind Terraform's back. The refresh has to notice
				// and the apply has to put them back.
				PreConfig: func() { setTenantTagsOutOfBand(t, &tenantID, firstTag) },
				Config:    testAccTenantWithTags(localName, "[octopusdeploy_tag.second.canonical_tag_name]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tenant_tags.#", "1"),
					testTenantTagsOnServer(resourceName, secondTag),
				),
			},
			{
				Config: testAccTenantWithTags(localName, "[]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tenant_tags.#", "0"),
					testTenantTagsOnServer(resourceName),
				),
			},
		},
	})
}

func testAccTenantWithTags(localName string, tenantTags string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_tag_set" "%[1]s" {
		  name = "%[1]s"
		}

		resource "octopusdeploy_tag" "first" {
		  name        = "first"
		  color       = "#6e6e6e"
		  description = "First tenant tag"
		  tag_set_id  = octopusdeploy_tag_set.%[1]s.id
		}

		resource "octopusdeploy_tag" "second" {
		  name        = "second"
		  color       = "#6e6e6e"
		  description = "Second tenant tag"
		  tag_set_id  = octopusdeploy_tag_set.%[1]s.id
		}

		resource "octopusdeploy_tenant" "%[1]s" {
		  name        = "%[1]s"
		  description = "Tenant tag test"

		  tenant_tags = %[2]s
		}
	`, localName, tenantTags)
}

func captureTenantID(prefix string, tenantID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("Not found: %s", prefix)
		}

		*tenantID = rs.Primary.ID

		return nil
	}
}

func setTenantTagsOutOfBand(t *testing.T, tenantID *string, tags ...string) {
	tenant, err := tenants.GetByID(octoClient, octoClient.GetSpaceID(), *tenantID)
	if err != nil {
		t.Fatalf("unable to load tenant (%s) to change its tags out of band: %s", *tenantID, err)
	}

	tenant.TenantTags = tags
	if _, err := tenants.Update(octoClient, tenant); err != nil {
		t.Fatalf("unable to change the tags on tenant (%s) out of band: %s", *tenantID, err)
	}
}

func testTenantTagsOnServer(prefix string, expected ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("Not found: %s", prefix)
		}

		tenant, err := tenants.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID)
		if err != nil {
			return err
		}

		actual := append([]string{}, tenant.TenantTags...)
		want := append([]string{}, expected...)
		sort.Strings(actual)
		sort.Strings(want)

		if strings.Join(actual, ",") != strings.Join(want, ",") {
			return fmt.Errorf("expected tenant tags %v on the server, got %v", want, actual)
		}

		return nil
	}
}

func testTenantExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("Not found: %s", prefix)
		}

		if _, err := tenants.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccTenantCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_tenant" {
			continue
		}

		if tenant, err := octoClient.Tenants.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("tenant (%s) still exists", tenant.GetID())
		}
	}

	return nil
}
