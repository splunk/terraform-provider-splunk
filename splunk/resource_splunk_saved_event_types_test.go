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
    name 		= "test"
    description = "Test New event description"
    disabled 	= "0"
    priority 	= 1
    search 		= "index=main"
    color		= "blue"
    tags 		= "tag"
    acl {
      owner = "admin"
      sharing = "app"
      app = "launcher"
    }
}
`

func TestAccSplunkSavedEventTypes(t *testing.T) {
	resourceName := "splunk_saved_event_types.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccSplunkSavedEventsDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newSavedEventTypes,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "search", "index=main"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test New event description"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "0"),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "color", "blue"),
					resource.TestCheckResourceAttr(resourceName, "tags", "tag"),
				),
			},
			{
				ResourceName:      "splunk_saved_event_types.test",
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
