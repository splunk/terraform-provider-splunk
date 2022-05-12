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

// UserContent defines the content of a User object.
type UserContent struct {
	// Read/Write
	DefaultApp            attributes.Explicit[string] `values:"defaultApp,omitzero"`
	Email                 attributes.Explicit[string] `values:"email,omitzero"`
	ForceChangePass       attributes.Explicit[bool]   `values:"force-change-pass,omitzero"`
	Password              attributes.Explicit[string] `values:"password,omitzero"`
	RealName              attributes.Explicit[string] `values:"realname,omitzero"`
	RestartBackgroundJobs attributes.Explicit[bool]   `values:"restart_background_jobs,omitzero"`
	Roles                 []string                    `values:"roles,omitzero,fillempty"`
	TZ                    attributes.Explicit[string] `values:"tz,omitzero"`

	// Read-only fields are populated by results returned by the Splunk API, but
	// are not settable by Create or Update operations.
	Capabilities []string                    `values:"-"`
	Type         attributes.Explicit[string] `values:"-"`
}

// User defines a Splunk user.
type User struct {
	ID      client.ID   `selective:"create" service:"authentication/users"`
	Content UserContent `json:"content" values:",anonymize"`
}
