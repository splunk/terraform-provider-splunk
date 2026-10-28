package splunk

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/splunk/terraform-provider-splunk/client"
	"github.com/splunk/terraform-provider-splunk/client/models"
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
    alert_track = true
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
    alert_track = true
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

const savedSearchesNamespacedImport = `
resource "splunk_saved_searches" "test" {
    name          = "Test Namespaced Import"
    search        = "index=main"
    cron_schedule = "*/5 * * * *"
    acl {
      owner   = "admin"
      sharing = "app"
      app     = "search"
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
    action_email_message_alert = "a non-default message"
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

// savedSearchesEmailIncludeLinksZero sets action_email_include_results_link and
// action_email_include_view_link to 0 to verify they are sent to the API and
// returned on read (and thus appear in savedsearches.conf).
const savedSearchesEmailIncludeLinksZero = `
resource "splunk_saved_searches" "test" {
    name = "Test Email Include Links Zero"
    search = "index=main"
    actions = "email"
    action_email_include_results_link = 0
    action_email_include_view_link = 0
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
    alert_track = true
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
const savedSearchesDispatchIndexWindow = `
resource "splunk_saved_searches" "test" {
    name = "Test Dispatch Index Window"
    search = "index=main"
    dispatch_index_earliest = "-10m"
    dispatch_index_latest = "-5m"
}
`

const savedSearchesDispatchIndexRemoved = `
resource "splunk_saved_searches" "test" {
    name = "Test Dispatch Index Window"
    search = "index=main"
}
`

const savedSearchesDispatchIndexEmpty = `
resource "splunk_saved_searches" "test" {
    name = "Test Dispatch Index Window"
    search = "index=main"
    dispatch_index_earliest = ""
    dispatch_index_latest = ""
}
`

const newSavedSearchesLogEvent = `
resource "splunk_saved_searches" "test" {
	name = "Test Log Event Alert"
	actions = "logevent"
	action_logevent = "1"
	action_logevent_param_event = "test"
	action_logevent_param_host = ""
	action_logevent_param_index = "main"
	action_logevent_param_sourcetype = "stash"
	action_logevent_param_source = "alert"
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

const newSavedSearchesXsoar = `
resource "splunk_saved_searches" "test" {
	name = "Test XSOAR Alert"
	actions = "create_xsoar_incident"
	action_create_xsoar_incident = 1
	action_create_xsoar_incident_param_send_all_servers = 1
	action_create_xsoar_incident_param_server_url = "https://xsoar.example.com"
	action_create_xsoar_incident_param_incident_name = "$name$"
	action_create_xsoar_incident_param_details = "This is a test alert."
	action_create_xsoar_incident_param_custom_fields = "logsource:Demisto,mycustomfield:Test"
	action_create_xsoar_incident_param_severity = 1
	action_create_xsoar_incident_param_occurred = "$trigger_time$"
	action_create_xsoar_incident_param_type = "Unclassified"
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
const newSavedSearchesSlackAppAlertIntegration = `
resource "splunk_saved_searches" "test" {
	name = "Test Slack App Alert"
	actions = "slack_app_alert"
	action_slack_app_alert_integration_param_auto_join_channel = "1"
	action_slack_app_alert_integration_param_bot_username = "SplunkBot"
	action_slack_app_alert_integration_param_channel = "channel"
	action_slack_app_alert_integration_param_emoji = ":splunk:"
	action_slack_app_alert_integration_param_message = "error message"
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

const newSavedSearchesVictorops = `
resource "splunk_saved_searches" "test" {
	name = "Test Victorops Alert"
	actions = "victorops"
	action_victorops_param_message_type = "CRITICAL"
	action_victorops_param_monitoring_tool = "test"
	action_victorops_param_entity_id = "test"
	action_victorops_param_state_message = "error message"
	action_victorops_param_record_id = "12345ab"
	action_victorops_param_routing_key_override = "ops"
	action_victorops_param_enable_recovery = "1"
	action_victorops_param_poll_interval = "5"
	action_victorops_param_inactive_polls = "10"
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

const newSavedSearchesBetterWebhook = `
resource "splunk_saved_searches" "test" {
	name = "Test BetterWebhook Alert"
	actions = "better_webhook"
	action_better_webhook_param_url = "https://example.com/webhook"
	action_better_webhook_param_body_format = "{'alert': '$$app$$'}"
	action_better_webhook_param_credential = "better_webhooks:test:"
	action_better_webhook_param_credentials = "test"
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

const newSavedSearchesPagerduty = `
resource "splunk_saved_searches" "test" {
	name = "Test Pagerduty Alert"
	actions = "pagerduty"
	action_pagerduty_integration_url = "abcd"
	action_pagerduty_integration_url_override = "efgh"
	action_pagerduty_custom_details = "ijkl"
	action_pagerduty_integration_key = "mnop"
	action_pagerduty_integration_key_override = "qrst"
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

const newSavedSearchesJiraServiceDesk = `
resource "splunk_saved_searches" "test" {
	name = "Test Jira Alert Ticket"
	actions = "jira_service_desk"
	action_jira_service_desk_param_account = "test_account"
	action_jira_service_desk_param_jira_project = "test_project"
	action_jira_service_desk_param_jira_issue_type = "Task"
	action_jira_service_desk_param_jira_summary = "error message"
	action_jira_service_desk_param_jira_priority = "Normal"
	action_jira_service_desk_param_jira_description = "test ticket creation"
	action_jira_service_desk_param_jira_dedup = "enabled"
	action_jira_service_desk_param_jira_customfields = "\"customfield_10058\":{\"value\":\"custom_field_value_1\"},\"customfield_10046\":{\"value\":\"custom_field_value_2\"}"
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

const newSavedSearchesWebhook = `
resource "splunk_saved_searches" "test" {
	name = "Test Webhook Alert"
	actions = "webhook"
	action_webhook_param_url = "http://localhost:1234"
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

const newSavedSearchesReport = `
resource "splunk_saved_searches" "test" {
    name = "Test Report"
    search = "index=main"
    actions = "email"
    action_email_include_search = 0
    action_email_include_trigger = 1
    action_email_format = "table"
    action_email_to = "splunk@splunk.com"
    is_scheduled = true
    action_email_message_report = "a non-default message"
    cron_schedule = "*/5 * * * *"
}
`

const newSavedSearchesServiceNowEvent = `
resource "splunk_saved_searches" "test" {
	name                                    = "Test ServiceNow Event"
	actions                                 = "snow_event"
	search                                  = "index=main"
	action_snow_event_param_account         = "account-test"
	action_snow_event_param_node            = "node-test"
	action_snow_event_param_type            = "type-test"
	action_snow_event_param_resource        = "resource-test"
	action_snow_event_param_severity        = 3
	action_snow_event_param_description     = "description-test"
	action_snow_event_param_ci_identifier   = "ci_identifier-test"
	action_snow_event_param_custom_fields   = "custom_fields-test"
	action_snow_event_param_additional_info = "additional_info-test"
	alert_digest_mode                       = true
	alert_expires                           = "24h"
	alert_severity                          = 4
	alert_comparator                        = "greater than"
	alert_threshold                         = "10"
	alert_type                              = "number of events"
	is_scheduled                            = true
	is_visible                              = true
	dispatch_earliest_time                  = "rt-15m"
	dispatch_latest_time                    = "rt-0m"
	cron_schedule                           = "*/15 * * * *"
  }
`

const newSavedSearchesRunAiagent = `
resource "splunk_saved_searches" "test" {
	name                                = "Test Run AI Agent Alert"
	actions                             = "run_aiagent"
	search                              = "index=main level=error"
	action_run_aiagent_param_agent_name = "test_agent"
	action_run_aiagent_param_prompt     = "Investigate the triggering events"
	disabled                            = true
	is_scheduled                        = false
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
					resource.TestCheckResourceAttr(resourceName, "alert_track", "true"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_earliest_time", "rt-15m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_latest_time", "rt-0m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_earliest", "-10m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_index_latest", "-5m"),
					resource.TestCheckResourceAttr(resourceName, "cron_schedule", "*/5 * * * *"),
					resource.TestCheckResourceAttr(resourceName, "is_visible", "true"),
				),
			},
			{
				// Recreate a saved search deleted outside Terraform, then keep the in-place update below.
				// The fixture ACL places the object in admin/launcher; nobody/search does not contain it.
				PreConfig: func() {
					deleteAndWaitForNotFound(t,
						func(c *client.Client) (*http.Response, error) {
							return c.DeleteSavedSearches("Test New Alert", "admin", "launcher")
						},
						func(c *client.Client) url.URL {
							return c.BuildSplunkURLWithEscapedPathPart(nil, "Test New Alert", "servicesNS", "admin", "launcher", "saved", "searches")
						},
					)
				},
				Config: newSavedSearches,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test New Alert"),
					resource.TestCheckResourceAttr(resourceName, "search", "index=main"),
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
					resource.TestCheckResourceAttr(resourceName, "alert_track", "true"),
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
					resource.TestCheckResourceAttr(resourceName, "action_email_message_alert", "a non-default message"),
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
				Config: newSavedSearchesLogEvent,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Log Event Alert"),
					resource.TestCheckResourceAttr(resourceName, "actions", "logevent"),
					resource.TestCheckResourceAttr(resourceName, "action_logevent", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_logevent_param_event", "test"),
					resource.TestCheckResourceAttr(resourceName, "action_logevent_param_host", ""),
					resource.TestCheckResourceAttr(resourceName, "action_logevent_param_index", "main"),
					resource.TestCheckResourceAttr(resourceName, "action_logevent_param_sourcetype", "stash"),
					resource.TestCheckResourceAttr(resourceName, "action_logevent_param_source", "alert"),
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
				Config: newSavedSearchesXsoar,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test XSOAR Alert"),
					resource.TestCheckResourceAttr(resourceName, "actions", "create_xsoar_incident"),
					resource.TestCheckResourceAttr(resourceName, "action_create_xsoar_incident", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_create_xsoar_incident_param_send_all_servers", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_create_xsoar_incident_param_server_url", "https://xsoar.example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_create_xsoar_incident_param_incident_name", "$name$"),
					resource.TestCheckResourceAttr(resourceName, "action_create_xsoar_incident_param_details", "This is a test alert."),
					resource.TestCheckResourceAttr(resourceName, "action_create_xsoar_incident_param_custom_fields", "logsource:Demisto,mycustomfield:Test"),
					resource.TestCheckResourceAttr(resourceName, "action_create_xsoar_incident_param_severity", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_create_xsoar_incident_param_occurred", "$trigger_time$"),
					resource.TestCheckResourceAttr(resourceName, "action_create_xsoar_incident_param_type", "Unclassified"),
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
				Config: newSavedSearchesSlackAppAlertIntegration,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Slack App Alert"),
					resource.TestCheckResourceAttr(resourceName, "actions", "slack_app_alert"),
					resource.TestCheckResourceAttr(resourceName, "action_slack_app_alert_integration_param_auto_join_channel", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_slack_app_alert_integration_param_bot_username", "SplunkBot"),
					resource.TestCheckResourceAttr(resourceName, "action_slack_app_alert_integration_param_channel", "channel"),
					resource.TestCheckResourceAttr(resourceName, "action_slack_app_alert_integration_param_emoji", ":splunk:"),
					resource.TestCheckResourceAttr(resourceName, "action_slack_app_alert_integration_param_message", "error message"),
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
				Config: newSavedSearchesVictorops,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Victorops Alert"),
					resource.TestCheckResourceAttr(resourceName, "actions", "victorops"),
					resource.TestCheckResourceAttr(resourceName, "action_victorops_param_message_type", "CRITICAL"),
					resource.TestCheckResourceAttr(resourceName, "action_victorops_param_monitoring_tool", "test"),
					resource.TestCheckResourceAttr(resourceName, "action_victorops_param_entity_id", "test"),
					resource.TestCheckResourceAttr(resourceName, "action_victorops_param_state_message", "error message"),
					resource.TestCheckResourceAttr(resourceName, "action_victorops_param_record_id", "12345ab"),
					resource.TestCheckResourceAttr(resourceName, "action_victorops_param_routing_key_override", "ops"),
					resource.TestCheckResourceAttr(resourceName, "action_victorops_param_enable_recovery", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_victorops_param_poll_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "action_victorops_param_inactive_polls", "10"),
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
				Config: newSavedSearchesBetterWebhook,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test BetterWebhook Alert"),
					resource.TestCheckResourceAttr(resourceName, "actions", "better_webhook"),
					resource.TestCheckResourceAttr(resourceName, "action_better_webhook_param_url", "https://example.com/webhook"),
					resource.TestCheckResourceAttr(resourceName, "action_better_webhook_param_body_format", "{'alert': '$$app$$'}"),
					resource.TestCheckResourceAttr(resourceName, "action_better_webhook_param_credential", "better_webhooks:test:"),
					resource.TestCheckResourceAttr(resourceName, "action_better_webhook_param_credentials", "test"),
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
				Config: newSavedSearchesPagerduty,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Pagerduty Alert"),
					resource.TestCheckResourceAttr(resourceName, "actions", "pagerduty"),
					resource.TestCheckResourceAttr(resourceName, "action_pagerduty_integration_url", "abcd"),
					resource.TestCheckResourceAttr(resourceName, "action_pagerduty_integration_url_override", "efgh"),
					resource.TestCheckResourceAttr(resourceName, "action_pagerduty_custom_details", "ijkl"),
					resource.TestCheckResourceAttr(resourceName, "action_pagerduty_integration_key", "mnop"),
					resource.TestCheckResourceAttr(resourceName, "action_pagerduty_integration_key_override", "qrst"),
				),
			},
			{
				Config: newSavedSearchesJiraServiceDesk,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Jira Alert Ticket"),
					resource.TestCheckResourceAttr(resourceName, "actions", "jira_service_desk"),
					resource.TestCheckResourceAttr(resourceName, "action_jira_service_desk_param_account", "test_account"),
					resource.TestCheckResourceAttr(resourceName, "action_jira_service_desk_param_jira_project", "test_project"),
					resource.TestCheckResourceAttr(resourceName, "action_jira_service_desk_param_jira_issue_type", "Task"),
					resource.TestCheckResourceAttr(resourceName, "action_jira_service_desk_param_jira_summary", "error message"),
					resource.TestCheckResourceAttr(resourceName, "action_jira_service_desk_param_jira_priority", "Normal"),
					resource.TestCheckResourceAttr(resourceName, "action_jira_service_desk_param_jira_description", "test ticket creation"),
					resource.TestCheckResourceAttr(resourceName, "action_jira_service_desk_param_jira_dedup", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "action_jira_service_desk_param_jira_customfields", "\"customfield_10058\":{\"value\":\"custom_field_value_1\"},\"customfield_10046\":{\"value\":\"custom_field_value_2\"}"),
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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: newSavedSearchesWebhook,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Webhook Alert"),
					resource.TestCheckResourceAttr(resourceName, "actions", "webhook"),
					resource.TestCheckResourceAttr(resourceName, "action_webhook_param_url", "http://localhost:1234"),
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
				Config: newSavedSearchesReport,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Report"),
					resource.TestCheckResourceAttr(resourceName, "search", "index=main"),
					resource.TestCheckResourceAttr(resourceName, "actions", "email"),
					resource.TestCheckResourceAttr(resourceName, "action_email", "true"),
					resource.TestCheckResourceAttr(resourceName, "action_email_include_search", "0"),
					resource.TestCheckResourceAttr(resourceName, "action_email_include_trigger", "1"),
					resource.TestCheckResourceAttr(resourceName, "action_email_message_report", "a non-default message"),
					resource.TestCheckResourceAttr(resourceName, "action_email_format", "table"),
					resource.TestCheckResourceAttr(resourceName, "action_email_to", "splunk@splunk.com"),
					resource.TestCheckResourceAttr(resourceName, "cron_schedule", "*/5 * * * *"),
				),
			},
			{
				Config: newSavedSearchesServiceNowEvent,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test ServiceNow Event"),
					resource.TestCheckResourceAttr(resourceName, "search", "index=main"),
					resource.TestCheckResourceAttr(resourceName, "actions", "snow_event"),
					resource.TestCheckResourceAttr(resourceName, "action_snow_event_param_account", "account-test"),
					resource.TestCheckResourceAttr(resourceName, "action_snow_event_param_node", "node-test"),
					resource.TestCheckResourceAttr(resourceName, "action_snow_event_param_type", "type-test"),
					resource.TestCheckResourceAttr(resourceName, "action_snow_event_param_resource", "resource-test"),
					resource.TestCheckResourceAttr(resourceName, "action_snow_event_param_severity", "3"),
					resource.TestCheckResourceAttr(resourceName, "action_snow_event_param_description", "description-test"),
					resource.TestCheckResourceAttr(resourceName, "action_snow_event_param_ci_identifier", "ci_identifier-test"),
					resource.TestCheckResourceAttr(resourceName, "action_snow_event_param_custom_fields", "custom_fields-test"),
					resource.TestCheckResourceAttr(resourceName, "alert_digest_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "alert_expires", "24h"),
					resource.TestCheckResourceAttr(resourceName, "alert_severity", "4"),
					resource.TestCheckResourceAttr(resourceName, "alert_comparator", "greater than"),
					resource.TestCheckResourceAttr(resourceName, "alert_threshold", "10"),
					resource.TestCheckResourceAttr(resourceName, "alert_type", "number of events"),
					resource.TestCheckResourceAttr(resourceName, "is_scheduled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_visible", "true"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_earliest_time", "rt-15m"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_latest_time", "rt-0m"),
					resource.TestCheckResourceAttr(resourceName, "cron_schedule", "*/15 * * * *"),
				),
			},
			{
				Config: newSavedSearchesRunAiagent,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Run AI Agent Alert"),
					resource.TestCheckResourceAttr(resourceName, "actions", "run_aiagent"),
					resource.TestCheckResourceAttr(resourceName, "action_run_aiagent_param_agent_name", "test_agent"),
					resource.TestCheckResourceAttr(resourceName, "action_run_aiagent_param_prompt", "Investigate the triggering events"),
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

func TestAccSplunkSavedSearchesNamespacedImport(t *testing.T) {
	resourceName := "splunk_saved_searches.test"
	importID := "/servicesNS/admin/search/saved/searches/" + url.PathEscape("Test Namespaced Import")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkSavedSearchesDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: savedSearchesNamespacedImport,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Namespaced Import"),
					resource.TestCheckResourceAttr(resourceName, "search", "index=main"),
					resource.TestCheckResourceAttr(resourceName, "cron_schedule", "*/5 * * * *"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "app"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     importID,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccSplunkSavedSearchesEmailIncludeLinksZero verifies that setting
// action_email_include_results_link and action_email_include_view_link to 0
// sends them to the API and they are returned on read (and appear in savedsearches.conf).
func TestAccSplunkSavedSearchesEmailIncludeLinksZero(t *testing.T) {
	resourceName := "splunk_saved_searches.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkSavedSearchesDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: savedSearchesEmailIncludeLinksZero,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Email Include Links Zero"),
					resource.TestCheckResourceAttr(resourceName, "action_email_include_results_link", "0"),
					resource.TestCheckResourceAttr(resourceName, "action_email_include_view_link", "0"),
				),
			},
		},
	})
}

func TestAccSplunkSavedSearchesDispatchIndexCleared(t *testing.T) {
	resourceName := "splunk_saved_searches.test"
	window := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "dispatch_index_earliest", "-10m"),
		resource.TestCheckResourceAttr(resourceName, "dispatch_index_latest", "-5m"),
	)
	cleared := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "dispatch_index_earliest", ""),
		resource.TestCheckResourceAttr(resourceName, "dispatch_index_latest", ""),
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkSavedSearchesDestroyResources,
		Steps: []resource.TestStep{
			{Config: savedSearchesDispatchIndexWindow, Check: window},
			{Config: savedSearchesDispatchIndexRemoved, Check: cleared},
			{Config: savedSearchesDispatchIndexWindow, Check: window},
			{Config: savedSearchesDispatchIndexEmpty, Check: cleared},
		},
	})
}

func testAccSplunkSavedSearchesDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
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

func testResourceAlertTrackStateDataV0() map[string]interface{} {
	return map[string]interface{}{
		"alert_track": "true",
	}
}

func testResourceAlertTrackStateDataV1() map[string]interface{} {
	v0 := testResourceAlertTrackStateDataV0()
	val, _ := strconv.ParseBool(v0["alert_track"].(string))
	return map[string]interface{}{
		"alert_track": val,
	}
}

func TestResourceExampleInstanceStateUpgradeV0(t *testing.T) {
	expected := testResourceAlertTrackStateDataV1()
	actual, err := resourceAlertTrackStateUpgradeV0(testResourceAlertTrackStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}

func TestGetSavedSearchesConfigSchedulePriority(t *testing.T) {
	enterpriseData := schema.TestResourceDataRaw(t, savedSearches().Schema, map[string]interface{}{
		"search":            "index=main",
		"schedule_priority": "high",
	})
	enterpriseConfig := getSavedSearchesConfig(enterpriseData)
	if got, want := enterpriseConfig.SchedulePriority, "high"; got != want {
		t.Errorf("SchedulePriority with ignore=false: got %q, want %q", got, want)
	}

	cloudData := schema.TestResourceDataRaw(t, savedSearches().Schema, map[string]interface{}{
		"search":                   "index=main",
		"schedule_priority":        "high",
		"ignore_schedule_priority": true,
	})
	cloudConfig := getSavedSearchesConfig(cloudData)
	if got, want := cloudConfig.SchedulePriority, ""; got != want {
		t.Errorf("SchedulePriority with ignore=true: got %q, want empty string", got)
	}
}

func TestSavedSearchesSchemaIgnoreSchedulePriorityDefault(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, savedSearches().Schema, map[string]interface{}{
		"search": "index=main",
	})
	if got, want := resourceData.Get("ignore_schedule_priority").(bool), false; got != want {
		t.Errorf("ignore_schedule_priority default: got %v, want %v", got, want)
	}
}

func TestGetSavedSearchesConfigAllowSkew(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, savedSearches().Schema, map[string]interface{}{
		"search":     "index=main",
		"allow_skew": "100%",
	})
	config := getSavedSearchesConfig(resourceData)
	if got, want := config.AllowSkew, "100%"; got != want {
		t.Errorf("AllowSkew: got %q, want %q", got, want)
	}
}

func TestGetSavedSearchesConfigByNameEmptySeverity(t *testing.T) {
	body := `{
		"entry": [{
			"name": "test_alert",
			"content": {
				"action.snow_event.param.severity": "",
				"search": "index=_internal | head 1"
			}
		}]
	}`
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	entry, err := getSavedSearchesConfigByName("test_alert", resp)
	if err != nil {
		t.Fatalf("getSavedSearchesConfigByName: %v", err)
	}
	if entry == nil {
		t.Fatal("getSavedSearchesConfigByName returned a nil entry")
	}
	if got, want := entry.Content.Search, "index=_internal | head 1"; got != want {
		t.Errorf("Search: got %q, want %q", got, want)
	}
	if got := entry.Content.ActionSnowEventParamSeverity; got != 0 {
		t.Errorf("ActionSnowEventParamSeverity: got %d, want 0", got)
	}
}

func TestGetSavedSearchesConfigByNameDecodeError(t *testing.T) {
	body := `{
		"entry": [{
			"name": "test_alert",
			"content": {
				"action.snow_event.param.severity": "abc",
				"search": "index=_internal | head 1"
			}
		}]
	}`
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	entry, err := getSavedSearchesConfigByName("test_alert", resp)
	if err == nil {
		t.Fatal("getSavedSearchesConfigByName error = nil, want a decode error")
	}
	if entry != nil {
		t.Fatalf("getSavedSearchesConfigByName entry = %#v, want nil", entry)
	}
}

func TestSavedSearchesDeleteReturnsClientError(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	tests := []struct {
		name string
		body string
	}{
		{
			name: "splunk message",
			body: `{"messages":[{"type":"ERROR","text":"Could not find object id=alert"}]}`,
		},
		{
			name: "empty messages",
			body: `{"messages":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = r.Method + " " + r.URL.EscapedPath()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = io.WriteString(w, tt.body)
			}))
			defer server.Close()

			d := savedSearchDeleteResourceData(t)
			err := savedSearchesDelete(d, httpTestProvider(t, server))
			if err == nil {
				t.Fatal("savedSearchesDelete error = nil, want the client status error")
			}
			if !strings.Contains(err.Error(), "400 Bad Request") || !strings.Contains(err.Error(), tt.body) {
				t.Fatalf("savedSearchesDelete error = %q, want status and response body", err)
			}
			const want = "DELETE /servicesNS/nobody/search/saved/searches/alert"
			if got != want {
				t.Fatalf("request = %q, want %q", got, want)
			}
		})
	}
}

