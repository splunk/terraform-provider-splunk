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

package client

import (
	"strings"

	"github.com/splunk/go-splunk-client/pkg/internal/paths"
)

// Namespace is a Splunk object's namespace, consisting of a User and App. It is valid only if both User
// and App are set, or both are unset.
type Namespace struct {
	User string
	App  string
}

// namespacePath returns the namespace path. If the resulting path is invalid, it will be returned
// along with an error.
func (ns Namespace) namespacePath() (string, error) {
	var path string

	// absence of both user/app indicates global context
	if (ns.User == "") && (ns.App == "") {
		path = "services"
	} else {
		path = paths.Join("servicesNS", ns.User, ns.App)
	}

	return path, ns.validate()
}

// GetServicePath implements custom GetServicePath encoding. It returns the given
// path back, which has the effect of using the ID field's struct tag as
// its GetServicePath.
func (ns Namespace) GetServicePath(path string) (string, error) {
	nsPath, err := ns.namespacePath()
	if err != nil {
		return "", err
	}

	return paths.Join(nsPath, path), nil
}

// parseNamespace parses a string ID into a Namespace object.
func parseNamespace(id string) (Namespace, []string, error) {
	// start with a clean IDFields
	newNS := Namespace{}

	// an empty id is a "clean slate", with no namespace info (user/app), endpoint, or title
	if id == "" {
		return newNS, nil, nil
	}

	reverseStrings := func(input []string) []string {
		output := make([]string, 0, len(input))

		for i := len(input) - 1; i >= 0; i-- {
			output = append(output, input[i])
		}

		return output
	}

	// we work backwards in the id path segments to find the namespace, so reverse idPartStrings
	idPartStrings := strings.Split(id, "/")
	idPartStrings = reverseStrings(idPartStrings)

	for i, idPartString := range idPartStrings {
		// if we've made it to the "services" segment we have an empty Namespace
		if idPartString == "services" {
			return newNS, reverseStrings(idPartStrings[0:i]), nil
		}

		if idPartString == "servicesNS" {
			// if we've made it to the "servicesNS" segment we have a user/app namespace

			if i < 2 {
				return Namespace{}, nil, wrapError(ErrorID, nil, "unable to parse ID, servicesNS found without user/app: %s", id)
			}

			newNS.User = idPartStrings[i-1]
			newNS.App = idPartStrings[i-2]

			return newNS, reverseStrings(idPartStrings[0 : i-2]), nil
		}
	}

	return Namespace{}, nil, wrapError(ErrorID, nil, "client: unable to parse Namespace, missing services or servicesNS: %s", id)
}

// validate returns an error if the Namespace is invalid.
func (ns Namespace) validate() error {
	if (ns.User == "") != (ns.App == "") {
		return wrapError(ErrorNamespace, nil, "client: invalid Namespace, user and app must both be empty or non-empty")
	}

	return nil
}
