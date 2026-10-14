terraform {
  required_providers {
    splunk = {
      source = "splunk/splunk"
    }
  }
}

provider "splunk" {
  # The password is read from SPLUNK_PASSWORD. Do not commit it.
  url      = "127.0.0.1:8089"
  username = "admin"
}

resource "splunk_apps_local" "lookup_file_editing" {
  filename         = true
  name             = "/workspace/terraform/apps/splunk-app-for-lookup-file-editing_404.tgz"
  explicit_appname = "splunk-app-for-lookup-file-editing"
}