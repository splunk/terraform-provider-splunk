# Resource: splunk_lookup_table_file
Create and manage lookup table files.

## Requirements
This resource uses the Splunk `data/lookup_edit` REST API to create and update lookup file contents. That API is not part of core Splunk Enterprise; it is typically provided by the **Splunk App for Lookup File Editing** (Lookup Editor app). If the API is not available, the provider will receive 404 errors. Install the app on your Splunk instance if you need to manage lookup table file contents with Terraform.

## Example Usage
```
resource "splunk_lookup_table_file" "lookup_table_file" {
  app           = "search"
  owner         = "nobody"
  file_name     = "lookup.csv"
  file_contents = [
    ["status", "status_description", "status_type"],
    ["100", "Continue", "Informational"],
    ["101", "Switching Protocols", "Informational"],
    ["200", "OK", "Successful"]
  ]
}
```

For large CSV files, use `file_path` to let the provider read the file without sending its contents through Terraform's provider RPC:

```
resource "splunk_lookup_table_file" "large_lookup" {
  app       = "search"
  owner     = "nobody"
  file_name = "large_lookup.csv"
  file_path = "${path.module}/large_lookup.csv"
}
```

The provider reads `file_path` during plan and apply on the machine where Terraform is running, including a remote agent. A relative path is resolved from the provider process working directory; use `path.module` or an absolute path, and keep the file in the configuration that the runner receives. Switching an existing resource from `file_contents` to `file_path` is an in-place update. `app`, `owner`, and `file_name` still force a new resource.

State stores `file_contents_hash` and does not store the CSV. That hash is the SHA-256 hex digest of the JSON encoding of the parsed rows, not a digest of the raw file bytes. The provider uploads those rows in one request, so Splunk's form-post size is the remaining limit. An edit made to the lookup in Splunk is left as-is until the local file changes.

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/Knowledge/LookupexampleinSplunkWeb

This resource block supports the following arguments:
* `app` - (Required) The app context for the resource.
* `owner` - (Required) User name of resource owner. Defaults to the resource creator. Required for updating any knowledge object ACL properties. nobody = All users may access the resource, but write access to the resource might be restricted.
* `file_name` - (Required) A name for the lookup table file. Generally ends with ".csv"
* `file_contents` - (Optional) The column header and row value contents for the lookup table file. Specify exactly one of `file_contents` or `file_path`.
* `file_path` - (Optional) Path to a local CSV file. Use an absolute path or `path.module`. The provider reads this path during plan and apply. A change to the file contents triggers an update. Specify exactly one of `file_contents` or `file_path`.

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the lookup table file resource
* `file_contents_hash` - SHA-256 hex digest of the JSON encoding of the parsed CSV rows when `file_path` is set. This is not the SHA-256 of the raw file. Empty when `file_contents` is used.
