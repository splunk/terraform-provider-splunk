# Example Splunk Terraform Configuration
# Set correct Splunk instance credentials under the provider splunk resource block
# This basic example.tf creates a Splunk user, setup global HEC configuration,
# creates a new HEC token with indexer acknowledgement enabled
# creates a saved search to search for events received with above token as source

provider "splunk" {
  url                  = "localhost:8089"
  username             = "admin"
  password             = "changeme"
  insecure_skip_verify = true
}

resource "splunk_authentication_user" "user01" {
  name              = "user01"
  email             = "user01@example.com"
  password          = "password01"
  force_change_pass = false
  roles             = ["power", "user"]
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
    splunk_authentication_user.user01,
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
}
