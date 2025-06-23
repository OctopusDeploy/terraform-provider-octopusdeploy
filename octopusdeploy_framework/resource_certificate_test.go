package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployCertificateBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_certificate." + localName

	certificateData := "MIIDiDCCAnACCQDXHofnqz05ITANBgkqhkiG9w0BAQsFADCBhTELMAkGA1UEBhMCVVMxETAPBgNVBAgMCE9rbGFob21hMQ8wDQYDVQQHDAZOb3JtYW4xEzARBgNVBAoMCk1vb25zd2l0Y2gxGTAXBgNVBAMMEGRlbW8ub2N0b3B1cy5jb20xIjAgBgkqhkiG9w0BCQEWE2plZmZAbW9vbnN3aXRjaC5jb20wHhcNMTkwNjE0MjExMzI1WhcNMjAwNjEzMjExMzI1WjCBhTELMAkGA1UEBhMCVVMxETAPBgNVBAgMCE9rbGFob21hMQ8wDQYDVQQHDAZOb3JtYW4xEzARBgNVBAoMCk1vb25zd2l0Y2gxGTAXBgNVBAMMEGRlbW8ub2N0b3B1cy5jb20xIjAgBgkqhkiG9w0BCQEWE2plZmZAbW9vbnN3aXRjaC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDSTiD0OHyFDMH9O+d/h3AiqcuvpvUgRkKjf+whZ6mVlQnGkvPddRTUY48xCEaQ4QD1MAVJcGaJ2PU4NxwhrQgHqWW8TQkAZESL4wfzSwIKO2NX/I2tWqyv7a0uA/WdtlWQye+2oPV5rCnS0kM75X+gjEwOTpFh/ryS6KhMPFDb0zeNGREdg6564FdxWSvN4ppUZMqhvMpfzM7rsDWqEzYsMaQ4CNJDFdWkG89D4j5qk4b4Qb4m+l7QINdmYIXf4qO/0LE1WcfIkCpAS65tjc/hefIHmYtj/E/ijoNJbWKZDK3WLZg3zq99Ipqv/9DFvSiMQFBhZT0jO2B5d5zBUuIHAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAKsa4gaW7GhByu8aq56h99DaIl1LauI5WMVH8Q9Qpapho2VLRIpfwGeI5eENFoXwuKrnJp1ncsCqGnMQnugQHS+SrruS3Yyl0Uog4Zak9GbbK6qn+olx7GNJbsckmD371lqQOaKITLqYzK6kTc7/v8Cv0BwHFCBda1OCrmeVBSaarucPxZhGxzLAielzHHdlkZFQT/oO2VR3thhURIqtni7jVQ2MoeZF1ccvmAfVbzr/QnlNe/jrcmyPYymuF2JyrezzIjlKuiDhalKqwqkCHpOOgzV4y6BFuS+0w3DS8pa07nUudZ6E0kZzvhjjiyAx/sBdX6ZDdUjP9TDJMM4f5YA="
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccCertificateCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testCertificateExists(prefix),
					resource.TestCheckResourceAttr(prefix, "certificate_data", certificateData),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "password", password),
				),
				Config: testCertificateBasic(localName, name, certificateData, password),
			},
		},
	})
}

func testCertificateBasic(localName string, name string, certificateData string, password string) string {
	return fmt.Sprintf(`
locals {
  environments = {
    alpha = {
      name       = "TFA1"
      sort_order = 1
    }
    development = {
      name       = "TFB1"
      sort_order = 2
    }
    test = {
      name       = "TFC1"
      sort_order = 3
    }
    prod = {
      name       = "TFD1"
      sort_order = 4
    }
  }
  certificate_scopes = ["alpha", "development"]
}

resource "octopusdeploy_environment" "environment" {
  for_each   = local.environments
  name       = each.value.name
  sort_order = each.value.sort_order

  lifecycle {
    prevent_destroy       = false
    create_before_destroy = true
  }
}

resource "octopusdeploy_certificate" "%s" {
  certificate_data = "%s"
  name             = "%s"
  password         = "%s"
  environments     = [for scope in local.certificate_scopes : octopusdeploy_environment.environment[scope].id]
}
`, localName, certificateData, name, password)
}

func testCertificateExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		certificateID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.Certificates.GetByID(certificateID); err != nil {
			return err
		}

		return nil
	}
}

func testAccCertificateCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_certficate" {
			continue
		}

		certificate, err := octoClient.Certificates.GetByID(rs.Primary.ID)
		if err == nil && certificate != nil {
			return fmt.Errorf("certificate (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
