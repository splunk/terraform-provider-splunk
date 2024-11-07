package splunk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const splunkLookupDefinitionObjectWithAcl = `
resource "splunk_lookup_definition" "example" {
  name     = "example_lookup_definition"
  filename = "example_lookup_file.csv"
  acl {
    owner = "admin"
    app = "search"
    sharing = "global"
    read = ["*"]
    write = ["admin"]
  }
}
`

const updateSplunkLookupDefinitionObjectWithAcl = `
resource "splunk_lookup_definition" "example" {
  name     = "example_lookup_definition"
  filename = "example_lookup_file2.csv"
  acl {
    owner = "admin"
    app = "search"
    read = ["*"]
    write = ["admin"]
  }
}
  `

const splunkLookupDefinitionObject = `
resource "splunk_lookup_definition" "example" {
  name     = "example_lookup_definition"
  filename = "example_lookup_file.csv"
  }
`

const updateSplunkLookupDefinitionObject = `
resource "splunk_lookup_definition" "example" {
  name     = "example_lookup_definition"
  filename = "example_lookup_file2.csv"
  }
`

func TestAccSplunkLookupDefinition(t *testing.T) {
	resourceName := "splunk_lookup_definition.example"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkLookupDefinitionDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: splunkLookupDefinitionObject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "example_lookup_definition"),
					resource.TestCheckResourceAttr(resourceName, "filename", "example_lookup_file.csv"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "user"),
				),
			},
			{
				Config: updateSplunkLookupDefinitionObject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "example_lookup_definition"),
					resource.TestCheckResourceAttr(resourceName, "filename", "example_lookup_file2.csv"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "user"),
				),
			},
			{
				ResourceName:      "splunk_lookup_definition.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSplunkLookupDefinitionWithAcl(t *testing.T) {
	resourceName := "splunk_lookup_definition.example"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkLookupDefinitionDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: splunkLookupDefinitionObjectWithAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "example_lookup_definition"),
					resource.TestCheckResourceAttr(resourceName, "filename", "example_lookup_file.csv"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "user"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.0", "admin"),
				),
			},
			{
				Config: updateSplunkLookupDefinitionObjectWithAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "example_lookup_definition"),
					resource.TestCheckResourceAttr(resourceName, "filename", "example_lookup_file2.csv"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "user"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.0", "admin"),
				),
			},
			{
				ResourceName:      "splunk_lookup_definition.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkLookupDefinitionDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_lookup_definition":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "admin", "search", "data", "transforms", "lookups", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
