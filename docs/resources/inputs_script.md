# Resource: splunk_inputs_script
Create or update scripted inputs.

## Example Usage
```
resource "splunk_inputs_script" "script" {
  name     = "opt/splunk/bin/scripts/readme.txt"
  interval = 360
}
```

```
terraform import splunk_inputs_script.bar "\$SPLUNK_HOME/etc/apps/splunk_instrumentation/bin/instrumentation.py"
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.1.0/RESTREF/RESTinput#data.2Finputs.2Fscript

This resource block supports the following arguments:
* `name` - (Required) Specify the name of the scripted input.
* `disabled` - (Optional) Specifies whether the input script is disabled.
* `index` - (Optional) Sets the index for events from this input. Defaults to the main index.
* `host` - (Optional) Sets the host for events from this input. Defaults to whatever host sent the event.
* `sourcetype` - (Optional) Sets the sourcetype key/field for events from this input. If unset, Splunk software picks a source type based on various aspects of the data. As a convenience, the chosen string is prepended with 'sourcetype::'. There is no hard-coded default.
                            Sets the sourcetype key initial value. The key is used during parsing/indexing, in particular to set the source type field during indexing. It is also the source type field used at search time.
* `source` - (Optional) Sets the source key/field for events from this input. Defaults to the input file path.
                        Sets the source key initial value. The key is used during parsing/indexing, in particular to set the source field during indexing. It is also the source field used at search time. As a convenience, the chosen string is prepended with 'source::'.
* `rename_source` - (Optional) Specify a new name for the source field for the script.
* `passauth` - (Optional) User to run the script as. If you provide a username, Splunk software generates an auth token for that user and passes it to the script.
* `interval` - (Optional) Specify an integer or cron schedule. This parameter specifies how often to execute the specified script, in seconds or a valid cron schedule. If you specify a cron schedule, the script is not executed on start-up.
* `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the http event collector resource
