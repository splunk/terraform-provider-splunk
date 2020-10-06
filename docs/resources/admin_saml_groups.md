# Resource: splunk_admin_saml_groups
Manage external groups in an IdP response to internal Splunk roles.

## Example Usage
```
resource "splunk_admin_saml_groups" "saml-group01" {
  name              = "mygroup"
  roles             = ["admin", "power"]
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.0.6/RESTREF/RESTaccess#admin.2FSAML-groups

This resource block supports the following arguments:
* `name` - (Required) The name of the external group.
* `roles` - (Required) List of internal roles assigned to the group.

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID (group_name) of the resource

## Import

SAML groups can be imported using the id, e.g.

```
terraform import splunk_admin_saml_groups.saml-group01 mygroup
```