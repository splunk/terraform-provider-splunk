package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

const inputsMonitorConfig = `
resource "splunk_inputs_monitor" "monitor" {
  name     = "%v"
  sourcetype = "text"
}
`

const inputsMonitorUpdateConfig = `
resource "splunk_inputs_monitor" "monitor" {
  name     = "%v"
  recursive = true
  sourcetype = "text"
}
`

func TestAccSplunkInputsMonitor(t *testing.T) {
	if os.Getenv("SKIP_DOCKER") != "" {
		t.Skip("Skipping not finished test")
	}
	// Required since monitor inputs are limited to certain paths within $SPLUNK_HOME
	path := os.Getenv("SPLUNK_HOME")
	f := filepath.Join(path, "/var/log/splunk/health.log")
	config := fmt.Sprintf(inputsMonitorConfig, f)
	updatedConfig := fmt.Sprintf(inputsMonitorUpdateConfig, f)

	resourceName := "splunk_inputs_monitor.monitor"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkInputsMonitorDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "recursive", "false"),
					resource.TestCheckResourceAttr(resourceName, "follow_tail", "false"),
					resource.TestCheckResourceAttr(resourceName, "sourcetype", "text"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "recursive", "true"),
					resource.TestCheckResourceAttr(resourceName, "follow_tail", "false"),
					resource.TestCheckResourceAttr(resourceName, "sourcetype", "text"),
				),
			},
			{
				ResourceName:      "splunk_inputs_monitor.monitor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkInputsMonitorDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		default:
			endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "monitor", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}

	return nil
}
