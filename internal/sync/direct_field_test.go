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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Test_directField_ResourceDataHandler(t *testing.T) {
	tests := syncResourceTestCases{
		{
			name:      "invalid key",
			sync:      NewDirectField("invalid_field", new(string)),
			wantError: true,
		},
		{
			name: "valid key, type mismatch",
			schema: map[string]*schema.Schema{
				"int_field": {
					Type: schema.TypeInt,
				},
			},
			schemaValues: map[string]interface{}{
				"int_field": 1,
			},
			sync:      NewDirectField("int_field", new(string)),
			checkFunc: checkResourceKeyEquals("int_field", 1),
			wantError: true,
		},
		{
			name: "valid key, type match (string)",
			schema: map[string]*schema.Schema{
				"string_field": {
					Type: schema.TypeString,
				},
			},
			schemaValues: map[string]interface{}{
				"string_field": "resource_value",
			},
			sync:      NewDirectField("string_field", pointerToValueOf("object_value")),
			checkFunc: checkResourceKeyEquals("string_field", "object_value"),
		},
		{
			name: "valid key, type match (int)",
			schema: map[string]*schema.Schema{
				"int_field": {
					Type: schema.TypeInt,
				},
			},
			schemaValues: map[string]interface{}{
				"int_field": 1,
			},
			sync:      NewDirectField("int_field", pointerToValueOf(2)),
			checkFunc: checkResourceKeyEquals("int_field", 2),
		},
		{
			// this is a confusing test, because a schema type List won't work for the object handler,
			// but it does work for the resource handler.
			name: "list of strings",
			schema: map[string]*schema.Schema{
				"list_string_field": {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			schemaValues: map[string]interface{}{
				"list_string_field": []interface{}{
					"resource_string_value_1",
					"resource_string_value_2",
				},
			},
			sync: NewDirectField("list_string_field", pointerToValueOf([]string{
				"object_string_value_1",
				"object_string_value_2",
			})),
			checkFunc: checkResourceKeyEquals("list_string_field", []interface{}{
				"object_string_value_1",
				"object_string_value_2",
			}),
		},
	}

	tests.test(t)
}

func Test_directField_ObjectHandler(t *testing.T) {
	tests := syncObjectTestCases{
		{
			name:      "invalid key",
			sync:      NewDirectField("invalid_field", pointerToValueOf("object_value")),
			checkFunc: checkGetObjectEquality(pointerToValueOf("object_value")),
			wantError: true,
		},
		{
			name: "valid key, type mismatch",
			schema: map[string]*schema.Schema{
				"int_field": {
					Type: schema.TypeInt,
				},
			},
			sync:      NewDirectField("int_field", pointerToValueOf("object_value")),
			checkFunc: checkGetObjectEquality(pointerToValueOf("object_value")),
			wantError: true,
		},
		{
			name: "valid key, type match (string)",
			schema: map[string]*schema.Schema{
				"string_field": {
					Type: schema.TypeString,
				},
			},
			schemaValues: map[string]interface{}{
				"string_field": "resource_value",
			},
			sync:      NewDirectField("string_field", pointerToValueOf("object_value")),
			checkFunc: checkGetObjectEquality(pointerToValueOf("resource_value")),
		},
		{
			name: "valid key, type match (int)",
			schema: map[string]*schema.Schema{
				"int_field": {
					Type: schema.TypeInt,
				},
			},
			schemaValues: map[string]interface{}{
				"int_field": 1,
			},
			sync:      NewDirectField("int_field", pointerToValueOf(2)),
			checkFunc: checkGetObjectEquality(pointerToValueOf(1)),
		},
		{
			// lists aren't supported
			name: "list of strings",
			schema: map[string]*schema.Schema{
				"list_string_field": {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			schemaValues: map[string]interface{}{
				// the field must be set to a non-zero value for this test to fail as expected,
				// because SyncObject() uses GetOk() which returns ok=false if the key is unset
				// or set to a zero value.
				"list_string_field": []interface{}{
					"resource_value_1",
					"resource_value_2",
				},
			},
			sync: NewDirectField("list_string_field", pointerToValueOf([]string{
				"object_value_1",
				"object_value_2",
			})),
			checkFunc: checkGetObjectEquality(pointerToValueOf([]string{
				"object_value_1",
				"object_value_2",
			})),
			wantError: true,
		},
		{
			name: "unset field, no default",
			schema: map[string]*schema.Schema{
				"string_field": {
					Type: schema.TypeString,
				},
			},
			sync:      NewDirectField("string_field", pointerToValueOf("object_value")),
			checkFunc: checkGetObjectEquality(pointerToValueOf("")),
		},
		{
			name: "unset field, non-zero default",
			schema: map[string]*schema.Schema{
				"string_field": {
					Type:    schema.TypeString,
					Default: "resource_default_value",
				},
			},
			sync:      NewDirectField("string_field", pointerToValueOf("object_value")),
			checkFunc: checkGetObjectEquality(pointerToValueOf("resource_default_value")),
		},
		{
			name: "explicitly zero field, non-zero default",
			schema: map[string]*schema.Schema{
				"string_field": {
					Type:    schema.TypeString,
					Default: "resource_default_value",
				},
			},
			schemaValues: map[string]interface{}{
				"string_field": "",
			},
			sync:      NewDirectField("string_field", pointerToValueOf("object_value")),
			checkFunc: checkGetObjectEquality(pointerToValueOf("")),
		},
	}

	tests.test(t)
}
