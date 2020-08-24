package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const newTCPGroupOutput = `
resource "splunk_outputs_tcp_group" "test" {
    name = "tcp-group"
    disabled = false
    drop_events_on_queue_full = 60
    send_cooked_data = true
    servers = ["1.1.1.1:1234", "2.2.2.2:1234"]
}
`

const updateTCPGroupOutput = `
resource "splunk_outputs_tcp_group" "test" {
    name = "tcp-group"
    disabled = false
    drop_events_on_queue_full = 60
    send_cooked_data = true
    servers = ["1.1.1.1:1234", "2.2.2.2:1234", "3.3.3.3:1234"]
}
`

func TestAccSplunkTCPGroupOutput(t *testing.T) {
	resourceName := "splunk_outputs_tcp_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkTCPGroupOutputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newTCPGroupOutput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tcp-group"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "drop_events_on_queue_full", "60"),
					resource.TestCheckResourceAttr(resourceName, "send_cooked_data", "true"),
					resource.TestCheckResourceAttr(resourceName, "servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "servers.0", "1.1.1.1:1234"),
					resource.TestCheckResourceAttr(resourceName, "servers.1", "2.2.2.2:1234"),
				),
			},
			{
				Config: updateTCPGroupOutput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tcp-group"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "drop_events_on_queue_full", "60"),
					resource.TestCheckResourceAttr(resourceName, "send_cooked_data", "true"),
					resource.TestCheckResourceAttr(resourceName, "servers.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "servers.0", "1.1.1.1:1234"),
					resource.TestCheckResourceAttr(resourceName, "servers.1", "2.2.2.2:1234"),
					resource.TestCheckResourceAttr(resourceName, "servers.2", "3.3.3.3:1234"),
				),
			},
			{
				ResourceName:      "splunk_outputs_tcp_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkTCPGroupOutputDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_outputs_tcp_group":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "search", "data", "outputs", "tcp", "group", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
