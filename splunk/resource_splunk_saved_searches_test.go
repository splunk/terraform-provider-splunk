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
    name = "Test New Alert"
    search = "index=main"
    actions = "email"
    action_email_format = "table"
    action_email_max_time = "5m"
    action_email_max_results = 10
    action_email_send_results = false
    action_email_subject = "Splunk Alert: $name$"
    action_email_to = "splunk@splunk.com"
    action_email_track_alert = true
    dispatch_earliest_time = "rt-15m"
    dispatch_latest_time = "rt-0m"
    dispatch_index_earliest = "-10m"
    dispatch_index_latest = "-5m"
    cron_schedule = "*/5 * * * *"
    acl {
      owner = "admin"
      sharing = "app"
      app = "launcher"
    }
}
`

const updatedSavedSearches = `
resource "splunk_saved_searches" "test" {
    name = "Test New Alert"
    search = "index=main"
    actions = "email"
    action_email_format = "table"
    action_email_max_time = "5m"
    action_email_max_results = 100
    action_email_send_results = false
    action_email_subject = "Splunk Alert: $name$"
    action_email_to = "splunk@splunk.com"
    action_email_track_alert = true
    dispatch_earliest_time = "rt-15m"
    dispatch_latest_time = "rt-0m"
    dispatch_index_earliest = "-20m"
    dispatch_index_latest = "-5m"
    cron_schedule = "*/15 * * * *"
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
					resource.TestCheckResourceAttr(resourceName, "name", "Test New Alert"),
					resource.TestCheckResourceAttr(resourceName, "search", "index=main"),
					resource.TestCheckResourceAttr(resourceName, "actions", "email"),
					resource.TestCheckResourceAttr(resourceName, "action_email", "true"),
					resource.TestCheckResourceAttr(resourceName, "action_email_format", "table"),
					resource.TestCheckResourceAttr(resourceName, "action_email_max_time", "5m"),
					resource.TestCheckResourceAttr(resourceName, "action_email_max_results", "10"),
					resource.TestCheckResourceAttr(resourceName, "action_email_send_results", "false"),
					resource.TestCheckResourceAttr(resourceName, "action_email_subject", "Splunk Alert: $name$"),
					resource.TestCheckResourceAttr(resourceName, "action_email_to", "splunk@splunk.com"),
					resource.TestCheckResourceAttr(resourceName, "action_email_track_alert", "true"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_earliest_time", "rt-15m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_latest_time", "rt-0m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_earliest", "-10m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_latest", "-5m"),
					resource.TestCheckResourceAttr(resourceName, "cron_schedule", "*/5 * * * *"),
				),
			},
			{
				Config: updatedSavedSearches,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test New Alert"),
					resource.TestCheckResourceAttr(resourceName, "search", "index=main"),
					resource.TestCheckResourceAttr(resourceName, "actions", "email"),
					resource.TestCheckResourceAttr(resourceName, "action_email", "true"),
					resource.TestCheckResourceAttr(resourceName, "action_email_format", "table"),
					resource.TestCheckResourceAttr(resourceName, "action_email_max_time", "5m"),
					resource.TestCheckResourceAttr(resourceName, "action_email_max_results", "100"),
					resource.TestCheckResourceAttr(resourceName, "action_email_send_results", "false"),
					resource.TestCheckResourceAttr(resourceName, "action_email_subject", "Splunk Alert: $name$"),
					resource.TestCheckResourceAttr(resourceName, "action_email_to", "splunk@splunk.com"),
					resource.TestCheckResourceAttr(resourceName, "action_email_track_alert", "true"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_earliest_time", "rt-15m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_latest_time", "rt-0m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_earliest", "-20m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_latest", "-5m"),
					resource.TestCheckResourceAttr(resourceName, "cron_schedule", "*/15 * * * *"),
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
