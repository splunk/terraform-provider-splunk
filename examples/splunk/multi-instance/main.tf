# Example Splunk Terraform Configuration
# Create provider resource blocks with alias for multi instance configuration
# This config creates a splunk index, some basic file monitoring inputs and installs splunk_TA_unix on a splunk indexer
# This config creates a new splunk user with a role, installs splunk_unix_app and an alert on a splunk search head.

# initialization
terraform {
  required_providers {
    splunk = {
      source  = "splunk/splunk"
      version = "1.4.2"
    }
  }
}

# provider block
provider "splunk" {
  alias = "indexer"
}

// search head configuration
provider "splunk" {
  alias = "search_head"
}

# resource blocks
# // On the indexer
resource "splunk_indexes" "monitoring-index" {
  provider               = splunk.indexer
  name                   = "file-monitoring"
  max_hot_buckets        = 6
  max_total_data_size_mb = 1000000
}

resource "splunk_inputs_monitor" "monitor_var_log" {
  provider      = splunk.indexer
  name          = "/var/log"
  blacklist     = ".gz$"
  recursive     = true
  rename_source = "/var/log"
  index         = splunk_indexes.monitoring-index.name
  depends_on = [
    splunk_indexes.monitoring-index
  ]
}

resource "splunk_inputs_monitor" "monitor_etc" {
  provider      = splunk.indexer
  name          = "/etc"
  blacklist     = ".gz$"
  recursive     = true
  rename_source = "/etc"
  index         = splunk_indexes.monitoring-index.name
  depends_on = [
    splunk_indexes.monitoring-index
  ]
}

resource "splunk_apps_local" "splunk_TA_nix" {
  provider         = splunk.indexer
  name             = "/splunk-add-on-for-unix-and-linux_820.tgz" // Location of the app on the remote instance
  filename         = true
  explicit_appname = "Splunk_TA_nix"
  update           = true

  depends_on = [
    splunk_indexes.monitoring-index,
    splunk_inputs_monitor.monitor_etc,
    splunk_inputs_monitor.monitor_var_log
  ]
}

// On the search head
resource "splunk_authorization_roles" "role01" {
  provider       = splunk.search_head
  name           = "terraform-user-role"
  default_app    = "search"
  imported_roles = ["power", "user"]
  capabilities   = ["accelerate_datamodel", "change_authentication", "restart_splunkd"]
}

resource "splunk_authentication_users" "user01" {
  provider          = splunk.search_head
  name              = "user01"
  email             = "user01@example.com"
  password          = "changeme"
  force_change_pass = false
  roles             = ["terraform-user-role"]
  depends_on = [
    splunk_authorization_roles.role01
  ]
}

resource "splunk_apps_local" "splunk_unix_app" {
  provider         = splunk.search_head
  name             = "/splunk-app-for-unix-and-linux_600.tgz" // Location of the app on the remote instance
  filename         = true
  explicit_appname = "splunk_app_for_nix"
  update           = true

  depends_on = [
    splunk_authentication_users.user01
  ]
}

resource "splunk_saved_searches" "syslog_alert" {
  provider            = splunk.search_head
  actions             = "email"
  action_email_format = "table"
  action_email_to     = "user01@splunk.com"
  alert_comparator    = "greater than"
  alert_digest_mode   = true
  alert_expires       = "30d"
  alert_threshold     = "0"
  alert_type          = "number of events"
  cron_schedule       = "*/1 * * * *"
  name                = "syslog error alert"
  disabled            = false
  is_scheduled        = true
  is_visible          = true
  realtime_schedule   = true
  search              = "index=${splunk_indexes.monitoring-index.name} sourcetype=syslog level=error"

  acl {
    app     = "search"
    owner   = "user01"
    sharing = "user"
  }

  depends_on = [
    splunk_authentication_users.user01,
    splunk_indexes.monitoring-index
  ]
}
