# Resource: splunk_inputs_tcp_cooked
Create or update cooked TCP input information and create new containers for managing cooked data.

## Example Usage
```
resource "splunk_inputs_tcp_cooked" "tcp_cooked" {
    name = "50000"
    disabled = false
    connection_host = "dns"
    restrict_to_host = "splunk"
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.1.0/RESTREF/RESTinput#data.2Finputs.2Ftcp.2Fcooked

This resource block supports the following arguments:
* `name` - (Required) The port number of this input.
* `disabled` - (Optional) Indicates if input is disabled.
* `host` - (Optional) Host from which the indexer gets data.
* `restrict_to_host` - (Optional) Restrict incoming connections on this port to the host specified here.
* `connection_host` - (Optional) Valid values: (ip | dns | none)
                                 Set the host for the remote server that is sending data.
                                 ip sets the host to the IP address of the remote server sending data.
                                 dns sets the host to the reverse DNS entry for the IP address of the remote server sending data.
                                 none leaves the host as specified in inputs.conf, which is typically the Splunk system hostname.
                                 Default value is dns.
* `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the http event collector resource
