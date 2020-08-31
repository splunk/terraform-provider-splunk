# Resource: splunk_global_http_event_collector
Update Global HTTP Event Collector input configuration.

## Example Usage
```
resource "splunk_global_http_event_collector" "http" {
  disabled   = false
  enable_ssl = true
  port       = 8088
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.1.0/RESTREF/RESTinput#data.2Finputs.2Fhttp

This resource block supports the following arguments:
* `disabled` - (Optional) Input disabled indicator.
* `port` - (Optional) HTTP data input IP port.
* `enable_ssl` - (Optional) Enable SSL protocol for HTTP data input. `true` = SSL enabled, `false` = SSL disabled.
* `dedicated_io_threads` - (Optional) Number of threads used by HTTP Input server.
* `max_sockets` - (Optional) Maximum number of simultaneous HTTP connections accepted. Adjusting this value may cause server performance issues and is not generally recommended. Possible values for this setting vary by OS.
* `max_threads` - (Optional) Maximum number of threads that can be used by active HTTP transactions. Adjusting this value may cause server performance issues and is not generally recommended. Possible values for this setting vary by OS.
* `use_deployment_server` - (Optional) Indicates whether the event collector input writes its configuration to a deployment server repository. When this setting is set to 1 (enabled), the input writes its configuration to the directory specified as repositoryLocation in serverclass.conf.
Copy the full contents of the splunk_httpinput app directory to this directory for the configuration to work. When enabled, only the tokens defined in the splunk_httpinput app in this repository are viewable and editable on the API and the Data Inputs page in Splunk Web. When disabled, the input writes its configuration to $SPLUNK_HOME/etc/apps by default. Defaults to 0 (disabled).

