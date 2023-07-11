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

// Package selective provides functionality to translate one struct type into another
// by selectively including fields with matching field tags.
package selective

import (
	"fmt"
	"reflect"
	"strings"
)

// hasTagValue returns true if an individual tag value is present in the comma-separated
// tag value. An empty csvTag indicates the field should not be filtered, and wantTag
// should be assumed to be valid for the field, thus an empty csvTag always returns true.
func hasTagValue(csvTag string, wantTag string) bool {
	if csvTag == "" {
		return true
	}

	parts := strings.Split(csvTag, ",")
	for _, part := range parts {
		if part == wantTag {
			return true
		}
	}

	return false
}

// Encode returns a value for an input struct value calculated by:
//
// • If no fields non-empty "selective" tag values, the input is returned directly, retaining
// unexported fields and methods.
//
// • Any field with a non-empty "selective" tag value is retained only if its tag contains
// the tag provided to Encode.
//
// • Struct field values are calculated by calling this function on the input field's value
// with the same tag.
//
// Returned values with changes (structs that contain at least one field with a non-empty
// "selective" tag) will lose unexported fields and methods. This is due to reflection being
// unable to retain them.
func Encode(i interface{}, tag string) (interface{}, error) {
	iV := reflect.ValueOf(i)

	// dereference as many times as needed
	for iV.Kind() == reflect.Ptr {
		iV = iV.Elem()
	}

	if iV.Kind() != reflect.Struct {
		return nil, fmt.Errorf("selective: attempted Encode on non-struct type (%T)", i)
	}

	var newStructFields []reflect.StructField
	var newStructValues []reflect.Value

	iT := iV.Type()
	hadTags := false
	for i := 0; i < iT.NumField(); i++ {
		field := iT.Field(i)
		if !field.IsExported() {
			continue
		}

		foundTag := field.Tag.Get("selective")
		if foundTag != "" {
			hadTags = true
		}

		if !hasTagValue(foundTag, tag) {
			continue
		}

		if field.Type.Kind() != reflect.Struct {
			newStructFields = append(newStructFields, field)
			newStructValues = append(newStructValues, iV.Field(i))
		} else {
			embeddedStruct, err := Encode(iV.Field(i).Interface(), tag)
			if err != nil {
				return nil, err
			}

			field.Type = reflect.TypeOf(embeddedStruct)
			newStructFields = append(newStructFields, field)
			newStructValues = append(newStructValues, reflect.ValueOf(embeddedStruct))
		}
	}

	if !hadTags {
		return i, nil
	}

	newStructT := reflect.StructOf(newStructFields)
	newStructVPtr := reflect.New(newStructT)
	newStructV := newStructVPtr.Elem()
	for i, newStructValue := range newStructValues {
		newStructV.Field(i).Set(newStructValue)
	}
	newStructI := newStructV.Interface()

	return newStructI, nil
}
