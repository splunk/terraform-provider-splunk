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

package entry

import (
	"net/url"

	"github.com/splunk/go-splunk-client/pkg/attributes"
	"github.com/splunk/go-splunk-client/pkg/client"
)

// SavedSearchContent defines the content of a Savedsearch object.
type SavedSearchContent struct {
	Actions                           attributes.NamedParametersCollection `json:"-"                                     values:"action,omitzero"                                named_parameters_collection:"action"`
	AlertDigestMode                   attributes.Explicit[bool]            `json:"alert.digest_mode"                     values:"alert.digest_mode,omitzero"`
	AlertExpires                      attributes.Explicit[string]          `json:"alert.expires"                         values:"alert.expires,omitzero"`
	AlertSeverity                     attributes.Explicit[string]          `json:"alert.severity"                        values:"alert.severity,omitzero"`
	AlertSuppress                     attributes.Explicit[bool]            `json:"alert.suppress"                        values:"alert.suppress,omitzero"`
	AlertSuppressFields               attributes.Explicit[string]          `json:"alert.suppress.fields"                 values:"alert.suppress.fields,omitzero"`
	AlertSuppressGroupName            attributes.Explicit[string]          `json:"alert.suppress.group_name"             values:"alert.suppress.group_name,omitzero"`
	AlertSuppressPeriod               attributes.Explicit[int]             `json:"alert.suppress.period"                 values:"alert.suppress.period,omitzero"`
	AlertTrack                        attributes.Explicit[string]          `json:"alert.track"                           values:"alert.track,omitzero"`
	AlertComparator                   attributes.Explicit[string]          `json:"alert_comparator"                      values:"alert_comparator,omitzero"`
	AlertCondition                    attributes.Explicit[string]          `json:"alert_condition"                       values:"alert_condition,omitzero"`
	AlertThreshold                    attributes.Explicit[int]             `json:"alert_threshold"                       values:"alert_threshold,omitzero"`
	AlertType                         attributes.Explicit[string]          `json:"alert_type"                            values:"alert_type,omitzero"`
	AllowSkew                         attributes.Explicit[string]          `json:"allow_skew"                            values:"allow_skew,omitzero"`
	AutoSummarize                     attributes.Explicit[bool]            `json:"auto_summarize"                        values:"auto_summarize,omitzero"`
	AutoSummarizeCommand              attributes.Explicit[string]          `json:"auto_summarize.command"                values:"auto_summarize.command,omitzero"`
	AutoSummarizeCronSchedule         attributes.Explicit[string]          `json:"auto_summarize.cron_schedule"          values:"auto_summarize.cron_schedule,omitzero"`
	AutoSummarizeDispatchEarliestTime attributes.Explicit[string]          `json:"auto_summarize.dispatch.earliest_time" values:"auto_summarize.dispatch.earliest_time,omitzero"`
	AutoSummarizeDispatchLatestTime   attributes.Explicit[string]          `json:"auto_summarize.dispatch.latest_time"   values:"auto_summarize.dispatch.latest_time,omitzero"`
	AutoSummarizeDispatchTimeFormat   attributes.Explicit[string]          `json:"auto_summarize.dispatch.time_format"   values:"auto_summarize.dispatch.time_format,omitzero"`
	AutoSummarizeDispatchTtl          attributes.Explicit[string]          `json:"auto_summarize.dispatch.ttl"           values:"auto_summarize.dispatch.ttl,omitzero"`
	AutoSummarizeMaxConcurrent        attributes.Explicit[int]             `json:"auto_summarize.max_concurrent"         values:"auto_summarize.max_concurrent,omitzero"`
	AutoSummarizeMaxDisabledBuckets   attributes.Explicit[int]             `json:"auto_summarize.max_disabled_buckets"   values:"auto_summarize.max_disabled_buckets,omitzero"`
	AutoSummarizeMaxSummaryRatio      attributes.Explicit[int]             `json:"auto_summarize.max_summary_ratio"      values:"auto_summarize.max_summary_ratio,omitzero"`
	AutoSummarizeMaxSummarySize       attributes.Explicit[int]             `json:"auto_summarize.max_summary_size"       values:"auto_summarize.max_summary_size,omitzero"`
	AutoSummarizeMaxTime              attributes.Explicit[int]             `json:"auto_summarize.max_time"               values:"auto_summarize.max_time,omitzero"`
	AutoSummarizeSuspendPeriod        attributes.Explicit[string]          `json:"auto_summarize.suspend_period"         values:"auto_summarize.suspend_period,omitzero"`
	AutoSummarizeTimespan             attributes.Explicit[string]          `json:"auto_summarize.timespan"               values:"auto_summarize.timespan,omitzero"`
	CronSchedule                      attributes.Explicit[string]          `json:"cron_schedule"                         values:"cron_schedule,omitzero"`
	Description                       attributes.Explicit[string]          `json:"description"                           values:"description,omitzero"`
	Disabled                          attributes.Explicit[bool]            `json:"disabled"                              values:"disabled,omitzero"`
	Dispatch                          attributes.NamedParametersCollection `json:"-"                                     values:"dispatch,omitzero"                              named_parameters_collection:"dispatch" `
	DispatchAs                        attributes.Explicit[string]          `json:"dispatchAs"                            values:"dispatchAs,omitzero"`
	Displayview                       attributes.Explicit[string]          `json:"displayview"                           values:"displayview,omitzero"`
	DurableBackfillType               attributes.Explicit[string]          `json:"durable.backfill_type"                 values:"durable.backfill_type,omitzero"`
	DurableLagTime                    attributes.Explicit[int]             `json:"durable.lag_time"                      values:"durable.lag_time,omitzero"`
	DurableMaxBackfillIntervals       attributes.Explicit[int]             `json:"durable.max_backfill_intervals"        values:"durable.max_backfill_intervals,omitzero"`
	DurableTrackTimeType              attributes.Explicit[string]          `json:"durable.track_time_type"               values:"durable.track_time_type,omitzero"`
	IsScheduled                       attributes.Explicit[bool]            `json:"is_scheduled"                          values:"is_scheduled,omitzero"`
	IsVisible                         attributes.Explicit[bool]            `json:"is_visible"                            values:"is_visible,omitzero"`
	MaxConcurrent                     attributes.Explicit[int]             `json:"max_concurrent"                        values:"max_concurrent,omitzero"`
	Name                              attributes.Explicit[string]          `json:"name"                                  values:"name,omitzero"`
	NextScheduledTime                 attributes.Explicit[string]          `json:"next_scheduled_time"                   values:"next_scheduled_time,omitzero"`
	QualifiedSearch                   attributes.Explicit[string]          `json:"qualifiedSearch"                       values:"qualifiedSearch,omitzero"`
	RealtimeSchedule                  attributes.Explicit[bool]            `json:"realtime_schedule"                     values:"realtime_schedule,omitzero"`
	RequestUiDispatchApp              attributes.Explicit[string]          `json:"request.ui_dispatch_app"               values:"request.ui_dispatch_app,omitzero"`
	RequestUiDispatchView             attributes.Explicit[string]          `json:"request.ui_dispatch_view"              values:"request.ui_dispatch_view,omitzero"`
	RestartOnSearchpeerAdd            attributes.Explicit[bool]            `json:"restart_on_searchpeer_add"             values:"restart_on_searchpeer_add,omitzero"`
	RunNTimes                         attributes.Explicit[int]             `json:"run_n_times"                           values:"run_n_times,omitzero"`
	RunOnStartup                      attributes.Explicit[bool]            `json:"run_on_startup"                        values:"run_on_startup,omitzero"`
	SchedulePriority                  attributes.Explicit[string]          `json:"schedule_priority"                     values:"schedule_priority,omitzero"`
	ScheduleWindow                    attributes.Explicit[string]          `json:"schedule_window"                       values:"schedule_window,omitzero"`
	Search                            attributes.Explicit[string]          `json:"search"                                values:"search,omitzero"`
	Vsid                              attributes.Explicit[string]          `json:"vsid"                                  values:"vsid"`
	WorkloadPool                      attributes.Explicit[string]          `json:"workload_pool"                         values:"workload_pool"`
}

