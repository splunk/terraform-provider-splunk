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

The file must be available to the provider process during plan and apply. Terraform stores a hash for change detection instead of the full CSV contents in state.

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/Knowledge/LookupexampleinSplunkWeb

This resource block supports the following arguments:
* `app` - (Required) The app context for the resource.
* `owner` - (Required) User name of resource owner. Defaults to the resource creator. Required for updating any knowledge object ACL properties. nobody = All users may access the resource, but write access to the resource might be restricted.
* `file_name` - (Required) A name for the lookup table file. Generally ends with ".csv"
* `file_contents` - (Optional) The column header and row value contents for the lookup table file. Specify exactly one of `file_contents` or `file_path`.
* `file_path` - (Optional) Path to a local CSV file. Changes to its contents trigger an update.

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the lookup table file resource
* `file_contents_hash` - Hash used to detect changes to a CSV provided through `file_path`.
