package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"terraform-provider-splunk/client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func index() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"block_sign_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Controls how many events make up a block for block signatures.
				If this is set to 0, block signing is disabled for this index.
				A recommended value is 100.`,
			},
			"bucket_rebuild_memory_hint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `Suggestion for the bucket rebuild process for the size of the time-series (tsidx) file to make.
				Caution: This is an advanced parameter. Inappropriate use of this parameter causes splunkd to not start if rebuild is required. Do not set this parameter unless instructed by Splunk Support.

				Default value, auto, varies by the amount of physical RAM on the host

				- less than 2GB RAM = 67108864 (64MB) tsidx
				- 2GB to 8GB RAM = 134217728 (128MB) tsidx
				- more than 8GB RAM = 268435456 (256MB) tsidx

				Values other than "auto" must be 16MB-1GB. Highest legal value (of the numerical part) is 4294967295

				You can specify the value using a size suffix: "16777216" or "16MB" are equivalent.`,
			},
			"cold_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `An absolute path that contains the colddbs for the index. The path must be readable and writable. Cold databases are opened as needed when searching. May be defined in terms of a volume definition (see volume section below).
				Required. Splunk software does not start if an index lacks a valid coldPath.`,
			},
			"cold_to_frozen_dir": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `Destination path for the frozen archive. Use as an alternative to a coldToFrozenScript. Splunk software automatically puts frozen buckets in this directory.
				Bucket freezing policy is as follows:

				- New style buckets (4.2 and on): removes all files but the rawdata
				To thaw, run splunk rebuild <bucket dir> on the bucket, then move to the thawed directory
				- Old style buckets (Pre-4.2): gzip all the .data and .tsidx files
				To thaw, gunzip the zipped files and move the bucket into the thawed directory

				If both coldToFrozenDir and coldToFrozenScript are specified, coldToFrozenDir takes precedence`,
			},
			"cold_to_frozen_script": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `Path to the archiving script.
				If your script requires a program to run it (for example, python), specify the program followed by the path. The script must be in $SPLUNK_HOME/bin or one of its subdirectories.

				Splunk software ships with an example archiving script in $SPLUNK_HOME/bin called coldToFrozenExample.py. DO NOT use this example script directly. It uses a default path, and if modified in place any changes are overwritten on upgrade.

				It is best to copy the example script to a new file in bin and modify it for your system. Most importantly, change the default archive path to an existing directory that fits your needs.

				If your new script in bin/ is named myColdToFrozen.py, set this key to the following:

				coldToFrozenScript = "$SPLUNK_HOME/bin/python" "$SPLUNK_HOME/bin/myColdToFrozen.py"

				By default, the example script has two possible behaviors when archiving:

				- For buckets created from version 4.2 and on, it removes all files except for rawdata. To thaw: cd to the frozen bucket and type splunk rebuild ., then copy the bucket to thawed for that index. We recommend using the coldToFrozenDir parameter unless you need to perform a more advanced operation upon freezing buckets.
				- For older-style buckets, we simply gzip all the .tsidx files. To thaw: cd to the frozen bucket and unzip the tsidx files, then copy the bucket to thawed for that index`,
			},
			"compress_rawdata": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `This parameter is ignored. The splunkd process always compresses raw data.`,
			},
			"datatype": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"event", "metric"}, false),
				Description:  `Valid values: (event | metric). Specifies the type of index.`,
			},
			"enable_online_bucket_repair": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: `Enables asynchronous "online fsck" bucket repair, which runs concurrently with Splunk software.
				When enabled, you do not have to wait until buckets are repaired to start the Splunk platform. However, you might observe a slight performance degratation.`,
			},
			"frozen_time_period_in_secs": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Number of seconds after which indexed data rolls to frozen. Defaults to 188697600 (6 years).
				Freezing data means it is removed from the index. If you need to archive your data, refer to coldToFrozenDir and coldToFrozenScript parameter documentation.`,
			},
			"home_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `An absolute path that contains the hot and warm buckets for the index.
				Required. Splunk software does not start if an index lacks a valid homePath.

				Caution: The path must be readable and writable.`,
			},
			"max_bloom_backfill_bucket_age": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Valid values are: Integer[m|s|h|d]." +
					"If a warm or cold bucket is older than the specified age, do not create or rebuild its bloomfilter. Specify 0 to never rebuild bloomfilters.",
			},
			"max_concurrent_optimizes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `The number of concurrent optimize processes that can run against a hot bucket.
				This number should be increased if instructed by Splunk Support. Typically the default value should suffice.`,
			},
			"max_data_size": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `The maximum size in MB for a hot DB to reach before a roll to warm is triggered. Specifying "auto" or "auto_high_volume" causes Splunk software to autotune this parameter (recommended).Use "auto_high_volume" for high volume indexes (such as the main index); otherwise, use "auto". A "high volume index" would typically be considered one that gets over 10GB of data per day.
				- "auto" sets the size to 750MB.
				- "auto_high_volume" sets the size to 10GB on 64-bit, and 1GB on 32-bit systems.

				Although the maximum value you can set this is 1048576 MB, which corresponds to 1 TB, a reasonable number ranges anywhere from 100 - 50000. Any number outside this range should be approved by Splunk Support before proceeding.

				If you specify an invalid number or string, maxDataSize is auto-tuned.

				Note: The precise size of your warm buckets may vary from maxDataSize, due to post-processing and timing issues with the rolling policy.`,
			},
			"max_hot_buckets": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Maximum hot buckets that can exist per index. Defaults to 3.
				When maxHotBuckets is exceeded, Splunk software rolls the least recently used (LRU) hot bucket to warm. Both normal hot buckets and quarantined hot buckets count towards this total. This setting operates independently of maxHotIdleSecs, which can also cause hot buckets to roll.`,
			},
			"max_hot_idle_secs": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Maximum life, in seconds, of a hot bucket. Defaults to 0.
				If a hot bucket exceeds maxHotIdleSecs, Splunk software rolls it to warm. This setting operates independently of maxHotBuckets, which can also cause hot buckets to roll. A value of 0 turns off the idle check (equivalent to INFINITE idle time).`,
			},
			"max_hot_span_secs": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Upper bound of target maximum timespan of hot/warm buckets in seconds. Defaults to 7776000 seconds (90 days).
				Note: If you set this too small, you can get an explosion of hot/warm buckets in the filesystem. The system sets a lower bound implicitly for this parameter at 3600, but this is an advanced parameter that should be set with care and understanding of the characteristics of your data.`,
			},
			"max_mem_mb": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `The amount of memory, expressed in MB, to allocate for buffering a single tsidx file into memory before flushing to disk. Defaults to 5. The default is recommended for all environments.
				IMPORTANT: Calculate this number carefully. Setting this number incorrectly may have adverse effects on your systems memory and/or splunkd stability/performance.`,
			},
			"max_meta_entries": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Sets the maximum number of unique lines in .data files in a bucket, which may help to reduce memory consumption. If set to 0, this setting is ignored (it is treated as infinite).
				If exceeded, a hot bucket is rolled to prevent further increase. If your buckets are rolling due to Strings.data hitting this limit, the culprit may be the punct field in your data. If you do not use punct, it may be best to simply disable this (see props.conf.spec in $SPLUNK_HOME/etc/system/README).

				There is a small time delta between when maximum is exceeded and bucket is rolled. This means a bucket may end up with epsilon more lines than specified, but this is not a major concern unless excess is significant.`,
			},
			"max_time_unreplicated_no_acks": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Upper limit, in seconds, on how long an event can sit in raw slice. Applies only if replication is enabled for this index. Otherwise ignored.
				If there are any acknowledged events sharing this raw slice, this paramater does not apply. In this case, maxTimeUnreplicatedWithAcks applies.

				Highest legal value is 2147483647. To disable this parameter, set to 0.

				Note: this is an advanced parameter. Understand the consequences before changing.`,
			},
			"max_time_unreplicated_with_acks": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Upper limit, in seconds, on how long events can sit unacknowledged in a raw slice. Applies only if you have enabled acks on forwarders and have replication enabled (with clustering).
				Note: This is an advanced parameter. Make sure you understand the settings on all forwarders before changing this. This number should not exceed ack timeout configured on any forwarder, and should actually be set to at most half of the minimum value of that timeout. You can find this setting in outputs.conf readTimeout setting under the tcpout stanza.

				To disable, set to 0, but this is NOT recommended. Highest legal value is 2147483647.`,
			},
			"max_total_data_size_mb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The maximum size of an index (in MB). If an index grows larger than the maximum size, the oldest data is frozen.`,
			},
			"max_warm_db_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The maximum number of warm buckets. If this number is exceeded, the warm bucket/s with the lowest value for their latest times is moved to cold.`,
			},
			"min_raw_file_sync_secs": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `Specify an integer (or "disable") for this parameter.
				This parameter sets how frequently splunkd forces a filesystem sync while compressing journal slices.

				During this period, uncompressed slices are left on disk even after they are compressed. Then splunkd forces a filesystem sync of the compressed journal and removes the accumulated uncompressed files.

				If 0 is specified, splunkd forces a filesystem sync after every slice completes compressing. Specifying "disable" disables syncing entirely: uncompressed slices are removed as soon as compression is complete.

				Note: Some filesystems are very inefficient at performing sync operations, so only enable this if you are sure it is needed`,
			},
			"min_stream_group_queue_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Minimum size of the queue that stores events in memory before committing them to a tsidx file.
				Caution: Do not set this value, except under advice from Splunk Support.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the index to create.`,
			},
			"partial_service_meta_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Related to serviceMetaPeriod. If set, it enables metadata sync every <integer> seconds, but only for records where the sync can be done efficiently in-place, without requiring a full re-write of the metadata file. Records that require full re-write are be sync'ed at serviceMetaPeriod.
				partialServiceMetaPeriod specifies, in seconds, how frequently it should sync. Zero means that this feature is turned off and serviceMetaPeriod is the only time when metadata sync happens.

				If the value of partialServiceMetaPeriod is greater than serviceMetaPeriod, this setting has no effect.

				By default it is turned off (zero).`,
			},
			"process_tracker_service_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Specifies, in seconds, how often the indexer checks the status of the child OS processes it launched to see if it can launch new processes for queued requests. Defaults to 15.
				If set to 0, the indexer checks child process status every second.

				Highest legal value is 4294967295.`,
			},
			"quarantine_future_secs": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Events with timestamp of quarantineFutureSecs newer than "now" are dropped into quarantine bucket. Defaults to 2592000 (30 days).
				This is a mechanism to prevent main hot buckets from being polluted with fringe events.`,
			},
			"quarantine_past_secs": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Events with timestamp of quarantinePastSecs older than "now" are dropped into quarantine bucket. Defaults to 77760000 (900 days).
				This is a mechanism to prevent the main hot buckets from being polluted with fringe events.`,
			},
			"raw_chunk_size_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Target uncompressed size in bytes for individual raw slice in the rawdata journal of the index. Defaults to 131072 (128KB). 0 is not a valid value. If 0 is specified, rawChunkSizeBytes is set to the default value.
				Note: rawChunkSizeBytes only specifies a target chunk size. The actual chunk size may be slightly larger by an amount proportional to an individual event size.

				WARNING: This is an advanced parameter. Only change it if you are instructed to do so by Splunk Support.`,
			},
			"rep_factor": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `Index replication control. This parameter applies to only clustering slaves.
				auto = Use the master index replication configuration value.

				0 = Turn off replication for this index.`,
			},
			"rotate_period_in_secs": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `How frequently (in seconds) to check if a new hot bucket needs to be created. Also, how frequently to check if there are any warm/cold buckets that should be rolled/frozen.`,
			},
			"service_meta_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Defines how frequently metadata is synced to disk, in seconds. Defaults to 25 (seconds).
				You may want to set this to a higher value if the sum of your metadata file sizes is larger than many tens of megabytes, to avoid the hit on I/O in the indexing fast path.`,
			},
			"sync_meta": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: `When true, a sync operation is called before file descriptor is closed on metadata file updates. This functionality improves integrity of metadata files, especially in regards to operating system crashes/machine failures.
				Note: Do not change this parameter without the input of a Splunk Support.`,
			},
			"thawed_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `An absolute path that contains the thawed (resurrected) databases for the index.
				Cannot be defined in terms of a volume definition.

				Required. Splunk software does not start if an index lacks a valid thawedPath.`,
			},
			"throttle_check_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Defines how frequently Splunk software checks for index throttling condition, in seconds. Defaults to 15 (seconds).
				Note: Do not change this parameter without the input of Splunk Support.`,
			},
			"tstats_home_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `Location to store datamodel acceleration TSIDX data for this index. Restart splunkd after changing this parameter.
				If specified, it must be defined in terms of a volume definition.

				Caution: Path must be writable.

				Default value: volume:_splunk_summaries/$_index_name/tstats`,
			},
			"warm_to_cold_script": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `Path to a script to run when moving data from warm to cold.
				This attribute is supported for backwards compatibility with Splunk software versions older than 4.0. Contact Splunk support if you need help configuring this setting.

				Caution: Migrating data across filesystems is now handled natively by splunkd. If you specify a script here, the script becomes responsible for moving the event data, and Splunk-native data migration is not used.`,
			},
			"acl": aclSchema(),
		},
		Read:   indexRead,
		Create: indexCreate,
		Delete: indexDelete,
		Update: indexUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func indexCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	indexConfigObj := getIndexConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.Owner = "nobody"
		aclObject.App = "system"
	}
	err := (*provider.Client).CreateIndexObject(name, aclObject.Owner, aclObject.App, indexConfigObj)
	if err != nil {
		return err
	}
	if _, ok := d.GetOk("acl"); ok {
		err := (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, name, aclObject, "data", "indexes")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return indexRead(d, meta)
}

func indexRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	// We first get list of indexes to get owner and app name for the specific index
	resp, err := (*provider.Client).ReadAllIndexObject()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getIndexConfigByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %s", name)
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadIndexObject(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getIndexConfigByName(name, resp)
	if err != nil {
		return err
	}

	if err = d.Set("block_sign_size", entry.Content.BlockSignSize); err != nil {
		return err
	}

	if err = d.Set("bucket_rebuild_memory_hint", entry.Content.BucketRebuildMemoryHint); err != nil {
		return err
	}

	if err = d.Set("cold_path", entry.Content.ColdPath); err != nil {
		return err
	}

	if err = d.Set("cold_to_frozen_dir", entry.Content.ColdToFrozenDir); err != nil {
		return err
	}

	if err = d.Set("cold_to_frozen_script", entry.Content.ColdToFrozenScript); err != nil {
		return err
	}

	if err = d.Set("compress_rawdata", entry.Content.CompressRawdata); err != nil {
		return err
	}

	if err = d.Set("datatype", entry.Content.Datatype); err != nil {
		return err
	}

	if err = d.Set("enable_online_bucket_repair", entry.Content.EnableOnlineBucketRepair); err != nil {
		return err
	}

	if err = d.Set("frozen_time_period_in_secs", entry.Content.FrozenTimePeriodInSecs); err != nil {
		return err
	}

	if err = d.Set("home_path", entry.Content.HomePath); err != nil {
		return err
	}

	if err = d.Set("max_bloom_backfill_bucket_age", entry.Content.MaxBloomBackfillBucketAge); err != nil {
		return err
	}

	if err = d.Set("max_concurrent_optimizes", entry.Content.MaxConcurrentOptimizes); err != nil {
		return err
	}

	if err = d.Set("max_data_size", entry.Content.MaxDataSize); err != nil {
		return err
	}

	if err = d.Set("max_hot_buckets", entry.Content.MaxHotBuckets); err != nil {
		return err
	}

	if err = d.Set("max_hot_idle_secs", entry.Content.MaxHotIdleSecs); err != nil {
		return err
	}

	if err = d.Set("max_hot_span_secs", entry.Content.MaxHotSpanSecs); err != nil {
		return err
	}

	if err = d.Set("max_mem_mb", entry.Content.MaxMemMB); err != nil {
		return err
	}

	if err = d.Set("max_meta_entries", entry.Content.MaxMetaEntries); err != nil {
		return err
	}

	if err = d.Set("max_time_unreplicated_no_acks", entry.Content.MaxTimeUnreplicatedNoAcks); err != nil {
		return err
	}

	if err = d.Set("max_time_unreplicated_with_acks", entry.Content.MaxTimeUnreplicatedWithAcks); err != nil {
		return err
	}

	if err = d.Set("max_total_data_size_mb", entry.Content.MaxTotalDataSizeMB); err != nil {
		return err
	}

	if err = d.Set("max_warm_db_count", entry.Content.MaxWarmDBCount); err != nil {
		return err
	}

	if err = d.Set("min_raw_file_sync_secs", entry.Content.MinRawFileSyncSecs); err != nil {
		return err
	}

	if err = d.Set("min_stream_group_queue_size", entry.Content.MinStreamGroupQueueSize); err != nil {
		return err
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}

	if err = d.Set("partial_service_meta_period", entry.Content.PartialServiceMetaPeriod); err != nil {
		return err
	}

	if err = d.Set("process_tracker_service_interval", entry.Content.ProcessTrackerServiceInterval); err != nil {
		return err
	}

	if err = d.Set("quarantine_future_secs", entry.Content.QuarantineFutureSecs); err != nil {
		return err
	}

	if err = d.Set("quarantine_past_secs", entry.Content.QuarantinePastSecs); err != nil {
		return err
	}

	if err = d.Set("raw_chunk_size_bytes", entry.Content.RawChunkSizeBytes); err != nil {
		return err
	}

	if err = d.Set("rep_factor", entry.Content.RepFactor); err != nil {
		return err
	}

	if err = d.Set("rotate_period_in_secs", entry.Content.RotatePeriodInSecs); err != nil {
		return err
	}

	if err = d.Set("service_meta_period", entry.Content.ServiceMetaPeriod); err != nil {
		return err
	}

	if err = d.Set("sync_meta", entry.Content.SyncMeta); err != nil {
		return err
	}

	if err = d.Set("thawed_path", entry.Content.ThawedPath); err != nil {
		return err
	}

	if err = d.Set("throttle_check_period", entry.Content.ThrottleCheckPeriod); err != nil {
		return err
	}

	if err = d.Set("tstats_home_path", entry.Content.TstatsHomePath); err != nil {
		return err
	}

	if err = d.Set("warm_to_cold_script", entry.Content.WarmToColdScript); err != nil {
		return err
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func indexUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	indexConfigObj := getIndexConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	err := (*provider.Client).UpdateIndexObject(d.Id(), aclObject.Owner, aclObject.App, indexConfigObj)
	if err != nil {
		return err
	}

	//ACL update
	err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, d.Id(), aclObject, "data", "indexes")
	if err != nil {
		return err
	}

	return indexRead(d, meta)
}

func indexDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).DeleteIndexObject(d.Id(), aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.IndexResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getIndexConfig(d *schema.ResourceData) (indexConfigObject *models.IndexObject) {
	indexConfigObject = &models.IndexObject{}
	indexConfigObject.BlockSignSize = d.Get("block_sign_size").(int)
	indexConfigObject.BucketRebuildMemoryHint = d.Get("bucket_rebuild_memory_hint").(string)
	indexConfigObject.ColdPath = d.Get("cold_path").(string)
	indexConfigObject.ColdToFrozenDir = d.Get("cold_to_frozen_dir").(string)
	indexConfigObject.ColdToFrozenScript = d.Get("cold_to_frozen_script").(string)
	indexConfigObject.CompressRawdata = d.Get("compress_rawdata").(bool)
	indexConfigObject.Datatype = d.Get("datatype").(string)
	indexConfigObject.EnableOnlineBucketRepair = d.Get("enable_online_bucket_repair").(bool)
	indexConfigObject.FrozenTimePeriodInSecs = d.Get("frozen_time_period_in_secs").(int)
	indexConfigObject.HomePath = d.Get("home_path").(string)
	indexConfigObject.MaxBloomBackfillBucketAge = d.Get("max_bloom_backfill_bucket_age").(string)
	indexConfigObject.MaxConcurrentOptimizes = d.Get("max_concurrent_optimizes").(int)
	indexConfigObject.MaxDataSize = d.Get("max_data_size").(string)
	indexConfigObject.MaxHotBuckets = d.Get("max_hot_buckets").(int)
	indexConfigObject.MaxHotIdleSecs = d.Get("max_hot_idle_secs").(int)
	indexConfigObject.MaxHotSpanSecs = d.Get("max_hot_span_secs").(int)
	indexConfigObject.MaxMemMB = d.Get("max_mem_mb").(int)
	indexConfigObject.MaxMetaEntries = d.Get("max_meta_entries").(int)
	indexConfigObject.MaxTimeUnreplicatedNoAcks = d.Get("max_time_unreplicated_no_acks").(int)
	indexConfigObject.MaxTimeUnreplicatedWithAcks = d.Get("max_time_unreplicated_with_acks").(int)
	indexConfigObject.MaxTotalDataSizeMB = d.Get("max_total_data_size_mb").(int)
	indexConfigObject.MaxWarmDBCount = d.Get("max_warm_db_count").(int)
	indexConfigObject.MinRawFileSyncSecs = d.Get("min_raw_file_sync_secs").(string)
	indexConfigObject.MinStreamGroupQueueSize = d.Get("min_stream_group_queue_size").(int)
	indexConfigObject.PartialServiceMetaPeriod = d.Get("partial_service_meta_period").(int)
	indexConfigObject.ProcessTrackerServiceInterval = d.Get("process_tracker_service_interval").(int)
	indexConfigObject.QuarantineFutureSecs = d.Get("quarantine_future_secs").(int)
	indexConfigObject.QuarantinePastSecs = d.Get("quarantine_past_secs").(int)
	indexConfigObject.RawChunkSizeBytes = d.Get("raw_chunk_size_bytes").(int)
	indexConfigObject.RepFactor = d.Get("rep_factor").(string)
	indexConfigObject.RotatePeriodInSecs = d.Get("rotate_period_in_secs").(int)
	indexConfigObject.ServiceMetaPeriod = d.Get("service_meta_period").(int)
	indexConfigObject.SyncMeta = d.Get("sync_meta").(bool)
	indexConfigObject.ThawedPath = d.Get("thawed_path").(string)
	indexConfigObject.ThrottleCheckPeriod = d.Get("throttle_check_period").(int)
	indexConfigObject.TstatsHomePath = d.Get("tstats_home_path").(string)
	indexConfigObject.WarmToColdScript = d.Get("warm_to_cold_script").(string)
	return indexConfigObject
}

func getIndexConfigByName(name string, httpResponse *http.Response) (indexEntry *models.IndexEntry, err error) {
	response := &models.IndexResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:
		_ = json.NewDecoder(httpResponse.Body).Decode(&response)
		re := regexp.MustCompile(`(.*)`)
		for _, entry := range response.Entry {
			if name == re.FindStringSubmatch(entry.Name)[1] {
				return &entry, nil
			}
		}

	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return indexEntry, err
	}

	return indexEntry, nil
}
