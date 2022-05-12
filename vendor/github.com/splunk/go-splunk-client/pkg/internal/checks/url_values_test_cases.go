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

package checks

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/splunk/go-splunk-client/pkg/values"
)

// QueryValuesTestCase defines a test case for query.Values.
type QueryValuesTestCase struct {
	Name      string
	Input     interface{}
	Want      url.Values
	WantError bool
}

// Test runs the Test.
func (test QueryValuesTestCase) Test(t *testing.T) {
	got, err := values.Encode(test.Input)
	gotError := err != nil

	if gotError != test.WantError {
		t.Errorf("%s values.Encode returned error? %v: %s", test.Name, gotError, err)
	}

	if !reflect.DeepEqual(got, test.Want) {
		t.Errorf("%s values.Encode got\n%#v, want\n%#v", test.Name, got, test.Want)
	}
}

// QueryValuesTestCases is a collection of queryValuesTestCase tests.
type QueryValuesTestCases []QueryValuesTestCase

// Test runs the Test defined for each item in the collection.
func (tests QueryValuesTestCases) Test(t *testing.T) {
	for _, test := range tests {
		test.Test(t)
	}
}
