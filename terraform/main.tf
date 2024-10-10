terraform {
  required_providers {
    splunk = {
      source = "splunk/splunk"
    }
  }
}

provider "splunk" {
  url      = "localhost:8089"
  username = "admin"
  password = "password"
}

resource "splunk_apps_local" "lookup_file_editing" {
  filename         = true
  name             = "/apps/splunk-app-for-lookup-file-editing_404.tgz"
  explicit_appname = "splunk-app-for-lookup-file-editing"
}