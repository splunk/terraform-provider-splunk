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

import (
	"encoding/xml"
	"reflect"
	"testing"
)

// xmlUnmarshalerTestCase represents an individual test case for unmarshaling from XML.
type xmlUnmarshalerTestCase struct {
	input            string
	gotInterfacePtr  interface{}
	wantInterfacePtr interface{}
	wantError        bool
}

// test performs the test for a xmlUnmarshalerTestCase definition.
func (test xmlUnmarshalerTestCase) test(t *testing.T) {
	err := xml.Unmarshal([]byte(test.input), test.gotInterfacePtr)
	gotError := err != nil

	if gotError != test.wantError {
		t.Errorf("xml.Unmarshal(%q) returned error? %v (%s)", test.input, gotError, err)
	}

	if !reflect.DeepEqual(test.gotInterfacePtr, test.wantInterfacePtr) {
		t.Errorf("xml.Unmarshal(%q) got\n%#v, want\n%#v", test.input, test.gotInterfacePtr, test.wantInterfacePtr)
	}
}

// xmlUnmarshalerTestCases is a slice of xmlUnmarshalerTestCase objects.
type xmlUnmarshalerTestCases []xmlUnmarshalerTestCase

// test performs each xmlUnmarshalerTestCase in xmlUnmarshalerTestCases.
func (tests xmlUnmarshalerTestCases) test(t *testing.T) {
	for _, test := range tests {
		test.test(t)
	}
}
