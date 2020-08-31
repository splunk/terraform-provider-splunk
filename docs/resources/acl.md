# Resource: acl
The ACL resource is a dependent resource. It is optional. The ACL context applies to other Splunk resources.
Typically, knowledge objects, such as saved searches or event types, have an app/user context that is the namespace.

## Example Usage
```
resource "splunk_inputs_http_event_collector" "hec-token-01" {
  name       = "hec-token-01"
  index      = "main"
  source     = "new:source"
  sourcetype = "new:sourcetype"
  acl {
    owner   = "user01"
    sharing = "global"
    read    = ["admin"]
    write   = ["admin"]
  }
}

resource "splunk_saved_searches" "new-search" {
  name                      = "new-search-01"
  search                    = "index=user01-index source=http:hec-token-01"
  acl {
    app     = "search"
    owner   = "nobody"
    sharing = "user"
  }
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.1.0/RESTUM/RESTusing#Access_Control_List

This resource block supports the following arguments:
* `app` - (Optional) The app context for the resource. Required for updating saved search ACL properties. Allowed values are: The name of an app and system.
* `owner` - (Optional) User name of resource owner. Defaults to the resource creator. Required for updating any knowledge object ACL properties. nobody = All users may access the resource, but write access to the resource might be restricted.
* `sharing` - (Optional) Indicates how the resource is shared. Required for updating any knowledge object ACL properties.
<br>app: Shared within a specific app<br>global: (Default) Shared globally to all apps<br>user: Private to a user
* `read` - (Optional) Properties that indicate resource read permissions
* `write` - (Optional) Properties that indicate write permissions of the resource
* `can_change_perms` - (Optional) Indicates if the active user can change permissions for this object. Defaults to true.
* `can_share_app` - (Optional) Indicates if the active user can change sharing to app level. Defaults to true.
* `can_share_global` - (Optional) Indicates if the active user can change sharing to system level. Defaults to true.
* `can_share_user` - (Optional) Indicates if the active user can change sharing to user level. Defaults to true.
* `can_write` - (Optional) Indicates if the active user can edit this object. Defaults to true.
* `removable` - (Optional) Indicates whether an admin or user with sufficient permissions can delete the entity.
