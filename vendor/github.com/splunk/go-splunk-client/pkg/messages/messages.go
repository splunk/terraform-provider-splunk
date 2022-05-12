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

// Package messages provides the Messages and Message type to represent messages
// returned by the Splunk REST API.
package messages

import (
	"strings"
)

// Messages represents the <messages> element of a <response> entry.
type Messages struct {
	XMLName string    `xml:"messages"`
	Items   []Message `json:"messages" xml:"msg"`
}

// String returns the string representation of Messages. If multiple Message items are
// present, they will be comma-separated.
func (m Messages) String() string {
	itemStrings := make([]string, len(m.Items))
	for i := 0; i < len(m.Items); i++ {
		itemStrings[i] = m.Items[i].String()
	}

	return strings.Join(itemStrings, ",")
}
