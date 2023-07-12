# Example Splunk Terraform Configuration
# Set correct Splunk instance credentials under the provider splunk resource block
# This basic example.tf creates a Splunk user and role, setup global HEC configuration,
# creates a new HEC token with indexer acknowledgement enabled
# creates a saved search to search for events received with above token as source
terraform {
  required_providers {
    splunk = {
      source  = "splunk/splunk"
      version = "1.4.7"
    }
  }
}

provider "splunk" {
  url                  = "localhost:8089"
  username             = "admin"
  password             = "changeme"
  insecure_skip_verify = true
  // Or use environment variables used:
  // SPLUNK_USERNAME
  // SPLUNK_PASSWORD
  // SPLUNK_URL
  // SPLUNK_INSECURE_SKIP_VERIFY (Defaults to true)
}

resource "splunk_admin_saml_groups" "saml-group01" {
  name  = "terraform"
  roles = ["admin", "power"]
}

resource "splunk_authorization_roles" "role01" {
  name                   = "terraform-user01-role"
  default_app            = "search"
  imported_roles         = ["power", "user"]
  capabilities           = ["accelerate_datamodel", "change_authentication", "restart_splunkd"]
  search_indexes_allowed = ["_audit", "_internal", "main"]
  search_indexes_default = ["_audit", "_internal", "main"]
}

resource "splunk_authentication_users" "user01" {
  name              = "user01"
  email             = "user01@example.com"
  password          = "password01"
  force_change_pass = false
  roles             = ["terraform-user01-role"]
  depends_on = [
    splunk_authorization_roles.role01
  ]
}

resource "splunk_indexes" "user01-index" {
  name                   = "user01-index"
  max_hot_buckets        = 6
  max_total_data_size_mb = 1000000
}

resource "splunk_global_http_event_collector" "http" {
  disabled   = false
  enable_ssl = true
  port       = 8088
}

resource "random_uuid" "hec_token" {
}

resource "splunk_inputs_http_event_collector" "hec-token-01" {
  name       = "hec-token-01"
  token      = random_uuid.hec_token.result
  index      = splunk_indexes.user01-index.name
  indexes    = [splunk_indexes.user01-index.name, "history", "summary"]
  source     = "new:source"
  sourcetype = "new:sourcetype"
  disabled   = false
  use_ack    = 0
  acl {
    owner   = "user01"
    sharing = "global"
    read    = ["admin"]
    write   = ["admin"]
  }
  depends_on = [
    splunk_indexes.user01-index,
    splunk_authentication_users.user01,
    splunk_global_http_event_collector.http,
  ]
}

resource "splunk_configs_conf" "terraform-stanza" {
  name = "props-example/terraform"
  variables = {
    "disabled" : "false"
    "custom_key" : "value"
  }

  acl {
    app = "search"
  }

  depends_on = [
    splunk_authentication_users.user01,
  ]
}

resource "splunk_data_ui_views" "dashboard" {
  name     = "Terraform_Sample_Dashboard"
  eai_data = "<dashboard version=\"1.1\"><label>Terraform</label><description>Terraform operations</description><row><panel><chart><search><query>index=_internal sourcetype=splunkd_access useragent=\"splunk-simple-go-client\" | timechart fixedrange=f values(status) by uri_path</query><earliest>-24h@h</earliest><latest>now</latest><sampleRatio>1</sampleRatio></search><option name=\"charting.axisLabelsX.majorLabelStyle.overflowMode\">ellipsisNone</option><option name=\"charting.axisLabelsX.majorLabelStyle.rotation\">0</option><option name=\"charting.axisTitleX.visibility\">collapsed</option><option name=\"charting.axisTitleY.text\">HTTP status codes</option><option name=\"charting.axisTitleY.visibility\">visible</option><option name=\"charting.axisTitleY2.visibility\">visible</option><option name=\"charting.axisX.abbreviation\">none</option><option name=\"charting.axisX.scale\">linear</option><option name=\"charting.axisY.abbreviation\">none</option><option name=\"charting.axisY.scale\">linear</option><option name=\"charting.axisY2.abbreviation\">none</option><option name=\"charting.axisY2.enabled\">0</option><option name=\"charting.axisY2.scale\">inherit</option><option name=\"charting.chart\">column</option><option name=\"charting.chart.bubbleMaximumSize\">50</option><option name=\"charting.chart.bubbleMinimumSize\">10</option><option name=\"charting.chart.bubbleSizeBy\">area</option><option name=\"charting.chart.nullValueMode\">connect</option><option name=\"charting.chart.showDataLabels\">none</option><option name=\"charting.chart.sliceCollapsingThreshold\">0.01</option><option name=\"charting.chart.stackMode\">default</option><option name=\"charting.chart.style\">shiny</option><option name=\"charting.drilldown\">none</option><option name=\"charting.layout.splitSeries\">0</option><option name=\"charting.layout.splitSeries.allowIndependentYRanges\">0</option><option name=\"charting.legend.labelStyle.overflowMode\">ellipsisMiddle</option><option name=\"charting.legend.mode\">standard</option><option name=\"charting.legend.placement\">right</option><option name=\"charting.lineWidth\">2</option><option name=\"trellis.enabled\">0</option><option name=\"trellis.scales.shared\">1</option><option name=\"trellis.size\">small</option><option name=\"trellis.splitBy\">_aggregation</option></chart></panel></row></dashboard>"

  acl {
    owner = "admin"
    app   = "search"
  }
}

resource "splunk_saved_searches" "new-search-01" {
  actions                   = "email"
  action_email_format       = "table"
  action_email_max_time     = "5m"
  action_email_send_results = true
  action_email_subject      = "Splunk Alert: $name$"
  action_email_to           = "user01@splunk.com"
  action_email_track_alert  = true
  alert_comparator          = "greater than"
  alert_digest_mode         = true
  alert_expires             = "30d"
  alert_threshold           = "0"
  alert_type                = "number of events"
  description               = "source=http:hec-token-01 is receiving events"
  dispatch_earliest_time    = "rt-15m"
  dispatch_latest_time      = "rt-0m"
  cron_schedule             = "*/15 * * * *"
  name                      = "new-search-01"
  schedule_priority         = "default"
  search                    = "index=${splunk_indexes.user01-index.name} source=${splunk_inputs_http_event_collector.hec-token-01.name}"

  acl {
    app     = "search"
    owner   = "user01"
    sharing = "user"
  }
  depends_on = [
    splunk_authentication_users.user01,
    splunk_indexes.user01-index
  ]
}
