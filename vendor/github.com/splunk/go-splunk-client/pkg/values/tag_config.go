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

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// tagConfig represents a struct tag's intended configuration.
type tagConfig struct {
	// fields are exported (though type is not) so they can be set via reflection
	Name      string
	Omitzero  bool
	Fillempty bool
	Anonymize bool
}

// parseTagConfig returns a new tagConfig by parsing a tag string.
func parseTagConfig(tag string) (tagConfig, error) {
	newConfig := tagConfig{}

	if tag == "" {
		return newConfig, nil
	}

	tagParts := strings.Split(tag, ",")

	newConfig.Name = tagParts[0]
	flags := tagParts[1:]

	// for newOptsV to be settable, it needs to be created from a pointer
	newOptsV := reflect.ValueOf(&newConfig).Elem()

	for _, flag := range flags {
		// a flag may be an empty string, specifically to permit `values:"-,"` to be used to
		// name a key "-" without omitting the field
		if flag == "" {
			continue
		}

		// capitalize the first letter of the flag so it can match a settable (exported) field
		flagFieldName := cases.Title(language.English).String(flag)

		flagField := newOptsV.FieldByName(flagFieldName)
		if !flagField.IsValid() && flagField.Kind() == reflect.Bool {
			return tagConfig{}, fmt.Errorf("values: unknown flag: %s", flag)
		}

		if !flagField.CanSet() {
			return tagConfig{}, fmt.Errorf("values: unexpected error, unable to set option field %s", flagFieldName)
		}

		// checked to be Bool above, so this should be safe
		flagField.SetBool(true)
	}

	return newConfig, nil
}
