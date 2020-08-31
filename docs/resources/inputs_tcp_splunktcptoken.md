# Resource: splunk_inputs_tcp_splunk_tcp_token
Manage receiver access using tokens.

## Example Usage
```
resource "splunk_inputs_tcp_splunk_tcp_token" "tcp_splunk_tcp_token" {
    name = "new-splunk-tcp-token"
    token = "D66C45B3-7C28-48A1-A13A-027914146501"
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.1.0/RESTREF/RESTinput#data.2Finputs.2Ftcp.2Fsplunktcptoken

This resource block supports the following arguments:
* `name` - (Required) Required. Name for the token to create.
* `token` - (Optional) Optional. Token value to use. If unspecified, a token is generated automatically.
* `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the http event collector resource
