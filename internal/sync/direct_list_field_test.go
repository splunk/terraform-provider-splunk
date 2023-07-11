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

func Test_directListField_ResourceDataHandler(t *testing.T) {
	tests := syncResourceTestCases{
		{
			name: "invalid key",
			sync: NewDirectListField("invalid_field", pointerToValueOf([]string{
				"object_value_1",
				"object_value_2",
			})),
			wantError: true,
		},
		{
			name: "valid key, type mismatch",
			schema: map[string]*schema.Schema{
				"list_int_field": {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
				},
			},
			schemaValues: map[string]interface{}{
				"list_int_field": []interface{}{
					0,
					1,
				},
			},
			sync: NewDirectListField("list_int_field", pointerToValueOf([]string{
				"object_value_1",
				"object_value_2",
			})),
			checkFunc: checkResourceKeyEquals("list_int_field", []interface{}{
				0,
				1,
			}),
			wantError: true,
		},
		{
			name: "valid key, type match (string)",
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
					"resource_value_1",
					"resource_value_2",
				},
			},
			sync: NewDirectListField("list_string_field", pointerToValueOf([]string{
				"object_value_1",
				"object_value_2",
			})),
			checkFunc: checkResourceKeyEquals("list_string_field", []interface{}{
				"object_value_1",
				"object_value_2",
			}),
		},
		{
			name: "valid key, type match (int)",
			schema: map[string]*schema.Schema{
				"list_int_field": {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
				},
			},
			schemaValues: map[string]interface{}{
				"list_int_field": []interface{}{
					0,
					1,
				},
			},
			sync: NewDirectListField("list_int_field", pointerToValueOf([]int{
				2,
				3,
			})),
			checkFunc: checkResourceKeyEquals("list_int_field", []interface{}{
				2,
				3,
			}),
		},
	}

	tests.test(t)
}

func Test_directListField_ObjectHandler(t *testing.T) {
	tests := syncObjectTestCases{
		{
			name: "invalid key",
			sync: NewDirectListField("invalid_key", pointerToValueOf([]string{
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
			name: "valid key, type mismatch",
			schema: map[string]*schema.Schema{
				"list_int_field": {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
				},
			},
			schemaValues: map[string]interface{}{
				"list_int_field": []interface{}{
					0,
					1,
				},
			},
			sync: NewDirectListField("list_int_field", pointerToValueOf([]string{
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
			name: "valid key, type match (string)",
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
					"resource_value_1",
					"resource_value_2",
				},
			},
			sync: NewDirectListField("list_string_field", pointerToValueOf([]string{
				"object_value_1",
				"object_value_2",
			})),
			checkFunc: checkGetObjectEquality(pointerToValueOf([]string{
				"resource_value_1",
				"resource_value_2",
			})),
		},
		{
			name: "valid key, type match (int)",
			schema: map[string]*schema.Schema{
				"list_int_field": {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
				},
			},
			schemaValues: map[string]interface{}{
				"list_int_field": []interface{}{
					0,
					1,
				},
			},
			sync: NewDirectListField("list_int_field", pointerToValueOf([]int{
				2,
				3,
			})),
			checkFunc: checkGetObjectEquality(pointerToValueOf([]int{
				0,
				1,
			})),
		},
		{
			name: "unset field",
			schema: map[string]*schema.Schema{
				"list_string_field": {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			sync: NewDirectListField("list_string_field", pointerToValueOf([]string{
				"object_value_1",
				"object_value_2",
			})),
			checkFunc: checkGetObjectEquality(pointerToValueOf([]string{})),
		},
	}

	tests.test(t)
}
