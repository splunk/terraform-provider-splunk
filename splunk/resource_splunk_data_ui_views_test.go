package splunk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const splunkDashboardsObject = `
resource "splunk_data_ui_views" "dashboard" {
  name     = "Terraform_Test_Dashboard"
  eai_data = "<dashboard><label>Terraform Test Dashboard</label></dashboard>"
  acl {
    owner = "admin"
    app = "search"
  }
}
`

const updateSplunkDashboardsObject = `

resource "splunk_data_ui_views" "dashboard" {
  name     = "Terraform_Test_Dashboard"
  eai_data = "<dashboard><label>Terraform Test Dashboard</label></dashboard>"
}
`

func TestAccSplunkDashboards(t *testing.T) {
	resourceName := "splunk_data_ui_views.dashboard"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkDashboardDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: splunkDashboardsObject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform_Test_Dashboard"),
					resource.TestCheckResourceAttr(resourceName, "eai_data", "<dashboard><label>Terraform Test Dashboard</label></dashboard>"),
				),
			},
			{
				Config: updateSplunkDashboardsObject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform_Test_Dashboard"),
					resource.TestCheckResourceAttr(resourceName, "eai_data", "<dashboard><label>Terraform Test Dashboard</label></dashboard>"),
				),
			},
			{
				ResourceName:      "splunk_data_ui_views.dashboard",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkDashboardDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_data_ui_views":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "admin", "search", "data", "ui", "views", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
