# Resource: splunk_outputs_tcp_server
Access data forwarding configurations.

## Example Usage
```
resource "splunk_outputs_tcp_server" "tcp_server" {
    name = "new-host:1234"
    ssl_alt_name_to_check = "old-host"
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.1.0/RESTREF/RESToutput#data.2Foutputs.2Ftcp.2Fserver

This resource block supports the following arguments:
* `name` - (Required) <host>:<port> of the Splunk receiver. <host> can be either an ip address or server name. <port> is the that port that the Splunk receiver is listening on.
* `disabled` - (Optional) If true, disables the group.
* `method` - Valid values: (clone | balance | autobalance)
             The data distribution method used when two or more servers exist in the same forwarder group.
* `ssl_alt_name_to_check` - (Optional) The alternate name to match in the remote server's SSL certificate.
* `ssl_cert_path` - (Optional) Path to the client certificate. If specified, connection uses SSL.
* `ssl_cipher` - (Optional) SSL Cipher in the form ALL:!aNULL:!eNULL:!LOW:!EXP:RC4+RSA:+HIGH:+MEDIUM
* `ssl_common_name_to_check` - (Optional) Check the common name of the server's certificate against this name.
                                          If there is no match, assume that Splunk Enterprise is not authenticated against this server. You must specify this setting if sslVerifyServerCert is true.
* `ssl_root_ca_path` - (Optional) The path to the root certificate authority file.
* `ssl_password` - (Optional) The password associated with the CAcert.
                              The default Splunk Enterprise CAcert uses the password "password."
* `ssl_verify_server_cert` - (Optional) If true, make sure that the server you are connecting to is a valid one (authenticated). Both the common name and the alternate name of the server are then checked for a match.
* `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the http event collector resource
