package splunk

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const splunkDashboardsObject = `

resource "splunk_dashboards" "dashboard" {
  name     = "test"
  eai_data = "<dashboard><label>the_new_label</label></dashboard>"
}
`

func TestAccSplunkDashboards(t *testing.T) {
	resourceName := "splunk_dashboards.dashboard"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: splunkDashboardsObject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "eai_data", "<dashboard><label>the_new_label</label></dashboard>"),
				),
			},
			{
				ResourceName:      "splunk_dashboards.dashboard",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