func TestSavedSearchesDeleteSucceeds(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()

	if err := savedSearchesDelete(savedSearchDeleteResourceData(t), httpTestProvider(t, server)); err != nil {
		t.Fatal(err)
	}
}

func TestSavedSearchesReadDropsStateWhenMissing(t *testing.T) {
	t.Setenv("HTTPScheme", "http")
	useRefreshTimeout(t, 200*time.Millisecond)

	tests := []struct {
		name   string
		status int
		body   string
	}{
		{name: "not found", status: http.StatusNotFound, body: `{"messages":[{"type":"ERROR","text":"Could not find object"}]}`},
		{name: "empty feed", status: http.StatusOK, body: `{"entry":[]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				_, _ = io.WriteString(w, tt.body)
			}))
			defer server.Close()

			d := savedSearchReadResourceData(t)
			if err := savedSearchesRead(d, httpTestProvider(t, server)); err != nil {
				t.Fatal(err)
			}
			if d.Id() != "" {
				t.Fatalf("id = %q, want empty", d.Id())
			}
		})
	}
}

func TestSavedSearchesReadKeepsStateWhenCreateCannotSeeSearch(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	d := savedSearchReadResourceData(t)
	d.MarkNewResource()
	createDeadlines.Store(d, time.Now().Add(time.Second))
	t.Cleanup(func() { createDeadlines.Delete(d) })
	err := savedSearchesRead(d, httpTestProvider(t, server))
	if err == nil {
		t.Fatal("savedSearchesRead error = nil, want a missing-search error during create")
	}
	if d.Id() != "deleted-search" {
		t.Fatalf("id = %q, want deleted-search", d.Id())
	}
}

func TestSavedSearchesReadKeepsStateOnServerError(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"messages":[{"type":"ERROR","text":"unavailable"}]}`, http.StatusInternalServerError)
	}))
	defer server.Close()

	d := savedSearchReadResourceData(t)
	if err := savedSearchesRead(d, httpTestProvider(t, server)); err == nil {
		t.Fatal("savedSearchesRead error = nil, want the server error")
	}
	if d.Id() != "deleted-search" {
		t.Fatalf("id = %q, want deleted-search", d.Id())
	}
}

func savedSearchReadResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, savedSearches().Schema, map[string]interface{}{
		"name":   "deleted-search",
		"search": "index=main",
	})
	d.SetId("deleted-search")
	return d
}

func savedSearchDeleteResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, savedSearches().Schema, map[string]interface{}{
		"name":   "alert",
		"search": "index=main",
		"acl": []interface{}{map[string]interface{}{
			"owner":   "admin",
			"app":     "search",
			"sharing": "app",
		}},
	})
	d.SetId("alert")
	return d
}

func TestGetSavedSearchesConfigActionEmailCommand(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, savedSearches().Schema, map[string]interface{}{
		"search":               "index=main",
		"action_email_command": "$name$",
	})
	config := getSavedSearchesConfig(resourceData)
	if got, want := config.ActionEmailCommand, "$name$"; got != want {
		t.Errorf("ActionEmailCommand: got %q, want %q", got, want)
	}
}

// builtinSavedSearchActions are alert actions shipped with Splunk Enterprise.
// Any other action.* argument belongs to a modular-alert add-on. Splunk rejects
// it when that add-on is not installed (#256).
var builtinSavedSearchActions = map[string]struct{}{
	"email":           {},
	"logevent":        {},
	"populate_lookup": {},
	"rss":             {},
	"script":          {},
	"summary_index":   {},
	"webhook":         {},
}

