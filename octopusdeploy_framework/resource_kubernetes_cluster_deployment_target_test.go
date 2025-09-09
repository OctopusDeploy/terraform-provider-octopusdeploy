package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployKubernetesClusterDeploymentTargetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_kubernetes_cluster_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	clusterUrl := "https://k8s-cluster.example.com"
	
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccKubernetesClusterDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesClusterDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "cluster_url", clusterUrl),
					resource.TestCheckResourceAttr(prefix, "environments.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.0", "k8s-cluster"),
					resource.TestCheckResourceAttr(prefix, "skip_tls_verification", "true"),
				),
				Config: testAccKubernetesClusterDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, clusterUrl),
			},
		},
	})
}

func TestAccOctopusDeployKubernetesClusterDeploymentTargetUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_kubernetes_cluster_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	clusterUrl := "https://k8s-cluster.example.com"
	newClusterUrl := "https://new-k8s-cluster.example.com"

	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccKubernetesClusterDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesClusterDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "cluster_url", clusterUrl),
				),
				Config: testAccKubernetesClusterDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, clusterUrl),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesClusterDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "cluster_url", newClusterUrl),
					resource.TestCheckResourceAttr(prefix, "roles.#", "2"),
					resource.TestCheckResourceAttr(prefix, "namespace", "production"),
					resource.TestCheckResourceAttr(prefix, "is_disabled", "true"),
				),
				Config: testAccKubernetesClusterDeploymentTargetUpdate(localName, environmentLocalName, environmentName, newName, newClusterUrl),
			},
		},
	})
}

func TestAccOctopusDeployKubernetesClusterDeploymentTargetWithCertAuth(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_kubernetes_cluster_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	clusterUrl := "https://k8s-cluster.example.com"
	
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	certificateLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	certificateName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccKubernetesClusterDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesClusterDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "certificate_authentication.#", "1"),
				),
				Config: testAccKubernetesClusterDeploymentTargetWithCertAuth(localName, environmentLocalName, environmentName, certificateLocalName, certificateName, name, clusterUrl),
			},
		},
	})
}

func TestAccOctopusDeployKubernetesClusterDeploymentTargetImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_kubernetes_cluster_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	clusterUrl := "https://k8s-cluster.example.com"
	
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccKubernetesClusterDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, clusterUrl),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"has_latest_calamari",
					"health_status",
					"is_in_process",
					"operating_system",
					"shell_name",
					"shell_version",
					"status",
					"status_summary",
					"uri",
				},
				ImportStateIdFunc: testAccKubernetesClusterDeploymentTargetImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccKubernetesClusterDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, clusterUrl string) string {
	accountLocalName := "acc" + localName
	accountName := "TestAccount-" + localName
	return fmt.Sprintf(`
	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_username_password_account" "%s" {
		name     = "%s"
		username = "test-user"
	}

	resource "octopusdeploy_kubernetes_cluster_deployment_target" "%s" {
		name         = "%s"
		cluster_url  = "%s"
		environments = [octopusdeploy_environment.%s.id]
		roles        = ["k8s-cluster"]
		skip_tls_verification = true
		
		authentication {
			account_id = octopusdeploy_username_password_account.%s.id
		}
	}`, environmentLocalName, environmentName, accountLocalName, accountName, localName, name, clusterUrl, environmentLocalName, accountLocalName)
}

func testAccKubernetesClusterDeploymentTargetUpdate(localName, environmentLocalName, environmentName, name, clusterUrl string) string {
	accountLocalName := "acc" + localName
	accountName := "TestAccount-" + localName
	return fmt.Sprintf(`
	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_username_password_account" "%s" {
		name     = "%s"
		username = "test-user"
	}

	resource "octopusdeploy_kubernetes_cluster_deployment_target" "%s" {
		name         = "%s"
		cluster_url  = "%s"
		environments = [octopusdeploy_environment.%s.id]
		roles        = ["k8s-cluster", "production"]
		namespace    = "production"
		is_disabled  = true
		skip_tls_verification = true
		machine_policy_id = "MachinePolicies-1"
		
		authentication {
			account_id = octopusdeploy_username_password_account.%s.id
		}
	}`, environmentLocalName, environmentName, accountLocalName, accountName, localName, name, clusterUrl, environmentLocalName, accountLocalName)
}

