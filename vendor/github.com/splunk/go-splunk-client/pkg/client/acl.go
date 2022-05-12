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

package client

import (
	"encoding/json"

	"github.com/splunk/go-splunk-client/pkg/attributes"
)

// Sharing represents the level of sharing of a Splunk object.
type Sharing string

const (
	SharingUndefined Sharing = ""
	SharingGlobal    Sharing = "global"
	SharingUser      Sharing = "user"
	SharingApp       Sharing = "app"
)

// validate returns an error if Sharing is a value other than the predefined Sharing constants.
func (sharing Sharing) validate() error {
	switch sharing {
	default:
		return wrapError(ErrorSharing, nil, "client: invalid Sharing value %q", sharing)
	case SharingUndefined, SharingGlobal, SharingUser, SharingApp:
		return nil
	}
}

// MarshalJSON implements custom marshaling.
func (sharing Sharing) MarshalJSON() ([]byte, error) {
	if err := sharing.validate(); err != nil {
		return nil, err
	}

	return json.Marshal(string(sharing))
}

// Permissions represents the read/write permissions of a Splunk object.
type Permissions struct {
	Read  []string `json:"read"  values:"read,omitzero,fillempty"`
	Write []string `json:"write" values:"write,omitzero,fillempty"`
}

// ACL represents the ACL of a Splunk object.
type ACL struct {
	Permissions Permissions                 `json:"perms"   values:"perms"`
	Owner       attributes.Explicit[string] `json:"owner"   values:"owner,omitzero"`
	Sharing     Sharing                     `json:"sharing" values:"sharing,omitzero"`
}
