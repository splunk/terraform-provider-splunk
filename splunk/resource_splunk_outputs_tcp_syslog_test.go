package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const newTCPSyslogOutput = `
resource "splunk_outputs_tcp_syslog" "test" {
    name = "new-syslog"
    server = "new-host-1:1234"
}
`

const updateTCPSyslogOutput = `
resource "splunk_outputs_tcp_syslog" "test" {
    name = "new-syslog"
    server = "new-host-1:1234"
    priority = 5
}
`

func TestAccSplunkTCPSyslogOutput(t *testing.T) {
	resourceName := "splunk_outputs_tcp_syslog.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkTCPSyslogOutputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newTCPSyslogOutput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-syslog"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "server", "new-host-1:1234"),
					resource.TestCheckResourceAttr(resourceName, "priority", "0"),
				),
			},
			{
				Config: updateTCPSyslogOutput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-syslog"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "server", "new-host-1:1234"),
					resource.TestCheckResourceAttr(resourceName, "priority", "5"),
				),
			},
			{
				ResourceName:      "splunk_outputs_tcp_syslog.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkTCPSyslogOutputDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_outputs_tcp_syslog":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "search", "data", "outputs", "tcp", "server", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
