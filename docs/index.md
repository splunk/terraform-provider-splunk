# <provider> Provider

The Splunk provider can interact with the resources supported by Splunk. The provider needs to be configured with the proper credentials before it can be used.

## Example Usage

```
provider "splunk" {
  url                  = "localhost:8089"
  username             = "admin"
  password             = "changeme"
  insecure_skip_verify = true
}
```

## Argument Reference

* `url` - (Required) The URL for the Splunk instance to be configured. (The provider uses `https` as the default schema as prefix to the URL)
* `username` - (Required) The username to access the Splunk instance to be configured.
* `password` - (Required) The password to access the Splunk instance to be configured.
* `insecure_skip_verify` - (Optional) Insecure skip verification flag (Defaults to `true`)
