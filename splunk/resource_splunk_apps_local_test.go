package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const newAppsLocal = `
resource "splunk_apps_local" "app" {
    name = "new-app"
    author = "Splunk Works"
	configured = false
	description = "New App"
	label = "New App"
	version = "0.0.4"
	visible = true
    update = true
}`

const updateAppsLocal = `
resource "splunk_apps_local" "app" {
    name = "new-app"
    author = "Splunk Works"
	configured = false
	description = "New App"
	label = "New App"
	version = "0.0.5"
	visible = true
}
`

func TestAccSplunkAppsLocal(t *testing.T) {
	resourceName := "splunk_apps_local.app"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkAppsLocalDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newAppsLocal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "author", "Splunk Works"),
					resource.TestCheckResourceAttr(resourceName, "configured", "false"),
					resource.TestCheckResourceAttr(resourceName, "version", "0.0.4"),
					resource.TestCheckResourceAttr(resourceName, "visible", "true"),
				),
			},
			{
				Config: updateAppsLocal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "author", "Splunk Works"),
					resource.TestCheckResourceAttr(resourceName, "configured", "false"),
					resource.TestCheckResourceAttr(resourceName, "version", "0.0.5"),
					resource.TestCheckResourceAttr(resourceName, "visible", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "filename", "explicit_appname", "update"},
			},
		},
	})
}

func testAccSplunkAppsLocalDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_apps_local":
			endpoint := client.BuildSplunkURL(nil, "services", "apps", "local", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
