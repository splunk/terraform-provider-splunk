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

package messages

import "fmt"

// Message represents the <msg> element of a <messages> entry.
type Message struct {
	Value string `json:"text" xml:",chardata"`
	Code  string `json:"type" xml:"code,attr"`
}

// String returns the string representation of a message. It will be in the form:
//
//   Code: Value
func (m Message) String() string {
	return fmt.Sprintf("%s: %s", m.Code, m.Value)
}
