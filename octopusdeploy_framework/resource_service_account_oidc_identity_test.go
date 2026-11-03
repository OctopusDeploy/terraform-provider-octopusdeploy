package octopusdeploy_framework

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/serviceaccounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOctopusDeployServiceAccountOIDCIdentity(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_service_account_oidc_identity." + localName

	localUserName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userPrefix := " octopusdeploy_user." + localUserName

	userData := users.User{
		DisplayName:  acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		EmailAddress: acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) + "@test.com",
		Username:     acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
	}

	data := serviceaccounts.OIDCIdentity{
		Name:             acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		ServiceAccountID: userPrefix + ".id",
		Issuer:           "https://token.actions.githubusercontent.com",
		Subject:          "repo:test/test:environment:test",
		Audience:         "oidc-machine-user-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testServiceAccountIdentityConfig(localName, localUserName, data, userData),
				Check: resource.ComposeTestCheckFunc(
					testScriptModuleExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", data.Name),
					resource.TestCheckResourceAttr(prefix, "issuer", data.Issuer),
					resource.TestCheckResourceAttr(prefix, "subject", data.Subject),
					resource.TestCheckResourceAttr(prefix, "audience", data.Audience),
					resource.TestCheckResourceAttrSet("octopusdeploy_user."+localUserName, "external_id"),
				),
			},
			{
				Config: testServiceAccountIdentityUpdate(localName, localUserName, data, userData),
				Check: resource.ComposeTestCheckFunc(
					testScriptModuleExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", data.Name+"-updated"),
					resource.TestCheckResourceAttr(prefix, "issuer", data.Issuer),
					resource.TestCheckResourceAttr(prefix, "subject", data.Subject),
					resource.TestCheckResourceAttr(prefix, "audience", data.Audience),
					resource.TestCheckResourceAttrSet("octopusdeploy_user."+localUserName, "external_id"),
				),
			},
		},
	})
}

func testServiceAccountIdentityConfig(localName string, localUserName string, data serviceaccounts.OIDCIdentity, userData users.User) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_user" "%s" {
		display_name = "%s"
		email_address = "%s"
		is_active = true
		is_service = true
		username = "%s"
	}
	resource "octopusdeploy_service_account_oidc_identity" "%s" {
		name        = "%s"
		service_account_id = %s
		issuer = "%s"
		subject = "%s"
		audience = "%s"
	}`,
		localUserName,
		userData.DisplayName,
		userData.EmailAddress,
		userData.Username,
		localName,
		data.Name,
		data.ServiceAccountID,
		data.Issuer,
		data.Subject,
		data.Audience)
}

func testServiceAccountIdentityUpdate(localName string, localUserName string, data serviceaccounts.OIDCIdentity, userData users.User) string {
	data.Name = data.Name + "-updated"
	return testServiceAccountIdentityConfig(localName, localUserName, data, userData)
}

// The oidc-claim-conditions feature toggle is off by default on Octopus Server.
// Set OCTOPUS_OIDC_CLAIM_CONDITIONS_ENABLED=true when the target server has the toggle on.
func TestAccOctopusDeployServiceAccountOIDCIdentityClaimConditions(t *testing.T) {
	if os.Getenv("OCTOPUS_OIDC_CLAIM_CONDITIONS_ENABLED") != "true" {
		t.Skip("Skipping claim conditions test: set OCTOPUS_OIDC_CLAIM_CONDITIONS_ENABLED=true on a server with the oidc-claim-conditions toggle on")
	}

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_service_account_oidc_identity." + localName
	localUserName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	identityName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	config := func(claimConditions string) string {
		return fmt.Sprintf(`
	resource "octopusdeploy_user" "%[1]s" {
		display_name  = "%[2]s"
		email_address = "%[2]s@test.com"
		is_active     = true
		is_service    = true
		username      = "%[2]s"
	}
	resource "octopusdeploy_service_account_oidc_identity" "%[3]s" {
		name               = "%[4]s"
		service_account_id = octopusdeploy_user.%[1]s.id
		issuer             = "https://token.actions.githubusercontent.com"
		subject            = "repo:test/test:*"
		%[5]s
	}`, localUserName, userName, localName, identityName, claimConditions)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config(`claim_conditions = [
			{ claim = "repository_owner", values = ["test"] },
			{ claim = "ref", values = ["refs/heads/main", "refs/tags/*"] },
		]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prefix, "claim_conditions.#", "2"),
					resource.TestCheckResourceAttr(prefix, "claim_conditions.0.claim", "repository_owner"),
					resource.TestCheckResourceAttr(prefix, "claim_conditions.0.values.0", "test"),
					resource.TestCheckResourceAttr(prefix, "claim_conditions.1.claim", "ref"),
					resource.TestCheckResourceAttr(prefix, "claim_conditions.1.values.#", "2"),
					resource.TestCheckResourceAttr(prefix, "claim_conditions.1.values.1", "refs/tags/*"),
				),
			},
			{
				Config: config(`claim_conditions = [
			{ claim = "environment", values = ["production"] },
		]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prefix, "claim_conditions.#", "1"),
					resource.TestCheckResourceAttr(prefix, "claim_conditions.0.claim", "environment"),
					resource.TestCheckResourceAttr(prefix, "claim_conditions.0.values.0", "production"),
				),
			},
			{
				// Removing the attribute must clear the conditions on the server.
				Config: config(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(prefix, "claim_conditions.#"),
				),
			},
		},
	})
}

func TestAccOctopusDeployServiceAccountOIDCIdentityClaimConditionsValidation(t *testing.T) {
	config := func(claimConditions string) string {
		return fmt.Sprintf(`
	resource "octopusdeploy_service_account_oidc_identity" "test" {
		name               = "test"
		service_account_id = "Users-1"
		issuer             = "https://token.actions.githubusercontent.com"
		subject            = "repo:test/test:*"
		claim_conditions   = %s
	}`, claimConditions)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      config(`[{ claim = "sub", values = ["x"] }]`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`value must be none of`),
			},
			{
				Config:      config(`[{ claim = "ref", values = [] }]`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`list must contain at least 1 elements`),
			},
			{
				Config:      config(`[{ claim = "ref", values = ["**"] }]`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`can't be made only of \* characters`),
			},
			{
				Config:      config(`[]`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`list must contain at least 1 elements`),
			},
		},
	})
}
