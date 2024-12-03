package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const newConfigsConf = `
resource "splunk_configs_conf" "tftest-stanza" {
	name = "tf_test/sqs://tftest_stanza"
	variables = {
        "disabled": "false"
		"key": "value"
	}
}
`

const updateConfigsConf = `
resource "splunk_configs_conf" "tftest-stanza" {
	name = "tf_test/sqs://tftest_stanza"
	variables = {
        "disabled": "false"
		"key": "new-value"
	}
}
`

func TestAccCreateSplunkConfigsConf(t *testing.T) {
	resourceName := "splunk_configs_conf.tftest-stanza"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkConfigsConfDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newConfigsConf,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "variables.key", "value"),
				),
			},
			{
				Config: updateConfigsConf,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "variables.key", "new-value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkConfigsConfDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_configs_conf":
			endpoint := client.BuildSplunkURL(nil, "services", "configs", "conf", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
