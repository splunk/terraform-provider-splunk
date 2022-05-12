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

// RoleContent defines the Content for a Role.
type RoleContent struct {
	Capabilities              []string                    `json:"capabilities"              values:"capabilities,omitzero,fillempty"`
	CumulativeRTSrchJobsQuota attributes.Explicit[int]    `json:"cumulativeRTSrchJobsQuota" values:"cumulativeRTSrchJobsQuota,omitzero"`
	CumulativeSrchJobsQuota   attributes.Explicit[int]    `json:"cumulativeSrchJobsQuota"   values:"cumulativeSrchJobsQuota,omitzero"`
	DefaultApp                attributes.Explicit[string] `json:"defaultApp"                values:"defaultApp,omitzero"`
	RtSrchJobsQuota           attributes.Explicit[int]    `json:"rtSrchJobsQuota"           values:"rtSrchJobsQuota,omitzero"`
	SrchDiskQuota             attributes.Explicit[int]    `json:"srchDiskQuota"             values:"srchDiskQuota,omitzero"`
	SrchFilter                attributes.Explicit[string] `json:"srchFilter"                values:"srchFilter,omitzero"`
	SrchIndexesAllowed        []string                    `json:"srchIndexesAllowed"        values:"srchIndexesAllowed,omitzero,fillempty"`
	SrchIndexesDefault        []string                    `json:"srchIndexesDefault"        values:"srchIndexesDefault,omitzero,fillempty"`
	SrchJobsQuota             attributes.Explicit[int]    `json:"srchJobsQuota"             values:"srchJobsQuota,omitzero"`
	SrchTimeWin               attributes.Explicit[int]    `json:"srchTimeWin"               values:"srchTimeWin,omitzero"`

	// Read-only fields are populated by results returned by the Splunk API, but
	// are not settable by Create or Update operations.
	ImportedCapabilities       []string                    `json:"imported_capabilities"       values:"-"`
	ImportedRoles              []string                    `json:"imported_roles"              values:"-"`
	ImportedRtSrchJobsQuota    attributes.Explicit[int]    `json:"imported_rtSrchJobsQuota"    values:"-"`
	ImportedRtSrchJObsQuota    attributes.Explicit[int]    `json:"imported_rtSrchJObsQuota"    values:"-"`
	ImportedSrchDiskQuota      attributes.Explicit[int]    `json:"imported_srchDiskQuota"      values:"-"`
	ImportedSrchFilter         attributes.Explicit[string] `json:"imported_srchFilter"         values:"-"`
	ImportedSrchIndexesAllowed []string                    `json:"imported_srchIndexesAllowed" values:"-"`
	ImportedSrchIndexesDefault []string                    `json:"imported_srchIndexesDefault" values:"-"`
	ImportedSrchJobsQuota      attributes.Explicit[int]    `json:"imported_srchJobsQuota"      values:"-"`
	ImportedSrchTimeWin        attributes.Explicit[int]    `json:"imported_srchTimeWin"        values:"-"`
}

// Role defines a Splunk role.
type Role struct {
	ID      client.ID   `selective:"create" service:"authorization/roles"`
	Content RoleContent `json:"content" values:",anonymize"`
}
