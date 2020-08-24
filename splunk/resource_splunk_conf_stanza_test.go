package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const newConfStanza = `
resource "splunk_conf_stanza" "tftest-stanza" {
	name = "tftest/tftest_stanza"
	variables = {
		"key": "value"
	}
}
`

const updateConfStanza = `
resource "splunk_conf_stanza" "tftest-stanza" {
	name = "tftest/tftest_stanza"
	variables = {
		"key": "new-value"
	}
}
`

func TestAccCreateSplunkConfStanza(t *testing.T) {
	resourceName := "splunk_conf_stanza.tftest-stanza"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkConfStanzaDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newConfStanza,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "variables.key", "value"),
				),
			},
			{
				Config: updateConfStanza,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "variables.key", "new-value"),
				),
			},
			{
				ResourceName:      "splunk_conf_stanza.tftest-stanza",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkConfStanzaDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_conf_stanza":
			endpoint := client.BuildSplunkURL(nil, "services", "configs", "conf", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
