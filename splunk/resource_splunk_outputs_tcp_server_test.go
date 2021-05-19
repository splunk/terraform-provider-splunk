package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const newTCPServerOutput = `
resource "splunk_outputs_tcp_server" "test" {
    name = "new-host:1234"
}
`

const updateTCPServerOutput = `
resource "splunk_outputs_tcp_server" "test" {
    name = "new-host:1234"
    ssl_alt_name_to_check = "old-host"
}
`

func TestAccSplunkTCPServerOutput(t *testing.T) {
	resourceName := "splunk_outputs_tcp_server.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkTCPServerOutputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newTCPServerOutput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-host:1234"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "method", "autobalance"),
				),
			},
			{
				Config: updateTCPServerOutput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-host:1234"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "method", "autobalance"),
					resource.TestCheckResourceAttr(resourceName, "ssl_alt_name_to_check", "old-host"),
				),
			},
			{
				ResourceName:      "splunk_outputs_tcp_server.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkTCPServerOutputDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_outputs_tcp_server":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "search", "data", "outputs", "tcp", "server", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
