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
	"sync"

	"github.com/splunk/go-splunk-client/pkg/client"
)

// Password defines password authentication to Splunk.
type Password struct {
	Username string `values:"username"`
	Password string `values:"password"`

	// UseBasicAuth can be set to true if Basic Authentication should always be used,
	// which causes Username/Password to be passed with each authenticated request.
	UseBasicAuth bool `values:"-"`

	// SessionKey holds the SessionKey after initial authentication occurs. Unless
	// UseBasicAuth is set to true, this SessionKey will be used to authenticate requests.
	SessionKey `url:"-"`

	// mu is used to enable locking to prevent race conditions when checking for and obtaining
	// a SessionKey.
	mu sync.Mutex

	// empty Namespace used to satisfy service.ServicePathGetter
	_ client.Namespace `service:"auth/login"`
}

// loginResponse represents the response returned from auth/login.
type loginResponse struct {
	SessionKey
}

// authenticate performs the authentication request and handles the response, storing the SessionKey
// if successful.
func (p *Password) authenticate(c *client.Client) error {
	lR := loginResponse{}

	if err := c.RequestAndHandle(
		client.ComposeRequestBuilder(
			client.BuildRequestMethod(http.MethodPost),
			client.BuildRequestServiceURL(c, p),
			client.BuildRequestBodyValues(p),
		),
		client.ComposeResponseHandler(
			client.HandleResponseCode(http.StatusUnauthorized, client.HandleResponseXMLMessagesCustomError(client.ErrorUnauthorized)),
			client.HandleResponseRequireCode(http.StatusOK, client.HandleResponseXMLMessagesError()),
			client.HandleResponseXML(&lR),
		),
	); err != nil {
		return err
	}

	p.SessionKey = lR.SessionKey

	return nil
}

// authenticateOnce calls authenticate only if currently unauthenticated.
func (p *Password) authenticateOnce(c *client.Client) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.authenticated() {
		return p.authenticate(c)
	}

	return nil
}

// AuthenticateRequest adds authentication to an http.Request.
func (p *Password) AuthenticateRequest(c *client.Client, r *http.Request) error {
	if err := p.authenticateOnce(c); err != nil {
		return err
	}

	return p.SessionKey.AuthenticateRequest(c, r)
}
