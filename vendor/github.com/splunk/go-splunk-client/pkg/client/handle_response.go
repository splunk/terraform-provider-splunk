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
	"encoding/json"
	"encoding/xml"
	"net/http"
	"reflect"

	"github.com/splunk/go-splunk-client/pkg/messages"
)

// ResponseHandler defines a function that performs an action on an http.Response.
type ResponseHandler func(*http.Response) error

// ComposeResponseHandler creates a new ResponseHandler that runs each ResponseHandler
// provided as an argument.
func ComposeResponseHandler(handlers ...ResponseHandler) ResponseHandler {
	return func(r *http.Response) error {
		for _, handler := range handlers {
			if err := handler(r); err != nil {
				return err
			}
		}

		return nil
	}
}

// HandleResponseXML returns a ResponseHandler that decodes an http.Response's Body
// as XML to the given interface.
func HandleResponseXML(i interface{}) ResponseHandler {
	return func(r *http.Response) error {
		if err := xml.NewDecoder(r.Body).Decode(i); err != nil {
			return wrapError(ErrorResponseBody, err, "unable to decode response XML: %s", err)
		}

		return nil
	}
}

// HandleResponseXMLMessagesCustomError returns a ResponseHandler that decode's an http.Response's
// Body as an XML document of Messages and returns the Messages as an error with the given ErrorCode.
func HandleResponseXMLMessagesCustomError(code ErrorCode) ResponseHandler {
	return func(r *http.Response) error {
		response := struct {
			Messages messages.Messages
		}{}

		if err := HandleResponseXML(&response)(r); err != nil {
			return err
		}

		return wrapError(code, nil, "response contained message: %s", response.Messages.String())
	}
}

// HandleResponseXMLMessagesError returns a ResponseHandler that decodes an http.Response's Body
// as an XML document of Messages and returns the Messages as an error.
func HandleResponseXMLMessagesError() ResponseHandler {
	return func(r *http.Response) error {
		return HandleResponseXMLMessagesCustomError(ErrorSplunkMessage)(r)
	}
}

// HandleResponseJSON returns a ResponseHandler that decodes an http.Response's Body
// as JSON to the given interface.
func HandleResponseJSON(i interface{}) ResponseHandler {
	return func(r *http.Response) error {
		if err := json.NewDecoder(r.Body).Decode(i); err != nil {
			return wrapError(ErrorResponseBody, err, "unable to decode response JSON: %s", err)
		}

		return nil
	}
}

// HandleResponseJSONMessagesCustomError returns a ResponseHandler that decodes an http.Response's
// Body as JSON document of Messages and returns the Messages as an error with the given ErrorCode.
func HandleResponseJSONMessagesCustomError(code ErrorCode) ResponseHandler {
	return func(r *http.Response) error {
		msg := messages.Messages{}
		if err := HandleResponseJSON(&msg)(r); err != nil {
			return err
		}

		return wrapError(code, nil, "response contained message: %s", msg.String())
	}
}

// HandleResponseJSONMessagesError returns a ResponseHandler that decode's an http.Response's Body
// as a JSON document of Messages and returns the Messages as an error with the Code ErrorSplunkMessage.
func HandleResponseJSONMessagesError() ResponseHandler {
	return func(r *http.Response) error {
		return HandleResponseJSONMessagesCustomError(ErrorSplunkMessage)(r)
	}
}

// HandleResponseCode returns a ResponseHandler that calls errorResponseHandler if an http.Response's
// StatusCode is equal to the provided code.
func HandleResponseCode(code int, errorResponseHandler ResponseHandler) ResponseHandler {
	return func(r *http.Response) error {
		if r.StatusCode != code {
			return nil
		}

		return errorResponseHandler(r)
	}
}

// HandleResponseRequireCode returns a ResponseHandler that checks for a given StatusCode. If
// the http.Response has a different StatusCode, the provided ResponseHandler will be called
// to return the appopriate error message.
func HandleResponseRequireCode(code int, errorResponseHandler ResponseHandler) ResponseHandler {
	return func(r *http.Response) error {
		if r.StatusCode == code {
			return nil
		}

		return errorResponseHandler(r)
	}
}

// HandleResponseEntries returns a ResponseHandler that parses the http.Response Body
// into the list of Entry reference provided.
func HandleResponseEntries(entries interface{}) ResponseHandler {
	return func(r *http.Response) error {
		entriesPtrV := reflect.ValueOf(entries)
		if entriesPtrV.Kind() != reflect.Ptr {
			return wrapError(ErrorPtr, nil, "attempted to read entries to non-pointer")
		}

		entriesV := reflect.Indirect(entriesPtrV)
		if entriesV.Kind() != reflect.Slice {
			return wrapError(ErrorSlice, nil, "attempted to read entries to non-slice")
		}

		responseT := reflect.StructOf([]reflect.StructField{
			{
				Name: "Entry",
				Type: entriesV.Type(),
			},
		})

		entriesResponsePtrV := reflect.New(responseT)
		entriesResponsePtrI := entriesResponsePtrV.Interface()

		d := json.NewDecoder(r.Body)
		if err := d.Decode(entriesResponsePtrI); err != nil {
			return wrapError(ErrorResponseBody, err, "unable to decode JSON: %s", err)
		}

		entriesResponseV := reflect.Indirect(entriesResponsePtrV)
		responseEntriesFieldV := entriesResponseV.FieldByName("Entry")

		entriesV.Set(responseEntriesFieldV)

		return nil
	}
}

// HandleResponseEntry returns a responseHaResponseHandlerndler that parses the http.Response Body
// into the given Entry.
func HandleResponseEntry(entry interface{}) ResponseHandler {
	return func(r *http.Response) error {
		entryPtrV := reflect.ValueOf(entry)
		if entryPtrV.Kind() != reflect.Ptr {
			return wrapError(ErrorPtr, nil, "attempted to read entry to non-pointer")
		}

		entryV := reflect.Indirect(entryPtrV)
		entryT := entryV.Type()

		entriesT := reflect.SliceOf(entryT)
		entriesPtrV := reflect.New(entriesT)
		entriesPtrI := entriesPtrV.Interface()

		if err := HandleResponseEntries(entriesPtrI)(r); err != nil {
			return err
		}

		entriesV := reflect.Indirect(entriesPtrV)

		if entriesV.Len() != 1 {
			return wrapError(ErrorResponseBody, nil, "expected exactly 1 entry, got %d", entriesV.Len())
		}

		entryV.Set(entriesV.Index(0))

		return nil
	}
}
