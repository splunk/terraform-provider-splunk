# Resource: splunk_saved_searches

Create and manage saved searches.

## Example Usage

```
resource "splunk_saved_searches" "saved_search" {
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
    cron_schedule = "*/5 * * * *"
    acl {
      owner = "admin"
      sharing = "app"
      app = "launcher"
    }
}
```

## Argument Reference

For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTsearch#saved.2Fsearches

This resource block supports the following arguments:

- `name` - (Required) A name for the search.
- `search` - (Required) Required when creating a new search.
- `action_email` - (Optional) The state of the email action. Read-only attribute. Value ignored on POST. Use actions to specify a list of enabled actions. Defaults to 0.
- `action_email_auth_password` - (Optional) The password to use when authenticating with the SMTP server. Normally this value is set when editing the email settings, however you can set a clear text password here and it is encrypted on the next platform restart.Defaults to empty string.
- `action_email_auth_username` - (Optional) The username to use when authenticating with the SMTP server. If this is empty string, no authentication is attempted. Defaults to empty stringNOTE: Your SMTP server might reject unauthenticated emails.
- `action_email_bcc` - (Optional) BCC email address to use if action.email is enabled.
- `action_email_cc` - (Optional) CC email address to use if action.email is enabled.
- `action_email_command` - (Optional) The search command (or pipeline) which is responsible for executing the action.Generally the command is a template search pipeline which is realized with values from the saved search. To reference saved search field values wrap them in $, for example to reference the savedsearch name use $name$, to reference the search use $search$.
- `action_email_format` - (Optional) Valid values: (table | plain | html | raw | csv)Specify the format of text in the email. This value also applies to any attachments.
- `action_email_from` - (Optional) Email address from which the email action originates.Defaults to splunk@$LOCALHOST or whatever value is set in alert_actions.conf.
- `action_email_hostname` - (Optional) Sets the hostname used in the web link (url) sent in email actions.This value accepts two forms:hostname (for example, splunkserver, splunkserver.example.com)
- `action_email_include_results_link` - (Optional) Specify whether to include a link to the results. Defaults to 0.
- `action_email_include_search` - (Optional) Specify whether to include the search that caused an email to be sent. Defaults to 0.
- `action_email_include_trigger` - (Optional) Specify whether to show the trigger condition that caused the alert to fire. Defaults to 0.
- `action_email_include_trigger_time` - (Optional) Specify whether to show the time that the alert was fired. Defaults to 0.
- `action_email_include_view_link` - (Optional) Specify whether to show the title and a link to enable the user to edit the saved search. Defaults to 0.
- `action_email_inline` - (Optional) Indicates whether the search results are contained in the body of the email.Results can be either inline or attached to an email.
- `action_email_mailserver` - (Optional) Set the address of the MTA server to be used to send the emails.Defaults to <LOCALHOST> or whatever is set in alert_actions.conf.
- `action_email_max_results` - (Optional) Sets the global maximum number of search results to send when email.action is enabled. Defaults to 100.
- `action_email_max_time` - (Optional) Valid values are Integer[m|s|h|d].Specifies the maximum amount of time the execution of an email action takes before the action is aborted. Defaults to 5m.
- `action_email_message_alert` - (Optional) Customize the message sent in the emailed alert. Defaults to: The alert condition for '$name$' was triggered.
- `action_email_message_report` - (Optional) Customize the message sent in the emailed report. Defaults to: The scheduled report '$name$' has run
- `action_email_pdfview` - (Optional) The name of the view to deliver if sendpdf is enabled
- `action_email_preprocess_results` - (Optional) Search string to preprocess results before emailing them. Defaults to empty string (no preprocessing).Usually the preprocessing consists of filtering out unwanted internal fields.
- `action_email_report_cid_font_list` - (Optional) Space-separated list. Specifies the set (and load order) of CID fonts for handling Simplified Chinese(gb), Traditional Chinese(cns), Japanese(jp), and Korean(kor) in Integrated PDF Rendering.If multiple fonts provide a glyph for a given character code, the glyph from the first font specified in the list is used.To skip loading any CID fonts, specify the empty string.Defaults to 'gb cns jp kor'
- `action_email_report_include_splunk_logo` - (Optional) Indicates whether to include the Splunk logo with the report.
- `action_email_report_paper_orientation` - (Optional) Valid values: (portrait | landscape)Specifies the paper orientation: portrait or landscape. Defaults to portrait.
- `action_email_report_paper_size` - (Optional) Valid values: (letter | legal | ledger | a2 | a3 | a4 | a5)Specifies the paper size for PDFs. Defaults to letter.
- `action_email_report_server_enabled` - (Optional) No Supported
- `action_email_report_server_url` - (Optional) Not supported.For a default locally installed report server, the URL is http://localhost:8091/
- `action_email_send_csv` - (Optional) Specify whether to send results as a CSV file. Defaults to 0.
- `action_email_send_pdf` - (Optional) Indicates whether to create and send the results as a PDF. Defaults to false.
- `action_email_send_results` - (Optional) Indicates whether to attach the search results in the email.Results can be either attached or inline. See action.email.inline.
- `action_email_subject` - (Optional) Specifies an alternate email subject.Defaults to SplunkAlert-<savedsearchname>.
- `action_email_to` - (Optional) A comma or semicolon separated list of recipient email addresses. Required if this search is scheduled and the email alert action is enabled.
- `action_email_track_alert` - (Optional) Indicates whether the execution of this action signifies a trackable alert.
- `action_email_ttl` - (Optional) Valid values are Integer[p].Specifies the minimum time-to-live in seconds of the search artifacts if this action is triggered. If p follows <Integer>, int is the number of scheduled periods. Defaults to 86400 (24 hours).If no actions are triggered, the artifacts have their ttl determined by dispatch.ttl in savedsearches.conf.
- `action_email_use_ssl` - (Optional) Indicates whether to use SSL when communicating with the SMTP server. Defaults to false.
- `action_email_use_tls` - (Optional) Indicates whether to use TLS (transport layer security) when communicating with the SMTP server (starttls).Defaults to false.
- `action_email_width_sort_columns` - (Optional) Indicates whether columns should be sorted from least wide to most wide, left to right.Only valid if format=text.
- `action_pagerduty_custom_details` - (Optional) The PagerDuty custom details information.
- `action_pagerduty_integration_key` - (Optional) The PagerDuty integration Key.
- `action_pagerduty_integration_key_override` - (Optional) The PagerDuty integration Key override.
- `action_pagerduty_integration_url` - (Optional) The pagerduty integration URL. This integration uses Splunk's native webhooks to send events to PagerDuty.
- `action_pagerduty_integration_url_override` - (Optional) The pagerduty integration URL override. This integration uses Splunk's native webhooks to send events to PagerDuty.
- `action_populate_lookup` - (Optional) The state of the populate lookup action. Read-only attribute. Value ignored on POST. Use actions to specify a list of enabled actions. Defaults to 0.
- `action_populate_lookup_command` - (Optional) The search command (or pipeline) which is responsible for executing the action.
- `action_populate_lookup_dest` - (Optional) Lookup name of path of the lookup to populate
- `action_populate_lookup_hostname` - (Optional) Sets the hostname used in the web link (url) sent in alert actions.This value accepts two forms: hostname (for example, splunkserver, splunkserver.example.com)\n\nprotocol://hostname:port (for example, http://splunkserver:8000, https://splunkserver.example.com:443)
- `action_populate_lookup_max_results` - (Optional) Sets the maximum number of search results sent using alerts. Defaults to 100.
- `action_populate_lookup_max_time` - (Optional) Valid values are: Integer[m|s|h|d]Sets the maximum amount of time the execution of an action takes before the action is aborted. Defaults to 5m.
- `action_populate_lookup_track_alert` - (Optional) Indicates whether the execution of this action signifies a trackable alert.
- `action_populate_lookup_ttl` - (Optional) Valid values are Integer[p]Specifies the minimum time-to-live in seconds of the search artifacts if this action is triggered. If p follows Integer, then this specifies the number of scheduled periods. Defaults to 10p.
- `action_rss` - (Optional) The state of the rss action. Read-only attribute. Value ignored on POST.Use actions to specify a list of enabled actions. Defaults to 0.
- `action_rss_command` - (Optional) The search command (or pipeline) which is responsible for executing the action.Generally the command is a template search pipeline which is realized with values from the saved search. To reference saved search field values wrap them in $, for example to reference the savedsearch name use $name$, to reference the search use $search$.
- `action_rss_hostname` - (Optional) Sets the hostname used in the web link (url) sent in alert actions.This value accepts two forms:hostname (for example, splunkserver, splunkserver.example.com)\n\nprotocol://hostname:port (for example, http://splunkserver:8000, https://splunkserver.example.com:443)
- `action_rss_max_results` - (Optional) Sets the maximum number of search results sent using alerts. Defaults to 100.
- `action_rss_max_time` - (Optional) Valid values are Integer[m|s|h|d].Sets the maximum amount of time the execution of an action takes before the action is aborted. Defaults to 1m.
- `action_rss_track_alert` - (Optional) Indicates whether the execution of this action signifies a trackable alert.
- `action_rss_ttl` - (Optional) Valid values are: Integer[p] Specifies the minimum time-to-live in seconds of the search artifacts if this action is triggered. If p follows Integer, specifies the number of scheduled periods. Defaults to 86400 (24 hours).
- `action_script` - (Optional) The state of the script action. Read-only attribute. Value ignored on POST. Use actions to specify a list of enabled actions. Defaults to 0.
- `action_script_command` - (Optional) The search command (or pipeline) which is responsible for executing the action.Generally the command is a template search pipeline which is realized with values from the saved search. To reference saved search field values wrap them in $, for example to reference the savedsearch name use $name$, to reference the search use $search$.
- `action_script_filename` - (Optional) File name of the script to call. Required if script action is enabled
- `action_script_hostname` - (Optional) Sets the hostname used in the web link (url) sent in alert actions.This value accepts two forms:hostname (for example, splunkserver, splunkserver.example.com)\n\nprotocol://hostname:port (for example, http://splunkserver:8000, https://splunkserver.example.com:443)
- `action_script_max_results` - (Optional) Sets the maximum number of search results sent using alerts. Defaults to 100.
- `action_script_max_time` - (Optional) Valid values are Integer[m|s|h|d].Sets the maximum amount of time the execution of an action takes before the action is aborted. Defaults to 1m.
- `action_script_track_alert` - (Optional) Indicates whether the execution of this action signifies a trackable alert.
- `action_script_ttl` - (Optional) Valid values are: Integer[p] Specifies the minimum time-to-live in seconds of the search artifacts if this action is triggered. If p follows Integer, specifies the number of scheduled periods. Defaults to 86400 (24 hours).
- `action_summary_index` - (Optional) The state of the summary index action. Read-only attribute. Value ignored on POST. Use actions to specify a list of enabled actions. Defaults to 0.
- `action_summary_index_command` - (Optional) The search command (or pipeline) which is responsible for executing the action.Generally the command is a template search pipeline which is realized with values from the saved search. To reference saved search field values wrap them in $, for example to reference the savedsearch name use $name$, to reference the search use $search$.
- `action_summary_index_hostname` - (Optional) Sets the hostname used in the web link (url) sent in summary-index alert actions.This value accepts two forms:hostname (for example, splunkserver, splunkserver.example.com)protocol://hostname:port (for example, http://splunkserver:8000, https://splunkserver.example.com:443)
- `action_summary_index_inline` - (Optional) Determines whether to execute the summary indexing action as part of the scheduled search.NOTE: This option is considered only if the summary index action is enabled and is always executed (in other words, if counttype = always).Defaults to true
- `action_summary_index_max_results` - (Optional) Sets the maximum number of search results sent using alerts. Defaults to 100.
- `action_summary_index_max_time` - (Optional) Valid values are Integer[m|s|h|d].Sets the maximum amount of time the execution of an action takes before the action is aborted. Defaults to 1m.
- `action_summary_index_name` - (Optional) Specifies the name of the summary index where the results of the scheduled search are saved.Defaults to summary.
- `action_summary_index_track_alert` - (Optional) Indicates whether the execution of this action signifies a trackable alert.
- `action_summary_index_ttl` - (Optional) Valid values are: Integer[p] Specifies the minimum time-to-live in seconds of the search artifacts if this action is triggered. If p follows Integer, specifies the number of scheduled periods. Defaults to 86400 (24 hours).
- `action_logevent` - (Optional) Enable log event action (Should be 1 (Enabled) or 0 (Disabled))
- `action_logevent_param_event` - (Optional) Specific event text for the logged event.
- `action_logevent_param_host` - (Optional) Value for the host field of the logged event.
- `action_logevent_param_index` - (Optional) Destination index where to store the logged event.
- `action_logevent_param_sourcetype` - (Optional) Destination sourcetype where to store the logged event.
- `action_logevent_param_source` - (Optional) Value for the source field of the logged event.
- `action_create_xsoar_incident` - (Optional) Enable XSOAR alerting (Should by 1 (Enabled) or 0 (Disabled))
- `action_create_xsoar_incident_param_send_all_servers` - (Optional) Enable XSOAR alerting sending to all servers (Should by 1 (Enabled) or 0 (Disabled)
- `action_create_xsoar_incident_param_server_url` - (Optional) XSOAR Server instance URL (Should start with https:// || http://)
- `action_create_xsoar_incident_param_incident_name` - (Optional) XSOAR incident name
- `action_create_xsoar_incident_param_details` - (Optional) XSOAR incident description
- `action_create_xsoar_incident_param_custom_fields` - (Optional) XSOAR custom incident fields (should be a comma separated list)
- `action_create_xsoar_incident_param_severity` - (Optional) XSOAR Severity (1 - Low, 2 - Medium, 3 - High, 4 - Critical)
- `action_create_xsoar_incident_param_occurred` - (Optional) XSOAR incident time
- `action_create_xsoar_incident_param_type` - (Optional) XSOAR incident type
- `action_slack_param_channel` - (Optional) Slack channel to send the message to (Should start with # or @)
- `action_slack_param_fields` - (Optional) Show one or more fields from the search results below your Slack message. Comma-separated list of field names. Allows wildcards. eg. index,source\*
- `action_slack_param_attachment` - (Optional) Include a message attachment. Valid values are message, none, or alert_link
- `action_slack_param_message` - (Optional) Enter the chat message to send to the Slack channel. The message can include tokens that insert text based on the results of the search.
- `action_slack_param_webhook_url_override` - (Optional) You can override the Slack webhook URL here if you need to send the alert message to a different Slack team
- `action_slack_app_alert_param_auto_join_channel` - (Optional) Automatically join the channel if the bot is not already a member (Should be 1 (Enabled) or 0 (Disabled))
- `action_slack_app_alert_param_bot_username` - (Optional) The bot username that will post the message
- `action_slack_app_alert_param_channel` - (Optional) Slack channel to send the message to (Should start with # or @)
- `action_slack_app_alert_param_emoji` - (Optional) Emoji icon to use as the bot's avatar (Should start and end with :)
- `action_slack_app_alert_param_message` - (Optional) Enter the chat message to send to the Slack channel. The message can include tokens that insert text based on the results of the
- `action_slack_app_alert_param_fields`
- `action_jira_service_desk_param_account` - (Optional) Jira Service Desk account name
- `action_jira_service_desk_param_jira_project` - (Optional) Jira Project name
- `action_jira_service_desk_param_jira_issue_type` - (Optional) Jira issue type name
- `action_jira_service_desk_param_jira_summary` - (Optional) Jira issue title/summary
- `action_jira_service_desk_param_jira_priority` - (Optional) Jira priority of issue
- `action_jira_service_desk_param_jira_description` - (Optional) Jira issue description
- `action_jira_service_desk_param_jira_customfields` - (Optional) Jira custom fields data (see https://ta-jira-service-desk-simple-addon.readthedocs.io/en/latest/userguide.html)

* `action_victorops_param_message_type` - (Optional) Type of VictorOps message. Valid values are info, warning, critical, recovery, ack
* `action_victorops_param_monitoring_tool` - (Optional) Name of the monitoring tool sending the alert
* `action_victorops_param_entity_id` - (Optional) Unique identifier for the affected system or service
* `action_victorops_param_state_message` - (Optional) Description of the alert condition
* `action_victorops_param_record_id` - (Optional) Identifier used to correlate related alerts
* `action_victorops_param_routing_key_override` - (Optional) You can override the VictorOps routing key here if you need to send the alert message to a different VictorOps team
* `action_victorops_param_enable_recovery` - (Optional) Enable sending of recovery messages (Should be 1 (Enabled) or 0 (Disabled))
* `action_victorops_param_poll_interval` - (Optional) Polling interval for checking the status of the alert (in minutes)
* `action_victorops_param_inactive_polls` - (Optional) Number of inactive polls before sending a recovery message

- `action_better_webhook_param_url` - (Optional) URL to send the HTTP POST request to. Must be accessible from the Splunk server
- `action_better_webhook_param_body_format` - (Optional) Format of the body content. Valid values are json, xml, form-urlencoded, or raw
- `action_better_webhook_param_credential` - (Optional) Name of the Splunk stored credential to use for authentication
- `action_better_webhook_param_credentials` - (Optional) Use the credentials defined in the webhook URL
- `action_webhook_param_url` - (Optional) URL to send the HTTP POST request to. Must be accessible from the Splunk server
- `actions` - (Optional) A comma-separated list of actions to enable. For example: rss,email
- `alert_comparator` - (Optional) One of the following strings: greater than, less than, equal to, rises by, drops by, rises by perc, drops by percUsed with alert_threshold to trigger alert actions.
- `alert_condition` - (Optional) Contains a conditional search that is evaluated against the results of the saved search. Defaults to an empty string.
- `alert_digest_mode` - (Optional) Specifies whether alert actions are applied to the entire result set or on each individual result.Defaults to 1 (true).
- `alert_expires` - (Optional) Valid values: [number][time-unit]Sets the period of time to show the alert in the dashboard. Defaults to 24h.Use [number][time-unit] to specify a time. For example: 60 = 60 seconds, 1m = 1 minute, 1h = 60 minutes = 1 hour.
- `alert_severity` - (Optional) Valid values: (1 | 2 | 3 | 4 | 5 | 6) Sets the alert severity level.Valid values are:1 DEBUG 2 INFO 3 WARN 4 ERROR 5 SEVERE 6 FATAL Defaults to 3.
- `alert_suppress` - (Optional) Indicates whether alert suppression is enabled for this scheduled search.
- `alert_suppress_fields` - (Optional) Comma delimited list of fields to use for suppression when doing per result alerting. Required if suppression is turned on and per result alerting is enabled.
- `alert_suppress_period` - (Optional) Valid values: [number][time-unit] Specifies the suppresion period. Only valid if alert.supress is enabled.Use [number][time-unit] to specify a time. For example: 60 = 60 seconds, 1m = 1 minute, 1h = 60 minutes = 1 hour.
- `alert_threshold` - (Optional) Valid values are: Integer[%]Specifies the value to compare (see alert_comparator) before triggering the alert actions. If expressed as a percentage, indicates value to use when alert_comparator is set to rises by perc or drops by perc.
- `alert_track` - (Optional) Valid values: (true | false | auto) Specifies whether to track the actions triggered by this scheduled search.auto - determine whether to track or not based on the tracking setting of each action, do not track scheduled searches that always trigger actions. Default value true - force alert tracking.false - disable alert tracking for this search.
- `alert_type` - (Optional) What to base the alert on, overriden by alert_condition if it is specified. Valid values are: always, custom, number of events, number of hosts, number of sources.
- `allow_skew` - (Optional) Allows the search scheduler to distribute scheduled searches randomly and more evenly over their specified search periods.
- `auto_summarize` - (Optional) Indicates whether the scheduler should ensure that the data for this search is automatically summarized. Defaults to 0.
- `auto_summarize_command` - (Optional) An auto summarization template for this search. See auto summarization options in savedsearches.conf for more details.
- `auto_summarize_cron_schedule` - (Optional) Cron schedule that probes and generates the summaries for this saved search.The default value is _/10 _ \* \* \* and corresponds to \`every ten hours\`.
- `auto_summarize_dispatch_earliest_time` - (Optional) A time string that specifies the earliest time for summarizing this search. Can be a relative or absolute time.If this value is an absolute time, use the dispatch.time_format to format the value.
- `auto_summarize_dispatch_latest_time` - (Optional) A time string that specifies the latest time for summarizing this saved search. Can be a relative or absolute time.If this value is an absolute time, use the dispatch.time_format to format the value.
- `auto_summarize_dispatch_time_format` - (Optional) Defines the time format that Splunk software uses to specify the earliest and latest time. Defaults to %FT%T.%Q%:z
- `auto_summarize_dispatch_ttl` - (Optional) Valid values: Integer[p]. Defaults to 60.Indicates the time to live (in seconds) for the artifacts of the summarization of the scheduled search.
- `auto_summarize_max_disabled_buckets` - (Optional) The maximum number of buckets with the suspended summarization before the summarization search is completely stopped, and the summarization of the search is suspended for auto_summarize.suspend_period. Defaults to 2.
- `auto_summarize_max_summary_ratio` - (Optional) The maximum ratio of summary_size/bucket_size, which specifies when to stop summarization and deem it unhelpful for a bucket. Defaults to 0.1 Note: The test is only performed if the summary size is larger than auto_summarize.max_summary_size.
- `auto_summarize_max_summary_size` - (Optional) The minimum summary size, in bytes, before testing whether the summarization is helpful.The default value is 52428800 and is equivalent to 5MB.
- `auto_summarize_max_time` - (Optional) Maximum time (in seconds) that the summary search is allowed to run. Defaults to 3600.Note: This is an approximate time. The summary search stops at clean bucket boundaries.
- `auto_summarize_suspend_period` - (Optional) Time specfier indicating when to suspend summarization of this search if the summarization is deemed unhelpful.Defaults to 24h.
- `auto_summarize_timespan` - (Optional) The list of time ranges that each summarized chunk should span. This comprises the list of available granularity levels for which summaries would be available. Specify a comma delimited list of time specifiers.For example a timechart over the last month whose granuality is at the day level should set this to 1d. If you need the same data summarized at the hour level for weekly charts, use: 1h,1d.
- `cron_schedule` - (Optional) Valid values: cron stringThe cron schedule to execute this search. For example: _/5 _ \* \* \* causes the search to execute every 5 minutes.
- `description` - (Optional) Human-readable description of this saved search. Defaults to empty string.
- `disabled` - (Optional) Indicates if the saved search is enabled. Defaults to 0.Disabled saved searches are not visible in Splunk Web.
- `dispatch_buckets` - (Optional) The maximum number of timeline buckets. Defaults to 0.
- `dispatch_earliest_time` - (Optional) A time string that specifies the earliest time for this search. Can be a relative or absolute time. If this value is an absolute time, use the dispatch.time_format to format the value.
- `dispatch_index_earliest` - (Optional) A time string that specifies the earliest index time for this search. Can be a relative or absolute time. If this value is an absolute time, use the dispatch.time_format to format the value.
- `dispatch_index_latest` - (Optional) A time string that specifies the latest index time for this search. Can be a relative or absolute time. If this value is an absolute time, use the dispatch.time_format to format the value.
- `dispatch_indexed_realtime` - (Optional) A time string that specifies the earliest time for this search. Can be a relative or absolute time. If this value is an absolute time, use the dispatch.time_format to format the value.
- `dispatch_indexed_realtime_minspan` - (Optional) Allows for a per-job override of the [search] indexed_realtime_disk_sync_delay setting in limits.conf.
- `dispatch_indexed_realtime_offset` - (Optional) Allows for a per-job override of the [search] indexed_realtime_disk_sync_delay setting in limits.conf.
- `dispatch_latest_time` - (Optional) A time string that specifies the latest time for this saved search. Can be a relative or absolute time.If this value is an absolute time, use the dispatch.time_format to format the value.
- `dispatch_lookups` - (Optional) Enables or disables the lookups for this search. Defaults to 1.
- `dispatch_max_count` - (Optional) The maximum number of results before finalizing the search. Defaults to 500000.
- `dispatch_max_time` - (Optional) Indicates the maximum amount of time (in seconds) before finalizing the search. Defaults to 0.
- `dispatch_reduce_freq` - (Optional) Specifies, in seconds, how frequently the MapReduce reduce phase runs on accumulated map values. Defaults to 10.
- `dispatch_rt_backfill` - (Optional) Whether to back fill the real time window for this search. Parameter valid only if this is a real time search. Defaults to 0.
- `dispatch_rt_maximum_span` - (Optional) Allows for a per-job override of the [search] indexed_realtime_maximum_span setting in limits.conf.
- `dispatch_spawn_process` - (Optional) Specifies whether a new search process spawns when this saved search is executed. Defaults to 1. Searches against indexes must run in a separate process.
- `dispatch_time_format` - (Optional) A time format string that defines the time format for specifying the earliest and latest time. Defaults to %FT%T.%Q%:z
- `dispatch_ttl` - (Optional) Valid values: Integer[p]. Defaults to 2p.Indicates the time to live (in seconds) for the artifacts of the scheduled search, if no actions are triggered.
- `display_view` - (Optional) Defines the default UI view name (not label) in which to load the results. Accessibility is subject to the user having sufficient permissions.
- `is_scheduled` - (Optional) Whether this search is to be run on a schedule
- `is_visible` - (Optional) Specifies whether this saved search should be listed in the visible saved search list. Defaults to 1.
- `max_concurrent` - (Optional) The maximum number of concurrent instances of this search the scheduler is allowed to run. Defaults to 1.
- `realtime_schedule` - (Optional) Defaults to 1. Controls the way the scheduler computes the next execution time of a scheduled search. If this value is set to 1, the scheduler bases its determination of the next scheduled search execution time on the current time. If this value is set to 0, the scheduler bases its determination of the next scheduled search on the last search execution time. This is called continuous scheduling. If set to 0, the scheduler never skips scheduled execution periods. However, the execution of the saved search might fall behind depending on the scheduler load. Use continuous scheduling whenever you enable the summary index option.
- `request_ui_dispatch_app` - (Optional) Specifies a field used by Splunk Web to denote the app this search should be dispatched in.
- `request_ui_dispatch_view` - (Optional) Specifies a field used by Splunk Web to denote the view this search should be displayed in.
- `restart_on_searchpeer_add` - (Optional) Specifies whether to restart a real-time search managed by the scheduler when a search peer becomes available for this saved search. Defaults to 1.
- `run_on_startup` - (Optional) Indicates whether this search runs at startup. If it does not run on startup, it runs at the next scheduled time. Defaults to 0. Set to 1 for scheduled searches that populate lookup tables.
- `schedule_priority` - (Optional) Raises the scheduling priority of the named search. Defaults to Default
- `schedule_window` - (Optional) Time window (in minutes) during which the search has lower priority. Defaults to 0. The scheduler can give higher priority to more critical searches during this window. The window must be smaller than the search period.Set to auto to let the scheduler determine the optimal window value automatically. Requires the edit_search_schedule_window capability to override auto.
- `vsid` - (Optional) Defines the viewstate id associated with the UI view listed in 'displayview'.
- `workload_pool` - (Optional) Specifies the new workload pool where the existing running search will be placed.`
- `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference

In addition to all arguments above, This resource block exports the following arguments:

- `id` - The ID of the saved search resource
