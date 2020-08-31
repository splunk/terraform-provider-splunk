# Resource: splunk_inputs_tcp_raw
Create or update raw TCP input information for managing raw tcp inputs from forwarders.

## Example Usage
```
resource "splunk_inputs_tcp_raw" "tcp_raw" {
    name = "41000"
    index = "main"
    queue = "indexQueue"
    source = "new"
    sourcetype = "new"
    disabled = false
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.1.0/RESTREF/RESTinput#data.2Finputs.2Ftcp.2Fraw

This resource block supports the following arguments:
* `name` - (Required) The input port which receives raw data.
* `disabled` - (Optional) Indicates if input is disabled.
* `index` - (Optional) Index to store generated events. Defaults to default.
* `host` - (Optional) Host from which the indexer gets data.
* `source` - (Optional) Sets the source key/field for events from this input. Defaults to the input file path.
                        Sets the source key initial value. The key is used during parsing/indexing, in particular to set the source field during indexing. It is also the source field used at search time. As a convenience, the chosen string is prepended with 'source::'.
* `sourcetype` - (Optional) Set the source type for events from this input.
                            "sourcetype=" is automatically prepended to <string>.
                            Defaults to audittrail (if signedaudit=true) or fschange (if signedaudit=false).
* `restrict_to_host` - (Optional) Allows for restricting this input to only accept data from the host specified here.
* `queue` - (Optional) Valid values: (parsingQueue | indexQueue)
                       Specifies where the input processor should deposit the events it reads. Defaults to parsingQueue.
                       Set queue to parsingQueue to apply props.conf and other parsing rules to your data. For more information about props.conf and rules for timestamping and linebreaking, refer to props.conf and the online documentation at "Monitor files and directories with inputs.conf"
                       Set queue to indexQueue to send your data directly into the index.
* `connection_host` - (Optional) Valid values: (ip | dns | none)
                                 Set the host for the remote server that is sending data.
                                 ip sets the host to the IP address of the remote server sending data.
                                 dns sets the host to the reverse DNS entry for the IP address of the remote server sending data.
                                 none leaves the host as specified in inputs.conf, which is typically the Splunk system hostname.
                                 Default value is dns.
* `raw_tcp_done_timeout` - (Optional) Specifies in seconds the timeout value for adding a Done-key. Default value is 10 seconds.
                                      If a connection over the port specified by name remains idle after receiving data for specified number of seconds, it adds a Done-key. This implies the last event is completely received.
* `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the http event collector resource
