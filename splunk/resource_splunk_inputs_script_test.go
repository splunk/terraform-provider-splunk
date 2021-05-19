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

const inputsScriptConfig = `
resource "splunk_inputs_script" "script" {
  name     = "%v"
  interval = 60
}
`

const inputsScriptUpdateConfig = `
resource "splunk_inputs_script" "script" {
  name     = "%v"
  interval = 120
}
`

func TestAccSplunkInputsScript(t *testing.T) {
	// Required since scripted inputs are limited to certain paths within $SPLUNK_HOME
	path := os.Getenv("SPLUNK_HOME")
	f := filepath.Join(path, "/bin/scripts/readme.txt")
	config := fmt.Sprintf(inputsScriptConfig, f)
	updatedConfig := fmt.Sprintf(inputsScriptUpdateConfig, f)

	resourceName := "splunk_inputs_script.script"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkInputsScriptDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "interval", "60"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "interval", "120"),
				),
			},
			{
				ResourceName:      "splunk_inputs_script.script",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkInputsScriptDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		default:
			endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "script", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}

	return nil
}
