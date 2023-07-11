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

// Package service provides functionality to determine the service and
// entry paths for a struct.
package service

import (
	"fmt"
	"reflect"
)

// ServicePathGetter is the interface for types that implement GetServicePath.
type ServicePathGetter interface {
	// GetServicePath returns the object's service path based on the provided value.
	GetServicePath(string) (string, error)
}

// EntryPathGetter is the interfacae for types that implement GetEntryPath.
type EntryPathGetter interface {
	// GetEntryPath returns the object's entry path based on the provided value.
	GetEntryPath(string) (string, error)
}

// ServicePath returns the relative service path for an object.
//
// The path will be the first of these that are present:
//
// • The value returned by the struct's GetServicePath method, if the struct is a service.ServicePathGetter.
//
// • The value returned by an exported field's GetServicePath method, if the struct is a service.ServicePathGetter.
// The field's "service" tag will be passed to its GetServicePath method..
func ServicePath(i interface{}) (string, error) {
	if path, err := structOrFieldPath(i, fieldServicePath); path != "" || err != nil {
		return path, err
	}

	return "", fmt.Errorf("service: unable to determine ServicePath for type %T", i)
}

// EntryPath returns the relative entry path for an object.
//
// The path will be the first of these that are present:
//
// • The value returned by the struct's EntryPath method, if implemented.
//
// • The value returned by an exported service.EntryPathGetter member field's EntryPath
// method, which will be passed that field's "service" tag value.
func EntryPath(i interface{}) (string, error) {
	if path, err := structOrFieldPath(i, fieldEntryPath); path != "" || err != nil {
		return path, err
	}

	return "", fmt.Errorf("service: unable to determine EntryPath")
}

// structOrFieldPath returns the path for a given struct. The path is determined as the first of:
//
// • the result of calling pathFunc directly for the input struct (does the struct implement pathFunc's required interface?) and an empty tag value
//
// • the result of caling getPathFromStructFields for the input struct
//
// If no error is encountered, no path was determined, an empty path and nil error will be returned.
func structOrFieldPath(i interface{}, f pathFunc) (string, error) {
	iV := reflect.ValueOf(i)

	if !iV.IsValid() {
		return "", fmt.Errorf("service: attempted operation on invalid value")
	}

	if path, err := f(iV, ""); path != "" || err != nil {
		return path, err
	}

	if path, err := getPathFromStructFields(iV, f); path != "" || err != nil {
		return path, err
	}

	return "", nil
}

// getPathFromStructFields returns the path for a given struct's fields. The path is determined by iterating
// through each of the struct's exported fields and running the given pathFunc for them with the field's
// "service" struct tag. If multiple fields return a non-empty value for the pathFunc, an error will be returned,
// as this is an ambiguous configuration.
func getPathFromStructFields(v reflect.Value, f pathFunc) (string, error) {
	derefV := derefValue(v)

	if derefV.Kind() != reflect.Struct {
		return "", fmt.Errorf("service: attempted ServicePath of non-struct type %T", derefV.Type().Name())
	}

	var servicePathGetterPath string

	for i := 0; i < derefV.Type().NumField(); i++ {
		field := derefV.Type().Field(i)

		var fieldV reflect.Value

		if !field.IsExported() {
			// unexported fields are treated like their type's zero value
			fieldV = reflect.New(field.Type).Elem()
		} else {
			fieldV = derefV.Field(i)
		}

		fieldTag := field.Tag.Get("service")

		gotPath, err := f(fieldV, fieldTag)
		if err != nil {
			return "", err
		}

		if gotPath != "" {
			if servicePathGetterPath != "" {
				return "", fmt.Errorf("service: multiple ServicePathGetter fields found in type %T", i)
			}

			servicePathGetterPath = gotPath
		}
	}

	return servicePathGetterPath, nil
}

// pathFunc describes a function that returns a path for a reflect.Value and tag value.
type pathFunc func(reflect.Value, string) (string, error)

// fieldServicePath is a pathFunc that returns the path for a ServicePathGetter.
var fieldServicePath pathFunc = func(fieldV reflect.Value, fieldTag string) (string, error) {
	if fieldPathGetter, ok := fieldV.Interface().(ServicePathGetter); ok {
		return fieldPathGetter.GetServicePath(fieldTag)
	}

	return "", nil
}

// fieldEntryPath is a pathFunc that returns the path for an EntryPathGetter.
var fieldEntryPath pathFunc = func(fieldV reflect.Value, fieldTag string) (string, error) {
	if entryPathGetter, ok := fieldV.Interface().(EntryPathGetter); ok {
		return entryPathGetter.GetEntryPath(fieldTag)
	}

	return "", nil
}

// derefValue returns a new reflect.Value by dereferencing v.
func derefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v
}
