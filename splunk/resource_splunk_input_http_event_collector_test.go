package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

const newHttpEventCollectorInput = `

resource "splunk_global_http_event_collector" "http" {
  disabled     = false
  enable_ssl   = false
  port         = 8088
}

resource "splunk_input_http_event_collector" "token" {
    name = "new_token"
    index = "main"
    source = "new"
    sourcetype = "new"
    disabled = false
    use_ack = false

    depends_on = ["splunk_global_http_event_collector.http"]
}
`

const updateHttpEventCollectorInput = `
resource "splunk_input_http_event_collector" "token" {
    name = "new_token"
    app_context = "search"
    index = "main"
    source = "new"
    sourcetype = "new"
    disabled = false
    use_ack = true
}
`

func TestAccCreateSplunkHttpEventCollectorInput(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: testAccSplunkHttpEventCollectorInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config:newHttpEventCollectorInput,
			},
			{
				Config:updateHttpEventCollectorInput,
			},
		},
	})
}

func testAccSplunkHttpEventCollectorInputDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		default:
			splunkHttpInputConfig, err := client.doGetHttpInput(&SplunkHttpInputConfig{
				Name: rs.Primary.ID,
			})
			if splunkHttpInputConfig != nil {
				return fmt.Errorf("found deleted token: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
