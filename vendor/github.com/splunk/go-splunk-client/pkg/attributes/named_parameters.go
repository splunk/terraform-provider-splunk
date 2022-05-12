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

package attributes

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// NamedParameters represent a set of Parameters that are associated with an overall Name.
type NamedParameters struct {
	// Name is the overall name of this set of Parameters. It is likely the leftmost segment
	// of a dotted parameter name, such as "actions" for "actions.email".
	Name string `values:"-"`

	// Status is the string representation of a NamedParameters' status. This is typically
	// true/false or 0/1, and is the value associated directly with the name segment, such as
	// email=true.
	Status Explicit[string] `values:",omitzero,anonymize"`

	Parameters Parameters `values:",omitzero,anonymize"`
}

// GetURLKey implements custom key encoding for url.Values.
func (params NamedParameters) GetURLKey(parentKey string, childKey string) (string, error) {
	if params.Name == "" {
		return "", fmt.Errorf("attributes: unable to determine url.Values key for empty NamedParameters.Name")
	}

	return fmt.Sprintf("%s.%s", parentKey, params.Name), nil
}

// NamedParametersCollection is a collection of NamedParameters.
type NamedParametersCollection []NamedParameters

// EnabledNames returns a list of Names of the member NamedParameters that have a true Status value.
func (collection NamedParametersCollection) EnabledNames() []string {
	var enabled []string

	for _, params := range collection {
		if isEnabled, ok := params.Status.Bool(); ok && isEnabled {
			enabled = append(enabled, params.Name)
		}
	}

	return enabled
}

// UnmarshalJSONForNamedParametersCollections unmarshals JSON data into the given dest interface. dest must be
// a pointer to a struct, and any struct fields with the "named_parameters_collection" tag must be of the
// NamedParametersCollection type.
//
// This method exists to enable unmarshaling of the same level of a JSON document to a struct and also
// to NamedParametersCollection fields of the same struct.
func UnmarshalJSONForNamedParametersCollections(data []byte, dest interface{}) error {
	destVPtr := reflect.ValueOf(dest)
	if destVPtr.Kind() != reflect.Ptr {
		return fmt.Errorf("attempted UnmarshalJSONForNamedParametersCollection on non-pointer type: %T", dest)
	}

	destV := destVPtr.Elem()
	destT := destV.Type()

	if destT.Kind() != reflect.Struct {
		return fmt.Errorf("attempted UnmarshalJSONForNamedParametersCollection on non-struct type: %T", dest)
	}

	for i := 0; i < destT.NumField(); i++ {
		fieldF := destT.Field(i)
		if !fieldF.IsExported() {
			continue
		}

		fieldTag := fieldF.Tag.Get("named_parameters_collection")
		if fieldTag == "" {
			continue
		}

		var collection NamedParametersCollection
		if fieldF.Type != reflect.TypeOf(collection) {
			return fmt.Errorf("attempted UnmarshalJSONForNamedParametersCollection on non-NamedParametersCollection type %T for field %s", destV.Field(i).Interface(), fieldF.Name)
		}

		var allParams Parameters
		if err := json.Unmarshal(data, &allParams); err != nil {
			return err
		}

		newParams := allParams.withDottedName(fieldTag)

		newCollection := newParams.namedParametersCollection()
		newCollectionV := reflect.ValueOf(newCollection)

		destV.Field(i).Set(newCollectionV)
	}

	return nil
}
