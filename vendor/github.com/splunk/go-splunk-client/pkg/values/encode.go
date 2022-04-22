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

// Package values implements encoding objects into url.Values.
package values

import (
	"fmt"
	"net/url"
	"reflect"
)

// Encode returns url.Values for a given input interface.
//
// A struct field's key defaults to the name of the field. Nested values default
// to having their key as a dotted suffix to the parent object's key. An empty parent
// key results in the child key being used directly. The parent key for the original
// value passed to Encode is an empty string.
//
// Anonymous struct fields also start with an empty parent key.
//
// Struct tags can be used to customize the encoding behavior of specific fields.
// Tags are formatted as "<key>,<option>...". If <key> is left empty the default
// key value will be used, which is the name of the field. The field is not encoded
// if the full value of the tag is "-".
//
// The name of the encoded value may also be determined by having it adhere to the
// values.KeyEncoder interface, of which the GetURLKey method accepts a parent and
// child key and returns a customized key.
//
// Option flags can be included as comma-separated values after <key>. Supported
// options are:
//
// • omitzero - Omit the field if it is the zero value for its type. Note that
// the zero value for a slice is nil. An empty slice is not treated as a nil slice.
//
// • fillempty - Populate empty slices (including nil) with a single value. The assigned
// value will be the zero type for the slice element. If both fillempty and omitzero flags
// are given, omitzero has precedence. This permits empty and nil slices to behave differently,
// as configured.
//
// • anonymize - Treat the field as if it were anonymous. This gives it an empty
// parent key.
//
// An error is returned if a value's key is an empty string at any level
// of encoding.
func Encode(i interface{}) (url.Values, error) {
	inputV := reflect.ValueOf(i)

	if !inputV.IsValid() {
		return nil, fmt.Errorf("values: attempted Encode() on invalid (likely nil) type %T", i)
	}

	newValues := url.Values{}

	if err := encodeValue("", &newValues, inputV); err != nil {
		return nil, err
	}

	return newValues, nil
}

// encodeStructValue adds the given input struct to url.Values for a given key.
func encodeStructValue(key string, inputV reflect.Value, values *url.Values) error {
	inputT := inputV.Type()
	for i := 0; i < inputV.NumField(); i++ {
		field := inputT.Field(i)

		if !field.IsExported() {
			continue
		}

		fieldV := inputV.Field(i)

		fieldName := field.Name

		fieldTag := field.Tag.Get("values")
		if fieldTag == "-" {
			continue
		}

		fieldOptions, err := parseTagConfig(fieldTag)
		if err != nil {
			return err
		}

		// override if tag specified name
		if fieldOptions.Name != "" {
			fieldName = fieldOptions.Name
		}

		// clear anonymous-like fields
		if field.Anonymous || fieldOptions.Anonymize {
			fieldName = ""
		}

		fieldName, err = computedKey(inputV, key, fieldName)
		if err != nil {
			return err
		}

		// skip empty values if configured to do so
		if fieldOptions.Omitzero && fieldV.IsZero() {
			continue
		}

		// fill empty slices with a single zero value if configured to do so
		if fieldV.Kind() == reflect.Slice && fieldV.Len() == 0 && fieldOptions.Fillempty {
			nestedKey, err := computedKey(fieldV, fieldName, 0)
			if err != nil {
				return err
			}

			iV := reflect.New(fieldV.Type().Elem()).Elem()
			if err := encodeValue(nestedKey, values, iV); err != nil {
				return err
			}
		}

		if err := encodeValue(fieldName, values, fieldV); err != nil {
			return err
		}
	}

	return nil
}

// encodeValue adds the given reflect.Value to url.Values for a given key and tagConfig.
func encodeValue(key string, values *url.Values, value reflect.Value) error {
	// use full custom encoding if implemented
	if valuesEncoder, ok := value.Interface().(URLValuesSetter); ok {
		return valuesEncoder.SetURLValues(key, values)
	}

	// use value-only custom encoding if implemented
	if valueGetter, ok := value.Interface().(URLValueGetter); ok {
		return encodeValue(key, values, reflect.ValueOf(valueGetter.GetURLValue()))
	}

	// fully dereference if needed
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {
	default:
		// refuse to encode to an empty key.
		// note this is done only where this function adds to url.Values, as custom encoding should be
		// trusted to choose the proper key, even if it was empty here.
		if key == "" {
			return fmt.Errorf("values: attempted to encode empty key for value %#v", value.Interface())
		}

		values.Add(key, fmt.Sprint(value.Interface()))

	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			nestedKey, err := computedKey(value, key, i)
			if err != nil {
				return err
			}

			iV := value.Index(i)
			if err := encodeValue(nestedKey, values, iV); err != nil {
				return err
			}
		}

	case reflect.Map:
		for _, keyV := range value.MapKeys() {
			nestedKey, err := computedKey(value, key, keyV.Interface())
			if err != nil {
				return err
			}

			nestedValue := value.MapIndex(keyV)
			if err := encodeValue(nestedKey, values, nestedValue); err != nil {
				return err
			}
		}

	case reflect.Struct:
		if err := encodeStructValue(key, value, values); err != nil {
			return err
		}

	}

	// use additional custom encoding if implemented
	if valuesAdder, ok := value.Interface().(URLValuesAdder); ok {
		return valuesAdder.AddURLValues(key, values)
	}

	return nil
}

// computedKey returns the key that should be used for the url.Values value. childKeyInterface
// is an interface because it may be a map key, which isn't guaranteed to be a string. If the
// input value doesn't implement the values.GetURLKey interface, the returned key will be
// <parentKey>.<childKey>, or just <childKey> if <parentKey> is empty.
func computedKey(value reflect.Value, parentKey string, childKeyInterface interface{}) (string, error) {
	childKey := fmt.Sprint(childKeyInterface)

	if keyEncoder, ok := value.Interface().(URLKeyGetter); ok {
		return keyEncoder.GetURLKey(parentKey, childKey)
	}

	// unless overridden by KeyEncoder (above), slices and arrays use the parent key directly
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		return parentKey, nil
	}

	// an empty parentKey indicates this is a top-level key, so the childKey should be returned directly
	if parentKey == "" {
		return childKey, nil
	}

	return fmt.Sprintf("%s.%s", parentKey, childKey), nil
}
