package splunk

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"terraform-provider-splunk/client/models"
	"testing"
)

const newTCPDefaultOutput = `
resource "splunk_outputs_tcp_default" "test" {
    disabled = false
    default_group = "test-indexers"
    drop_events_on_queue_full = 60
    index_and_forward = true
    send_cooked_data = true
    max_queue_size = "100KB"
}
`

const updateTCPDefaultOutput = `
resource "splunk_outputs_tcp_default" "test" {
    disabled = false
    default_group = "test-indexers"
    drop_events_on_queue_full = 60
    index_and_forward = true
    send_cooked_data = true
    max_queue_size = "100KB"
}
`

func TestAccSplunkTCPDefaultOutput(t *testing.T) {
	resourceName := "splunk_outputs_tcp_default.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkTCPDefaultOutputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newTCPDefaultOutput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tcpout"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_group", "test-indexers"),
					resource.TestCheckResourceAttr(resourceName, "drop_events_on_queue_full", "60"),
					resource.TestCheckResourceAttr(resourceName, "index_and_forward", "true"),
					resource.TestCheckResourceAttr(resourceName, "max_queue_size", "100KB"),
					resource.TestCheckResourceAttr(resourceName, "send_cooked_data", "true"),
				),
			},
			{
				Config: updateTCPDefaultOutput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tcpout"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_group", "test-indexers"),
					resource.TestCheckResourceAttr(resourceName, "drop_events_on_queue_full", "60"),
					resource.TestCheckResourceAttr(resourceName, "index_and_forward", "true"),
					resource.TestCheckResourceAttr(resourceName, "max_queue_size", "100KB"),
					resource.TestCheckResourceAttr(resourceName, "send_cooked_data", "true"),
				),
			},
			{
				ResourceName:      "splunk_outputs_tcp_default.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkTCPDefaultOutputDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_outputs_tcp_default":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "search", "data", "outputs", "tcp", "default", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if err != nil {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
			defer resp.Body.Close()
			v := &models.OutputsTCPDefaultResponse{}
			_ = json.NewDecoder(resp.Body).Decode(&v)
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
			if v.Entry[0].Content.Disabled == false {
				return fmt.Errorf("TCP outputs should be disabled: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
