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

package paths

import "strings"

// Join returns a string crafted by trimming leading and trailing slashes
// from the input paths and joining them with slashes. Empty paths are retained,
// because each passed path is expected to be a component to be passed to the API,
// which shouldn't be removed just because they are empty.
func Join(path ...string) string {
	parts := make([]string, len(path))

	for i, part := range path {
		parts[i] = strings.Trim(part, "/")
	}

	return strings.Join(parts, "/")
}
