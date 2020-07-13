package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

const globalHttpInputConfig = `

resource "splunk_global_http_event_collector" "http" {
  disabled     = false
  enable_ssl   = false
  port         = 8088
}
`

func TestAccGlobalSplunkHttpEventCollectorInput(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: testAccGlobalSplunkHttpInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config:globalHttpInputConfig,
			},
		},
	})
}

func testAccGlobalSplunkHttpInputDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		default:
			_, err := client.doGetGlobalHttpInput(&SplunkGlobalHttpInputConfig{Name: rs.Primary.ID})
			if err != nil {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}