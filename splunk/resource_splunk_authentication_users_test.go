package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const authenticationUserInput = `
resource "splunk_authentication_users" "user" {
  name = "new-user"
  password = "changeme"
  force_change_pass = true
  roles = ["admin"]
}
`

const updateAuthenticationUserInput = `
resource "splunk_authentication_users" "user" {
  name = "new-user"
  default_app = "search"
  force_change_pass = true
  password = "changeme"
  roles = ["admin", "user"]
}
`

func TestAccSplunkAuthenticationUsers(t *testing.T) {
	resourceName := "splunk_authentication_users.user"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkAuthenticationUsersInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: authenticationUserInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-user"),
					resource.TestCheckResourceAttr(resourceName, "force_change_pass", "true"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "roles.0", "admin"),
				),
			},
			{
				Config: updateAuthenticationUserInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-user"),
					resource.TestCheckResourceAttr(resourceName, "default_app", "search"),
					resource.TestCheckResourceAttr(resourceName, "force_change_pass", "true"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "roles.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.1", "user"),
				),
			},
			{
				ResourceName:            "splunk_authentication_users.user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "force_change_pass"},
			},
		},
	})
}

func testAccSplunkAuthenticationUsersInputDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_authentication_users":
			endpoint := client.BuildSplunkURL(nil, "services", "authentication", "user", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