// TestVanillaSavedSearchPostsNoAddOnActions is the saved-search create/update
// regression for a scheduled alert that does not use an add-on. It fails when
// any non-built-in action.* argument is present, including a new field added
// with a schema default or without omitempty on its url tag.
func TestVanillaSavedSearchPostsNoAddOnActions(t *testing.T) {
	const search = "index=main EventCode=1116 | table _time, host, user"
	resourceData := schema.TestResourceDataRaw(t, savedSearches().Schema, map[string]interface{}{
		"name":                   "Highly_Relevant_Antivirus_Alert",
		"search":                 search,
		"disabled":               false,
		"is_scheduled":           true,
		"cron_schedule":          "*/5 * * * *",
		"dispatch_earliest_time": "-5m",
		"dispatch_latest_time":   "now",
		"alert_type":             "number of events",
		"alert_comparator":       "greater than",
		"alert_threshold":        "0",
		"alert_track":            true,
		"alert_severity":         4,
	})
	config := getSavedSearchesConfig(resourceData)
	forms := []struct {
		name   string
		values url.Values
	}{
		{"create", savedSearchFormValues(t, config)},
		{"update", savedSearchUpdateFormValues(t, config)},
	}
	for _, form := range forms {
		t.Run(form.name, func(t *testing.T) {
			var posted []string
			for key, vals := range form.values {
				action, ok := savedSearchActionName(key)
				if !ok {
					continue
				}
				if _, builtin := builtinSavedSearchActions[action]; builtin {
					continue
				}
				posted = append(posted, fmt.Sprintf("%s=%q", key, vals))
			}
			if len(posted) > 0 {
				t.Errorf("vanilla saved search posts add-on arguments:\n%s", strings.Join(posted, "\n"))
			}
			if got := form.values.Get("search"); got != search {
				t.Errorf("search = %q, want %q", got, search)
			}
			if got := form.values.Get("cron_schedule"); got != "*/5 * * * *" {
				t.Errorf("cron_schedule = %q", got)
			}
			if got := form.values.Get("alert_threshold"); got != "0" {
				t.Errorf("alert_threshold = %q, want 0", got)
			}
			if !form.values.Has("action.email") {
				t.Error("built-in action.email was omitted")
			}
		})
	}
}

