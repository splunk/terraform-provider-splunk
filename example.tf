# Example Splunk Terraform Configuration
# Set correct Splunk instance credentials under the provider splunk resource block
# This basic example.tf creates a Splunk user and role, setup global HEC configuration,
# creates a new HEC token with indexer acknowledgement enabled
# creates a saved search to search for events received with above token as source

provider "splunk" {
  url                  = "localhost:8089"
  username             = "admin"
  password             = "password"
  insecure_skip_verify = true
}

resource "splunk_authorization_roles" "role01" {
  name           = "terraform-user01-role"
  default_app    = "search"
  imported_roles = ["power", "user"]
  capabilities   = ["accelerate_datamodel", "change_authentication", "restart_splunkd"]
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

resource "splunk_inputs_http_event_collector" "hec-token-01" {
  name       = "hec-token-01"
  index      = "user01-index"
  indexes    = ["user01-index", "history", "summary"]
  source     = "new:source"
  sourcetype = "new:sourcetype"
  disabled   = false
  use_ack    = false
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

resource "splunk_saved_searches" "new-search-01" {
  actions                   = "email"
  action_email_format       = "table"
  action_email_max_time     = "5m"
  action_email_send_results = false
  action_email_subject      = "Splunk Alert: $name$"
  action_email_to           = "user01@splunk.com"
  action_email_track_alert  = true
  description               = "New search for user01"
  dispatch_earliest_time    = "rt-15m"
  dispatch_latest_time      = "rt-0m"
  cron_schedule             = "*/15 * * * *"
  name                      = "new-search-01"
  search                    = "index=user01-index source=http:hec-token-01"

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

resource "splunk_configs_conf" "new-conf-stanza" {
  name = "internaltf/custom"
  variables = {
    "custom_key" : "value"
  }
  acl {
    app = "search"
  }
}

