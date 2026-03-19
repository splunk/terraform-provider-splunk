# One-off config to create a saved search with all action_email_include_* params
# for manual verification in Splunk UI. Uses SPLUNK_URL, SPLUNK_USERNAME,
# SPLUNK_PASSWORD (and SPLUNK_HOME if needed). Run: terraform init && terraform apply
# Then in Splunk: Settings → Saved searches → "Test Email Include Links Zero".
# Clean up with: terraform destroy

terraform {
  required_providers {
    splunk = {
      source  = "splunk/splunk"
      version = ">= 1.4.0"
    }
  }
}

provider "splunk" {}

resource "splunk_saved_searches" "test" {
  name                             = "Test Email Include Links Zero"
  search                           = "index=main"
  actions                          = "email"
  action_email_include_results_link = 0
  action_email_include_view_link    = 0
  action_email_include_search       = 0
  action_email_include_trigger      = 1
  action_email_include_trigger_time = 1
  action_email_format               = "table"
  action_email_max_time             = "5m"
  action_email_max_results          = 10
  action_email_send_csv             = 1
  action_email_send_results         = false
  action_email_subject              = "Splunk Alert: $name$"
  action_email_to                   = "splunk@splunk.com"
  action_email_track_alert          = true
  alert_track                      = true
  dispatch_earliest_time            = "rt-15m"
  dispatch_latest_time              = "rt-0m"
  dispatch_index_earliest           = "-10m"
  dispatch_index_latest             = "-5m"
  cron_schedule                    = "*/5 * * * *"
  acl {
    owner   = "admin"
    sharing = "app"
    app     = "launcher"
  }
}
