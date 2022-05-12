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

package authenticators

import (
	"net/http"
	"testing"

	"github.com/splunk/go-splunk-client/pkg/client"
	"github.com/splunk/go-splunk-client/pkg/internal/checks"
)

// AuthenticatorTestCase defines a test against a specific Authenticator and Client.
type AuthenticatorTestCase struct {
	name               string
	inputAuthenticator client.Authenticator
	inputClient        *client.Client
	wantError          bool
	requestCheck       checks.CheckRequestFunc
}

// test performs a AuthenticatorTestCase's defined test.
func (test AuthenticatorTestCase) test(t *testing.T) {
	r := &http.Request{}
	err := test.inputAuthenticator.AuthenticateRequest(test.inputClient, r)
	gotError := err != nil

	if gotError != test.wantError {
		t.Errorf("%s AuthenticateRequest returned error? %v", test.name, gotError)
	}

	if test.requestCheck != nil {
		test.requestCheck(r, t)
	}
}

// AuthenticatorTestCases is a list of AuthenticatorTestCase instances.
type AuthenticatorTestCases []AuthenticatorTestCase

// test performs the test for each of its AuthenticatorTestCase items.
func (tests AuthenticatorTestCases) test(t *testing.T) {
	for _, test := range tests {
		test.test(t)
	}
}
