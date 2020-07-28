package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
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
    disabled = true
    use_ack = false

    acl {
      owner = "splunker"
      sharing = "app"
      read = ["admin", "splunker"]
      write = ["admin"]
    }

    depends_on = ["splunk_global_http_event_collector.http"]
}
`

const updateHttpEventCollectorInput = `
resource "splunk_global_http_event_collector" "http" {
  disabled     = false
  enable_ssl   = false
  port         = 8088
}

resource "splunk_input_http_event_collector" "token" {
    name = "new_token"
    index = "main"
    indexes = ["main", "history"]
    source = "new"
    sourcetype = "new"
    disabled = false
    use_ack = true

    acl {
      owner = "splunker"
      sharing = "global"
      read = ["admin", "splunker"]
      write = ["admin", "splunker"]
    }

    depends_on = ["splunk_global_http_event_collector.http"]
}
`

func TestAccCreateSplunkHttpEventCollectorInput(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkHttpEventCollectorInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newHttpEventCollectorInput,
			},
			{
				Config: updateHttpEventCollectorInput,
			},
		},
	})
}

func testAccSplunkHttpEventCollectorInputDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_input_http_event_collector":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "splunk_httpinput", "data", "inputs", "http", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
