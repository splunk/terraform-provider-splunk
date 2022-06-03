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

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// indirectObject stores a pointer to an object, and a Sync that indirectly manages
// the stored object.
type indirectObject[T any] struct {
	object *T
	sync   Sync
}

// NewIndirectObject returns a SyncGetter for the given value and sync. sync is expected
// to indirectly manage value, typically via operating on its member fields.
//
// For example, the stored value may be a full entry.SAMLGroup, with the given sync
// operating on its Roles field.
//
// The stored object is never managed directly by the returned SyncGetter, but only by
// the provided Sync.
func NewIndirectObject[T any](value *T, sync Sync) SyncGetter {
	return indirectObject[T]{
		object: value,
		sync:   sync,
	}
}

// SyncResource runs the indirect Sync's SyncResource method.
func (o indirectObject[T]) SyncResource(d *schema.ResourceData) error {
	return o.sync.SyncResource(d)
}

// SyncResource runs the indirect Sync's SyncObject method.
func (o indirectObject[T]) SyncObject(d *schema.ResourceData) error {
	return o.sync.SyncObject(d)
}

// GetObject returns the stored object.
func (o indirectObject[T]) GetObject() interface{} {
	return o.object
}
