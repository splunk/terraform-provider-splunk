package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

const globalHttpInputConfig = `

resource "splunk_global_http_event_collector" "http" {
  disabled     = false
  enable_ssl   = false
}
`

const globalUpdateHttpInputConfig = `

resource "splunk_global_http_event_collector" "http" {
  disabled     = false
  enable_ssl   = true
  port         = 8088
  dedicated_io_threads = 2
}
`

func TestAccGlobalSplunkHttpEventCollectorInput(t *testing.T) {
	resourceName := "splunk_global_http_event_collector.http"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkGlobalHttpInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: globalHttpInputConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "8088"),
					resource.TestCheckResourceAttr(resourceName, "enable_ssl", "false"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "max_threads", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_sockets", "0"),
					resource.TestCheckResourceAttr(resourceName, "use_deployment_server", "0"),
				),
			},
			{
				Config: globalUpdateHttpInputConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "8088"),
					resource.TestCheckResourceAttr(resourceName, "enable_ssl", "true"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "dedicated_io_threads", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_threads", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_sockets", "0"),
					resource.TestCheckResourceAttr(resourceName, "use_deployment_server", "0"),
				),
			},
			{
				ResourceName:      "splunk_global_http_event_collector.http",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkGlobalHttpInputDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		default:
			endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "http", rs.Primary.ID)
			_, err := client.Get(endpoint)
			if err != nil {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
