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

package sync

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// directListField implements the Sync interface for list types that are synchronized directly
// between the resource and object.
type directListField[T any] struct {
	value *[]T
	key   string
}

// NewDirectListField returns a Sync that directly associates a list (slice) value with a list
// resource field.
func NewDirectListField[T any](key string, value *[]T) Sync {
	return directListField[T]{
		key:   key,
		value: value,
	}
}

// SyncResource synchronizes schema.ResourceData from the locally stored value.
func (field directListField[T]) SyncResource(d *schema.ResourceData) error {
	return d.Set(field.key, *field.value)
}

// SyncObject synchronizes the locally stored value from schema.ResourceData.
func (field directListField[T]) SyncObject(d *schema.ResourceData) error {
	resourceValueI := d.Get(field.key)
	// nil interface = unknown key
	if resourceValueI == nil {
		return fmt.Errorf("resource: unknown key %q", field.key)
	}

	resourceValue, ok := resourceValueI.([]interface{})
	if !ok {
		return fmt.Errorf("resource: key %q not a slice (%T)", field.key, resourceValueI)
	}

	newValues := make([]T, len(resourceValue))
	for i, resourceValue := range resourceValue {
		newValue, ok := resourceValue.(T)
		if !ok {
			return fmt.Errorf("resource: key %q not a slice of type %T", field.key, *new(T))
		}

		newValues[i] = newValue
	}

	*field.value = newValues

	return nil
}

// GetObject returns the locally stored value.
func (field directListField[T]) GetObject() interface{} {
	return field.value
}
