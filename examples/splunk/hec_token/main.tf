# Example: HEC token creation (no indexer acknowledgement).
# Uses SPLUNK_URL, SPLUNK_USERNAME, SPLUNK_PASSWORD. Run: terraform init && terraform apply

terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = ">= 3.0"
    }
    splunk = {
      source  = "splunk/splunk"
      version = ">= 1.4.0"
    }
  }
}

provider "splunk" {}

resource "splunk_indexes" "test_index_1" {
  name                   = "test_index_1"
  max_hot_buckets        = 6
  max_total_data_size_mb = 1000000
}

resource "splunk_global_http_event_collector" "http" {
  disabled   = false
  enable_ssl = true
  port       = 8088
}

resource "random_uuid" "no_ack" {}

resource "splunk_inputs_http_event_collector" "no_ack" {
  name     = "some-name"
  token    = random_uuid.no_ack.result
  index    = "test_index_1"
  indexes  = ["test_index_1"]
  disabled = false
  use_ack  = 0
  depends_on = [
    splunk_indexes.test_index_1,
    splunk_global_http_event_collector.http,
  ]
}
