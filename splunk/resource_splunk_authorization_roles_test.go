package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const authorizationRolesInput = `
resource "splunk_authorization_roles" "role" {
  name = "new-role"
  default_app = "search"
  imported_roles = ["user"]
  capabilities = ["accelerate_datamodel", "change_authentication"]
}
`

const updateAuthorizationRolesInput = `
resource "splunk_authorization_roles" "role" {
  name = "new-role"
  default_app = "search"
  imported_roles = ["power", "user"]
  capabilities = ["accelerate_datamodel", "change_authentication", "restart_splunkd"]
}
`

func TestAccSplunkAuthorizationRoles(t *testing.T) {
	resourceName := "splunk_authorization_roles.role"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkAuthorizationRolesInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: authorizationRolesInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-role"),
					resource.TestCheckResourceAttr(resourceName, "default_app", "search"),
					resource.TestCheckResourceAttr(resourceName, "imported_roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "imported_roles.0", "user"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.0", "accelerate_datamodel"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.1", "change_authentication"),
				),
			},
			{
				Config: updateAuthorizationRolesInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-role"),
					resource.TestCheckResourceAttr(resourceName, "default_app", "search"),
					resource.TestCheckResourceAttr(resourceName, "imported_roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "imported_roles.0", "power"),
					resource.TestCheckResourceAttr(resourceName, "imported_roles.1", "user"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.0", "accelerate_datamodel"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.1", "change_authentication"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.2", "restart_splunkd"),
				),
			},
			{
				ResourceName:      "splunk_authorization_roles.role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkAuthorizationRolesInputDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_authorization_roles":
			endpoint := client.BuildSplunkURL(nil, "services", "authorization", "roles", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
