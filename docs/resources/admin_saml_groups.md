# Resource: splunk_admin_saml_groups
Manage external groups in an IdP response to internal Splunk roles.

## Example Usage
```
resource "splunk_admin_saml_groups" "saml-group" {
  name              = "mygroup"
  roles             = ["admin", "power"]
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTaccess#admin.2FSAML-groups

This resource block supports the following arguments:
* `name` - (Required) The name of the external group.
* `roles` - (Required) List of internal roles assigned to the group.
* `use_client` - (Optional) Set to explicitly specify which client to use for this resource. Leave unset to use the provider's default. Permitted non-empty values are legacy and external.
The legacy client is being replaced by a standalone Splunk client with improved error and drift handling. The legacy client will be deprecated in a future version.

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID (group_name) of the resource

## Import

SAML groups can be imported using the id, e.g.

```
terraform import splunk_admin_saml_groups.saml-group mygroup
```
