package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const newSavedSearches = `
resource "splunk_saved_searches" "test" {
    name = "test_search"
    search = "index=main"
    acl {
      owner = "admin"
      sharing = "app"
      app = "launcher"
    }
}
`

const updatedSavedSearches = `
resource "splunk_saved_searches" "test" {
    name = "test_search"
    search = "index=main foo=bar"
    acl {
      owner = "admin"
      sharing = "app"
      app = "launcher"
    }
}
`

func TestAccSplunkSavedSearches(t *testing.T) {
	resourceName := "splunk_saved_searches.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkSavedSearchesDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newSavedSearches,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "search", "index=main"),
				),
			},
			{
				Config: updatedSavedSearches,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "search", "index=main foo=bar"),
				),
			},
			{
				ResourceName:      "splunk_saved_searches.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkSavedSearchesDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_saved_searches":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "search", "saved", "searches", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