// TestAddOnActionValuesAreEncodedWhenSet checks every non-built-in action
// field, so omitempty cannot drop a value the user configured.
func TestAddOnActionValuesAreEncodedWhenSet(t *testing.T) {
	base := getSavedSearchesConfig(schema.TestResourceDataRaw(t, savedSearches().Schema, map[string]interface{}{
		"search": "index=main",
	}))
	typ := reflect.TypeOf(*base)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		formKey, ok := formKeyFromURLTag(field.Tag.Get("url"))
		if !ok {
			continue
		}
		action, ok := savedSearchActionName(formKey)
		if !ok {
			continue
		}
		if _, builtin := builtinSavedSearchActions[action]; builtin {
			continue
		}
		t.Run(formKey, func(t *testing.T) {
			configured := *base
			fv := reflect.ValueOf(&configured).Elem().Field(i)
			switch fv.Kind() {
			case reflect.String:
				fv.SetString("configured")
			case reflect.Int:
				fv.SetInt(1)
			default:
				t.Fatalf("unsupported add-on field kind %s", fv.Kind())
			}
			for _, values := range []url.Values{savedSearchFormValues(t, &configured), savedSearchUpdateFormValues(t, &configured)} {
				if !values.Has(formKey) || values.Get(formKey) == "" {
					t.Errorf("configured %s was omitted", formKey)
				}
			}
		})
	}
}

