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
	"fmt"
	"net/http"

	"github.com/splunk/go-splunk-client/pkg/client"
)

// SessionKey provides authentication to Splunk via a session key.
type SessionKey struct {
	// SessionKey is the session key that will be used to authenticate to Splunk.
	SessionKey string `xml:"sessionKey"`
}

// authenticated returns true if SessionKey is not empty.
func (s SessionKey) authenticated() bool {
	return s.SessionKey != ""
}

// AuthenticateRequest adds the SessionKey to the http.Request's Header.
func (s SessionKey) AuthenticateRequest(c *client.Client, r *http.Request) error {
	if !s.authenticated() {
		return fmt.Errorf("attempted to authenticate request with empty SessionKey")
	}

	if r.Header == nil {
		r.Header = http.Header{}
	}

	r.Header.Add("Authorization", fmt.Sprintf("Splunk %s", s.SessionKey))

	return nil
}
