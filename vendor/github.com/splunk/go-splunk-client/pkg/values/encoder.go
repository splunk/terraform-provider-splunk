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

package values

import "net/url"

// URLValuesSetter is the interface for types that implement custom encoding to url.Values
// for a given key.
type URLValuesSetter interface {
	SetURLValues(string, *url.Values) error
}

// URLValuesAdder is the interface for types that implement custom encoding to url.Values
// for a given key in addition to default encoding. AddURLValues will be called in addition
// to the default encoding methods. It is not called if the encoded type is a URLValuesSetter
// or URLValueGetter.
type URLValuesAdder interface {
	AddURLValues(string, *url.Values) error
}

// URLKeyGetter is the interface for types that implement custom key calculation prior
// to encoding to url.Values.
type URLKeyGetter interface {
	GetURLKey(parentKey string, childKeystring string) (string, error)
}

// URLValueGetter is the interface for types that implement custom value calculation
// prior to encoding to url.Values.
type URLValueGetter interface {
	GetURLValue() interface{}
}