func savedSearchActionName(formKey string) (string, bool) {
	rest, ok := strings.CutPrefix(formKey, "action.")
	if !ok || rest == "" {
		return "", false
	}
	name, _, _ := strings.Cut(rest, ".")
	return name, name != ""
}

func formKeyFromURLTag(tag string) (string, bool) {
	if tag == "" || tag == "-" {
		return "", false
	}
	name, _, _ := strings.Cut(tag, ",")
	if name == "" || name == "-" {
		return "", false
	}
	return name, true
}

func savedSearchFormValues(t *testing.T, config *models.SavedSearchObject) url.Values {
	t.Helper()
	values, err := query.Values(config)
	if err != nil {
		t.Fatalf("query.Values: %v", err)
	}
	return values
}

func savedSearchUpdateFormValues(t *testing.T, config *models.SavedSearchObject) url.Values {
	t.Helper()
	// UpdateSavedSearches encodes a pointer to the saved-search pointer.
	values, err := query.Values(&config)
	if err != nil {
		t.Fatalf("query.Values update: %v", err)
	}
	return values
}

func TestSavedSearchesNamespaceForOtherOwners(t *testing.T) {
	tests := []struct {
		sharing  string
		wantPath string
	}{
		{"app", "/servicesNS/nobody/search/saved/searches/alert"},
		{"global", "/servicesNS/nobody/search/saved/searches/alert"},
		{"user", "/servicesNS/jdoe/search/saved/searches/alert"},
	}
	for _, tt := range tests {
		t.Run(tt.sharing, func(t *testing.T) {
			// A positive refresh timeout still retries a 404 once, because the wait is capped
			// by the time remaining. This test checks the namespace path, so do not retry.
			useRefreshTimeout(t, 0)
			var paths []string
			meta := newRoundTripProvider(t, func(r *http.Request) (*http.Response, error) {
				paths = append(paths, r.Method+" "+r.URL.Path)
				return notFoundResponse(), nil
			})
			d := schema.TestResourceDataRaw(t, savedSearches().Schema, map[string]interface{}{
				"name":   "alert",
				"search": "index=main",
				"acl":    []interface{}{map[string]interface{}{"owner": "jdoe", "app": "search", "sharing": tt.sharing}},
			})
			d.SetId("alert")

			if err := savedSearchesRead(d, meta); err != nil {
				t.Fatal(err)
			}
			if d.Id() != "" {
				t.Errorf("id = %q, want the resource removed from state", d.Id())
			}
			d.SetId("alert")
			if err := savedSearchesDelete(d, meta); err != nil {
				t.Fatal(err)
			}

			want := []string{"GET " + tt.wantPath, "DELETE " + tt.wantPath}
			if !reflect.DeepEqual(paths, want) {
				t.Errorf("requests = %q, want %q", paths, want)
			}
		})
	}
}

func TestSavedSearchesCreateWaitsForReplication(t *testing.T) {
	var creates, acls, reads int
	meta := newRoundTripProvider(t, func(r *http.Request) (*http.Response, error) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/servicesNS/admin/search/saved/searches":
			creates++
			if creates == 1 {
				return jsonResponse(http.StatusConflict, `{"messages":[{"type":"ERROR","text":"An object with name=alert already exists"}]}`), nil
			}
			return jsonResponse(http.StatusCreated, "{}"), nil
		case r.Method == http.MethodPost && r.URL.Path == "/servicesNS/admin/search/saved/searches/alert/acl":
			acls++
			if acls == 1 {
				return notFoundResponse(), nil
			}
			return jsonResponse(http.StatusOK, "{}"), nil
		case r.Method == http.MethodGet && r.URL.Path == "/servicesNS/nobody/search/saved/searches/alert":
			reads++
			if reads == 1 {
				return notFoundResponse(), nil
			}
			return jsonResponse(http.StatusOK, `{"entry":[{"name":"alert","acl":{"app":"search","owner":"admin","sharing":"app"},"content":{"search":"index=main"}}]}`), nil
		}
		t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		return nil, nil
	})

	savedSearch := savedSearches()
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":   "alert",
		"search": "index=main",
		"acl":    []interface{}{map[string]interface{}{"owner": "admin", "app": "search", "sharing": "app"}},
	})
	diff, err := savedSearch.Diff(nil, config, meta)
	if err != nil {
		t.Fatal(err)
	}
	state, err := savedSearch.Apply(nil, diff, meta)
	if err != nil {
		t.Fatal(err)
	}
	if state.ID != "alert" || state.Attributes["search"] != "index=main" {
		t.Errorf("state = %v", state)
	}
	if creates != 2 || acls != 2 || reads != 2 {
		t.Errorf("creates = %d, acls = %d, reads = %d, want 2 each", creates, acls, reads)
	}
}

