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
	"bytes"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

// CheckRequestFunc functions perform a check against an http.Request.
type CheckRequestFunc func(*http.Request, *testing.T)

// ComposeCheckRequestFunc returns a new CheckRequestFunc from an arbitrary number of other
// CheckRequestFunc functions.
func ComposeCheckRequestFunc(checks ...CheckRequestFunc) CheckRequestFunc {
	return func(r *http.Request, t *testing.T) {
		for _, check := range checks {
			check(r, t)
		}
	}
}

// CheckRequestHeaderKeyValue checks that an http.Request's header has a given value for
// a given key.
func CheckRequestHeaderKeyValue(key string, value ...string) CheckRequestFunc {
	return func(r *http.Request, t *testing.T) {
		if r.Header == nil {
			t.Errorf("CheckRequestHeaderKeyValue: Header not set")
			return
		}

		got, ok := r.Header[key]
		if !ok {
			t.Errorf("CheckRequestHeaderKeyValue: Key %s not set", key)
			return
		}

		if !reflect.DeepEqual(got, value) {
			t.Errorf("CheckRequestHeaderKeyValue: Key %s = %#v, want %#v", key, got, value)
			return
		}
	}
}

// CheckRequestMethod checks that a requests method matches the given method.
func CheckRequestMethod(method string) CheckRequestFunc {
	return func(r *http.Request, t *testing.T) {
		if r.Method != method {
			t.Errorf("CheckRequestMethod: got %s, want %s", r.Method, method)
		}
	}
}

// CheckRequestURL checks that a request's URL matches the given URL.
func CheckRequestURL(url string) CheckRequestFunc {
	return func(r *http.Request, t *testing.T) {
		gotURL := r.URL.String()

		if gotURL != url {
			t.Errorf("CheckRequestURL: got\n%s, want\n%s", gotURL, url)
		}
	}
}

// CheckRequestBodyValue checks that a requests URL-encoded body has the given
// key and value.
func CheckRequestBodyValue(key string, value ...string) CheckRequestFunc {
	return func(r *http.Request, t *testing.T) {
		if r.Body == nil {
			t.Errorf("CheckRequestBodyValue: Body is nil")
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("CheckRequestBodyValue: unable to read body: %s", err)
			return
		}

		// because reading is a one-time deal, we need to put back the data we
		// read from r.Body
		r.Body = io.NopCloser(bytes.NewReader(data))

		v, err := url.ParseQuery(string(data))
		if err != nil {
			t.Errorf("CheckRequestBodyValue: unable to parse query: %s", err)
			return
		}

		got, ok := v[key]
		if !ok {
			t.Errorf("CheckRequestBodyValue: key %s not present", key)
			return
		}

		if !reflect.DeepEqual(got, value) {
			t.Errorf("CheckRequestBodyValue: key %s got %v, want %v", key, got, value)
			return
		}
	}
}