// AddURLValues implements custom additional encoding to url.Values.
func (content SavedSearchContent) AddURLValues(key string, v *url.Values) error {
	// The Splunk REST API returns savedsearch action status like "action.email=1", but doesn't honor that format
	// for setting the action statuses. To set an action status, you must pass "actions=action1,action2" formatted
	// values. Here we iterate through the enabled actions and add a url.Values entry for all enabled actions.
	//
	// If Actions is empty (not nil), we "clear" the enabled actions list by setting a single empty value for "actions".

	if content.Actions != nil && len(content.Actions) == 0 {
		v.Add("actions", "")
	}

	for _, enabledActionName := range content.Actions.EnabledNames() {
		v.Add("actions", enabledActionName)
	}

	return nil
}

// UnmarshalJSON implements custom JSON unmarshaling.
func (content *SavedSearchContent) UnmarshalJSON(data []byte) error {
	type contentAlias SavedSearchContent
	var newAliasedContent contentAlias

	if err := attributes.UnmarshalJSONForNamedParametersCollections(data, &newAliasedContent); err != nil {
		return err
	}

	*content = SavedSearchContent(newAliasedContent)

	return nil
}

// SavedSearch defines a Splunk savedsearch.
type SavedSearch struct {
	ID      client.ID          `service:"saved/searches" selective:"create"`
	Content SavedSearchContent `json:"content" values:",anonymize"`
}
