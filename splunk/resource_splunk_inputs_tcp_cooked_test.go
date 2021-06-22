package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const newTCPCookedInput = `
resource "splunk_inputs_tcp_cooked" "test" {
    name = "50000"
    disabled = false
    restrict_to_host = "splunk"
}
`

const updateTCPCookedInput = `

resource "splunk_inputs_tcp_cooked" "test" {
    name = "50000"
    disabled = false
    connection_host = "dns"
}
`

func TestAccSplunkTCPCookedInput(t *testing.T) {
	resourceName := "splunk_inputs_tcp_cooked.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkTCPCookedInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newTCPCookedInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "connection_host", "ip"),
					resource.TestCheckResourceAttr(resourceName, "restrict_to_host", "splunk"),
				),
			},
			{
				Config: updateTCPCookedInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "connection_host", "dns"),
					resource.TestCheckResourceAttr(resourceName, "restrict_to_host", "splunk"),
				),
			},
			{
				ResourceName:      "splunk_inputs_tcp_cooked.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkTCPCookedInputDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_inputs_tcp_cooked":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "search", "data", "inputs", "tcp", "splunk:50000")
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
