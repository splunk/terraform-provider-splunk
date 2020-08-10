package models

type IndexResponse struct {
	Entry    []IndexEntry     `json:"entry"`
	Messages []ErrorMessage `json:"messages"`
}

type IndexEntry struct {
	Name    string                   `json:"name"`
	ACL     ACLObject                `json:"acl"`
	Content IndexObject `json:"content"`
}

type IndexObject struct {
	BlockSignSize                 int    `json:"block_sign_size,omitempty" url:"block_sign_size,omitempty"`
	BucketRebuildMemoryHint       string `json:"bucket_rebuild_memory_hint,omitempty" url:"bucket_rebuild_memory_hint,omitempty"`
	ColdPath                      string `json:"cold_path,omitempty" url:"cold_path,omitempty"`
	ColdToFrozenDir               string `json:"cold_to_frozen_dir,omitempty" url:"cold_to_frozen_dir,omitempty"`
	ColdToFrozenScript            string `json:"cold_to_frozen_script,omitempty" url:"cold_to_frozen_script,omitempty"`
	CompressRawdata               bool   `json:"compress_rawdata,omitempty" url:"compress_rawdata,omitempty"`
	Datatype                      string `json:"datatype,omitempty" url:"datatype,omitempty"`
	EnableOnlineBucketRepair      bool   `json:"enable_online_bucket_repair,omitempty" url:"enable_online_bucket_repair,omitempty"`
	FrozenTimePeriodInSecs        int    `json:"frozen_time_period_in_secs,omitempty" url:"frozeny_time_period_in_secs,omitempty"`
	HomePath                      string `json:"home_path,omitempty" url:"home_path,omitempty"`
	MaxBloomBackfillBucketAge     string `json:"max_bloom_backfill_bucket_age,omitempty" url:"max_bloom_backfill_bucket_age,omitempty"`
	MaxConcurrentOptimizes        int    `json:"max_concurrent_optimizes,omitempty" url:"max_concurrent_optimizes,omitempty"`
	MaxDataSize                   string `json:"max_data_size,omitempty" url:"max_data_size,omitempty"`
	MaxHotBuckets                 int    `json:"max_hot_buckets,omitempty" url:"max_hot_buckets,omitempty"`
	MaxHotIdleSecs                int    `json:"max_hot_idle_secs,omitempty" url:"max_hot_idle_secs,omitempty"`
	MaxHotSpanSecs                int    `json:"max_hot_span_secs,omitempty" url:"max_hot_span_secs,omitempty"`
	MaxMemMB                      int    `json:"max_mem_mb,omitempty" url:"max_mem_mb,omitempty"`
	MaxMetaEntries                int    `json:"max_meta_entries,omitempty" url:"max_meta_entries,omitempty"`
	MaxTimeUnreplicatedNoAcks     int    `json:"max_time_unreplicated_no_acks,omitempty" url:"max_time_unreplicated_no_acks,omitempty"`
	MaxTimeUnreplicatedWithAcks   int    `json:"max_time_unreplicated_with_acks,omitempty" url:"max_time_unreplicated_with_acks,omitempty"`
	MaxTotalDataSizeMB            int    `json:"max_total_data_size_mb,omitempty" url:"max_total_data_size_mb,omitempty"`
	MaxWarmDBCount                int    `json:"max_warm_db_count,omitempty" url:"max_warm_db_count,omitempty"`
	MinRawFileSyncSecs            string `json:"min_raw_file_sync_secs,omitempty" url:"min_raw_file_sync_secs,omitempty"`
	MinStreamGroupQueueSize       int    `json:"min_stream_group_queue_size,omitempty" url:"min_stream_group_queue_size,omitempty"`
	Name                          string `json:"name,omitempty" url:"name,omitempty"`
	PartialServiceMetaPeriod      int    `json:"partial_service_meta_period,omitempty" url:"partial_service_meta_period,omitempty"`
	ProcessTrackerServiceInterval int    `json:"process_tracker_service_interval,omitempty" url:"process_tracker_service_interval,omitempty"`
	QuarantineFutureSecs          int    `json:"quarantine_future_secs,omitempty" url:"quarantine_future_secs,omitempty"`
	QuarantinePastSecs            int    `json:"quarantine_past_secs,omitempty" url:"quarantine_past_secs,omitempty"`
	RawChunkSizeBytes             int    `json:"raw_chunk_size_bytes,omitempty" url:"raw_chunk_size_bytes,omitempty"`
	RepFactor                     string `json:"rep_factor,omitempty" url:"rep_factor,omitempty"`
	RotatePeriodInSecs            int    `json:"rotate_period_in_secs,omitempty" url:"rotate_period_in_secs,omitempty"`
	ServiceMetaPeriod             int    `json:"service_meta_period,omitempty" url:"service_meta_period,omitempty"`
	SyncMeta                      bool   `json:"sync_meta,omitempty" url:"sync_meta,omitempty"`
	ThawedPath                    string `json:"thawed_path,omitempty" url:"thawed_path,omitempty"`
	ThrottleCheckPeriod           int    `json:"throttle_check_period,omitempty" url:"throttle_check_period,omitempty"`
	TstatsHomePath                string `json:"tstats_home_path,omitempty" url:"tstats_home_path,omitempty"`
	WarmToColdScript              string `json:"warm_to_cold_script,omitempty" url:"warm_to_cold_script,omitempty"`
}
