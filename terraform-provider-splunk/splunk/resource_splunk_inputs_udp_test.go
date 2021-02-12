package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const newUDPInput = `
resource "splunk_inputs_udp" "test" {
    name = "41000"
    source = "new"
    disabled = false
}
`

const updateUDPInput = `

resource "splunk_inputs_udp" "test" {
    name = "41000"
    index = "main"
    source = "new"
    sourcetype = "new"
    disabled = false
}
`

func TestAccSplunkUDPInput(t *testing.T) {
	resourceName := "splunk_inputs_udp.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkUDPInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newUDPInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source", "new"),
					resource.TestCheckResourceAttr(resourceName, "index", "default"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "connection_host", "ip"),
					resource.TestCheckResourceAttr(resourceName, "queue", ""),
					resource.TestCheckResourceAttr(resourceName, "restrict_to_host", ""),
					resource.TestCheckResourceAttr(resourceName, "no_appending_timestamp", "false"),
					resource.TestCheckResourceAttr(resourceName, "no_priority_stripping", "false"),
				),
			},
			{
				Config: updateUDPInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source", "new"),
					resource.TestCheckResourceAttr(resourceName, "index", "main"),
					resource.TestCheckResourceAttr(resourceName, "sourcetype", "new"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "queue", ""),
					resource.TestCheckResourceAttr(resourceName, "connection_host", "ip"),
					resource.TestCheckResourceAttr(resourceName, "restrict_to_host", ""),
					resource.TestCheckResourceAttr(resourceName, "no_priority_stripping", "false"),
				),
			},
			{
				ResourceName:      "splunk_inputs_udp.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkUDPInputDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_inputs_udp":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "search", "data", "inputs", "udp", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
