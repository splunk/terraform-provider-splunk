package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const adminSAMLGroupsInput = `
resource "splunk_admin_saml_groups" "saml_group" {
  name = "new-saml-group"
  roles = ["admin", "power"]
}
`

const updateAdminSAMLGroupsInput = `
resource "splunk_admin_saml_groups" "saml_group" {
  name = "new-saml-group"
  roles = ["admin", "power", "user"]
}
`

func TestAccSplunkAdminSAMLGroups(t *testing.T) {
	resourceName := "splunk_admin_saml_groups.saml_group"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkAdminSAMLGroupsInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: adminSAMLGroupsInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-saml-group"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "roles.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.1", "power"),
				),
			},
			{
				// to test re-creation of remotely deleted or missing resources, delete the new saml group before updating it
				PreConfig: func() {
					client, _ := newTestClient()
					if _, err := client.DeleteAdminSAMLGroups("new-saml-group"); err != nil {
						t.Error("PreConfig deletion of new-saml-group failed")
					}
				},
				Config: updateAdminSAMLGroupsInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-saml-group"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "roles.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.1", "power"),
					resource.TestCheckResourceAttr(resourceName, "roles.2", "user"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkAdminSAMLGroupsInputDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_authorization_roles":
			endpoint := client.BuildSplunkURL(nil, "services", "admin", "SAML-groups", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
