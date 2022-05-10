# Resource: splunk_indexes
Create and manage data indexes.

## Authorization and authentication
By default, all users can list all indexes. However, if the indexes_list_all capability is enabled in authorize.conf, access to all indexes is limited to only those roles with this capability.
To enable indexes_list_all capability restrictions on the data/indexes endpoint, create a [capability::indexes_list_all] stanza in authorize.conf. Specify indexes_list_all=enabled for any role permitted to list all indexes from this endpoint.

## Example Usage
```
resource "splunk_indexes" "user01-index" {
  name                   = "user01-index"
  max_hot_buckets        = 6
  max_total_data_size_mb = 1000000
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTintrospect#data.2Findexes

This resource block supports the following arguments:
* `name` - (Required) The name of the index to create.
* `block_sign_size` - (Optional) Controls how many events make up a block for block signatures. If this is set to 0, block signing is disabled for this index. <br>A recommended value is 100.
* `bucket_rebuild_memory_hint` - (Optional) Suggestion for the bucket rebuild process for the size of the time-series (tsidx) file to make.
                                            <be>Caution: This is an advanced parameter. Inappropriate use of this parameter causes splunkd to not start if rebuild is required. Do not set this parameter unless instructed by Splunk Support.
                                            Default value, auto, varies by the amount of physical RAM on the host<br>
                                                less than 2GB RAM = 67108864 (64MB) tsidx
                                                2GB to 8GB RAM = 134217728 (128MB) tsidx
                                                more than 8GB RAM = 268435456 (256MB) tsidx<br>
                                            Values other than "auto" must be 16MB-1GB. Highest legal value (of the numerical part) is 4294967295 You can specify the value using a size suffix: "16777216" or "16MB" are equivalent.
* `cold_path` - (Optional) An absolute path that contains the colddbs for the index. The path must be readable and writable. Cold databases are opened as needed when searching.
* `cold_to_frozen_dir` - (Optional) Destination path for the frozen archive. Use as an alternative to a coldToFrozenScript. Splunk software automatically puts frozen buckets in this directory.
                                    <br>
                                    Bucket freezing policy is as follows:<br>
                                    New style buckets (4.2 and on): removes all files but the rawdata<br>
                                    To thaw, run splunk rebuild <bucket dir> on the bucket, then move to the thawed directory<br>
                                    Old style buckets (Pre-4.2): gzip all the .data and .tsidx files<br>
                                    To thaw, gunzip the zipped files and move the bucket into the thawed directory<br>
                                    If both coldToFrozenDir and coldToFrozenScript are specified, coldToFrozenDir takes precedence
* `cold_to_frozen_script` - (Optional) Path to the archiving script.
                                       <br>If your script requires a program to run it (for example, python), specify the program followed by the path. The script must be in $SPLUNK_HOME/bin or one of its subdirectories.
                                       <br>Splunk software ships with an example archiving script in $SPLUNK_HOME/bin called coldToFrozenExample.py. DO NOT use this example script directly. It uses a default path, and if modified in place any changes are overwritten on upgrade.
                                       <br>It is best to copy the example script to a new file in bin and modify it for your system. Most importantly, change the default archive path to an existing directory that fits your needs.
* `compress_rawdata` - (Optional) This parameter is ignored. The splunkd process always compresses raw data.
* `datatype` - (Optional)  	Valid values: (event | metric). Specifies the type of index.
* `enable_online_bucket_repair` - (Optional) Enables asynchronous "online fsck" bucket repair, which runs concurrently with Splunk software.
When enabled, you do not have to wait until buckets are repaired to start the Splunk platform. However, you might observe a slight performance degratation.
* `frozen_time_period_in_secs` - (Optional) Number of seconds after which indexed data rolls to frozen.
Defaults to 188697600 (6 years).Freezing data means it is removed from the index. If you need to archive your data, refer to coldToFrozenDir and coldToFrozenScript parameter documentation.
* `home_path` - (Optional) An absolute path that contains the hot and warm buckets for the index.
                          Required. Splunk software does not start if an index lacks a valid homePath.
                          <br>Caution: The path must be readable and writable.
* `max_bloom_backfill_bucket_age` - (Optional) Valid values are: Integer[m|s|h|d].
<br>If a warm or cold bucket is older than the specified age, do not create or rebuild its bloomfilter. Specify 0 to never rebuild bloomfilters.
* `max_concurrent_optimizes` - (Optional) The number of concurrent optimize processes that can run against a hot bucket.
This number should be increased if instructed by Splunk Support. Typically the default value should suffice.
* `max_data_size` - (Optional) The maximum size in MB for a hot DB to reach before a roll to warm is triggered. Specifying "auto" or "auto_high_volume" causes Splunk software to autotune this parameter (recommended).
Use "auto_high_volume" for high volume indexes (such as the main index); otherwise, use "auto". A "high volume index" would typically be considered one that gets over 10GB of data per day.
* `max_hot_buckets` - (Optional) Maximum hot buckets that can exist per index. Defaults to 3.
<br>When maxHotBuckets is exceeded, Splunk software rolls the least recently used (LRU) hot bucket to warm. Both normal hot buckets and quarantined hot buckets count towards this total. This setting operates independently of maxHotIdleSecs, which can also cause hot buckets to roll.
* `max_hot_idle_secs` - (Optional) Maximum life, in seconds, of a hot bucket. Defaults to 0. If a hot bucket exceeds maxHotIdleSecs, Splunk software rolls it to warm. This setting operates independently of maxHotBuckets, which can also cause hot buckets to roll. A value of 0 turns off the idle check (equivalent to INFINITE idle time).
* `max_hot_span_secs` - (Optional) Upper bound of target maximum timespan of hot/warm buckets in seconds. Defaults to 7776000 seconds (90 days).
* `max_mem_mb` - (Optional) The amount of memory, expressed in MB, to allocate for buffering a single tsidx file into memory before flushing to disk. Defaults to 5. The default is recommended for all environments.
* `max_meta_entries` - (Optional) Sets the maximum number of unique lines in .data files in a bucket, which may help to reduce memory consumption. If set to 0, this setting is ignored (it is treated as infinite).
If exceeded, a hot bucket is rolled to prevent further increase. If your buckets are rolling due to Strings.data hitting this limit, the culprit may be the punct field in your data. If you do not use punct, it may be best to simply disable this (see props.conf.spec in $SPLUNK_HOME/etc/system/README).
* `max_meta_entries` - (Optional) Upper limit, in seconds, on how long an event can sit in raw slice. Applies only if replication is enabled for this index. Otherwise ignored. If there are any acknowledged events sharing this raw slice, this paramater does not apply. In this case, maxTimeUnreplicatedWithAcks applies. Highest legal value is 2147483647. To disable this parameter, set to 0.
* `max_time_unreplicated_no_acks` - (Optional) Upper limit, in seconds, on how long an event can sit in raw slice. Applies only if replication is enabled for this index. Otherwise ignored.
                                               If there are any acknowledged events sharing this raw slice, this paramater does not apply. In this case, maxTimeUnreplicatedWithAcks applies.
                                               Highest legal value is 2147483647. To disable this parameter, set to 0.
* `max_time_unreplicated_with_acks` - (Optional) Upper limit, in seconds, on how long events can sit unacknowledged in a raw slice. Applies only if you have enabled acks on forwarders and have replication enabled (with clustering).
Note: This is an advanced parameter. Make sure you understand the settings on all forwarders before changing this. This number should not exceed ack timeout configured on any forwarder, and should actually be set to at most half of the minimum value of that timeout. You can find this setting in outputs.conf readTimeout setting under the tcpout stanza.
To disable, set to 0, but this is NOT recommended. Highest legal value is 2147483647.
* `max_total_data_size_mb` - (Optional) The maximum size of an index (in MB). If an index grows larger than the maximum size, the oldest data is frozen.
* `max_warm_db_count` - (Optional) The maximum number of warm buckets. If this number is exceeded, the warm bucket/s with the lowest value for their latest times is moved to cold.
* `min_raw_file_sync_secs` - (Optional) Specify an integer (or "disable") for this parameter.
                                        This parameter sets how frequently splunkd forces a filesystem sync while compressing journal slices.
                                        During this period, uncompressed slices are left on disk even after they are compressed. Then splunkd forces a filesystem sync of the compressed journal and removes the accumulated uncompressed files.
                                        If 0 is specified, splunkd forces a filesystem sync after every slice completes compressing. Specifying "disable" disables syncing entirely: uncompressed slices are removed as soon as compression is complete.
* `min_stream_group_queue_size` - (Optional) Minimum size of the queue that stores events in memory before committing them to a tsidx file.
* `partial_service_meta_period` - (Optional) Related to serviceMetaPeriod. If set, it enables metadata sync every <integer> seconds, but only for records where the sync can be done efficiently in-place, without requiring a full re-write of the metadata file. Records that require full re-write are be sync'ed at serviceMetaPeriod.
                                             partialServiceMetaPeriod specifies, in seconds, how frequently it should sync. Zero means that this feature is turned off and serviceMetaPeriod is the only time when metadata sync happens.
                                             If the value of partialServiceMetaPeriod is greater than serviceMetaPeriod, this setting has no effect.
                                             By default it is turned off (zero).
* `process_tracker_service_interval` - (Optional) Specifies, in seconds, how often the indexer checks the status of the child OS processes it launched to see if it can launch new processes for queued requests. Defaults to 15.
                                                  If set to 0, the indexer checks child process status every second.
                                                  Highest legal value is 4294967295.
* `quarantine_future_secs` - (Optional) Events with timestamp of quarantineFutureSecs newer than "now" are dropped into quarantine bucket. Defaults to 2592000 (30 days).
                                        This is a mechanism to prevent main hot buckets from being polluted with fringe events.
* `quarantine_past_secs` - (Optional) Events with timestamp of quarantinePastSecs older than "now" are dropped into quarantine bucket. Defaults to 77760000 (900 days). This is a mechanism to prevent the main hot buckets from being polluted with fringe events.
* `raw_chunk_size_bytes` - (Optional) Target uncompressed size in bytes for individual raw slice in the rawdata journal of the index. Defaults to 131072 (128KB). 0 is not a valid value. If 0 is specified, rawChunkSizeBytes is set to the default value.
* `rep_factor` - (Deprecated, Optional) Index replication control. This parameter applies to only clustering slaves.
                            auto = Use the master index replication configuration value.
                            0 = Turn off replication for this index.

  `rep_factor` is deprecated in this Terraform Provider.
				
  The REST API returns a 0 for both `repFactor = 0` and `repFactor = auto`. These are the two valid values for `repFactor`, yet they cannot be detected as different from the API's response.
				
  Additionally, `repFactor` only has meaning on clustered indexes, which should be configured by the Indexer Cluster Manager, not via REST.
* `rotate_period_in_secs` - (Optional) How frequently (in seconds) to check if a new hot bucket needs to be created. Also, how frequently to check if there are any warm/cold buckets that should be rolled/frozen.
* `service_meta_period` - (Optional) Defines how frequently metadata is synced to disk, in seconds. Defaults to 25 (seconds).
                                     You may want to set this to a higher value if the sum of your metadata file sizes is larger than many tens of megabytes, to avoid the hit on I/O in the indexing fast path.
* `sync_meta` - (Optional) When true, a sync operation is called before file descriptor is closed on metadata file updates. This functionality improves integrity of metadata files, especially in regards to operating system crashes/machine failures.
* `thawed_path` - (Optional) An absolute path that contains the thawed (resurrected) databases for the index.
                             Cannot be defined in terms of a volume definition.
                             Required. Splunk software does not start if an index lacks a valid thawedPath.
* `throttle_check_period` - (Optional) Defines how frequently Splunk software checks for index throttling condition, in seconds. Defaults to 15 (seconds).
* `tstats_home_path` - (Optional) Location to store datamodel acceleration TSIDX data for this index. Restart splunkd after changing this parameter.
                                  If specified, it must be defined in terms of a volume definition.
* `warm_to_cold_script` - (Optional)  Path to a script to run when moving data from warm to cold.
                                     This attribute is supported for backwards compatibility with Splunk software versions older than 4.0. Contact Splunk support if you need help configuring this setting.
* `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the http event collector resource
