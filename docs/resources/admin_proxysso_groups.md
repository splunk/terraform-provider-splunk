# Resource: splunk_admin_proxysso_groups

Manage Proxy SSO group-to-role mappings on Splunk Enterprise. Proxy SSO is an Enterprise authentication method and is not available on Splunk Cloud.

## Example Usage
```
resource "splunk_admin_proxysso_groups" "proxy_group" {
  name  = "mygroup"
  roles = ["admin", "power"]
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTaccess#admin.2FProxySSO-groups

This resource block supports the following arguments:
* `name` - (Required) The name of the external group. Changing this forces a new resource.
* `roles` - (Required) Set of internal roles assigned to the group. Role order does not matter.

## Attribute Reference
In addition to all arguments above, this resource exports the following attribute:

* `id` - The group name.

## Import

Proxy SSO groups can be imported using the group name, for example:

```
terraform import splunk_admin_proxysso_groups.proxy_group mygroup
```
