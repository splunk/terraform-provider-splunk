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
	BlockSignSize                 int    `json:"blockSignSize,omitempty" url:"blockSignSize,omitempty"`
	BucketRebuildMemoryHint       string `json:"bucketRebuildMemoryHint,omitempty" url:"bucketRebuildMemoryHint,omitempty"`
	ColdPath                      string `json:"coldPath,omitempty" url:"coldPath,omitempty"`
	ColdToFrozenDir               string `json:"coldToFrozenDir,omitempty" url:"coldToFrozenDir,omitempty"`
	ColdToFrozenScript            string `json:"coldToFrozenScript,omitempty" url:"coldToFrozenScript,omitempty"`
	CompressRawdata               bool   `json:"compressRawdata,omitempty" url:"compressRawdata,omitempty"`
	Datatype                      string `json:"datatype,omitempty" url:"datatype,omitempty"`
	EnableOnlineBucketRepair      bool   `json:"enableOnlineBucketRepair,omitempty" url:"enableOnlineBucketRepair,omitempty"`
	FrozenTimePeriodInSecs        int    `json:"frozenTimePeriodInSecs,omitempty" url:"frozenTimePeriodInSecs,omitempty"`
	HomePath                      string `json:"homePath,omitempty" url:"homePath,omitempty"`
	MaxBloomBackfillBucketAge     string `json:"maxBloomBackfillBucketAge,omitempty" url:"maxBloomBackfillBucketAge,omitempty"`
	MaxConcurrentOptimizes        int    `json:"maxConcurrentOptimizes,omitempty" url:"maxConcurrentOptimizes,omitempty"`
	MaxDataSize                   string `json:"maxDataSize,omitempty" url:"maxDataSize,omitempty"`
	MaxHotBuckets                 int    `json:"maxHotBuckets,omitempty" url:"maxHotBuckets,omitempty"`
	MaxHotIdleSecs                int    `json:"maxHotIdleSecs,omitempty" url:"maxHotIdleSecs,omitempty"`
	MaxHotSpanSecs                int    `json:"maxHotSpanSecs,omitempty" url:"maxHotSpanSecs,omitempty"`
	MaxMemMB                      int    `json:"maxMemMB,omitempty" url:"maxMemMB,omitempty"`
	MaxMetaEntries                int    `json:"maxMetaEntries,omitempty" url:"maxMetaEntries,omitempty"`
	MaxTimeUnreplicatedNoAcks     int    `json:"maxTimeUnreplicatedNoAcks,omitempty" url:"maxTimeUnreplicatedNoAcks,omitempty"`
	MaxTimeUnreplicatedWithAcks   int    `json:"maxTimeUnreplicatedWithAcks,omitempty" url:"maxTimeUnreplicatedWithAcks,omitempty"`
	MaxTotalDataSizeMB            int    `json:"maxTotalDataSizeMB,omitempty" url:"maxTotalDataSizeMB,omitempty"`
	MaxWarmDBCount                int    `json:"maxWarmDBCount,omitempty" url:"maxWarmDBCount,omitempty"`
	MinRawFileSyncSecs            string `json:"minRawFileSyncSecs,omitempty" url:"minRawFileSyncSecs,omitempty"`
	MinStreamGroupQueueSize       int    `json:"minStreamGroupQueueSize,omitempty" url:"minStreamGroupQueueSize,omitempty"`
	PartialServiceMetaPeriod      int    `json:"partialServiceMetaPeriod,omitempty" url:"partialServiceMetaPeriod,omitempty"`
	ProcessTrackerServiceInterval int    `json:"processTrackerServiceInterval,omitempty" url:"processTrackerServiceInterval,omitempty"`
	QuarantineFutureSecs          int    `json:"quarantineFutureSecs,omitempty" url:"quarantineFutureSecs,omitempty"`
	QuarantinePastSecs            int    `json:"quarantinePastSecs,omitempty" url:"quarantinePastSecs,omitempty"`
	RawChunkSizeBytes             int    `json:"rawChunkSizeBytes,omitempty" url:"rawChunkSizeBytes,omitempty"`
	RepFactor                     string `json:"repFactor,omitempty" url:"repFactor,omitempty"`
	RotatePeriodInSecs            int    `json:"rotatePeriodInSecs,omitempty" url:"rotatePeriodInSecs,omitempty"`
	ServiceMetaPeriod             int    `json:"serviceMetaPeriod,omitempty" url:"serviceMetaPeriod,omitempty"`
	SyncMeta                      bool   `json:"syncMeta,omitempty" url:"syncMeta,omitempty"`
	ThawedPath                    string `json:"thawedPath,omitempty" url:"thawedPath,omitempty"`
	ThrottleCheckPeriod           int    `json:"throttleCheckPeriod,omitempty" url:"throttleCheckPeriod,omitempty"`
	TstatsHomePath                string `json:"tstatsHomePath,omitempty" url:"tstatsHomePath,omitempty"`
	WarmToColdScript              string `json:"warmToColdScript,omitempty" url:"warmToColdScript,omitempty"`
}
