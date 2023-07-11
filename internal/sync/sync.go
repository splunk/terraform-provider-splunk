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

// Package sync provides a synchronization abstraction between object types
// and schema.ResourceData.
package sync

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// Sync is the interface for types that synchronize between objects and schema.ResourceData.
type Sync interface {
	// SyncResource synchronizes schema.ResourceData from the locally stored object.
	SyncResource(*schema.ResourceData) error

	// SyncObject synchronizes the locally stored object from schema.ResourceData.
	SyncObject(*schema.ResourceData) error
}

// SyncGetter is the interface for types that implement Sync and GetObject().
type SyncGetter interface {
	Sync

	// GetObject returns the locally stored object.
	GetObject() interface{}
}

// syncs is a slice of Syncs.
type syncs []Sync

// SyncResource runs SyncResource for each member Sync.
func (s syncs) SyncResource(d *schema.ResourceData) error {
	for _, sync := range s {
		if err := sync.SyncResource(d); err != nil {
			return err
		}
	}

	return nil
}

// SyncObject runs SyncObject for each member Sync.
func (s syncs) SyncObject(d *schema.ResourceData) error {
	for _, sync := range s {
		if err := sync.SyncObject(d); err != nil {
			return err
		}
	}

	return nil
}

// ComposeSync returns a new Sync that is composed of the given Syncs.
func ComposeSync(s ...Sync) Sync {
	newSyncs := make(syncs, len(s))

	for i, sync := range s {
		newSyncs[i] = sync
	}

	return newSyncs
}
