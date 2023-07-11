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
	"strconv"
	"strings"
)

// StatusCodes define the expected status code returned for a given operation.
type StatusCodes struct {
	Created  int
	Read     int
	Updated  int
	Deleted  int
	NotFound int
}

// WithDefaults returns a new StatusCodes with values from defaults applied for any
// field with the zero value.
func (codes StatusCodes) WithDefaults(defaults StatusCodes) StatusCodes {
	// for codesV to be settable, it needs to come from a pointer
	// codes wasn't passed as a pointer, so this still only permits changing the local copy of codes
	codesV := reflect.ValueOf(&codes).Elem()

	defaultsV := reflect.ValueOf(defaults)

	for i := 0; i < codesV.NumField(); i++ {
		field := codesV.Field(i)
		if !field.CanSet() || field.Kind() != reflect.Int {
			continue
		}

		if field.IsZero() {
			field.Set(defaultsV.Field(i))
		}
	}

	return codes
}

// withTagDefaults returns a new StatusCodes with values from the tag's configuration
// applied for any field with the zero value.
func (codes StatusCodes) withTagDefaults(tag string) (StatusCodes, error) {
	if tag == "" {
		return codes, nil
	}

	actions := strings.Split(tag, ",")

	tagDefaultCodes := StatusCodes{}
	// for codesV to be settable, it needs to come from a pointer
	codesV := reflect.ValueOf(&tagDefaultCodes).Elem()

	for _, action := range actions {
		actionParts := strings.Split(action, "=")
		if len(actionParts) != 2 {
			return StatusCodes{}, fmt.Errorf("service: unable to parse action code: %q", action)
		}

		actionName := actionParts[0]
		actionCode, err := strconv.ParseInt(actionParts[1], 10, 0)
		if err != nil {
			return StatusCodes{}, fmt.Errorf("service: unable to parse status code for action: %q: %s", action, err)
		}

		actionField := codesV.FieldByName(actionName)
		if !actionField.IsValid() {
			return StatusCodes{}, fmt.Errorf("service: unknown action name: %s", actionName)
		}

		if actionField.Kind() != reflect.Int {
			return StatusCodes{}, fmt.Errorf("service: field not int: %s", actionName)
		}

		if !actionField.CanSet() {
			return StatusCodes{}, fmt.Errorf("service: unable to set field: %s", actionName)
		}
		actionField.SetInt(actionCode)
	}

	return codes.WithDefaults(tagDefaultCodes), nil
}

// WithTagDefaults returns a new StatusCodes by applying the given tag and defaults.
func (codes StatusCodes) WithTagDefaults(tag string, defaults StatusCodes) (StatusCodes, error) {
	codes, err := codes.withTagDefaults(tag)
	if err != nil {
		return StatusCodes{}, err
	}

	return codes.WithDefaults(defaults), nil
}
