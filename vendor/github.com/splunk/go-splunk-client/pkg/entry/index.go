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
	ColdToFrozenDir        attributes.Explicit[string] `json:"coldToFrozenDir"        values:"coldToFrozenDir,omitzero"`
	ColdToFrozenScript     attributes.Explicit[string] `json:"coldToFrozenScript"     values:"coldToFrozenScript,omitzero"`
	DataType               attributes.Explicit[string] `json:"datatype"               values:"datatype,omitzero"               selective:"create"`
	Disabled               attributes.Explicit[bool]   `json:"disabled"               values:"disabled,omitzero"               selective:"read"`
	FrozenTimePeriodInSecs attributes.Explicit[int]    `json:"frozenTimePeriodInSecs" values:"frozenTimePeriodInSecs,omitzero"`
	HomePath               attributes.Explicit[string] `json:"homePath"               values:"homePath,omitzero"               selective:"create"`
	MaxDataSize            attributes.Explicit[string] `json:"maxDataSize"            values:"maxDataSize,omitzero"`
	MaxHotBuckets          attributes.Explicit[string] `json:"maxHotBuckets"          values:"maxHotBuckets,omitzero"`
	MaxHotIdleSecs         attributes.Explicit[int]    `json:"maxHotIdleSecs"         values:"maxHotIdleSecs,omitzero"`
	MaxHotSpanSecs         attributes.Explicit[int]    `json:"maxHotSpanSecs"         values:"maxHotSpanSecs,omitzero"`
	MaxTotalDataSizeMB     attributes.Explicit[int]    `json:"maxTotalDataSizeMB"     values:"maxTotalDataSizeMB,omitzero"`
	MaxWarmDBCount         attributes.Explicit[int]    `json:"maxWarmDBCount"         values:"maxWarmDBCount,omitzero"`
	QuarantineFutureSecs   attributes.Explicit[int]    `json:"quarantineFutureSecs"   values:"quarantineFutureSecs,omitzero"`
	QuarantinePastSecs     attributes.Explicit[int]    `json:"quarantinePastSecs"     values:"quarantinePastSecs,omitzero"`
	ThawedPath             attributes.Explicit[string] `json:"thawedPath"             values:"thawedPath,omitzero"            selective:"create"`
}

// Index is a Splunk Index.
type Index struct {
	ID      client.ID    `selective:"create" service:"data/indexes"`
	Content IndexContent `json:"content" values:",anonymize"`
}