func TestSavedSearchesCreateACLUsesOwner(t *testing.T) {
	const entry = `{"entry":[{"name":"alert","acl":{"app":"search","owner":"jdoe","sharing":"app"},"content":{"search":"index=main"}}]}`
	tests := []struct {
		sharing  string
		wantACL  string
		wantRead string
	}{
		{"app", "/servicesNS/jdoe/search/saved/searches/alert/acl", "/servicesNS/nobody/search/saved/searches/alert"},
		{"global", "/servicesNS/jdoe/search/saved/searches/alert/acl", "/servicesNS/nobody/search/saved/searches/alert"},
		{"user", "/servicesNS/jdoe/search/saved/searches/alert/acl", "/servicesNS/jdoe/search/saved/searches/alert"},
	}
	for _, tt := range tests {
		t.Run(tt.sharing, func(t *testing.T) {
			var aclPath, readPath string
			meta := newRoundTripProvider(t, func(r *http.Request) (*http.Response, error) {
				switch r.Method {
				case http.MethodPost:
					if strings.HasSuffix(r.URL.Path, "/acl") {
						aclPath = r.URL.Path
					}
					return jsonResponse(http.StatusCreated, "{}"), nil
				case http.MethodGet:
					readPath = r.URL.Path
					return jsonResponse(http.StatusOK, entry), nil
				}
				t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
				return nil, nil
			})

			savedSearch := savedSearches()
			config := terraform.NewResourceConfigRaw(map[string]interface{}{
				"name":   "alert",
				"search": "index=main",
				"acl":    []interface{}{map[string]interface{}{"owner": "jdoe", "app": "search", "sharing": tt.sharing}},
			})
			diff, err := savedSearch.Diff(nil, config, meta)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := savedSearch.Apply(nil, diff, meta); err != nil {
				t.Fatal(err)
			}
			if aclPath != tt.wantACL {
				t.Errorf("acl path = %q, want %q", aclPath, tt.wantACL)
			}
			if readPath != tt.wantRead {
				t.Errorf("read path = %q, want %q", readPath, tt.wantRead)
			}
		})
	}
}

func TestSavedSearchesUpdateSendsACLOnlyWhenChanged(t *testing.T) {
	search, write := "index=main", "admin"
	var aclRequests int
	var aclPath string
	meta := newRoundTripProvider(t, func(r *http.Request) (*http.Response, error) {
		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			form, err := url.ParseQuery(string(body))
			if err != nil {
				t.Fatal(err)
			}
			if strings.HasSuffix(r.URL.Path, "/acl") {
				aclRequests++
				aclPath = r.URL.Path
				write = form.Get("perms.write")
			} else if form.Has("search") {
				search = form.Get("search")
			}
			return jsonResponse(http.StatusOK, "{}"), nil
		}
		return jsonResponse(http.StatusOK, fmt.Sprintf(`{"entry":[{"name":"alert","acl":{"app":"search","owner":"jdoe","sharing":"app","perms":{"read":["*"],"write":[%q]}},"content":{"search":%q}}]}`, write, search)), nil
	})

	savedSearch := savedSearches()
	apply := func(state *terraform.InstanceState, search, write string) *terraform.InstanceState {
		t.Helper()
		config := terraform.NewResourceConfigRaw(map[string]interface{}{
			"name":   "alert",
			"search": search,
			"acl":    []interface{}{map[string]interface{}{"owner": "jdoe", "app": "search", "sharing": "app", "read": []interface{}{"*"}, "write": []interface{}{write}}},
		})
		diff, err := savedSearch.Diff(state, config, meta)
		if err != nil {
			t.Fatal(err)
		}
		state, err = savedSearch.Apply(state, diff, meta)
		if err != nil {
			t.Fatal(err)
		}
		return state
	}

	state := apply(nil, "index=main", "admin")
	aclRequests = 0

	state = apply(state, "index=main | head 1", "admin")
	if aclRequests != 0 {
		t.Errorf("acl requests = %d after a search change, want 0", aclRequests)
	}

	apply(state, "index=main | head 1", "power")
	if aclRequests != 1 {
		t.Errorf("acl requests = %d after an acl change, want 1", aclRequests)
	}
	if aclPath != "/servicesNS/nobody/search/saved/searches/alert/acl" {
		t.Errorf("acl path = %q, want the nobody namespace", aclPath)
	}
}

