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
	"net/url"
	"reflect"
	"sort"
	"strings"
)

// dottedParameterNameParts parses a dotted parameter name and returns the first segment as name,
// and the remaining segments as paramName. If there are not multiple segments, name will be empty.
//
// For example:
//
//   dottedParameterNameParts("name")
//   # ("name", "")
//
//   dottedParameterNameParts("actions.email.to")
//   # ("actions", "email.to")
func dottedParameterNameParts(fullFieldName string) (name string, paramName string) {
	parts := strings.Split(fullFieldName, ".")

	return parts[0], strings.Join(parts[1:], ".")
}

// Parameters is a map of parameter names to string values.
type Parameters map[string]string

// withDottedName returns a new Parameters object containing the nested parameters
// for the given name. The new Parameters name field will have this name prefix removed.
//
// For example:
//
//   Parameters{"action.email": "true", "action.email.to": "whoever@example.com"}.withDottedName("action")
//   # Parameters{"email": "true", "email.to": "whoever@example.com"}
func (p Parameters) withDottedName(name string) Parameters {
	var newParameters Parameters

	for key, value := range p {
		fieldName, fieldParamName := dottedParameterNameParts(key)

		if fieldName == name && fieldParamName != "" {
			if newParameters == nil {
				newParameters = Parameters{}
			}

			newParameters[fieldParamName] = value
		}
	}

	return newParameters
}

// namedParametersWithDottedName returns a NamedParameters with the given name and Status and Parameter values
// as calculated from the input Parameters.
//
//   Parameters{"email":"true","email.to":"whoever@example.com"}.namedParametersWithDottedName("email")
//   # NamedParameters{Name: "email", Status: "true", Parameters{"to": "whoever@example.com"}}
func (p Parameters) namedParametersWithDottedName(name string) NamedParameters {
	newParams := NamedParameters{
		Name:       name,
		Parameters: p.withDottedName(name),
	}

	// Status only set explicitly if there was a <name> key is present in Parameters
	if statusValue, ok := p[name]; ok {
		newParams.Status = NewExplicit(statusValue)
	}

	return newParams
}

// dottedNames returns the list of top-level names of fields in Parameters.
func (p Parameters) dottedNames() []string {
	foundNamesMap := map[string]bool{}
	var foundNames []string

	for key := range p {
		fieldName, _ := dottedParameterNameParts(key)

		if _, ok := foundNamesMap[fieldName]; !ok {
			foundNames = append(foundNames, fieldName)
		}
		foundNamesMap[fieldName] = true
	}

	sort.Strings(foundNames)

	return foundNames
}

// namedParametersCollection returns a NamedParametersCollection containing a NamedParameters object
// for each top-level name of Parameters.
func (p Parameters) namedParametersCollection() NamedParametersCollection {
	names := p.dottedNames()
	var newCollection NamedParametersCollection

	for _, name := range names {
		newCollection = append(newCollection, p.namedParametersWithDottedName(name))
	}

	return newCollection
}

// UnmarshalJSON implements custom JSON unmarshaling which assumes the content being unmarshaled is a simple map of strings
// to a single value (string, bool, float, int). It returns an error if a value other than these types is encountered.
func (p *Parameters) UnmarshalJSON(data []byte) error {
	interfaceMap := map[string]interface{}{}
	if err := json.Unmarshal(data, &interfaceMap); err != nil {
		return err
	}

	newP := Parameters{}
	for key, value := range interfaceMap {
		valueV := reflect.ValueOf(value)
		if !valueV.IsValid() {
			// invalid reflect.Value indicates a nil value, which we can safely ignore
			continue
		}

		switch valueV.Kind() {
		default:
			return fmt.Errorf("unable to unmarshal unhandled type %T into Parameters for key %s", value, key)
		case reflect.String, reflect.Bool, reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			newP[key] = fmt.Sprintf("%v", value)
		}
	}

	if len(newP) == 0 {
		return nil
	}

	*p = newP

	return nil
}

// SetURLValues implements custom encoding into url.Values.
func (p Parameters) SetURLValues(key string, v *url.Values) error {
	for paramKey, paramValue := range p {
		valueKey := strings.Join([]string{key, paramKey}, ".")

		v.Set(valueKey, paramValue)
	}

	return nil
}

// UnmarshalJSONForParameters unmarshals JSON data into the given dest interface. dest must be a pointer
// to a struct, and any struct fields with the "parameters" tag must be of the Parameters type.
//
// This method exists to enable unmarshaling of the same level of a JSON document to a struct and also
// to Parameters fields of the same struct.
func UnmarshalJSONForParameters(data []byte, dest interface{}) error {
	destVPtr := reflect.ValueOf(dest)
	if destVPtr.Kind() != reflect.Ptr {
		return fmt.Errorf("attempted UnmarshalJSONForParameters on non-pointer type: %T", dest)
	}

	destV := destVPtr.Elem()
	destT := destV.Type()

	if destT.Kind() != reflect.Struct {
		return fmt.Errorf("attempted UnmarshalJSONForParameters on non-struct type: %T", dest)
	}

	for i := 0; i < destT.NumField(); i++ {
		fieldF := destT.Field(i)
		if !fieldF.IsExported() {
			continue
		}

		fieldTag := fieldF.Tag.Get("parameters")
		if fieldTag == "" {
			continue
		}

		allParams := make(Parameters)
		if fieldF.Type != reflect.TypeOf(allParams) {
			return fmt.Errorf("attempted UnmarshalJSONForParameters on non-Parameters type %T for field %s", destV.Field(i).Interface(), fieldF.Name)
		}

		if err := json.Unmarshal(data, &allParams); err != nil {
			return err
		}

		newParams := allParams.withDottedName(fieldTag)
		newParamsV := reflect.ValueOf(newParams)

		destV.Field(i).Set(newParamsV)
	}

	return nil
}
