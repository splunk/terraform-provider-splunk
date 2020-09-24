# Resource: splunk_authorization_roles
Create and update user information or delete the user.

## Example Usage
```
resource "splunk_authorization_roles" "user01" {
  name              = "user01"
  email             = "user01@example.com"
  password          = "password01"
  force_change_pass = false
  roles             = ["terraform-user01-role"]
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTaccess#authorization.2Froles

This resource block supports the following arguments:
* `name` - (Required) The name of the user role to create.
* `capabilities` - (Optional) List of capabilities assigned to role.
* `cumulative_realtime_search_jobs_quota` - (Optional) Maximum number of concurrently running real-time searches that all members of this role can have.
* `cumulative_search_jobs_quota` - (Optional) Maximum number of concurrently running searches for all role members. Warning message logged when limit is reached.
* `default_app` - Specify the folder name of the default app to use for this role. A user-specific default app overrides this.
* `imported_roles` - (Optional) List of imported roles for this role. <br>Importing other roles imports all aspects of that role, such as capabilities and allowed indexes to search. In combining multiple roles, the effective value for each attribute is value with the broadest permissions.
* `realtime_search_jobs_quota` - (Optional) Specify the maximum number of concurrent real-time search jobs for this role. This count is independent from the normal search jobs limit.
* `search_disk_quota` - (Optional) Specifies the maximum disk space in MB that can be used by a user's search jobs. For example, a value of 100 limits this role to 100 MB total.
* `search_filter` - (Optional) Specify a search string that restricts the scope of searches run by this role. Search results for this role only show events that also match the search string you specify. In the case that a user has multiple roles with different search filters, they are combined with an OR.
* `search_indexes_allowed` - (Optional) Index that this role has permissions to search. Pass this argument once for each index that you want to specify. These may be wildcarded, but the index name must begin with an underscore to match internal indexes.
* `search_indexes_default` - (Optional) For this role, indexes to search when no index is specified. These indexes can be wildcarded, with the exception that '*' does not match internal indexes. To match internal indexes, start with '_'. All internal indexes are represented by '_*'. A user with this role can search other indexes using "index= "
* `search_jobs_quota` - (Optional) The maximum number of concurrent searches a user with this role is allowed to run. For users with multiple roles, the maximum quota value among all of the roles applies.
* `search_time_win` - (Optional) Maximum time span of a search, in seconds. By default, searches are not limited to any specific time window. To override any search time windows from imported roles, set srchTimeWin to '0', as the 'admin' role does.

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the resource
