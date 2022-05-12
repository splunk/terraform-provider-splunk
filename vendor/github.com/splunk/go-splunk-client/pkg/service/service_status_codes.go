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

package service

import (
	"fmt"
	"reflect"
)

// TagDefaultsStatusCoder is the interface for types that return a StatusCodes object from a service
// tag configuration and default StatusCodes value.
type TagDefaultsStatusCoder interface {
	WithTagDefaults(string, StatusCodes) (StatusCodes, error)
}

// ServiceStatusCodes returns the StatusCodes configuration for the given input struct.
func ServiceStatusCodes(input interface{}, defaults StatusCodes) (StatusCodes, error) {
	var defaultsCoder TagDefaultsStatusCoder
	var serviceTag string

	if coder, ok := input.(TagDefaultsStatusCoder); ok {
		defaultsCoder = coder
	}

	if defaultsCoder == nil {
		inputV := reflect.ValueOf(input)

		for inputV.Kind() == reflect.Ptr {
			inputV = inputV.Elem()
		}

		inputT := inputV.Type()
		if inputV.Kind() != reflect.Struct {
			return StatusCodes{}, fmt.Errorf("service: ServiceStatusCodes attempted on non-struct type %T", input)
		}

		for i := 0; i < inputV.NumField(); i++ {
			fieldV := inputV.Field(i)
			field := inputT.Field(i)
			fieldT := field.Type
			if !field.IsExported() {
				// unexported fields are treated like the zero values of their type
				fieldV = reflect.New(fieldT).Elem()
			}

			if coder, ok := fieldV.Interface().(TagDefaultsStatusCoder); ok {
				if defaultsCoder != nil {
					return StatusCodes{}, fmt.Errorf("service: ServiceStatusCodes attempted on struct type %T with multiple DefaultsStatusCoder fields", input)
				}

				defaultsCoder = coder
				serviceTag = inputV.Type().Field(i).Tag.Get("service")
			}
		}
	}

	newCodes := defaults
	if defaultsCoder != nil {
		withTagDefaults, err := defaultsCoder.WithTagDefaults(serviceTag, defaults)
		if err != nil {
			return StatusCodes{}, err
		}

		newCodes = withTagDefaults
	}

	return newCodes, nil
}
