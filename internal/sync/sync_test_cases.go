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

package sync

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// resourceDataFunc is a function that performs an operation against ResourceData.
type resourceDataFunc func(*schema.ResourceData) error

// withId returns a resourceDataFunc that calls schema.ResourceData.SetId for the given value.
func withId(id string) resourceDataFunc {
	return func(d *schema.ResourceData) error {
		d.SetId(id)

		return nil
	}
}

// resourceCheckFunc is a function that performs a check against schema.ResourceData.
type resourceCheckFunc func(t *testing.T, name string, d *schema.ResourceData)

// syncResourceTestCase defines a test case for SyncResource().
type syncResourceTestCase struct {
	// name is the name of the test. It is logged during failures.
	name string

	// schema is the schema.Schema for the test. It may be nil if the test needs no schema.
	schema map[string]*schema.Schema

	// schemaValues is a map of field names to values to be set directly to the schema prior
	// to the test being run. It may be nil if the test needs no initial schema values.
	schemaValues map[string]interface{}

	// sync is the Sync object that is being tested. It must be set.
	sync Sync

	// checkFunc is a resourceCheckFunc that will be run to perform a check against schema.ResourceData
	// after sync is applied. It may be nil if no check is needed (such as checking only for a returned error).
	checkFunc resourceCheckFunc

	// wantError defines if the test case is expected to return an error when sync.SyncResource() is called.
	wantError bool
}

// test performs the defined test.
func (test syncResourceTestCase) test(t *testing.T) {
	if test.sync == nil {
		t.Fatalf("%s: sync must be set", test.name)
	}

	d := schema.TestResourceDataRaw(t, test.schema, test.schemaValues)

	err := test.sync.SyncResource(d)
	gotError := err != nil

	if gotError != test.wantError {
		t.Errorf("%s: SyncResource() returned error? %v (%s)", test.name, gotError, err)
	}

	if test.checkFunc != nil {
		test.checkFunc(t, test.name, d)
	}
}

// syncResourceTestCases is a slice of syncResourceTestCases.
type syncResourceTestCases []syncResourceTestCase

// test runs test() for each syncResourceTestCase.
func (tests syncResourceTestCases) test(t *testing.T) {
	for _, test := range tests {
		test.test(t)
	}
}

// composeResourceCheckFunc returns a resourceCheckFunc composed of the given checks.
func composeResourceCheckFunc(checks ...resourceCheckFunc) resourceCheckFunc {
	return func(t *testing.T, name string, d *schema.ResourceData) {
		for _, check := range checks {
			check(t, name, d)
		}
	}
}

// checkResourceIdEquals returns a syncResourceTestFunc that checks the schema.ResourceData's
// Id value.
func checkResourceIdEquals(want string) resourceCheckFunc {
	return func(t *testing.T, name string, d *schema.ResourceData) {
		got := d.Id()

		if got != want {
			t.Errorf("%s: Id() got\n%s, want\n%s", name, got, want)
		}
	}
}

// checkResourceKeyEquals returns a syncResourceTestFunc that checks the value returned by
// schema.ResourceData's Get() method for the given key.
func checkResourceKeyEquals(key string, want interface{}) resourceCheckFunc {
	return func(t *testing.T, name string, d *schema.ResourceData) {
		got := d.Get(key)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: Get(key) got\n%#v, want\n%#v", name, got, want)
		}
	}
}

// objectCheckFunc is a function that performs a check against a Sync.
type objectCheckFunc func(t *testing.T, name string, s Sync)

// syncObjectTestCase defines a test case for SyncObject.
type syncObjectTestCase struct {
	// name is the name of the test. It is logged during failures.
	name string

	// schema is the schema.Schema for the test. It may be nil if the test needs no schema.
	schema map[string]*schema.Schema

	// schemaValues is a map of field names to values to be set directly to the schema prior
	// to the test being run. It may be nil if the test needs no initial schema values.
	schemaValues map[string]interface{}

	// prepFunc is a resourceDataFunc that is called prior to the test being run. It can be used
	// when schema.ResourceData needs to be manipulated in a way that schema.TestResourceDataRaw()
	// isn't capable of handling, such as calling SetId().
	prepFunc resourceDataFunc

	// sync is the Sync object that is being tested. It must be set.
	sync Sync

	// checkFunc is a objectCheckFunc that will be run to perform a check against the value returned by sync.GetObject()
	// after sync is applied. It may be nil if no check is needed (such as checking only for a returned error).
	checkFunc objectCheckFunc

	// wantError defines if the test case is expected to return an error when sync.SyncResource() is called.
	wantError bool
}

// test performs the defined test.
func (test syncObjectTestCase) test(t *testing.T) {
	d := schema.TestResourceDataRaw(t, test.schema, test.schemaValues)

	if test.prepFunc != nil {
		if err := test.prepFunc(d); err != nil {
			t.Fatalf("%s: prepFunc returned error: %s", test.name, err)
		}
	}

	if test.sync == nil {
		t.Fatalf("%s: sync must be set", test.name)
	}

	err := test.sync.SyncObject(d)
	gotError := err != nil

	if gotError != test.wantError {
		t.Errorf("%s: SyncObject() returned error? %v (%s)", test.name, gotError, err)
	}

	if test.checkFunc != nil {
		test.checkFunc(t, test.name, test.sync)
	}
}

// syncObjectTestCases is a slice of syncObjectTestCases.
type syncObjectTestCases []syncObjectTestCase

// test runs test() for each syncObjectTestCase.
func (tests syncObjectTestCases) test(t *testing.T) {
	for _, test := range tests {
		test.test(t)
	}
}

// checkGetObjectEquality returns a syncObjectTestFunc that tests equality between want and
// the value retured by Sync's getObject(). The given Sync object must also implement SyncGetter.
func checkGetObjectEquality(want interface{}) objectCheckFunc {
	return func(t *testing.T, name string, s Sync) {
		getter, ok := s.(SyncGetter)
		if !ok {
			t.Fatalf("%s: checkGetObjectEquality passed Sync that doesn't implement syncGetter: %#v", name, s)
		}

		got := getter.GetObject()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: GetObject() got\n%#v (%s), want\n%#v (%s)", name, got, mustJSON(got), want, mustJSON(want))
		}
	}
}

// mustJSON returns a JSOn string for the given interface. It panics if json.Marshal returns an error.
// It exists to ease crafting of JSON for test case error logging.
func mustJSON(i interface{}) string {
	out, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	return string(out)
}
