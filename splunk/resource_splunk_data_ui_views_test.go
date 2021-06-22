package splunk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const splunkDashboardsObjectWithAcl = `
resource "splunk_data_ui_views" "dashboard" {
  name     = "terraform"
  eai_data = "<dashboard><label>Terraform Test Dashboard</label></dashboard>"
  acl {
    owner = "admin"
    app = "search"
    sharing = "global"
    read = ["*"]
    write = ["admin"]
  }
}
`

const updateSplunkDashboardsObjectWithAcl = `
resource "splunk_data_ui_views" "dashboard" {
  name     = "terraform"
  eai_data = "<dashboard><label>Terraform Test Dashboard</label></dashboard>"
}
`

const splunkDashboardsObject = `
resource "splunk_data_ui_views" "dashboard" {
  name     = "terraform"
  eai_data = "<dashboard><label>Terraform Test Dashboard</label></dashboard>"
}
`

const updateSplunkDashboardsObject = `

resource "splunk_data_ui_views" "dashboard" {
  name     = "terraform"
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
					resource.TestCheckResourceAttr(resourceName, "name", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "eai_data", "<dashboard><label>Terraform Test Dashboard</label></dashboard>"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "user"),
				),
			},
			{
				Config: updateSplunkDashboardsObject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "eai_data", "<dashboard><label>Terraform Test Dashboard</label></dashboard>"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "user"),
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

func TestAccSplunkDashboardsWithAcl(t *testing.T) {
	resourceName := "splunk_data_ui_views.dashboard"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkDashboardDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: splunkDashboardsObjectWithAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "eai_data", "<dashboard><label>Terraform Test Dashboard</label></dashboard>"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "global"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.0", "admin"),
				),
			},
			{
				Config: updateSplunkDashboardsObjectWithAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "eai_data", "<dashboard><label>Terraform Test Dashboard</label></dashboard>"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "global"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.0", "admin"),
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
	client, err := newTestClient()
	if err != nil {
		return err
	}
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
