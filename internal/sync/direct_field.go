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

// directField implements the Sync interface for simple types that are synchronized directly
// between the resource and object.
type directField[T any] struct {
	key   string
	value *T
}

// NewDirectField returns a Sync that directly associates a value with a resource field.
//
// This direct association is generally not sufficient for lists, sets, or maps,
// as the resource value is stored with stored elements as interface{}, not
// the concrete types the referenced value likely contains.
func NewDirectField[T any](key string, value *T) Sync {
	return directField[T]{
		key:   key,
		value: value,
	}
}

// SyncResource synchronizes schema.ResourceData from the locally stored value.
func (field directField[T]) SyncResource(d *schema.ResourceData) error {
	return d.Set(field.key, *field.value)
}

// SyncObject synchronizes the locally stored value from schema.ResourceData.
func (field directField[T]) SyncObject(d *schema.ResourceData) error {
	resourceValueI := d.Get(field.key)
	// nil interface = unknown key
	if resourceValueI == nil {
		return fmt.Errorf("resource: likely unknown key %q", field.key)
	}

	resourceValueT, ok := resourceValueI.(T)
	if !ok {
		return fmt.Errorf("resource: unable to assign resource type %T to object type %T", resourceValueI, *new(T))
	}

	*field.value = resourceValueT

	return nil
}

// GetObject returns the locally stored value.
func (field directField[T]) GetObject() interface{} {
	return field.value
}

// pointerToValueOf returns a new pointer to a copy of the given value. It is used to ease testing
// without needing an existing variable to reference.
func pointerToValueOf[T any](value T) *T {
	return &value
}
