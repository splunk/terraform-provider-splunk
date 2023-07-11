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

package splunk

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/splunk/go-splunk-client/pkg/client"
	"github.com/splunk/terraform-provider-splunk/internal/sync"
)

// useLegacyClient returns the first of:
//
// * true if the resource's use_client field set to "legacy"
//
// * false if the resource's use_client field is not empty and set to anything other than "legacy" (only "external" is possible)
//
// * true if the provider's use_client_default field set to "legacy"
//
// * false if the provider's use_client_default field set to anything other than "legacy" (only "external" is possible)
func useLegacyClient(provider *SplunkProvider, d *schema.ResourceData) bool {
	if resourceUseClientI, ok := d.GetOk("use_client"); ok {
		resourceUseClient := resourceUseClientI.(string)
		return resourceUseClient == useClientLegacy
	}

	return provider.useClientDefault == useClientLegacy
}

// createFunc returns a schema.CreateFunc for the Sync returned by the given function.
//
// During the transition between the legacy and external Splunk clients, it will return legacyFunction
// if the provider configuration or the resource configuration sets use_legacy_client=true.
func createFunc(syncFn func() sync.SyncGetter, legacyFunction schema.CreateFunc) schema.CreateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		provider := meta.(*SplunkProvider)

		if useLegacyClient(provider, d) {
			return legacyFunction(d, meta)
		}

		c := provider.ExternalClient

		s := syncFn()

		if err := s.SyncObject(d); err != nil {
			return err
		}

		if err := c.Create(s.GetObject()); err != nil {
			return err
		}

		return readFunc(syncFn, legacyFunction)(d, meta)
	}
}

// readFunc returns a schema.CreateFunc for the Sync returned by the given function.
//
// During the transition between the legacy and external Splunk clients, it will return legacyFunction
// if the provider configuration or the resource configuration sets use_legacy_client=true.
func readFunc(syncFn func() sync.SyncGetter, legacyFunction schema.CreateFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		provider := meta.(*SplunkProvider)

		if useLegacyClient(provider, d) {
			return legacyFunction(d, meta)
		}

		c := provider.ExternalClient

		s := syncFn()

		if err := s.SyncObject(d); err != nil {
			return err
		}

		if err := c.Read(s.GetObject()); err != nil {
			if clientErr, ok := err.(client.Error); ok {
				if clientErr.Code == client.ErrorNotFound {
					d.SetId("")

					return nil
				}
			}
			return err
		}

		if err := s.SyncResource(d); err != nil {
			return err
		}

		return nil
	}
}

// updateFunc returns a schema.UpdateFunc for the Sync returned by the given function.
//
// During the transition between the legacy and external Splunk clients, it will return legacyFunction
// if the provider configuration or the resource configuration sets use_legacy_client=true.
func updateFunc(syncFn func() sync.SyncGetter, legacyFunction schema.CreateFunc) schema.UpdateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		provider := meta.(*SplunkProvider)

		if useLegacyClient(provider, d) {
			return legacyFunction(d, meta)
		}

		c := provider.ExternalClient

		s := syncFn()

		if err := s.SyncObject(d); err != nil {
			return err
		}

		if err := c.Update(s.GetObject()); err != nil {
			return err
		}

		return readFunc(syncFn, legacyFunction)(d, meta)
	}
}

// deleteFunc returns a schema.DeleteFunc for the Sync returned by the given function.
//
// During the transition between the legacy and external Splunk clients, it will return legacyFunction
// if the provider configuration or the resource configuration sets use_legacy_client=true.
func deleteFunc(syncFn func() sync.SyncGetter, legacyFunction schema.CreateFunc) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		provider := meta.(*SplunkProvider)

		if useLegacyClient(provider, d) {
			return legacyFunction(d, meta)
		}

		c := provider.ExternalClient

		s := syncFn()

		if err := s.SyncObject(d); err != nil {
			return err
		}

		if err := c.Delete(s.GetObject()); err != nil {
			return err
		}

		return nil
	}
}
