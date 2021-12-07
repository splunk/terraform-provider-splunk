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
  # capabilities are explicitly listed in a different order than returned by the Splunk API
  capabilities = ["change_authentication", "accelerate_datamodel"]
  search_indexes_allowed = ["_audit", "_internal"]
  search_indexes_default = ["_audit", "_internal"]
}
`

const updateAuthorizationRolesInput = `
resource "splunk_authorization_roles" "role" {
  name = "new-role"
  default_app = "search"
  imported_roles = ["power", "user"]
  # capabilities are explicitly listed in a different order than returned by the Splunk API
  capabilities = ["restart_splunkd", "change_authentication", "accelerate_datamodel"]
  search_indexes_allowed = ["_audit", "_internal", "main"]
  search_indexes_default = ["_audit", "_internal", "main"]
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
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.0", "_audit"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.1", "_internal"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.0", "_audit"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.1", "_internal"),
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
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.0", "_audit"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.1", "_internal"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.2", "main"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.0", "_audit"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.1", "_internal"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.2", "main"),
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
	client, err := newTestClient()
	if err != nil {
		return err
	}
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
