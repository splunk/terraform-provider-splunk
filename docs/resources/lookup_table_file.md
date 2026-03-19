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

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/Knowledge/LookupexampleinSplunkWeb

This resource block supports the following arguments:
* `app` - (Required) The app context for the resource.
* `owner` - (Required) User name of resource owner. Defaults to the resource creator. Required for updating any knowledge object ACL properties. nobody = All users may access the resource, but write access to the resource might be restricted.
* `file_name` - (Required) A name for the lookup table file. Generally ends with ".csv"
* `file_contents` - (Required) The column header and row value contents for the lookup table file.

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the lookup table file resource
