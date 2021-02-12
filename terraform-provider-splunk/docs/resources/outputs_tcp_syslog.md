# Resource: splunk_outputs_tcp_syslog
Access the configuration of a forwarded server configured to provide data in standard syslog format.

## Example Usage
```
resource "splunk_outputs_tcp_syslog" "tcp_syslog" {
    name = "new-syslog"
    server = "new-host-1:1234"
    priority = 5
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESToutput#data.2Foutputs.2Ftcp.2Fsyslog

This resource block supports the following arguments:
* `name` - (Required) Name of the syslog output group. This is name used when creating syslog configuration in outputs.conf.
* `disabled` - (Optional) If true, disables global syslog settings.
* `priority` - Sets syslog priority value. The priority value should specified as an integer. See $SPLUNK_HOME/etc/system/README/outputs.conf.spec for details.
* `server` - (Optional) host:port of the server where syslog data should be sent
* `syslog_sourcetype` - (Optional) Specifies a rule for handling data in addition to that provided by the "syslog" sourcetype. By default, there is no value for syslogSourceType.
                                   <br>This string is used as a substring match against the sourcetype key. For example, if the string is set to 'syslog', then all source types containing the string "syslog" receives this special treatment.
                                   To match a source type explicitly, use the pattern "sourcetype::sourcetype_name." For example
                                       syslogSourcetype = sourcetype::apache_common
                                   Data that is "syslog" or matches this setting is assumed to already be in syslog format.
                                   Data that does not match the rules has a header, potentially a timestamp, and a hostname added to the front of the event. This is how Splunk software causes arbitrary log data to match syslog expectations.
* `timestamp_format` - (Optional) Format of timestamp to add at start of the events to be forwarded.
                                  The format is a strftime-style timestamp formatting string. See $SPLUNK_HOME/etc/system/README/outputs.conf.spec for details.
* `type` - (Optional) Protocol to use to send syslog data. Valid values: (tcp | udp ).
* `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the http event collector resource