func testAccKubernetesClusterDeploymentTargetWithCertAuth(localName, environmentLocalName, environmentName, certificateLocalName, certificateName, name, clusterUrl string) string {
	// Note: This test creates a dummy certificate. In production, a real certificate would be needed.
	return fmt.Sprintf(`
	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_certificate" "%s" {
		name = "%s"
		certificate_data = "MIIDiDCCAnACCQDXHofnqz05ITANBgkqhkiG9w0BAQsFADCBhTELMAkGA1UEBhMCVVMxETAPBgNVBAgMCE9rbGFob21hMQ8wDQYDVQQHDAZOb3JtYW4xEzARBgNVBAoMCk1vb25zd2l0Y2gxGTAXBgNVBAMMEGRlbW8ub2N0b3B1cy5jb20xIjAgBgkqhkiG9w0BCQEWE2plZmZAbW9vbnN3aXRjaC5jb20wHhcNMTkwNjE0MjExMzI1WhcNMjAwNjEzMjExMzI1WjCBhTELMAkGA1UEBhMCVVMxETAPBgNVBAgMCE9rbGFob21hMQ8wDQYDVQQHDAZOb3JtYW4xEzARBgNVBAoMCk1vb25zd2l0Y2gxGTAXBgNVBAMMEGRlbW8ub2N0b3B1cy5jb20xIjAgBgkqhkiG9w0BCQEWE2plZmZAbW9vbnN3aXRjaC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDSTiD0OHyFDMH9O+d/h3AiqcuvpvUgRkKjf+whZ6mVlQnGkvPddRTUY48xCEaQ4QD1MAVJcGaJ2PU4NxwhrQgHqWW8TQkAZESL4wfzSwIKO2NX/I2tWqyv7a0uA/WdtlWQye+2oPV5rCnS0kM75X+gjEwOTpFh/ryS6KhMPFDb0zeNGREdg6564FdxWSvN4ppUZMqhvMpfzM7rsDWqEzYsMaQ4CNJDFdWkG89D4j5qk4b4Qb4m+l7QINdmYIXf4qO/0LE1WcfIkCpAS65tjc/hefIHmYtj/E/ijoNJbWKZDK3WLZg3zq99Ipqv/9DFvSiMQFBhZT0jO2B5d5zBUuIHAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAKsa4gaW7GhByu8aq56h99DaIl1LauI5WMVH8Q9Qpapho2VLRIpfwGeI5eENFoXwuKrnJp1ncsCqGnMQnugQHS+SrruS3Yyl0Uog4Zak9GbbK6qn+olx7GNJbsckmD371lqQOaKITLqYzK6kTc7/v8Cv0BwHFCBda1OCrmeVBSaarucPxZhGxzLAielzHHdlkZFQT/oO2VR3thhURIqtni7jVQ2MoeZF1ccvmAfVbzr/QnlNe/jrcmyPYymuF2JyrezzIjlKuiDhalKqwqkCHpOOgzV4y6BFuS+0w3DS8pa07nUudZ6E0kZzvhjjiyAx/sBdX6ZDdUjP9TDJMM4f5YA="
		password = "test-password"
	}

	resource "octopusdeploy_kubernetes_cluster_deployment_target" "%s" {
		name         = "%s"
		cluster_url  = "%s"
		environments = [octopusdeploy_environment.%s.id]
		roles        = ["k8s-cluster"]
		
		certificate_authentication {
			client_certificate = octopusdeploy_certificate.%s.id
		}
		
		skip_tls_verification = true
	}`, environmentLocalName, environmentName, certificateLocalName, certificateName, localName, name, clusterUrl, environmentLocalName, certificateLocalName)
}

func testAccKubernetesClusterDeploymentTargetExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		targetID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), targetID); err != nil {
			return err
		}

		return nil
	}
}

func testAccKubernetesClusterDeploymentTargetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_kubernetes_cluster_deployment_target" {
			continue
		}

		if target, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("kubernetes cluster deployment target (%s) still exists", target.GetID())
		}
	}

	return nil
}

func testAccKubernetesClusterDeploymentTargetImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}