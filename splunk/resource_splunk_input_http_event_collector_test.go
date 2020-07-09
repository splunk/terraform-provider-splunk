package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

const newHttpEventCollectorInput = `
resource "splunk_input_http_event_collector" "token" {
    name = "new_token"
    index = "main"
    source = "new"
    sourcetype = "new"
    disabled = false
    use_ack = false
}
`

const updateHttpEventCollectorInput = `
resource "splunk_input_http_event_collector" "token" {
    name = "new_token"
    index = "main"
    source = "new"
    sourcetype = "new"
    disabled = false
    use_ack = true
}
`

func TestAccCreateSplunkHttpEventCollectorInput(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: testAccDestroyResources,
		Steps: []resource.TestStep{
			{
				Config:newHttpEventCollectorInput,
			},
			{
				Config:updateHttpEventCollectorInput,
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SPLUNK_URL"); v == "" {
		t.Fatal("SPLUNK_URL must be set for acceptance tests")
	}
	if v := os.Getenv("SPLUNK_USERNAME"); v == "" {
		t.Fatal("SPLUNK_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("SPLUNK_PASSWORD"); v == "" {
		t.Fatal("SPLUNK_PASSWORD must be set for acceptance tests")
	}
}

func testAccDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		default:
			splunkHttpInputConfig, err := client.doGetHttpInput(&SplunkHttpInputConfig{Name: rs.Primary.ID})
			if splunkHttpInputConfig != nil {
				return fmt.Errorf("found deleted token: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
