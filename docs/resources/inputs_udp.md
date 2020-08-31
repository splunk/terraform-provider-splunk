# Resource: splunk_inputs_tcp_raw
Create and manage UDP data inputs.

## Example Usage
```
resource "splunk_inputs_udp" "udp" {
    name = "41000"
    index = "main"
    source = "new"
    sourcetype = "new"
    disabled = false
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.1.0/RESTREF/RESTinput#data.2Finputs.2Fudp

This resource block supports the following arguments:
* `name` - (Required) The UDP port that this input should listen on.
* `disabled` - (Optional) Indicates if input is disabled.
* `index` - (Optional) Which index events from this input should be stored in. Defaults to default.
* `host` - (Optional) The value to populate in the host field for incoming events. This is used during parsing/indexing, in particular to set the host field. It is also the host field used at search time.
* `source` - (Optional) The value to populate in the source field for incoming events. The same source should not be used for multiple data inputs.
* `sourcetype` - (Optional) The value to populate in the sourcetype field for incoming events.
* `restrict_to_host` - (Optional) Restrict incoming connections on this port to the host specified here.
                                  If this is not set, the value specified in [udp://<remote server>:<port>] in inputs.conf is used.
* `queue` - (Optional) Which queue events from this input should be sent to. Generally this does not need to be changed.
* `connection_host` - (Optional) Valid values: (ip | dns | none)
                                 Set the host for the remote server that is sending data.
                                 ip sets the host to the IP address of the remote server sending data.
                                 dns sets the host to the reverse DNS entry for the IP address of the remote server sending data.
                                 none leaves the host as specified in inputs.conf, which is typically the Splunk system hostname.
                                 Default value is dns.
* `no_appending_timestamp` - (Optional) If set to true, prevents Splunk software from prepending a timestamp and hostname to incoming events.
* `no_priority_stripping` - (Optional) If set to true, Splunk software does not remove the priority field from incoming syslog events.
* `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the http event collector resource
