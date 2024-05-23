package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const AdminProxyssoGroupsInput = `
resource "splunk_admin_proxysso_groups" "proxy_sso_group" {
  name = "newrole"
  roles = ["admin", "power"]
}
`

const updateAdminProxyssoGroupsInput = `
resource "splunk_admin_proxysso_groups" "proxy_sso_group" {
  name = "newrole"
  roles = ["admin", "power", "user"]
}
`

func TestAccSplunkAdminProxyssoGroups(t *testing.T) {
	resourceName := "splunk_admin_proxysso_groups.proxy_sso_group"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkAdminProxyssoGroupsInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: AdminProxyssoGroupsInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "newrole"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "roles.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.1", "power"),
				),
			},
			{
				Config: updateAdminProxyssoGroupsInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "newrole"),
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

func testAccSplunkAdminProxyssoGroupsInputDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_authorization_roles":
			endpoint := client.BuildSplunkURL(nil, "services", "admin", "ProxySSO-groups", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
