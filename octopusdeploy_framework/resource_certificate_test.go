package octopusdeploy_framework

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"math/big"
	"testing"
	"time"

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

func TestAccOctopusDeployCertificateReplace(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_certificate." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	originalCommonName := "original-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) + ".example.com"
	replacementCommonName := "replacement-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) + ".example.com"
	originalCertificateData := generateTestCertificateData(t, originalCommonName)
	replacementCertificateData := generateTestCertificateData(t, replacementCommonName)

	var originalID string

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccCertificateCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testCertificateExists(prefix),
					resource.TestCheckResourceAttr(prefix, "certificate_data", originalCertificateData),
					resource.TestCheckResourceAttr(prefix, "subject_common_name", originalCommonName),
					func(s *terraform.State) error {
						originalID = s.RootModule().Resources[prefix].Primary.ID
						return nil
					},
				),
				Config: testCertificateBasic(localName, name, originalCertificateData, password),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testCertificateExists(prefix),
					resource.TestCheckResourceAttr(prefix, "certificate_data", replacementCertificateData),
					resource.TestCheckResourceAttr(prefix, "subject_common_name", replacementCommonName),
					func(s *terraform.State) error {
						if id := s.RootModule().Resources[prefix].Primary.ID; id != originalID {
							return fmt.Errorf("expected certificate to be replaced in place (%s), but got %s", originalID, id)
						}

						certificate, err := octoClient.Certificates.GetByID(originalID)
						if err != nil {
							return err
						}
						if certificate.SubjectCommonName != replacementCommonName {
							return fmt.Errorf("expected server certificate common name %s, but got %s", replacementCommonName, certificate.SubjectCommonName)
						}

						return nil
					},
				),
				Config: testCertificateBasic(localName, name, replacementCertificateData, password),
			},
		},
	})
}

// generateTestCertificateData returns a base64 encoded, DER formatted, self-signed certificate
func generateTestCertificateData(t *testing.T, commonName string) string {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject:      pkix.Name{CommonName: commonName},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(der)
}

func testCertificateBasic(localName string, name string, certificateData string, password string) string {
	return fmt.Sprintf(`
locals {
  environments = [
    {
      key        = "alpha"
      name       = "TFA1"
      sort_order = 1
    },
    {
      key        = "development"
      name       = "TFB1"
      sort_order = 2
    },
    {
      key        = "test"
      name       = "TFC1"
      sort_order = 3
    },
    {
      key        = "prod"
      name       = "TFD1"
      sort_order = 4
    }
  ]
  certificate_scopes = ["alpha", "development"]
}

resource "octopusdeploy_environment" "environment" {
  count      = length(local.environments)
  name       = local.environments[count.index].name
  sort_order = local.environments[count.index].sort_order

  lifecycle {
    prevent_destroy       = false
    create_before_destroy = true
  }
}

resource "octopusdeploy_certificate" "%s" {
  certificate_data = "%s"
  name             = "%s"
  password         = "%s"
  environments     = [
    for scope in local.certificate_scopes : 
    octopusdeploy_environment.environment[
      index(local.environments[*].key, scope)
    ].id
  ]
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
