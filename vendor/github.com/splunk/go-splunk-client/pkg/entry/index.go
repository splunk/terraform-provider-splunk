// Copyright 2022 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package entry

import (
	"github.com/splunk/go-splunk-client/pkg/attributes"
	"github.com/splunk/go-splunk-client/pkg/client"
)

// IndexContent is the content for an Index.
type IndexContent struct {
	BlockSignSize                 attributes.Explicit[int]    `json:"blockSignSize"                 values:"blockSignSize,omitzero"`
	BucketRebuildMemoryHint       attributes.Explicit[int]    `json:"bucketRebuildMemoryHint"       values:"bucketRebuildMemoryHint,omitzero"`
	ColdPath                      attributes.Explicit[string] `json:"coldPath"                      values:"coldPath,omitzero"                      selective:"create"`
	ColdToFrozenDir               attributes.Explicit[string] `json:"coldToFrozenDir"               values:"coldToFrozenDir,omitzero"`
	ColdToFrozenScript            attributes.Explicit[string] `json:"coldToFrozenScript"            values:"coldToFrozenScript,omitzero"`
	CompressRawdata               attributes.Explicit[bool]   `json:"compressRawdata"               values:"compressRawdata,omitzero"`
	DataType                      attributes.Explicit[string] `json:"datatype"                      values:"datatype,omitzero"                      selective:"create"`
	Disabled                      attributes.Explicit[bool]   `json:"disabled"                      values:"disabled,omitzero"                      selective:"read"`
	EnableOnlineBucketRepair      attributes.Explicit[bool]   `json:"enableOnlineBucketRepair"      values:"enableOnlineBucketRepair,omitzero"`
	FrozenTimePeriodInSecs        attributes.Explicit[int]    `json:"frozenTimePeriodInSecs"        values:"frozenTimePeriodInSecs,omitzero"`
	HomePath                      attributes.Explicit[string] `json:"homePath"                      values:"homePath,omitzero"                      selective:"create"`
	MaxBloomBackfillBucketAge     attributes.Explicit[string] `json:"maxBloomBackfillBucketAge"     values:"maxBloomBackfillBucketAge,omitzero"`
	MaxConcurrentOptimizes        attributes.Explicit[int]    `json:"maxConcurrentOptimizes"        values:"maxConcurrentOptimizes,omitzero"`
	MaxDataSize                   attributes.Explicit[string] `json:"maxDataSize"                   values:"maxDataSize,omitzero"`
	MaxHotBuckets                 attributes.Explicit[string] `json:"maxHotBuckets"                 values:"maxHotBuckets,omitzero"`
	MaxHotIdleSecs                attributes.Explicit[int]    `json:"maxHotIdleSecs"                values:"maxHotIdleSecs,omitzero"`
	MaxHotSpanSecs                attributes.Explicit[int]    `json:"maxHotSpanSecs"                values:"maxHotSpanSecs,omitzero"`
	MaxMemMB                      attributes.Explicit[int]    `json:"maxMemMB"                      values:"maxMemMB,omitzero"`
	MaxMetaEntries                attributes.Explicit[int]    `json:"maxMetaEntries"                values:"maxMetaEntries,omitzero"`
	MaxTimeUnreplicatedNoAcks     attributes.Explicit[int]    `json:"maxTimeUnreplicatedNoAcks"     values:"maxTimeUnreplicatedNoAcks,omitzero"`
	MaxTimeUnreplicatedWithAcks   attributes.Explicit[int]    `json:"maxTimeUnreplicatedWithAcks"   values:"maxTimeUnreplicatedWithAcks,omitzero"`
	MaxTotalDataSizeMB            attributes.Explicit[int]    `json:"maxTotalDataSizeMB"            values:"maxTotalDataSizeMB,omitzero"`
	MaxWarmDBCount                attributes.Explicit[int]    `json:"maxWarmDBCount"                values:"maxWarmDBCount,omitzero"`
	MinRawFileSyncSecs            attributes.Explicit[string] `json:"minRawFileSyncSecs"            values:"minRawFileSyncSecs,omitzero"`
	MinStreamGroupQueueSize       attributes.Explicit[int]    `json:"minStreamGroupQueueSize"       values:"minStreamGroupQueueSize,omitzero"`
	PartialServiceMetaPeriod      attributes.Explicit[int]    `json:"partialServiceMetaPeriod"      values:"partialServiceMetaPeriod,omitzero"`
	ProcessTrackerServiceInterval attributes.Explicit[int]    `json:"processTrackerServiceInterval" values:"processTrackerServiceInterval,omitzero"`
	QuarantineFutureSecs          attributes.Explicit[int]    `json:"quarantineFutureSecs"          values:"quarantineFutureSecs,omitzero"`
	QuarantinePastSecs            attributes.Explicit[int]    `json:"quarantinePastSecs"            values:"quarantinePastSecs,omitzero"`
	RawChunkSizeBytes             attributes.Explicit[int]    `json:"rawChunkSizeBytes"             values:"rawChunkSizeBytes,omitzero"`
	RepFactor                     attributes.Explicit[int]    `json:"repFactor"                     values:"repFactor,omitzero"`
	RotatePeriodInSecs            attributes.Explicit[int]    `json:"rotatePeriodInSecs"            values:"rotatePeriodInSecs,omitzero"`
	ServiceMetaPeriod             attributes.Explicit[int]    `json:"serviceMetaPeriod"             values:"serviceMetaPeriod,omitzero"`
	SyncMeta                      attributes.Explicit[bool]   `json:"syncMeta"                      values:"syncMeta,omitzero"`
	ThawedPath                    attributes.Explicit[string] `json:"thawedPath"                    values:"thawedPath,omitzero"                    selective:"create"`
	ThrottleCheckPeriod           attributes.Explicit[int]    `json:"throttleCheckPeriod"           values:"throttleCheckPeriod,omitzero"`
	TstatsHomePath                attributes.Explicit[string] `json:"tstatsHomePath"                values:"tstatsHomePath,omitzero"                selective:"create"`
	WarmToColdScript              attributes.Explicit[string] `json:"warmToColdScript"              values:"warmToColdScript,omitzero"`
}

// Index is a Splunk Index.
type Index struct {
	ID      client.ID    `selective:"create" service:"data/indexes"`
	Content IndexContent `json:"content" values:",anonymize"`
}
