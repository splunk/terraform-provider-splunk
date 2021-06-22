package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const newSplunkTCPTokenInput = `
resource "splunk_inputs_tcp_splunk_tcp_token" "test" {
    name = "new-splunk-tcp-token"
    token = "D66C45B3-7C28-48A1-A13A-027914146501"
}
`

func TestAccSplunkTCPTokenInput(t *testing.T) {
	resourceName := "splunk_inputs_tcp_splunk_tcp_token.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkTCPTokenInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newSplunkTCPTokenInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-splunk-tcp-token"),
					resource.TestCheckResourceAttr(resourceName, "token", "D66C45B3-7C28-48A1-A13A-027914146501"),
				),
			},
			{
				ResourceName:      "splunk_inputs_tcp_splunk_tcp_token.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkTCPTokenInputDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_inputs_tcp_splunk_tcp_token":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "search", "data", "inputs", "tcp", "splunktcptoken", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
