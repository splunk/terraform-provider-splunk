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
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/splunk/go-splunk-client/pkg/attributes"
	"github.com/splunk/go-splunk-client/pkg/client"
)

// StanzaContent defines the content for a Stanza.
type StanzaContent struct {
	Disabled attributes.Explicit[bool] `json:"disabled" values:"disabled,omitzero"`
	Values   map[string]string         `json:"-"        values:",anonymize"`
}

// Stanza is a Splunk configs/conf-<file> stanza.
type Stanza struct {
	ID      client.ConfID `selective:"create" service:"configs"`
	Content StanzaContent `json:"content"     values:",anonymize"`
}

// UnmarshalJSON implements custom JSON unmarshaling.
func (content *StanzaContent) UnmarshalJSON(data []byte) error {
	type aliasType StanzaContent
	var aliasValue aliasType

	if err := json.Unmarshal(data, &aliasValue); err != nil {
		return err
	}

	*content = StanzaContent(aliasValue)

	var allValues map[string]interface{}
	if err := json.Unmarshal(data, &allValues); err != nil {
		return err
	}

	eaiR := regexp.MustCompile("^eai:")

	for key, value := range allValues {
		if eaiR.MatchString(key) {
			continue
		}

		if key == "disabled" {
			if valueBool, ok := value.(bool); ok {
				content.Disabled = attributes.NewExplicit(valueBool)
				continue
			} else {
				return fmt.Errorf("entry: unable to decode StanzaContent, disabled value %#v, expected bool", value)
			}
		}

		if valueString, ok := value.(string); ok {
			if content.Values == nil {
				content.Values = map[string]string{}
			}

			content.Values[key] = valueString
			continue
		} else {
			return fmt.Errorf("entry: unexpected non-string type in StanzaContent: %#v", value)
		}
	}

	return nil
}
