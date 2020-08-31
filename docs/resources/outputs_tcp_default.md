# Resource: splunk_outputs_tcp_default
Manage to global tcpout properties.

## Example Usage
```
resource "splunk_outputs_tcp_default" "tcp_default" {
    name = "tcpout"
    disabled = false
    default_group = "test-indexers"
    drop_events_on_queue_full = 60
    index_and_forward = true
    send_cooked_data = true
    max_queue_size = "100KB"
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.1.0/RESTREF/RESToutput#data.2Foutputs.2Ftcp.2Fdefault

This resource block supports the following arguments:
* `name` - (Required) Configuration to be edited. The only valid value is "tcpout".
* `default_group` - (Optional) Comma-separated list of one or more target group names, specified later in [tcpout:<target_group>] stanzas of outputs.conf.spec file.
                               The forwarder sends all data to the specified groups. If you do not want to forward data automatically, do not set this attribute. Can be overridden by an inputs.conf _TCP_ROUTING setting, which in turn can be overridden by a props.conf/transforms.conf modifier.
* `disabled` - (Optional) Disables default tcpout settings
* `drop_events_on_queue_full` - (Optional) If set to a positive number, wait the specified number of seconds before throwing out all new events until the output queue has space. Defaults to -1 (do not drop events).
                                           <br>CAUTION: Do not set this value to a positive integer if you are monitoring files.
                                           Setting this to -1 or 0 causes the output queue to block when it gets full, which causes further blocking up the processing chain. If any target group queue is blocked, no more data reaches any other target group.
                                           Using auto load-balancing is the best way to minimize this condition, because, in that case, multiple receivers must be down (or jammed up) before queue blocking can occur.
* `heartbeat_frequency` - (Optional) How often (in seconds) to send a heartbeat packet to the receiving server.
                                     Heartbeats are only sent if sendCookedData=true. Defaults to 30 seconds.
* `index_and_forward` - (Optional) Specifies whether to index all data locally, in addition to forwarding it. Defaults to false.
                                   This is known as an "index-and-forward" configuration. This attribute is only available for heavy forwarders. It is available only at the top level [tcpout] stanza in outputs.conf. It cannot be overridden in a target group.
* `max_queue_size` - (Optional) Specify an integer or integer[KB|MB|GB].
                                <br>Sets the maximum size of the forwarder output queue. It also sets the maximum size of the wait queue to 3x this value, if you have enabled indexer acknowledgment (useACK=true).
                                Although the wait queue and the output queues are both configured by this attribute, they are separate queues. The setting determines the maximum size of the queue in-memory (RAM) buffer.
                                For heavy forwarders sending parsed data, maxQueueSize is the maximum number of events. Since events are typically much shorter than data blocks, the memory consumed by the queue on a parsing forwarder is likely to be much smaller than on a non-parsing forwarder, if you use this version of the setting.
                                If specified as a lone integer (for example, maxQueueSize=100), maxQueueSize indicates the maximum number of queued events (for parsed data) or blocks of data (for unparsed data). A block of data is approximately 64KB. For non-parsing forwarders, such as universal forwarders, that send unparsed data, maxQueueSize is the maximum number of data blocks.
                                If specified as an integer followed by KB, MB, or GB (for example, maxQueueSize=100MB), maxQueueSize indicates the maximum RAM allocated to the queue buffer. Defaults to 500KB (which means a maximum size of 500KB for the output queue and 1500KB for the wait queue, if any).
* `send_cooked_data` - (Optional) If true, events are cooked (processed by Splunk software). If false, events are raw and untouched prior to sending. Defaults to true.
                                  Set to false if you are sending to a third-party system.
* `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the http event collector resource
