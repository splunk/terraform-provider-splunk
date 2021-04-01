# Resource: splunk_sh_indexes_manager
Create indexes on Splunk Cloud instances. [BETA]

## Authorization and authentication
As of now there is no support to create indexes in user-specified workspaces on Splunk Cloud.

## Example Usage
```
resource "splunk_sh_indexes_manager" "tf-index" {
    name = "tf-test-index-0"
    datatype = "event"
    frozen_time_period_in_secs = "94608000"
    max_global_raw_data_size_mb = "100"
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/RESTTUT/RESTandCloud#Example_use_cases

This resource block supports the following arguments:
* `name` - (Required) The name of the index to create.
* `datatype` - (Optional)  	Valid values: (event | metric). Specifies the type of index.
* `frozen_time_period_in_secs` - (Optional) Number of seconds after which indexed data rolls to frozen.
Defaults to 94608000 (3 years).Freezing data means it is removed from the index. If you need to archive your data, refer to coldToFrozenDir and coldToFrozenScript parameter documentation.
* `max_global_raw_data_size_mb` - (Optional) The maximum size of an index (in MB). If an index grows larger than the maximum size, the oldest data is frozen.
  Defaults to 100 MB.

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the splunk_sh_indexes_manager resource
