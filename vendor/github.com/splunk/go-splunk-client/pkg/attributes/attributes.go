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

// Package attributes provides types that can represent explicitly set values,
// including the zero values of their stored types, to enable HTTP requests to
// account for all intended parameter values.
package attributes

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// explicitlyStorableValue defines the types storable by ExplicitValue.
type explicitlyStorableValue interface {
	bool | int | string
}

// Explicit is permits storing a value such that it can be explicitly empty/zero.
type Explicit[T explicitlyStorableValue] struct {
	set   bool
	value T
}

// GetURLValue implements custom encoding of its url.Values value.
func (e Explicit[T]) GetURLValue() interface{} {
	return e.value
}

// Set explicitly sets the value.
func (e *Explicit[T]) Set(value T) {
	e.set = true
	e.value = value
}

// NewExplicit returns a new Explicit with its value explicitly set.
func NewExplicit[T explicitlyStorableValue](value T) (newExplicitValue Explicit[T]) {
	newExplicitValue.Set(value)

	return
}

// UnmarshalJSON implements custom JSON unmarshaling. The unmarshaled is explicitly set.
func (e *Explicit[T]) UnmarshalJSON(data []byte) error {
	var newValue T

	if err := json.Unmarshal(data, &newValue); err != nil {
		return err
	}

	e.Set(newValue)

	return nil
}

// Bool returns a value indicating the boolean representation of the stored value, and
// another boolean that will be true only if the value was explicitly set, and it can
// be parsed by strconv.ParseBool without error.
func (e Explicit[T]) Bool() (value bool, ok bool) {
	if !e.set {
		return
	}

	strValue := fmt.Sprint(e.value)
	value, err := strconv.ParseBool(strValue)
	ok = err == nil

	return
}

// Value returns the stored value.
func (e Explicit[T]) Value() T {
	return e.value
}

// String returns the defualt string representation of the stored value.
func (e Explicit[T]) String() string {
	return fmt.Sprint(e.value)
}
