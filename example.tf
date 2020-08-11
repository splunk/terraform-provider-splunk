# Example Splunk Terraform Configuration
# Set correct Splunk instance credentials under the provider splunk resource block
# This basic example.tf creates a HEC token with indexer acknowledgement enabled

provider "splunk" {
  url                  = "localhost:8089"
  username             = "admin"
  password             = "changeme"
  insecure_skip_verify = true
}

resource "splunk_global_http_event_collector" "http" {
  disabled    = false
  enable_ssl  = true
  port        = 8088
}

resource "splunk_inputs_http_event_collector" "hec" {
  name       = "new-token"
  index      = "main"
  indexes    = ["main", "history"]
  source     = "new-source"
  sourcetype = "new-sourcetype"
  disabled   = false
  use_ack    = 1

  acl {
    sharing = "global"
    read = ["admin"]
    write = ["admin"]
  }

  depends_on = ["splunk_global_http_event_collector.http"]
}
