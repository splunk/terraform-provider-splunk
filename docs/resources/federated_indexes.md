# Resource: splunk_federated_indexes

Creates and manages a local federated index that maps to a dataset on a federated provider.

## Example Usage

```
resource "splunk_federated_indexes" "remote_main" {
  name               = "remote-main"
  federated_provider = splunk_federated_providers.remote.name
  dataset_type       = "index"
  dataset_name       = "main"
}
```

## Argument Reference

* `name` - (Required) Local federated index name without the `federated:` prefix. Use lowercase letters, numbers, underscores, or hyphens, starting with a letter or number. The name must be at most 2048 characters and must not contain `kvstore`. Changing it creates a new index.
* `federated_provider` - (Required) Name of the federated provider containing the remote dataset.
* `dataset_type` - (Optional) Remote dataset type. Accepts only `index`, `savedsearch`, `lastjob`, or `datamodel`. Defaults to `index`.
* `dataset_name` - (Required) Remote index, saved search, scheduled search (for `lastjob`), or data model name, without a type prefix. For example, use `main`, not `index:main`. Changes to either dataset field update the mapping in place.

Federated indexes require a standard-mode provider. Remote saved searches and data models must be readable by the provider's service account and shared globally or in the provider's app context. See [Splunk's remote dataset documentation](https://help.splunk.com/en/splunk-enterprise/search/search-manual/9.1/run-federated-searches-across-multiple-splunk-deployments/create-a-federated-index).

## Attribute Reference

* `id` - The full Splunk index name, including `federated:`, used for REST operations and federated searches.

## Import

Import using the short name or the full Splunk index name. In either case, set `name = "remote-main"` in configuration:

```
terraform import splunk_federated_indexes.remote_main federated:remote-main
```

Alternatively:

```
terraform import splunk_federated_indexes.remote_main remote-main
```
