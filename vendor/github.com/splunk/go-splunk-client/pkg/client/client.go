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

// Package client implements a client to the Splunk REST API.
package client

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/splunk/go-splunk-client/pkg/deepset"
	"github.com/splunk/go-splunk-client/pkg/internal/paths"
	"github.com/splunk/go-splunk-client/pkg/service"
	"golang.org/x/net/publicsuffix"
)

const (
	defaultTimeout = time.Minute * 5
)

// Client defines connectivity and authentication to a Splunk REST API.
type Client struct {
	// URL is the URL to the Splunk REST API. It should include the scheme and port number.
	//
	// Example:
	//   https://localhost:8089
	URL string

	// Authenticator defines which authentication method and credentials to use.
	//
	// Example:
	//   authenticators.Password{Username: "admin", Password: "changeme"}
	Authenticator

	// Set TLSInsecureSkipVerify to true to skip TLS verification.
	TLSInsecureSkipVerify bool

	// Timeout configures the timeout of requests. If unspecified, defaults to 5 minutes.
	Timeout time.Duration

	httpClient *http.Client
	mu         sync.Mutex
}

// urlForPath returns a url.URL for path, relative to Client's URL.
func (c *Client) urlForPath(path ...string) (*url.URL, error) {
	if c.URL == "" {
		return nil, wrapError(ErrorMissingURL, nil, "Client has empty URL")
	}

	combinedPath := paths.Join(path...)

	u := paths.Join(c.URL, combinedPath)

	return url.Parse(u)
}

// ServiceURL returns a url.URL for a Service, relative to the Client's URL.
func (c *Client) ServiceURL(s interface{}) (*url.URL, error) {
	servicePath, err := service.ServicePath(s)
	if err != nil {
		return nil, err
	}

	return c.urlForPath(servicePath)
}

// EntryURL returns a url.URL for an Entry, relative to the Client's URL.
func (c *Client) EntryURL(e interface{}) (*url.URL, error) {
	entryPath, err := service.EntryPath(e)
	if err != nil {
		return nil, err
	}

	return c.urlForPath(entryPath)
}

func (c *Client) EntryACLURL(e interface{}) (*url.URL, error) {
	entryPath, err := service.EntryPath(e)
	if err != nil {
		return nil, err
	}

	return c.urlForPath(entryPath, "acl")
}

// httpClientPrep prepares the Client's http.Client.
func (c *Client) httpClientPrep() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.httpClient == nil {
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			return wrapError(ErrorHTTPClient, err, "unable to create new cookiejar: %s", err)
		}

		timeout := c.Timeout
		if timeout == 0 {
			timeout = defaultTimeout
		}

		c.httpClient = &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: c.TLSInsecureSkipVerify,
				},
			},
			Jar: jar,
		}
	}

	return nil
}

// do performs a given http.Request via the Client's http.Client.
func (c *Client) do(r *http.Request) (*http.Response, error) {
	if err := c.httpClientPrep(); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, wrapError(ErrorHTTPClient, err, "error encountered performing request: %s", err)
	}

	return resp, nil
}

// RequestAndHandle creates a new http.Request from the given RequestBuilder, performs the
// request, and handles the http.Response with the given ResponseHandler.
func (c *Client) RequestAndHandle(builder RequestBuilder, handler ResponseHandler) error {
	req, err := buildRequest(builder)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return handler(resp)
}

// Create performs a Create action for the given Entry.
func (client *Client) Create(entry interface{}) error {
	var codes service.StatusCodes

	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestGetServiceStatusCodes(entry, &codes),
			BuildRequestMethod(http.MethodPost),
			BuildRequestServiceURL(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestBodyValuesSelective(entry, "create"),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(codes.Created, HandleResponseJSONMessagesError()),
		),
	)
}

// Read performs a Read action for the given Entry. It modifies entry in-place,
// so entry must be a pointer.
func (client *Client) Read(entry interface{}) error {
	var codes service.StatusCodes

	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestGetServiceStatusCodes(entry, &codes),
			BuildRequestMethod(http.MethodGet),
			BuildRequestEntryURL(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseCode(codes.NotFound, HandleResponseJSONMessagesCustomError(ErrorNotFound)),
			HandleResponseRequireCode(codes.Read, HandleResponseJSONMessagesError()),
			HandleResponseEntry(entry),
		),
	)
}

// Update performs an Update action for the given Entry.
func (client *Client) Update(entry interface{}) error {
	var codes service.StatusCodes

	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestGetServiceStatusCodes(entry, &codes),
			BuildRequestMethod(http.MethodPost),
			BuildRequestEntryURL(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestBodyValuesSelective(entry, "update"),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(codes.Updated, HandleResponseJSONMessagesError()),
		),
	)
}

// Delete performs a Delete action for the given Entry.
func (client *Client) Delete(entry interface{}) error {
	var codes service.StatusCodes

	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestGetServiceStatusCodes(entry, &codes),
			BuildRequestMethod(http.MethodDelete),
			BuildRequestEntryURL(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(codes.Deleted, HandleResponseJSONMessagesError()),
		),
	)
}

func (client *Client) listModified(entries interface{}, modifier interface{}) error {
	entriesPtrV := reflect.ValueOf(entries)
	if entriesPtrV.Kind() != reflect.Ptr {
		return wrapError(ErrorPtr, nil, "client: List attempted on on-pointer value")
	}

	entriesV := reflect.Indirect(entriesPtrV)
	if entriesV.Kind() != reflect.Slice {
		return wrapError(ErrorSlice, nil, "client: List attempted on non-slice value")
	}
	entryT := entriesV.Type().Elem()
	entryI := reflect.New(entryT).Interface()

	if modifier != nil {
		if err := deepset.Set(entryI, modifier); err != nil {
			return err
		}
	}

	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodGet),
			BuildRequestEntryURL(client, entryI),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(http.StatusOK, HandleResponseJSONMessagesError()),
			HandleResponseEntries(entries),
		),
	)
}

// ListNamespace populates entries in place for a Namespace.
func (client *Client) ListNamespace(entries interface{}, ns Namespace) error {
	return client.listModified(entries, ns)
}

// ListNamespace populates entries in place for an ID.
func (client *Client) ListID(entries interface{}, id ID) error {
	return client.listModified(entries, id)
}

// ListNamespace populates entries in place without any ID or Namespace context.
func (client *Client) List(entries interface{}) error {
	return client.listModified(entries, nil)
}

// ReadACL performs a ReadACL action for the given Entry. It modifies acl in-place,
// so acl must be a pointer.
func (client *Client) ReadACL(entry interface{}, acl *ACL) error {
	var aclResponse struct {
		ACL ACL `json:"acl"`
	}

	if err := client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodGet),
			BuildRequestEntryACLURL(client, entry, *acl),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseCode(http.StatusNotFound, HandleResponseJSONMessagesCustomError(ErrorNotFound)),
			HandleResponseRequireCode(http.StatusOK, HandleResponseJSONMessagesError()),
			HandleResponseEntry(&aclResponse),
		),
	); err != nil {
		return err
	}

	*acl = aclResponse.ACL

	return nil
}

// UpdateACL performs an UpdateACL action for the given Entry.
func (client *Client) UpdateACL(entry interface{}, acl ACL) error {
	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodPost),
			BuildRequestEntryACLURL(client, entry, acl),
			BuildRequestBodyValues(acl),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(http.StatusOK, HandleResponseJSONMessagesError()),
		),
	)
}
