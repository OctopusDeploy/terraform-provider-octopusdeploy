package octopusdeploy_framework

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccOctopusDeployBuiltInRateLimitingPolicyBasic(t *testing.T) {
	resourceName := "octopusdeploy_built_in_rate_limiting_policy.anon"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: builtInRateLimitingPolicyConfig("anon", "anon", true, 5000, 100, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "slug", "anon"),
					resource.TestCheckResourceAttr(resourceName, "scope_type", "Unauthenticated"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "requests_per_hour", "5000"),
					resource.TestCheckResourceAttr(resourceName, "burst_limit", "100"),
					resource.TestCheckResourceAttr(resourceName, "audit_mode", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "anon",
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOctopusDeployBuiltInRateLimitingPolicyUpdate(t *testing.T) {
	resourceName := "octopusdeploy_built_in_rate_limiting_policy.anon"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: builtInRateLimitingPolicyConfig("anon", "anon", true, 5000, 100, false),
				Check:  resource.TestCheckResourceAttr(resourceName, "requests_per_hour", "5000"),
			},
			{
				Config: builtInRateLimitingPolicyConfig("anon", "anon", true, 8000, 200, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "requests_per_hour", "8000"),
					resource.TestCheckResourceAttr(resourceName, "burst_limit", "200"),
					resource.TestCheckResourceAttr(resourceName, "audit_mode", "true"),
				),
			},
		},
	})
}

func TestAccOctopusDeployBuiltInRateLimitingPolicyInvalidSlug(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      builtInRateLimitingPolicyConfig("nonexistent", "nonexistent", true, 5000, 100, false),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestAccOctopusDeployBuiltInRateLimitingPolicyReplaceOnSlugChange(t *testing.T) {
	resourceName := "octopusdeploy_built_in_rate_limiting_policy.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: builtInRateLimitingPolicyConfig("test", "anon", false, 5000, 100, false),
				Check:  resource.TestCheckResourceAttr(resourceName, "slug", "anon"),
			},
			{
				Config: builtInRateLimitingPolicyConfig("test", "user", false, 5000, 100, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// slug is the resource's identity, so changing it triggers a replace rather than an update;
						// because built-in policies can't be created or destroyed, destroy is a no-op and create is
						// really a modify, so the replace just changes the target to the "user" policy
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "slug", "user"),
					resource.TestCheckResourceAttr(resourceName, "scope_type", "AuthenticatedHuman"),
				),
			},
		},
	})
}

func builtInRateLimitingPolicyConfig(label, slug string, isEnabled bool, requestsPerHour, burstLimit int, auditMode bool) string {
	return fmt.Sprintf(
		`resource "octopusdeploy_built_in_rate_limiting_policy" "%[1]s" {
			slug              = "%[2]s"
			is_enabled        = %[3]t
			requests_per_hour = %[4]d
			burst_limit       = %[5]d
			audit_mode        = %[6]t
		}`,
		label,
		slug,
		isEnabled,
		requestsPerHour,
		burstLimit,
		auditMode)
}
