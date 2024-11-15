
# Resource: splunk_lookup_definition

Manage lookup definitions in Splunk. For more information on lookup definitions, refer to the official Splunk documentation: https://docs.splunk.com/Documentation/Splunk/latest/Knowledge/Aboutlookupsandfieldactions

## Example Usage
```hcl
resource "splunk_lookup_definition" "example" {
  name     = "example_lookup_definition"
  filename = "example_lookup_file.csv"
  acl {
    owner   = "admin"
    app     = "search"
    sharing = "app"
    read    = ["*"]
    write   = ["admin"]
  }
}
```

## Argument Reference
This resource block supports the following arguments:
* `name` - (Required) A unique name for the lookup definition within the app context.
* `filename` - (Required) The filename for the lookup table, usually ending in `.csv`.
* `acl` - (Optional) Defines the access control list (ACL) for the lookup definition. See [acl.md](acl.md) for more details.

## Validation Rules
When `acl.sharing` is set to `user`, the `acl.read` and `acl.write` fields must not be explicitly set. Setting them will trigger a validation error.

## Attribute Reference
In addition to the arguments listed above, this resource exports the following attributes:

* `id` - The ID of the lookup table file resource.
