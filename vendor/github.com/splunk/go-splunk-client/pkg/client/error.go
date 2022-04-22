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
	"fmt"
	"runtime/debug"
)

// ErrorCode identifies the type of error that was encountered.
type ErrorCode int

const (
	// ErrorUndefined is the zero-value, indicating that the error type
	// has not been defined.
	ErrorUndefined ErrorCode = iota

	// ErrorNamespace indicates an error with the Namespace.
	ErrorNamespace

	// ErrorEndpoint indicates an error with the Endpoint configuration.
	ErrorEndpoint

	// ErrorValues indicates an error was encountered while trying to encode
	// to url.Values.
	ErrorValues

	// ErrorNilValue indicates an attempt to perform an action against a nil value,
	// such as attempting to set RawQuery on http.Request with a nil URL.
	ErrorNilValue

	// ErrorOverwriteValue indicates an attempt to overwrite an existing value,
	// such as attempting to set RawQuery multiple times on a URL.
	ErrorOverwriteValue

	// ErrorMissingTitle indicates an operation that required a non-empty Title was
	// attempted with an empty Title.
	ErrorMissingTitle

	// ErrorMissingURL indicates the Client's URL value is missing.
	ErrorMissingURL

	// ErrorHTTPClient indicates an error related to the http.Client was encountered.
	ErrorHTTPClient

	// ErrorResponseBody indicates an error encountered while trying to parse the Body
	// from an http.Response.
	ErrorResponseBody

	// ErrorSplunkMessage indicates the Splunk REST API returned an error message.
	ErrorSplunkMessage

	// ErrorUnauthorized indicates a request was unauthorized.
	ErrorUnauthorized

	// ErrorNotFound indicates an attempt was made against an object that count not be found.
	ErrorNotFound

	// ErrorPtr indicates an operation requiring a pointer was passed a non-pointer.
	ErrorPtr

	// ErrorSlice indicates an operation requiring a slice was passed a non-slice, or a slice
	// of the wrong type.
	ErrorSlice

	// ErrorID indicates an error was encountered related to an object's ID.
	ErrorID

	// ErrorSharing indicates an error was encountered related to a Sharing value.
	ErrorSharing
)

// Error represents an encountered error. It adheres to the "error" interface,
// so will be returned as a standard error.
//
// Returned errors can be handled as this Error type:
//
//   if err := c.RequestAndHandle(...); err != nil {
// 	  if clientErr, ok := err.(client.Error) {
// 		  // check clientErr.Code to determine appropriate action
// 	  }
//   }
type Error struct {
	Code       ErrorCode
	Message    string
	Wrapped    error
	StackTrace string
}

// Wrap returns a new Error with the given code, error, and message. The error value
// may be nil if this is a new error.
func wrapError(code ErrorCode, err error, messagef string, messageArgs ...interface{}) Error {
	return Error{
		Code:       code,
		Message:    fmt.Sprintf(messagef, messageArgs...),
		Wrapped:    err,
		StackTrace: string(debug.Stack()),
	}
}

// Error returns the Error's message.
func (err Error) Error() string {
	return err.Message
}
