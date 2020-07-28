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
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkGlobalHttpInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: globalHttpInputConfig,
			},
		},
	})
}

func testAccSplunkGlobalHttpInputDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		default:
			endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "http", rs.Primary.ID)
			_, err := client.Get(endpoint)
			if err != nil {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