// TestSavedSearchesUpdateClearsDispatchIndexWindow runs against a server that,
// like Splunk, keeps any key an update leaves out.
func TestSavedSearchesUpdateClearsDispatchIndexWindow(t *testing.T) {
	keys := []string{"dispatch.index_earliest", "dispatch.index_latest"}
	tests := []struct {
		name    string
		cleared map[string]interface{}
	}{
		{"removed", map[string]interface{}{}},
		{"empty", map[string]interface{}{"dispatch_index_earliest": "", "dispatch_index_latest": ""}},
		{"null", map[string]interface{}{"dispatch_index_earliest": nil, "dispatch_index_latest": nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stored := map[string]string{}
			var updates []url.Values
			meta := newRoundTripProvider(t, func(r *http.Request) (*http.Response, error) {
				if r.Method == http.MethodPost {
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Fatal(err)
					}
					form, err := url.ParseQuery(string(body))
					if err != nil {
						t.Fatal(err)
					}
					if strings.HasSuffix(r.URL.Path, "/saved/searches/alert") {
						updates = append(updates, form)
					}
					for _, key := range append([]string{"search"}, keys...) {
						if form.Has(key) {
							stored[key] = form.Get(key)
						}
					}
					return jsonResponse(http.StatusOK, "{}"), nil
				}
				return jsonResponse(http.StatusOK, fmt.Sprintf(`{"entry":[{"name":"alert","acl":{"app":"search","owner":"admin","sharing":"app"},"content":{"search":%q,"dispatch.index_earliest":%q,"dispatch.index_latest":%q}}]}`,
					stored["search"], stored["dispatch.index_earliest"], stored["dispatch.index_latest"])), nil
			})

			savedSearch := savedSearches()
			config := func(search string, attrs map[string]interface{}) *terraform.ResourceConfig {
				raw := map[string]interface{}{"name": "alert", "search": search}
				for k, v := range attrs {
					raw[k] = v
				}
				return terraform.NewResourceConfigRaw(raw)
			}
			apply := func(state *terraform.InstanceState, config *terraform.ResourceConfig) *terraform.InstanceState {
				t.Helper()
				diff, err := savedSearch.Diff(state, config, meta)
				if err != nil {
					t.Fatal(err)
				}
				if diff.Empty() {
					return state
				}
				state, err = savedSearch.Apply(state, diff, meta)
				if err != nil {
					t.Fatal(err)
				}
				return state
			}

			state := apply(nil, config("index=main", map[string]interface{}{"dispatch_index_earliest": "-10m", "dispatch_index_latest": "-5m"}))
			if stored[keys[0]] != "-10m" || stored[keys[1]] != "-5m" {
				t.Fatalf("index window not created: %v", stored)
			}

			updates = nil
			state = apply(state, config("index=main", tt.cleared))
			for _, key := range keys {
				if got := stored[key]; got != "" {
					t.Errorf("%s = %q after the update, want it cleared", key, got)
				}
			}
			if len(updates) != 1 {
				t.Fatalf("updates = %d, want 1", len(updates))
			}
			for key, vals := range updates[0] {
				if vals[0] == "" && key != keys[0] && key != keys[1] {
					t.Errorf("update posted an empty %s", key)
				}
			}

			diff, err := savedSearch.Diff(state, config("index=main", tt.cleared), meta)
			if err != nil {
				t.Fatal(err)
			}
			if !diff.Empty() {
				for _, attr := range []string{"dispatch_index_earliest", "dispatch_index_latest"} {
					if _, ok := diff.Attributes[attr]; ok {
						t.Errorf("plan after the update still changes %s", attr)
					}
				}
			}

			updates = nil
			state = apply(state, config("index=main | head 1", tt.cleared))
			if len(updates) != 1 {
				t.Fatalf("updates = %d, want 1", len(updates))
			}
			for _, key := range keys {
				if updates[0].Has(key) {
					t.Errorf("a later search change posted %s", key)
				}
			}

			updates = nil
			apply(state, config("index=main | head 1", map[string]interface{}{"dispatch_index_earliest": "-10m", "dispatch_index_latest": "-5m"}))
			if len(updates) != 1 {
				t.Fatalf("updates = %d, want 1", len(updates))
			}
			want := map[string]string{keys[0]: "-10m", keys[1]: "-5m"}
			for _, key := range keys {
				if got := updates[0][key]; len(got) != 1 || got[0] != want[key] {
					t.Errorf("%s = %q when set again, want [%s]", key, got, want[key])
				}
				if stored[key] != want[key] {
					t.Errorf("%s stored %q when set again, want %q", key, stored[key], want[key])
				}
			}
		})
	}
}

func TestSavedSearchesReadRetriesNotFound(t *testing.T) {
	reads := 0
	meta := newRoundTripProvider(t, func(r *http.Request) (*http.Response, error) {
		reads++
		if reads == 1 {
			return notFoundResponse(), nil
		}
		return jsonResponse(http.StatusOK, `{"entry":[{"name":"alert","acl":{"app":"search","owner":"jdoe","sharing":"app"},"content":{"search":"index=main"}}]}`), nil
	})
	d := schema.TestResourceDataRaw(t, savedSearches().Schema, map[string]interface{}{
		"name":   "alert",
		"search": "index=main",
		"acl":    []interface{}{map[string]interface{}{"owner": "jdoe", "app": "search", "sharing": "app"}},
	})
	d.SetId("alert")

	if err := savedSearchesRead(d, meta); err != nil {
		t.Fatal(err)
	}
	if d.Id() != "alert" {
		t.Fatalf("id = %q, want the saved search kept after a retried 404", d.Id())
	}
	if reads < 2 {
		t.Fatalf("reads = %d, want a retry after the first 404", reads)
	}
}

func TestSavedSearchesDeleteRetriesNotFound(t *testing.T) {
	deletes := 0
	meta := newRoundTripProvider(t, func(r *http.Request) (*http.Response, error) {
		deletes++
		if r.URL.Path != "/servicesNS/nobody/search/saved/searches/alert" {
			t.Errorf("path = %s, want the nobody namespace", r.URL.Path)
		}
		if deletes == 1 {
			return notFoundResponse(), nil
		}
		return jsonResponse(http.StatusOK, "{}"), nil
	})
	d := schema.TestResourceDataRaw(t, savedSearches().Schema, map[string]interface{}{
		"name":   "alert",
		"search": "index=main",
		"acl":    []interface{}{map[string]interface{}{"owner": "jdoe", "app": "search", "sharing": "app"}},
	})
	d.SetId("alert")

	if err := savedSearchesDelete(d, meta); err != nil {
		t.Fatal(err)
	}
	if deletes < 2 {
		t.Fatalf("deletes = %d, want a retry after the first 404", deletes)
	}
}

func TestSavedSearchesCreateACLClientErrorDoesNotRetry(t *testing.T) {
	for _, code := range []int{http.StatusBadRequest, http.StatusForbidden} {
		t.Run(http.StatusText(code), func(t *testing.T) {
			var acls int
			meta := newRoundTripProvider(t, func(r *http.Request) (*http.Response, error) {
				switch {
				case r.Method == http.MethodPost && r.URL.Path == "/servicesNS/admin/search/saved/searches":
					return jsonResponse(http.StatusCreated, "{}"), nil
				case r.Method == http.MethodPost && r.URL.Path == "/servicesNS/admin/search/saved/searches/alert/acl":
					acls++
					return jsonResponse(code, `{"messages":[{"type":"ERROR","text":"denied"}]}`), nil
				}
				t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
				return nil, nil
			})

			savedSearch := savedSearches()
			config := terraform.NewResourceConfigRaw(map[string]interface{}{
				"name":   "alert",
				"search": "index=main",
				"acl":    []interface{}{map[string]interface{}{"owner": "admin", "app": "search", "sharing": "app"}},
			})
			diff, err := savedSearch.Diff(nil, config, meta)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := savedSearch.Apply(nil, diff, meta); err == nil {
				t.Fatal("expected the ACL error to fail the create")
			}
			if acls != 1 {
				t.Fatalf("acl requests = %d, want 1", acls)
			}
		})
	}
}
