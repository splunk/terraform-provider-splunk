package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const newHttpEventCollectorInput = `
resource "splunk_global_http_event_collector" "http" {
  disabled     = false
  enable_ssl   = true
}

resource "splunk_inputs_http_event_collector" "new-token" {
    name = "new-token"
    source = "new"
    disabled = false
    use_ack = 0

    acl {
      app = "launcher"
      sharing = "global"
    }

    depends_on = ["splunk_global_http_event_collector.http"]
}
`

const updateHttpEventCollectorInput = `
resource "splunk_global_http_event_collector" "http" {
  disabled     = false
  enable_ssl   = true
  port         = 8088
}

resource "splunk_inputs_http_event_collector" "new-token" {
    name = "new-token"
    index = "main"
    indexes = ["main", "history"]
    source = "new"
    sourcetype = "new"
    disabled = false
    use_ack = 1

    depends_on = ["splunk_global_http_event_collector.http"]
}
`

const newHttpEventCollectorInputWithToken = `
resource "splunk_global_http_event_collector" "http" {
  disabled     = false
  enable_ssl   = true
}

resource "splunk_inputs_http_event_collector" "new-token" {
    name = "new-token"
    token = "12345678-0000-0000-0000-123456780000"
    source = "new"
    disabled = false
    use_ack = 0

    acl {
      app = "launcher"
      sharing = "global"
    }

    depends_on = ["splunk_global_http_event_collector.http"]
}
`

func TestAccSplunkHttpEventCollectorInput(t *testing.T) {
	resourceName := "splunk_inputs_http_event_collector.new-token"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkHttpEventCollectorInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newHttpEventCollectorInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source", "new"),
					resource.TestCheckResourceAttr(resourceName, "index", "default"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "use_ack", "0"),
					resource.TestCheckResourceAttr(resourceName, "acl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "launcher"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "nobody"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "global"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.0", "admin"),
				),
			},
			{
				Config: updateHttpEventCollectorInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source", "new"),
					resource.TestCheckResourceAttr(resourceName, "index", "main"),
					resource.TestCheckResourceAttr(resourceName, "sourcetype", "new"),
					resource.TestCheckResourceAttr(resourceName, "indexes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "indexes.0", "main"),
					resource.TestCheckResourceAttr(resourceName, "indexes.1", "history"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "use_ack", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "launcher"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "nobody"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "global"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.0", "admin"),
				),
			},
			{
				ResourceName:      "splunk_inputs_http_event_collector.new-token",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSplunkHttpEventCollectorInputWithToken(t *testing.T) {
	resourceName := "splunk_inputs_http_event_collector.new-token"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkHttpEventCollectorInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newHttpEventCollectorInputWithToken,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source", "new"),
					resource.TestCheckResourceAttr(resourceName, "index", "default"),
					resource.TestCheckResourceAttr(resourceName, "token", "12345678-0000-0000-0000-123456780000"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "use_ack", "0"),
					resource.TestCheckResourceAttr(resourceName, "acl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "launcher"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "nobody"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "global"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.0", "admin"),
				),
			},
			{
				ResourceName:      "splunk_inputs_http_event_collector.new-token",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkHttpEventCollectorInputDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_inputs_http_event_collector":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "splunk_httpinput", "data", "inputs", "http", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
