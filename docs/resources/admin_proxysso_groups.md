# Resource: splunk_admin_proxysso_groups
Manage external groups in an proxy-sso response to internal Splunk roles.

## Example Usage
```
resource "splunk_admin_proxysso_groups" "proxy-group" {
  name              = "mygroup"
  roles             = ["admin", "power"]
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/9.2.1/RESTREF/RESTaccess#admin/ProxySSO-groups

This resource block supports the following arguments:
* `name` - (Required) The name of the external group.
* `roles` - (Required) List of internal roles assigned to the group.

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID (group_name) of the resource

## Import

Proxy-SSO groups can be imported using the id, e.g.

```
terraform import splunk_admin_proxysso_groups.proxy-group mygroup
```
