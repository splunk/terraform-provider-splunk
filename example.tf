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

resource "splunk_global_http_event_collector" "http" {
  disabled   = false
  enable_ssl = true
  port       = 8088
}

resource "splunk_inputs_http_event_collector" "hec-token-01" {
  name       = "hec-token-01"
  index      = "main"
  indexes    = ["main", "history", "summary"]
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
    splunk_authentication_users.user01,
    splunk_global_http_event_collector.http
  ]
}

resource "splunk_saved_searches" "new-search-01" {
  name   = "new-search-01"
  search = "index=main source=http:hec-token-01"
  acl {
    app     = "search"
    owner   = "user01"
    sharing = "global"
  }
  depends_on = [
    splunk_authentication_users.user01,
  ]
}

resource "splunk_index" "foo" {
  name = "foo"
  max_hot_buckets = 4
}

resource "splunk_index" "bar" {
  name = "summary"
  max_hot_buckets = 4
}
