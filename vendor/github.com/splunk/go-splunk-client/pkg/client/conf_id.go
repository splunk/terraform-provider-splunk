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
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/splunk/go-splunk-client/pkg/internal/paths"
)

// ConfID represents the ID of configs/conf-<file> resources.
type ConfID struct {
	Namespace Namespace

	File   string
	Stanza string
}

// parseConfID returns a new ConfsID by parsing the ID URL string.
func parseConfID(idURL string) (ConfID, error) {
	newNS, remnants, err := parseNamespace(idURL)
	if err != nil {
		return ConfID{}, err
	}

	if len(remnants) < 2 {
		return ConfID{}, wrapError(ErrorID, nil, "client: parseNamespace didn't return remnants for ConfID.File and ConfID.Stanza")
	}

	fileRemnant := remnants[len(remnants)-2]
	stanzaRemnant := remnants[len(remnants)-1]

	r := regexp.MustCompile("^conf-(.+)$")
	// remnants[0]
	foundFileStrings := r.FindStringSubmatch(fileRemnant)
	if foundFileStrings == nil {
		return ConfID{}, wrapError(ErrorID, nil, "client: unable to parse %q for ConfID.File", fileRemnant)
	}
	// [0] is the full match, [1] is the first capture group
	foundFile := foundFileStrings[1]

	return ConfID{
		Namespace: newNS,
		File:      foundFile,
		Stanza:    stanzaRemnant,
	}, nil
}

// Parse sets the ID's value to match what is parsed from the given ID URL.
func (confID *ConfID) Parse(idURL string) error {
	newConfID, err := parseConfID(idURL)
	if err != nil {
		return err
	}

	*confID = newConfID

	return nil
}

// GetServicePath implements custom GetServicePath encoding.
func (confID ConfID) GetServicePath(path string) (string, error) {
	if confID.File == "" {
		return "", wrapError(ErrorID, nil, "client: attempted ConfID.GetServicePath() with empty File")
	}

	nsServicePath, err := confID.Namespace.GetServicePath(path)
	if err != nil {
		return "", err
	}

	return paths.Join(
		nsServicePath,
		fmt.Sprintf("conf-%s", confID.File),
	), nil
}

// GetEntryPath implements custom GetEntryPath encoding.
func (confID ConfID) GetEntryPath(path string) (string, error) {
	servicePath, err := confID.GetServicePath(path)
	if err != nil {
		return "", err
	}

	return paths.Join(servicePath, url.PathEscape(confID.Stanza)), nil
}

// UnmarshalJSON implements custom JSON unmarshaling for ConfID.
func (confID *ConfID) UnmarshalJSON(data []byte) error {
	idString := ""
	if err := json.Unmarshal(data, &idString); err != nil {
		return wrapError(ErrorID, err, "client: unable to unmarshal %q as string", data)
	}

	if err := confID.Parse(idString); err != nil {
		return err
	}

	return nil
}

// SetURLValues implements custom url.Query encoding of ConfID. It adds a field "name" for the ConfID's
// Stanza. If the Title value is empty, it returns an error, as there are no scenarios where a ConfID
// object is expected to be POSTed with an empty Stanza.
func (confID ConfID) SetURLValues(key string, v *url.Values) error {
	if confID.Stanza == "" {
		return wrapError(ErrorID, nil, "client: attempted SetURLValues on ConfID with empty Stanza")
	}

	v.Add("name", confID.Stanza)

	return nil
}
