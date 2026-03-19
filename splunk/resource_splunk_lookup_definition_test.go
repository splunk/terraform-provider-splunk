package splunk

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

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

// if sharing=user there is no read and write attributes stored in perms
const splunkLookupDefinitionObjectWithAcl = `
resource "splunk_lookup_definition" "example" {
  name     = "example_lookup_definition"
  filename = "example_lookup_file.csv"
  acl {
    owner = "admin"
    app = "search"
	sharing = "app"
    read = ["*"]
    write = ["admin"]
  }
}
`

const splunkLookupDefinitionObjectWithAclFailedValidation = `
resource "splunk_lookup_definition" "example" {
  name     = "example_lookup_definition"
  filename = "example_lookup_file.csv"
  acl {
    owner = "admin"
    app = "search"
	sharing = "user"
    read = ["*"]
    write = ["admin"]
  }
}
`

const updateSplunkLookupDefinitionObjectWithAclSharingUser = `
  resource "splunk_lookup_definition" "example" {
	name     = "example_lookup_definition"
	filename = "example_lookup_file2.csv"
	acl {
	  owner = "admin"
	  app = "search"
	  sharing = "user"
	}
  }
`

const updateSplunkLookupDefinitionObjectWithAcl = `
resource "splunk_lookup_definition" "example" {
  name     = "example_lookup_definition"
  filename = "example_lookup_file3.csv"
  acl {
    owner = "admin"
    app = "search"
	sharing = "global"
    read = ["*"]
    write = ["admin"]
  }
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
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "app"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.0", "admin"),
				),
			},
			{
				Config: updateSplunkLookupDefinitionObjectWithAclSharingUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "example_lookup_definition"),
					resource.TestCheckResourceAttr(resourceName, "filename", "example_lookup_file2.csv"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "user"),
					resource.TestCheckNoResourceAttr(resourceName, "acl.0.read.0"),
					resource.TestCheckNoResourceAttr(resourceName, "acl.0.write.0"),
				),
			},
			{
				Config: updateSplunkLookupDefinitionObjectWithAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "example_lookup_definition"),
					resource.TestCheckResourceAttr(resourceName, "filename", "example_lookup_file3.csv"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "global"),
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

func TestAccSplunkLookupDefinitionValidationNewResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkLookupDefinitionDestroyResources,
		Steps: []resource.TestStep{
			{
				Config:      splunkLookupDefinitionObjectWithAclFailedValidation,
				ExpectError: regexp.MustCompile("errors during plan: `acl.read` cannot be set when `acl.sharing` is `user`"),
			},
		},
	})
}

func TestAccSplunkLookupDefinitionValidationUpdateResource(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "app"),
				),
			},
			{
				Config:      splunkLookupDefinitionObjectWithAclFailedValidation,
				ExpectError: regexp.MustCompile("errors during plan: `acl.read` cannot be set when `acl.sharing` is `user`"),
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
