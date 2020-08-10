package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const newIndex = `
resource "splunk_index" "new-index" {
    name = "new-index"
    acl {
      app = "launcher"
      sharing = "global"
    }
}
`

const updateIndex = `
resource "splunk_index" "new-index" {
    max_time_unreplicated_no_acks = 301

    acl {
      app = "launcher"
      sharing = "global"
      read = ["admin"]
      write = ["admin"]
    }
}
`

func TestAccCreateSplunkIndex(t *testing.T) {
	resourceName := "splunk_index.new-index"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkIndexDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newIndex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "maxTimeUnreplicatedNoAcks", "300"),
				),
			},
			{
				Config: updateIndex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "maxTimeUnreplicatedNoAcks", "301"),
				),
			},
			{
				ResourceName:      "splunk_index.new-index",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkIndexDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_index":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "search", "data", "indexes", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
