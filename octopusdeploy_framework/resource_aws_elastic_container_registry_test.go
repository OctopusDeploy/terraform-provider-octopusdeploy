package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

type awsEcrFeedTestData struct {
	name      string
	region    string
	accessKey string
	secretKey string
}

func TestAccOctopusDeployAwsEcrFeed(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_aws_elastic_container_registry." + localName
	createData := awsEcrFeedTestData{
		name:      acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		region:    "us-east-1",
		accessKey: acctest.RandStringFromCharSet(16, acctest.CharSetAlpha),
		secretKey: acctest.RandStringFromCharSet(32, acctest.CharSetAlpha),
	}
	updateData := awsEcrFeedTestData{
		name:      createData.name + "-updated",
		region:    "us-west-2",
		accessKey: createData.accessKey + "-updated",
		secretKey: createData.secretKey + "-updated",
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testAwsEcrFeedCheckDestroy(s) },
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAwsEcrFeedBasic(createData, localName),
				Check:  testAssertAwsEcrFeedAttributes(createData, prefix),
			},
			{
				Config: testAwsEcrFeedBasic(updateData, localName),
				Check:  testAssertAwsEcrFeedAttributes(updateData, prefix),
			},
		},
	})
}

func testAwsEcrFeedBasic(data awsEcrFeedTestData, localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_aws_elastic_container_registry" "%s" {
			name       = "%s"
			region     = "%s"
			access_key = "%s"
			secret_key = "%s"
		}
	`,
		localName,
		data.name,
		data.region,
		data.accessKey,
		data.secretKey,
	)
}

func testAssertAwsEcrFeedAttributes(expected awsEcrFeedTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttr(prefix, "region", expected.region),
		resource.TestCheckResourceAttr(prefix, "access_key", expected.accessKey),
		resource.TestCheckResourceAttr(prefix, "secret_key", expected.secretKey),
	)
}

func testAwsEcrFeedCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_aws_elastic_container_registry" {
			continue
		}

		feed, err := feeds.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("aws elastic container registry feed (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func TestAccOctopusDeployAwsEcrFeedWithOIDC(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_aws_elastic_container_registry." + localName

	accessKeySecretKeyConfig := `
resource "octopusdeploy_aws_elastic_container_registry" "%s" {
  name       = "AWS ECR With AccessKey"
  region     = "us-east-1"
  access_key = "testaccesskey"
  secret_key = "testsecretkey"
}
`
	oidcConfig := `
resource "octopusdeploy_aws_elastic_container_registry" "%s" {
  name       = "AWS ECR OIDC Registry"
  region     = "us-east-1"
  oidc_authentication = {
    session_duration = "3600"
    role_arn = "role_arn_value"
    audience = "audience"
    subject_keys = ["feed", "space"]
  }
}
`
	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testAwsEcrFeedCheckDestroy(s) },
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(accessKeySecretKeyConfig, localName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prefix, "name", "AWS ECR With AccessKey"),
					resource.TestCheckResourceAttr(prefix, "region", "us-east-1"),
					resource.TestCheckResourceAttr(prefix, "access_key", "testaccesskey"),
					resource.TestCheckResourceAttr(prefix, "secret_key", "testsecretkey"),
					resource.TestCheckNoResourceAttr(prefix, "oidc_authentication"),
				),
			},
			{
				Config: fmt.Sprintf(oidcConfig, localName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prefix, "name", "AWS ECR OIDC Registry"),
					resource.TestCheckResourceAttr(prefix, "region", "us-east-1"),
					resource.TestCheckResourceAttr(prefix, "oidc_authentication.session_duration", "3600"),
					resource.TestCheckResourceAttr(prefix, "oidc_authentication.role_arn", "role_arn_value"),
					resource.TestCheckResourceAttr(prefix, "oidc_authentication.audience", "audience"),
					resource.TestCheckTypeSetElemAttr(prefix, "oidc_authentication.subject_keys.*", "feed"),
					resource.TestCheckTypeSetElemAttr(prefix, "oidc_authentication.subject_keys.*", "space"),
					resource.TestCheckNoResourceAttr(prefix, "access_key"),
					resource.TestCheckNoResourceAttr(prefix, "secret_key"),
				),
			},
		},
	})
}
