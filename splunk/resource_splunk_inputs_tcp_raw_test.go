package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const newTCPRawInput = `
resource "splunk_inputs_tcp_raw" "test" {
    name = "41000"
    source = "new"
    queue = "parsingQueue"
    disabled = false
}
`

const updateTCPRawInput = `
resource "splunk_inputs_tcp_raw" "test" {
    name = "41000"
    index = "main"
    queue = "indexQueue"
    source = "new"
    sourcetype = "new"
    disabled = false
}
`

func TestAccSplunkTCPRawInput(t *testing.T) {
	resourceName := "splunk_inputs_tcp_raw.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkTCPRawInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newTCPRawInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source", "new"),
					resource.TestCheckResourceAttr(resourceName, "index", "default"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "connection_host", "dns"),
					resource.TestCheckResourceAttr(resourceName, "queue", "parsingQueue"),
					resource.TestCheckResourceAttr(resourceName, "restrict_to_host", ""),
				),
			},
			{
				Config: updateTCPRawInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source", "new"),
					resource.TestCheckResourceAttr(resourceName, "index", "main"),
					resource.TestCheckResourceAttr(resourceName, "sourcetype", "new"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "queue", "indexQueue"),
					resource.TestCheckResourceAttr(resourceName, "connection_host", "dns"),
					resource.TestCheckResourceAttr(resourceName, "restrict_to_host", ""),
				),
			},
			{
				ResourceName:      "splunk_inputs_tcp_raw.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkTCPRawInputDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_inputs_tcp_raw":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "search", "data", "inputs", "tcp", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
