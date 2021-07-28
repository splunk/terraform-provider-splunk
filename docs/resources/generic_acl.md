# Resource: splunk_generic_acl
Manage the ACL of any Splunk object.

Note: This resource doesn't actually create any remote resources, because ACLs can only exist (and always exist) for
knowledge objects. They can, however, be managed separately.

## Example Usage
```
resource "splunk_generic_acl" "my_app" {
  # apps are managed via the apps/local/<app> endpoint
  path = "apps/local/my_app"
  acl {
    # use app=system, owner=nobody when managing apps, as they have no owner or app context
    app   = "system"
    owner = "nobody"
    read  = ["*"]
    write = ["admin", "power"]
  }
}

resource "splunk_generic_acl" "my_dashboard" {
  path = "data/ui/views/my_dashboard"
  acl {
    app   = "my_app"
    owner = "joe_user"
    read  = ["team_joe"]
    write = ["team_joe"]
  }
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTapps#apps.2Flocal

This resource block supports the following arguments:
* `path` - (Required) REST API Endpoint path to the object, relative to servicesNS/<owner>/<app>
* `acl` - (Optional) The ACL to apply to the object, including app/owner to properly identify the object.
  Though technically optional, it should be explicitly set for this resource to really be valid. Some objects, such as
  apps, require specific values for app and owner. Consult the REST API documentation regarding which values to use for
  app and owner for objects that don't fit in the normal namespace.

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the resource
