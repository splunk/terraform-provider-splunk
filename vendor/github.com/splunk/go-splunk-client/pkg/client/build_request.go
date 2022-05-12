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

package client

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/splunk/go-splunk-client/pkg/selective"
	"github.com/splunk/go-splunk-client/pkg/service"
	"github.com/splunk/go-splunk-client/pkg/values"
)

// defaultStatusCodes set the expected StatusCodes for most CRUD operations. If
// a given type doesn't override them, responses will be checked against these
// codes.
var defaultStatusCodes = service.StatusCodes{
	Created:  http.StatusCreated,
	Read:     http.StatusOK,
	Updated:  http.StatusOK,
	Deleted:  http.StatusOK,
	NotFound: http.StatusNotFound,
}

// RequestBuilder defines a function that performs an operation on an http.Request.
type RequestBuilder func(*http.Request) error

// ComposeRequestBuilder creates a new RequestBuilder that performs each RequestBuilder
// provided as an argument, returning the first error encountered, if any.
func ComposeRequestBuilder(builders ...RequestBuilder) RequestBuilder {
	return func(r *http.Request) error {
		for _, builder := range builders {
			if err := builder(r); err != nil {
				return err
			}
		}

		return nil
	}
}

// buildRequest creates a new http.Request and applies the provided RequestBuilder.
func buildRequest(builder RequestBuilder) (*http.Request, error) {
	r := &http.Request{}

	if err := builder(r); err != nil {
		return nil, err
	}

	return r, nil
}

// BuildRequestMethod returns a RequestBuilder that sets the given method.
func BuildRequestMethod(method string) RequestBuilder {
	return func(r *http.Request) error {
		r.Method = method

		return nil
	}
}

// BuildRequestServiceURL returns a RequestBuilder that sets the URL to the ServiceURL
// for a given Service.
func BuildRequestServiceURL(c *Client, service interface{}) RequestBuilder {
	return func(r *http.Request) error {
		u, err := c.ServiceURL(service)
		if err != nil {
			return err
		}

		r.URL = u

		return nil
	}
}

// BuildRequestBodyValues returns a RequestBuilder that sets the Body to the encoded url.Values for
// a given interface.
func BuildRequestBodyValues(i interface{}) RequestBuilder {
	return func(r *http.Request) error {
		v, err := values.Encode(i)
		if err != nil {
			return wrapError(ErrorValues, err, err.Error())
		}

		r.Body = io.NopCloser(strings.NewReader(v.Encode()))

		return nil
	}
}

// BuildRequestOutputModeJSON returns a RequestBuilder that sets the URL's RawQuery to output_mode=json.
// It checks that the URL is already set, so it must be applied after setting the URL. It overwrites
// any existing RawQuery Values.
func BuildRequestOutputModeJSON() RequestBuilder {
	return func(r *http.Request) error {
		if r.URL == nil {
			return wrapError(ErrorNilValue, nil, "unable to set output mode on nil URL")
		}

		if r.URL.RawQuery != "" {
			return wrapError(ErrorOverwriteValue, nil, "attempted to set output_mode after RawQuery already set")
		}

		r.URL.RawQuery = url.Values{
			"output_mode": []string{"json"},
		}.Encode()

		return nil
	}
}

// BuildRequestBodyValuesSelective returns a RequestBuilder that sets the Body to the encoded url.Values
// for a given interface and selective tag.
func BuildRequestBodyValuesSelective(c interface{}, tag string) RequestBuilder {
	return func(r *http.Request) error {
		selected, err := selective.Encode(c, tag)
		if err != nil {
			return err
		}

		return BuildRequestBodyValues(selected)(r)
	}
}

// BuildRequestCollectionURL returns a RequestBuilder that sets the URL to the EntryURL
// for a given Entry.
func BuildRequestEntryURL(c *Client, entry interface{}) RequestBuilder {
	return func(r *http.Request) error {
		u, err := c.EntryURL(entry)
		if err != nil {
			return err
		}

		r.URL = u

		return nil
	}
}

// BuildRequestEntryACLURL returns a RequestBuilder that sets the URL to the ACL URL
// for a given Entry.
func BuildRequestEntryACLURL(c *Client, entry interface{}, acl ACL) RequestBuilder {
	return func(r *http.Request) error {
		u, err := c.EntryACLURL(entry)
		if err != nil {
			return err
		}

		r.URL = u

		return nil
	}
}

// BuildRequestAuthenticate returns a RequestBuilder that authenticates a request for a given Client.
func BuildRequestAuthenticate(c *Client) RequestBuilder {
	return func(r *http.Request) error {
		return c.Authenticator.AuthenticateRequest(c, r)
	}
}

// BuildRequestGetServiceStatusCodes updates codes for the given entry. It returns a RequestBuilder
// that returns the error (if any) returned by service.ServiceStatusCodes.
func BuildRequestGetServiceStatusCodes(entry interface{}, codes *service.StatusCodes) RequestBuilder {
	newCodes, err := service.ServiceStatusCodes(entry, defaultStatusCodes)
	*codes = newCodes

	return func(r *http.Request) error {
		return err
	}
}
