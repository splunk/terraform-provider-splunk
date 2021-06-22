package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

const createInputsTCPSSLConfig = `

resource "splunk_inputs_tcp_ssl" "test" {
  disabled     = false
  require_client_cert = true
}
`

const updateInputsTCPSSLConfig = `

resource "splunk_inputs_tcp_ssl" "test" {
  disabled     = false
  require_client_cert = false
}
`

func TestAccInputsTCPSSL(t *testing.T) {
	resourceName := "splunk_inputs_tcp_ssl.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccTCPSSLInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: createInputsTCPSSLConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "require_client_cert", "true"),
				),
			},
			{
				Config: updateInputsTCPSSLConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "require_client_cert", "false"),
				),
			},
			{
				ResourceName:      "splunk_inputs_tcp_ssl.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTCPSSLInputDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		default:
			endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "tcp", "ssl", rs.Primary.ID)
			_, err := client.Get(endpoint)
			if err != nil {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
