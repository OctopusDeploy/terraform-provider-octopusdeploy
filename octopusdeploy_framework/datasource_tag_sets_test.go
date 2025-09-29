package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceTagSets(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagSetName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localTagName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagSetResourceName := fmt.Sprintf("octopusdeploy_tag_set.%s", localName)
	tagResourceName := fmt.Sprintf("octopusdeploy_tag.%s", localTagName)
	dataSourceName := fmt.Sprintf("data.octopusdeploy_tag_sets.%s", localName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Create a tag set
			{
				Config: testAccTagSetConfig(localName, tagSetName, localTagName, tagName),
				Check: resource.ComposeTestCheckFunc(
					testTagSetExists(tagSetResourceName),
					resource.TestCheckResourceAttr(tagSetResourceName, "name", tagSetName),
					resource.TestCheckResourceAttr(tagSetResourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(tagSetResourceName, "scopes.0", "Tenant"),
					resource.TestCheckResourceAttr(tagSetResourceName, "type", "MultiSelect"),
				),
			},
			// Query the created tag set using the data source with scope filter
			{
				Config: testAccTagSetConfig(localName, tagSetName, localTagName, tagName) + testAccDataSourceTagSetsConfig(localName, tagSetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagSetsDataSourceID(dataSourceName),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag_sets.0.id", tagSetResourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag_sets.0.name", tagSetResourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag_sets.0.scopes.#", tagSetResourceName, "scopes.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag_sets.0.scopes.0", tagSetResourceName, "scopes.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag_sets.0.type", tagSetResourceName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag_sets.0.tags.0.name", tagResourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag_sets.0.tags.0.color", tagResourceName, "color"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag_sets.0.tags.0.description", tagResourceName, "description"),
				),
			},
		},
	})
}

func testAccCheckTagSetsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cannot find TagSets data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("TagSets data source ID not set")
		}
		return nil
	}
}

func testAccTagSetConfig(localName, tagSetName, localTagName, tagName string) string {
	return fmt.Sprintf(`
resource "octopusdeploy_tag_set" "%s" {
    name        = "%s"
    description = "Test tag set"
    scopes      = ["Tenant"]
    type        = "MultiSelect"
}

resource "octopusdeploy_tag" "%s" {
  name        = "%s"
  tag_set_id  = octopusdeploy_tag_set.%s.id
  color       = "#333333"
  description = "a description"
}

`, localName, tagSetName, localTagName, tagName, localName)
}

func testAccDataSourceTagSetsConfig(localName, tagSetName string) string {
	return fmt.Sprintf(`
data "octopusdeploy_tag_sets" "%s" {
    partial_name = "%s"
    scopes       = ["Tenant"]
    skip         = 0
    take         = 10
}
`, localName, tagSetName)
}
