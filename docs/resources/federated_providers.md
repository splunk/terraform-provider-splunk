# Resource: splunk_federated_providers

Creates and manages a Splunk-to-Splunk federated search provider on a Splunk Enterprise search head.

The provider definition is stored in the local deployment's `federated.conf`. It requires an account on the remote deployment that is dedicated to federated searching.

## Security

Use a dedicated, least-privileged remote service account and a trusted TLS certificate on both Splunk management endpoints. For production Terraform provider connections, set `insecure_skip_verify = false`. This resource stores a SHA-256 password hash in state. Protect state and saved plans: hashes permit offline password guessing, and saved plans can contain the configured plaintext password. Historical state from earlier builds may still contain plaintext passwords.

## Example Usage

```
resource "splunk_federated_providers" "remote" {
  name            = "remote-splunk"
  host_port       = "remote-splunk:8089"
  service_account = "federated_search"
  password        = var.remote_splunk_password
  mode            = "standard"
  app_context     = "search"
}
```

## Argument Reference

* `name` - (Required) Unique name for the federated provider. Changing it creates a new provider.
* `host_port` - (Required) Remote Splunk management host and port, such as `remote-splunk:8089`.
* `service_account` - (Required) Username of the service account on the remote deployment.
* `password` - (Required, Sensitive) Password for the remote service account. State stores its SHA-256 hash. Creation and configured password changes send the original password to Splunk; password changes update the provider in place without replacement. Other updates omit the password. Splunk does not return this secret, so changes made directly in Splunk cannot be detected.
* `mode` - (Optional) Federated search mode: `standard` (default) or `transparent`. A local deployment must not mix provider modes.
* `app_context` - (Optional) App on the remote search head used by standard-mode searches. Defaults to `search`; transparent mode ignores this value.

## Import

Import using the provider name:

```
terraform import splunk_federated_providers.remote remote-splunk
```

After import, set `password` in configuration before updating the resource.

## Acceptance testing

The optional federated acceptance test uses the standard repository acceptance-test
variables plus `SPLUNK_FEDERATED_PROVIDER_HOST_PORT`,
`SPLUNK_FEDERATED_PROVIDER_USERNAME`, and `SPLUNK_FEDERATED_PROVIDER_PASSWORD`.
The remote test deployment must contain the built-in `search` and `launcher`
apps and the `main` and `_internal` indexes. Use a disposable test environment.
If any federated credential variable is missing, the test skips with an
explanation, including when `TF_ACC=1`. Existing CI needs no additional secrets,
services, or workflow changes. Mocked lifecycle tests run without Splunk or
credentials. Ordinary unit-test runs leave `TF_ACC` unset or empty.
