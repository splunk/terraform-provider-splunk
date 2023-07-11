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

// Package deepset provides functionality to set the an input object's value directly,
// or set it on one of the object's struct fields, based on the type of the value being
// set.
package deepset

import (
	"fmt"
	"reflect"
)

// Set attempts to set the given value on dest.
//
// If dest and value are of the same underlying type, dest will be set directly.
//
// If dest is a struct, and exactly one of its fields are of the same underlying type as value, that field will be set to value.
//
// Any dest field that is a struct will be followed recursively.
//
// If multiple struct fields at any given level match the value type, an error is returned.
func Set(dest interface{}, value interface{}) error {
	destV := reflect.ValueOf(dest)

	for destV.Kind() == reflect.Ptr {
		destV = destV.Elem()
	}

	valueV := reflect.ValueOf(value)

	if !(destV.IsValid() && valueV.IsValid()) {
		return fmt.Errorf("deepset: dest (%T) or input (%T) invalid", dest, value)
	}

	if destV.Type() == valueV.Type() {
		if !destV.CanSet() {
			return fmt.Errorf("deepset: dest type (%T) is not settable", dest)
		}

		destV.Set(valueV)
		return nil
	}

	if destV.Kind() == reflect.Struct {
		setFieldV, err := structSetFieldValue(destV, valueV)
		if err != nil {
			return err
		}

		if setFieldV.CanSet() {
			setFieldV.Set(valueV)
			return nil
		}
	}

	return fmt.Errorf("deepset: unable to set value (%T) on dest (%T)", value, dest)
}

func structSetFieldValue(destV reflect.Value, valueV reflect.Value) (reflect.Value, error) {
	var foundSetFieldV reflect.Value

	fieldsV := make([]reflect.Value, 0, destV.NumField())

	for i := 0; i < destV.NumField(); i++ {
		fieldV := destV.Field(i)
		if !fieldV.CanSet() {
			continue
		}

		fieldsV = append(fieldsV, fieldV)
	}

	for _, fieldV := range fieldsV {
		if fieldV.Type() == valueV.Type() {
			if foundSetFieldV.IsValid() {
				return reflect.Value{}, fmt.Errorf("deepset: dest type (%T) has multiple fields of value type (%T)", destV, valueV)
			}

			foundSetFieldV = fieldV
		}
	}

	if foundSetFieldV.IsValid() {
		return foundSetFieldV, nil
	}

	for _, fieldV := range fieldsV {
		if fieldV.Kind() == reflect.Struct {
			foundEmbeddedSetFieldV, err := structSetFieldValue(fieldV, valueV)
			// error will be returned if the embedded struct is ambiguous
			if err != nil {
				return reflect.Value{}, err
			}

			// check if *this* struct is ambiguous
			if foundSetFieldV.IsValid() {
				return reflect.Value{}, fmt.Errorf("deepset: dest type (%T) has multiple fields of value type (%T)", destV, valueV)
			}

			foundSetFieldV = foundEmbeddedSetFieldV
		}
	}

	return foundSetFieldV, nil
}
