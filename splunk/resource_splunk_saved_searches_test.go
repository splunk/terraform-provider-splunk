package splunk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const newSavedSearches = `
resource "splunk_saved_searches" "test" {
    name = "Test New Alert"
    search = "index=main"
    actions = "email"
    action_email_include_search = 0
    action_email_include_trigger = 1
    action_email_format = "table"
    action_email_max_time = "5m"
    action_email_max_results = 10
    action_email_send_csv = 1
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
    action_email_include_search = 1
    action_email_include_trigger = 1
    action_email_format = "table"
    action_email_max_time = "5m"
    action_email_max_results = 100
    action_email_send_csv = 1
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

const newSavedSearchesBracket = `
resource "splunk_saved_searches" "test" {
    name = "[Test New Alert]"
    search = "index=main"
    actions = "email"
    action_email_include_search = 0
    action_email_include_trigger = 1
    action_email_format = "table"
    action_email_max_time = "5m"
    action_email_max_results = 10
    action_email_send_csv = 1
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

const updateSavedSearchesBracket = `
resource "splunk_saved_searches" "test" {
    name = "[Test New Alert]"
    search = "index=main"
    actions = "email"
    action_email_include_search = 1
    action_email_include_trigger = 1
    action_email_format = "table"
    action_email_max_time = "5m"
    action_email_max_results = 100
    action_email_send_csv = 1
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

const newSavedSearchesSlack = `
resource "splunk_saved_searches" "test" {
	name = "Test Slack Alert"
	actions = "slack"
	action_slack_param_attachment = "alert_link"
	action_slack_param_channel = "#channel"
	action_slack_param_message = "error message"
	alert_comparator    = "greater than"
	alert_digest_mode   = true
	alert_expires       = "30d"
	alert_threshold     = "0"
	alert_type          = "number of events"
	cron_schedule       = "*/1 * * * *"
	disabled            = false
	is_scheduled        = true
	is_visible          = true
	realtime_schedule   = true
	search              = "index=main level=error"
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
					resource.TestCheckResourceAttr(resourceName, "action_email_include_search", "0"),
					resource.TestCheckResourceAttr(resourceName, "action_email_include_trigger", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_email_format", "table"),
					resource.TestCheckResourceAttr(resourceName, "action_email_max_time", "5m"),
					resource.TestCheckResourceAttr(resourceName, "action_email_max_results", "10"),
					resource.TestCheckResourceAttr(resourceName, "action_email_send_csv", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_email_send_results", "false"),
					resource.TestCheckResourceAttr(resourceName, "action_email_subject", "Splunk Alert: $name$"),
					resource.TestCheckResourceAttr(resourceName, "action_email_to", "splunk@splunk.com"),
					resource.TestCheckResourceAttr(resourceName, "action_email_track_alert", "true"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_earliest_time", "rt-15m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_latest_time", "rt-0m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_earliest", "-10m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_latest", "-5m"),
					resource.TestCheckResourceAttr(resourceName, "cron_schedule", "*/5 * * * *"),
					resource.TestCheckResourceAttr(resourceName, "is_visible", "true"),
				),
			},
			{
				Config: updatedSavedSearches,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test New Alert"),
					resource.TestCheckResourceAttr(resourceName, "search", "index=main"),
					resource.TestCheckResourceAttr(resourceName, "actions", "email"),
					resource.TestCheckResourceAttr(resourceName, "action_email", "true"),
					resource.TestCheckResourceAttr(resourceName, "action_email_include_search", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_email_include_trigger", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_email_format", "table"),
					resource.TestCheckResourceAttr(resourceName, "action_email_max_time", "5m"),
					resource.TestCheckResourceAttr(resourceName, "action_email_max_results", "100"),
					resource.TestCheckResourceAttr(resourceName, "action_email_send_csv", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_email_send_results", "false"),
					resource.TestCheckResourceAttr(resourceName, "action_email_subject", "Splunk Alert: $name$"),
					resource.TestCheckResourceAttr(resourceName, "action_email_to", "splunk@splunk.com"),
					resource.TestCheckResourceAttr(resourceName, "action_email_track_alert", "true"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_earliest_time", "rt-15m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_latest_time", "rt-0m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_earliest", "-20m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_latest", "-5m"),
					resource.TestCheckResourceAttr(resourceName, "cron_schedule", "*/15 * * * *"),
					resource.TestCheckResourceAttr(resourceName, "is_visible", "true"),
				),
			},
			{
				Config: newSavedSearchesBracket,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "[Test New Alert]"),
					resource.TestCheckResourceAttr(resourceName, "search", "index=main"),
					resource.TestCheckResourceAttr(resourceName, "actions", "email"),
					resource.TestCheckResourceAttr(resourceName, "action_email", "true"),
					resource.TestCheckResourceAttr(resourceName, "action_email_include_search", "0"),
					resource.TestCheckResourceAttr(resourceName, "action_email_include_trigger", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_email_format", "table"),
					resource.TestCheckResourceAttr(resourceName, "action_email_max_time", "5m"),
					resource.TestCheckResourceAttr(resourceName, "action_email_max_results", "10"),
					resource.TestCheckResourceAttr(resourceName, "action_email_send_csv", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_email_send_results", "false"),
					resource.TestCheckResourceAttr(resourceName, "action_email_subject", "Splunk Alert: $name$"),
					resource.TestCheckResourceAttr(resourceName, "action_email_to", "splunk@splunk.com"),
					resource.TestCheckResourceAttr(resourceName, "action_email_track_alert", "true"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_earliest_time", "rt-15m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_latest_time", "rt-0m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_earliest", "-10m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_latest", "-5m"),
					resource.TestCheckResourceAttr(resourceName, "cron_schedule", "*/5 * * * *"),
					resource.TestCheckResourceAttr(resourceName, "is_visible", "true"),
				),
			},
			{
				Config: updateSavedSearchesBracket,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "[Test New Alert]"),
					resource.TestCheckResourceAttr(resourceName, "search", "index=main"),
					resource.TestCheckResourceAttr(resourceName, "actions", "email"),
					resource.TestCheckResourceAttr(resourceName, "action_email", "true"),
					resource.TestCheckResourceAttr(resourceName, "action_email_include_search", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_email_include_trigger", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_email_format", "table"),
					resource.TestCheckResourceAttr(resourceName, "action_email_max_time", "5m"),
					resource.TestCheckResourceAttr(resourceName, "action_email_max_results", "100"),
					resource.TestCheckResourceAttr(resourceName, "action_email_send_csv", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_email_send_results", "false"),
					resource.TestCheckResourceAttr(resourceName, "action_email_subject", "Splunk Alert: $name$"),
					resource.TestCheckResourceAttr(resourceName, "action_email_to", "splunk@splunk.com"),
					resource.TestCheckResourceAttr(resourceName, "action_email_track_alert", "true"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_earliest_time", "rt-15m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_latest_time", "rt-0m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_earliest", "-20m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_latest", "-5m"),
					resource.TestCheckResourceAttr(resourceName, "cron_schedule", "*/15 * * * *"),
					resource.TestCheckResourceAttr(resourceName, "is_visible", "true"),
				),
			},
			{
				Config: newSavedSearchesSlack,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Slack Alert"),
					resource.TestCheckResourceAttr(resourceName, "actions", "slack"),
					resource.TestCheckResourceAttr(resourceName, "action_slack_param_attachment", "alert_link"),
					resource.TestCheckResourceAttr(resourceName, "action_slack_param_channel", "#channel"),
					resource.TestCheckResourceAttr(resourceName, "action_slack_param_message", "error message"),
					resource.TestCheckResourceAttr(resourceName, "alert_comparator", "greater than"),
					resource.TestCheckResourceAttr(resourceName, "alert_digest_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "alert_expires", "30d"),
					resource.TestCheckResourceAttr(resourceName, "alert_threshold", "0"),
					resource.TestCheckResourceAttr(resourceName, "alert_type", "number of events"),
					resource.TestCheckResourceAttr(resourceName, "cron_schedule", "*/1 * * * *"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_scheduled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_visible", "true"),
					resource.TestCheckResourceAttr(resourceName, "realtime_schedule", "true"),
					resource.TestCheckResourceAttr(resourceName, "search", "index=main level=error"),
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
