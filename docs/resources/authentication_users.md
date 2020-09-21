# Resource: splunk_authentication_users
Create and update user information or delete the user.

## Example Usage
```
resource "splunk_authentication_users" "user01" {
  name              = "user01"
  email             = "user01@example.com"
  password          = "password01"
  force_change_pass = false
  roles             = ["terraform-user01-role"]
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTaccess#authentication.2Fusers

This resource block supports the following arguments:
* `name` - (Required) Unique user login name.
* `default_app` - (Optional) User default app. Overrides the default app inherited from the user roles.
* `email` - (Optional) User email address.
* `force_change_pass` - Force user to change password indication
* `password` - (Optional) User login password.
* `restart_background_jobs` - (Optional) Restart background search job that has not completed when Splunk restarts indication.
* `realname` - (Optional) Full user name.
* `roles` - (Optional) Role to assign to this user. At least one existing role is required.
* `tz` - (Optional) User timezone.

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the resource
