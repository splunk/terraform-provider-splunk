# Resource: splunk_lookup_table_file
Create and manage lookup table files.

## Example Usage
```
resource "splunk_lookup_table_file" "lookup_table_file" {
  app           = "search"
  owner         = "nobody"
  file_name     = "lookup.csv"
  file_contents = <<-EOT
[
  ["status", "status_description", "status_type"],
  ["100", "Continue", "Informational"],
  ["101", "Switching Protocols", "Informational"],
  ["200", "OK", "Successful"]
]
EOT
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
