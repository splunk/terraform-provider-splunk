variable "remote_splunk_password" {
  type      = string
  sensitive = true
}

resource "splunk_federated_providers" "remote" {
  name            = "remote-splunk"
  host_port       = "remote-splunk:8089"
  service_account = "federated_search"
  password        = var.remote_splunk_password
}

resource "splunk_federated_indexes" "remote_main" {
  name               = "remote-main"
  federated_provider = splunk_federated_providers.remote.name
  dataset_type       = "index"
  dataset_name       = "main"
}
