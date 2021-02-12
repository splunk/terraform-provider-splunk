# Resource: splunk_inputs_monitor
Create or update a new file or directory monitor input.

## Example Usage
```
resource "splunk_inputs_monitor" "monitor" {
  name     = "opt/splunk/var/log/splunk/health.log"
  recursive = true
  sourcetype = "text"
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTinput#data.2Finputs.2Fmonitor

This resource block supports the following arguments:
* `name` - (Required) The file or directory path to monitor on the system.
* `index` - (Optional) Which index events from this input should be stored in. Defaults to default.
* `host` - (Optional) The value to populate in the host field for events from this data input.
* `sourcetype` - (Optional) The value to populate in the sourcetype field for incoming events.
* `disabled` - (Optional) Indicates if input monitoring is disabled.
* `rename_source` - (Optional) The value to populate in the source field for events from this data input. The same source should not be used for multiple data inputs.
* `blacklist` - (Optional) Specify a regular expression for a file path. The file path that matches this regular expression is not indexed.
* `whitelist` - (Optional) Specify a regular expression for a file path. Only file paths that match this regular expression are indexed.
* `crc_salt` - (Optional) A string that modifies the file tracking identity for files in this input. The magic value <SOURCE> invokes special behavior.
* `follow_tail` - (Optional) If set to true, files that are seen for the first time is read from the end.
* `recursive` - (Optional) Setting this to false prevents monitoring of any subdirectories encountered within this data input.
* `host_regex` - (Optional) Specify a regular expression for a file path. If the path for a file matches this regular expression, the captured value is used to populate the host field for events from this data input. The regular expression must have one capture group.
* `host_segment` - (Optional) Use the specified slash-separate segment of the filepath as the host field value.
* `ignore_older_than` - (Optional) Specify a time value. If the modification time of a file being monitored falls outside of this rolling time window, the file is no longer being monitored.
* `time_before_close` - (Optional) When Splunk software reaches the end of a file that is being read, the file is kept open for a minimum of the number of seconds specified in this value. After this period has elapsed, the file is checked again for more data.
* `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the http event collector resource
