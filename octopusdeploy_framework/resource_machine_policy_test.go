package octopusdeploy_framework

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployMachinePolicyBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccMachinePolicyCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "space_id", space.ID),
					resource.TestCheckResourceAttr(prefix, "is_default", "false"),
				),
				Config: testAccMachinePolicyBasic(localName, name, space.ID),
			},
		},
	})
}

func TestAccOctopusDeployMachinePolicyUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccMachinePolicyCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
				),
				Config: testAccMachinePolicyWithDescription(localName, name, description, space.ID),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
				),
				Config: testAccMachinePolicyWithDescription(localName, newName, newDescription, space.ID),
			},
		},
	})
}

func TestAccOctopusDeployMachinePolicyWithConnectivitySettings(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccMachinePolicyCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "connection_connect_timeout", "60000000000"),
					resource.TestCheckResourceAttr(prefix, "connection_retry_count_limit", "5"),
					resource.TestCheckResourceAttr(prefix, "connection_retry_sleep_interval", "1000000000"),
					resource.TestCheckResourceAttr(prefix, "connection_retry_time_limit", "300000000000"),
				),
				Config: testAccMachinePolicyWithConnectivitySettings(localName, name, space.ID),
			},
		},
	})
}

func TestAccOctopusDeployMachinePolicyImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccMachinePolicyCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccMachinePolicyBasic(localName, name, space.ID),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMachinePolicyImportStateIdFunc(resourceName),
			},
		},
	})
}

// A health check interval of 0 means the schedule is "Never". The server represents this as the
// absence of an interval, so the assertions here check the value that came back from the server as
// well as the value in state.
func TestAccOctopusDeployMachinePolicyHealthCheckNever(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccMachinePolicyCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					resource.TestCheckTypeSetElemNestedAttrs(prefix, "machine_health_check_policy.*", map[string]string{
						"health_check_interval": "0",
					}),
					testAccMachinePolicyHealthCheckIntervalOnServer(prefix, 0),
					testAccMachinePolicyHealthCheckScheduleIsNever(prefix),
				),
				Config: testAccMachinePolicyWithHealthCheckInterval(localName, name, space.ID, 0),
			},
			{
				// A second plan against the same configuration must be empty. The interval carries a
				// schema default of 24 hours, and this guards against that default reappearing over
				// the explicit 0.
				Config:   testAccMachinePolicyWithHealthCheckInterval(localName, name, space.ID, 0),
				PlanOnly: true,
			},
		},
	})
}

func TestAccOctopusDeployMachinePolicyHealthCheckIntervalToNever(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccMachinePolicyCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					testAccMachinePolicyHealthCheckIntervalOnServer(prefix, 24*time.Hour),
				),
				Config: testAccMachinePolicyWithHealthCheckInterval(localName, name, space.ID, 24*time.Hour),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					resource.TestCheckTypeSetElemNestedAttrs(prefix, "machine_health_check_policy.*", map[string]string{
						"health_check_interval": "0",
					}),
					testAccMachinePolicyHealthCheckIntervalOnServer(prefix, 0),
				),
				Config: testAccMachinePolicyWithHealthCheckInterval(localName, name, space.ID, 0),
			},
		},
	})
}

// The server writes intervals of a day or longer as "d.hh:mm:ss". This covers that round trip, which
// is the shape that has to keep working for 0 to be usable as the "Never" sentinel.
func TestAccOctopusDeployMachinePolicyHealthCheckMultiDayInterval(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccMachinePolicyCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					resource.TestCheckTypeSetElemNestedAttrs(prefix, "machine_health_check_policy.*", map[string]string{
						"health_check_interval": "604800000000000",
					}),
					testAccMachinePolicyHealthCheckIntervalOnServer(prefix, 7*24*time.Hour),
				),
				Config: testAccMachinePolicyWithHealthCheckInterval(localName, name, space.ID, 7*24*time.Hour),
			},
			{
				Config:   testAccMachinePolicyWithHealthCheckInterval(localName, name, space.ID, 7*24*time.Hour),
				PlanOnly: true,
			},
		},
	})
}

