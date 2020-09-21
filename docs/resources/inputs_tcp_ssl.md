# Resource: splunk_inputs_tcp_ssl
Access or update the SSL configuration for the host.

## Example Usage
```
resource "splunk_inputs_tcp_ssl" "test" {
  disabled     = false
  require_client_cert = true
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTinput#data.2Finputs.2Ftcp.2Fssl

This resource block supports the following arguments:
* `disabled` - (Optional) Indicates if input is disabled.
* `root_ca` - (Optional) Certificate authority list (root file)
* `server_cert` - (Optional) Full path to the server certificate.
* `password` - (Optional) Server certificate password, if any.
* `require_client_cert` - (Optional) Determines whether a client must authenticate.

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the http event collector resource
