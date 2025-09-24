package splunk

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
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
				ResourceName:      "splunk_saved_searches.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
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
