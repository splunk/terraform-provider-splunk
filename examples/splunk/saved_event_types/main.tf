# Example: saved event type (knowledge object).
# Uses SPLUNK_URL, SPLUNK_USERNAME, SPLUNK_PASSWORD. Run: terraform init && terraform apply
# See: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTknowledge#saved.2Feventtypes

terraform {
  required_providers {
    splunk = {
      source  = "splunk/splunk"
      version = ">= 1.4.0"
    }
  }
}

provider "splunk" {}

resource "splunk_saved_event_types" "example" {
  name        = "terraform-example-event-type"
  search      = "index=main sourcetype=access_combined"
  description = "Example event type created by Terraform"
  disabled    = false
  priority    = 1
  color       = "et_green"
  acl {
    owner   = "admin"
    sharing = "app"
    app     = "search"
  }
}