// A policy that never configured a health check block must keep its server-side interval when some
// unrelated attribute changes. The block is Optional and Computed, so the interval read back from
// the server is what gets sent on update, and it must not be mistaken for "Never".
func TestAccOctopusDeployMachinePolicyHealthCheckDefaultSurvivesUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccMachinePolicyCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					testAccMachinePolicyHealthCheckIntervalOnServer(prefix, 24*time.Hour),
				),
				Config: testAccMachinePolicyBasic(localName, name, space.ID),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					resource.TestCheckResourceAttr(prefix, "description", description),
					testAccMachinePolicyHealthCheckIntervalOnServer(prefix, 24*time.Hour),
				),
				Config: testAccMachinePolicyWithDescription(localName, name, description, space.ID),
			},
		},
	})
}

func testAccMachinePolicyWithHealthCheckInterval(localName, name, spaceID string, interval time.Duration) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`resource "octopusdeploy_machine_policy" "%s" {
		name     = "%s"
		space_id = "%s"

		machine_health_check_policy {
			health_check_interval = %d

			bash_health_check_policy {}
			powershell_health_check_policy {}
		}
	}`, localName, name, spaceID, interval)
}

func testAccMachinePolicyHealthCheckIntervalOnServer(prefix string, expected time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("Not found: %s", prefix)
		}

		machinePolicy, err := machinepolicies.GetByID(octoClient, rs.Primary.Attributes["space_id"], rs.Primary.ID)
		if err != nil {
			return err
		}

		if actual := machinePolicy.MachineHealthCheckPolicy.HealthCheckInterval; actual != expected {
			return fmt.Errorf("health check interval on the server is %s, expected %s", actual, expected)
		}

		return nil
	}
}

// "Never" is the absence of a health check interval, which is not something the parsed duration can
// show: an interval of "00:00:00" also reads back as zero, and the server treats that as an interval
// of zero rather than as "Never". This reads the raw payload so the distinction is asserted.
func testAccMachinePolicyHealthCheckScheduleIsNever(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("Not found: %s", prefix)
		}

		url := fmt.Sprintf("%s/api/%s/machinepolicies/%s", strings.TrimSuffix(os.Getenv("OCTOPUS_URL"), "/"), rs.Primary.Attributes["space_id"], rs.Primary.ID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		request.Header.Set("X-Octopus-ApiKey", os.Getenv("OCTOPUS_APIKEY"))

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("GET %s returned %s", url, response.Status)
		}

		var payload struct {
			MachineHealthCheckPolicy struct {
				HealthCheckInterval *string `json:"HealthCheckInterval"`
				HealthCheckCron     *string `json:"HealthCheckCron"`
			} `json:"MachineHealthCheckPolicy"`
		}
		if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
			return err
		}

		if interval := payload.MachineHealthCheckPolicy.HealthCheckInterval; interval != nil {
			return fmt.Errorf("health check interval on the server is %q, expected it to be absent", *interval)
		}
		if cron := payload.MachineHealthCheckPolicy.HealthCheckCron; cron != nil && *cron != "" {
			return fmt.Errorf("health check cron on the server is %q, expected it to be absent", *cron)
		}

		return nil
	}
}

func testAccMachinePolicyBasic(localName, name, spaceID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`resource "octopusdeploy_machine_policy" "%s" {
		name     = "%s"
		space_id = "%s"
	}`, localName, name, spaceID)
}

func testAccMachinePolicyWithDescription(localName, name, description, spaceID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`resource "octopusdeploy_machine_policy" "%s" {
		name        = "%s"
		description = "%s"
		space_id    = "%s"
	}`, localName, name, description, spaceID)
}

func testAccMachinePolicyWithConnectivitySettings(localName, name, spaceID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`resource "octopusdeploy_machine_policy" "%s" {
		name                           = "%s"
		space_id                       = "%s"
		connection_connect_timeout     = 60000000000
		connection_retry_count_limit   = 5
		connection_retry_sleep_interval = 1000000000
		connection_retry_time_limit     = 300000000000
	}`, localName, name, spaceID)
}

func testAccMachinePolicyExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[prefix]
		machinePolicyID := rs.Primary.ID
		if _, err := machinepolicies.GetByID(octoClient, rs.Primary.Attributes["space_id"], machinePolicyID); err != nil {
			return err
		}

		return nil
	}
}

func testAccMachinePolicyCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_machine_policy" {
			continue
		}

		if machinePolicy, err := machinepolicies.GetByID(octoClient, rs.Primary.Attributes["space_id"], rs.Primary.ID); err == nil {
			return fmt.Errorf("machine policy (%s) still exists", machinePolicy.GetID())
		}
	}

	return nil
}

func testAccMachinePolicyImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
