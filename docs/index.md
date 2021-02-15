# Splunk Provider

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

Terraform 0.13 and later must add:
```
terraform {
   required_providers {
    splunk = {
      source  = "splunk/splunk"
    }
  }
}
```

## Argument Reference

Below arguments for the provider can also be set as environment variables.

* `url` or `SPLUNK_URL` - (Required) The URL for the Splunk instance to be configured. (The provider uses `https` as the default schema as prefix to the URL)
* `username` or `SPLUNK_USERNAME`  - (Optional) The username to access the Splunk instance to be configured.
* `password` or `SPLUNK_PASSWORD` - (Optional) The password to access the Splunk instance to be configured.
* `auth_token` or `SPLUNK_AUTH_TOKEN` - (Optional) Use auth token instead of username and password to configure Splunk instance.
If specified, auth token takes priority over username/password.
* `insecure_skip_verify` or `SPLUNK_INSECURE_SKIP_VERIFY` - (Optional) Insecure skip verification flag (Defaults to `true`)
* `timeout` or `SPLUNK_TIMEOUT` -  (Optional) Timeout when making calls to Splunk server. (Defaults to `60 seconds`)

(NOTE: Auth token can only be used with certain type of Splunk deployments.
Read more on authentication with tokens here: https://docs.splunk.com/Documentation/Splunk/latest/Security/Setupauthenticationwithtokens)
