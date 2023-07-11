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
	"encoding/json"
	"reflect"
	"testing"
)

// JSONUnmarshalTestCase defines a test case for json.Unmarshal.
type JSONUnmarshalTestCase struct {
	Name        string
	InputString string
	Want        interface{}
	WantError   bool
}

// Test runs the Test case.
func (test JSONUnmarshalTestCase) Test(t *testing.T) {
	// create a new pointer to a zero value of test.want
	gotT := reflect.TypeOf(test.Want)
	if gotT == nil {
		t.Fatalf("%s attempted with nil want type", test.Name)
	}
	gotV := reflect.New(gotT)
	gotP := gotV.Interface()

	// create a new pointer to a the same type as test.want,
	// and set its data to match test.want
	wantT := reflect.TypeOf(test.Want)
	wantV := reflect.New(wantT)
	wantV.Elem().Set(reflect.ValueOf(test.Want))
	wantP := wantV.Interface()

	err := json.Unmarshal([]byte(test.InputString), gotP)
	gotError := err != nil
	if gotError != test.WantError {
		t.Fatalf("%s json.Unmarshal returned error? %v (%s)", test.Name, gotError, err)
	}

	if !reflect.DeepEqual(gotP, wantP) {
		t.Errorf("%s json.Unmarshal got\n%#v, want\n%#v", test.Name, gotP, wantP)
	}
}

// JSONUnmarshalTestCases is a collection of JSONUnmarshalTestCases tests.
type JSONUnmarshalTestCases []JSONUnmarshalTestCase

// Test runs the Test defined for each item in the collection.
func (tests JSONUnmarshalTestCases) Test(t *testing.T) {
	for _, test := range tests {
		test.Test(t)
	}
}
