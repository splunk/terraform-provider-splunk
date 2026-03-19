package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const newSavedEventTypes = `
resource "splunk_saved_event_types" "event-type" {
  name        = "test"
  description = "Test New event description"
  disabled    = false
  priority    = 1
  search      = "index=main"
  color       = "et_blue"
  tags        = ["tag"]
  acl {
    owner   = "admin"
    sharing = "app"
    app     = "launcher"
  }
}
`

func TestAccSplunkSavedEventTypes(t *testing.T) {
	resourceName := "splunk_saved_event_types.event-type"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkSavedEventsDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newSavedEventTypes,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "search", "index=main"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test New event description"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "color", "et_blue"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag"),
				),
			},
			{
				ResourceName:      "splunk_saved_event_types.event-type",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkSavedEventsDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_saved_event_types":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "search", "saved", "eventtypes", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
