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
	"github.com/splunk/go-splunk-client/pkg/client"
)

func Test_clientIdResourceDataHandler(t *testing.T) {
	tests := syncResourceTestCases{
		{
			name: "empty",
			sync: NewClientID(&client.ID{}, "name"),
			schema: map[string]*schema.Schema{
				"name": {
					Type: schema.TypeString,
				},
			},
			checkFunc: composeResourceCheckFunc(
				checkResourceIdEquals(""),
				checkResourceKeyEquals("name", ""),
			),
		},
		{
			name: "set",
			sync: NewClientID(mustParseID("https://localhost:8089/services/authentication/users/testuser"), "name"),
			schema: map[string]*schema.Schema{
				"name": {
					Type: schema.TypeString,
				},
			},
			checkFunc: composeResourceCheckFunc(
				checkResourceIdEquals("https://localhost:8089/services/authentication/users/testuser"),
				checkResourceKeyEquals("name", "testuser"),
			),
		},
	}

	tests.test(t)
}

func Test_clientIdObjectValueHandler(t *testing.T) {
	tests := syncObjectTestCases{
		{
			name: "empty",
			sync: NewClientID(&client.ID{}, "name"),
			schema: map[string]*schema.Schema{
				"name": {
					Type: schema.TypeString,
				},
			},
			checkFunc: checkGetObjectEquality(&client.ID{}),
		},
		{
			name:     "invalid resource Id",
			prepFunc: withId("invalid"),
			sync:     NewClientID(&client.ID{}, "name"),
			schema: map[string]*schema.Schema{
				"name": {
					Type: schema.TypeString,
				},
			},
			// clientId.SyncObject doesn't actually return any errors, as invalid URLs
			// are assumed to be due to migration from the legacy client.
			checkFunc: checkGetObjectEquality(&client.ID{}),
		},
		{
			name: "valid resource name",
			sync: NewClientID(&client.ID{}, "name"),
			schema: map[string]*schema.Schema{
				"name": {
					Type: schema.TypeString,
				},
			},
			schemaValues: map[string]interface{}{
				"name": "valid_name",
			},
			checkFunc: checkGetObjectEquality(&client.ID{Title: "valid_name"}),
		},
		{
			name:     "valid resource Id",
			prepFunc: withId("https://localhost:8089/services/authentication/users/testuser"),
			sync:     NewClientID(&client.ID{}, "name"),
			schema: map[string]*schema.Schema{
				"name": {
					Type: schema.TypeString,
				},
			},
			checkFunc: checkGetObjectEquality(mustParseID("https://localhost:8089/services/authentication/users/testuser")),
		},
	}

	tests.test(t)
}
