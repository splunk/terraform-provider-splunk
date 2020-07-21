package model

// SearchConfig represents the search job post params
type SearchConfig struct {

	//f this is not documented on docs.splunk.com
	Fields []string `json:"f"`

	//sample_ratio this is not documented on docs.splunk.com
	SampleRatio string `json:"sample_ratio"`

	//adhoc_search_level    String        Use one of the following search modes.
	//[ verbose | fast | smart ]
	AdhocSearchLevel string `json:"adhoc_search_level"`

	//auto_cancel    Number    0    If specified, the job automatically cancels after this many seconds of inactivity. (0 means never auto-cancel)
	AutoCancel *uint `json:"auto_cancel"`

	//auto_finalize_ec    Number    0    Auto-finalize the search after at least this many events are processed.
	//Specify 0 to indicate no limit.
	AutoFinalizeEventCount *uint `json:"auto_finalize_ec"`

	//auto_pause    Number    0    If specified, the search job pauses after this many seconds of inactivity. (0 means never auto-pause.)
	//To restart a paused search job, specify unpause as an action to POST search/jobs/{search_id}/control.
	//auto_pause only goes into effect once. Unpausing after auto_pause does not put auto_pause into effect again.
	AutoPause *uint `json:"auto_pause"`

	//custom fields
	CustomFields map[string]interface{}

	//earliest_time    String        Specify a time string. Sets the earliest (inclusive), respectively, time bounds for the search.
	//The time string can be either a UTC time (with fractional seconds), a relative time specifier (to now) or a formatted time string. Refer to Time modifiers for search for information and examples of specifying a time string.
	//Compare to index_earliest parameter. Also see comment for the search_mode parameter.
	EarliestTime string `json:"earliest_time"`

	//enable_lookups    Boolean    true    Indicates whether lookups should be applied to events.
	//Specifying true (the default) may slow searches significantly depending on the nature of the lookups.
	EnableLookUps *bool `json:"enable_lookups"`

	//exec_mode    Enum    normal    Valid values: (blocking | oneshot | normal)
	//If set to normal, runs an asynchronous search.
	//If set to blocking, returns the sid when the job is complete.
	//If set to oneshot, returns results in the same call. In this case, you can specify the format for the output (for example, json output) using the output_mode parameter as described in GET search/jobs/export. Default format for output is xml.
	ExecuteMode string `json:"exec_mode"`

	//force_bundle_replication    Boolean    false    Specifies whether this search should cause (and wait depending on the value of sync_bundle_replication) for bundle synchronization with all search peers.
	ForceBundleReplication *bool `json:"force_bundle_replication"`

	//id    String        Optional string to specify the search ID (<sid>). If unspecified, a random ID is generated.
	ID string `json:"id"`

	//index_earliest    String        Specify a time string. Sets the earliest (inclusive), respectively, time bounds for the search, based on the index time bounds.
	//The time string can be either a UTC time (with fractional seconds), a relative time specifier (to now) or a formatted time string. Compare to earliest_time parameter. Also see comment for the search_mode parameter.
	//Refer to Time modifiers for search for information and examples of specifying a time string.
	IndexEarliestTime string `json:"index_earliest"`

	//index_latest    String        Specify a time string. Sets the latest (exclusive), respectively, time bounds for the search, based on the index time bounds.
	//The time string can be either a UTC time (with fractional seconds), a relative time specifier (to now) or a formatted time string.
	//Refer to Time modifiers for search for information and examples of specifying a time string.
	//Compare to latest_time parameter. Also see comment for the search_mode parameter.
	IndexLatestTime string `json:"index_latest"`

	//indexedRealtime    Boolean        Indicate whether or not to used indexed-realtime mode for real-time searches.
	IndexedRealTime *bool `json:"indexedRealtime"`

	//indexedRealtimeOffset    Number        Set disk sync delay for indexed real-time search (seconds).
	IndexedRealTimeOffset *uint `json:"indexedRealtimeOffset"`

	//latest_time    String        Specify a time string. Sets the latest (exclusive), respectively, time bounds for the search.
	//The time string can be either a UTC time (with fractional seconds), a relative time specifier (to now) or a formatted time string.
	//Refer to Time modifiers for search for information and examples of specifying a time string.
	//Compare to index_latest parameter. Also see comment for the search_mode parameter.
	LatestTime string `json:"latest_time"`

	//max_count    Number    10000    The number of events that can be accessible in any given status bucket.
	//Also, in transforming mode, the maximum number of results to store. Specifically, in all calls, codeoffset+count max_count.
	MaxCount *uint `json:"max_count"`

	//max_time    Number    0    The number of seconds to run this search before finalizing. Specify 0 to never finalize.
	MaxTime *uint `json:"max_time"`

	//namespace    String        The application namespace in which to restrict searches.
	//The namespace corresponds to the identifier recognized in the /services/apps/local endpoint.
	Namespace string `json:"namespace"`

	//now    String    current system time    Specify a time string to set the absolute time used for any relative time specifier in the search. Defaults to the current system time.
	//You can specify a relative time modifier for this parameter. For example, specify +2d to specify the current time plus two days.
	//If you specify a relative time modifier both in this parameter and in the search string, the search string modifier takes precedence.
	//Refer to Time modifiers for search for details on specifying relative time modifiers.
	Now string `json:"now"`

	//reduce_freq    Number    0    Determines how frequently to run the MapReduce reduce phase on accumulated map values.
	ReduceFrequency *uint `json:"reduce_freq"`

	//reload_macros    Boolean    true    Specifies whether to reload macro definitions from macros.conf.
	//Default is true.
	ReloadMacros *bool `json:"reload_macros"`

	//remote_server_list    String    empty list    Comma-separated list of (possibly wildcarded) servers from which raw events should be pulled. This same server list is to be used in subsearches.
	RemoteServerList string `json:"remote_server_list"`

	//replay_speed    Number greater than 0        Indicate a real-time search replay speed factor. For example, 1 indicates normal speed. 0.5 indicates half of normal speed, and 2 indicates twice as fast as normal.
	//earliest_time and latest_time arguments must indicate a real-time time range to use replay options.
	//Use replay_speed with replay_et and replay_lt relative times to indicate a speed and time range for the replay. For example,
	//replay_speed = 10
	//replay_et = -d@d
	//replay_lt = -@d
	//specifies a replay at 10x speed, as if the "wall clock" time starts yesterday at midnight and ends when it reaches today at midnight.
	//For more information about using relative time modifiers, see Search time modifiers in the Search reference.
	ReplaySpeed *uint `json:"replay_speed"`

	//replay_et    Time modifier string        Relative "wall clock" start time for the replay.
	ReplayEarliestTime string `json:"replay_et"`

	//replay_lt    Time modifier string.        Relative end time for the replay clock. The replay stops when clock time reaches this time.
	ReplayLatestTime string `json:"replay_lt"`

	//reuse_max_seconds_ago    Number        Specifies the number of seconds ago to check when an identical search is started and return the job's search ID instead of starting a new job.
	ReuseMaxSecondsAgo *uint `json:"reuse_max_seconds_ago"`

	//rf    String        Adds a required field to the search. There can be multiple rf POST arguments to the search.
	//These fields, even if not referenced or used directly by the search, are still included by the events and summary endpoints. Splunk Web uses these fields to prepopulate panels in the Search view.
	//Consider using this form of passing the required fields to the search instead of the deprecated required_field_list. If both rf and required_field_list are provided, the union of the two lists is used.
	RequiredFields []string `json:"rf"`

	//rt_blocking    Boolean    false    For a real-time search, indicates if the indexer blocks if the queue for this search is full.
	RealTimeBlocking *bool `json:"rt_blocking"`

	//rt_indexfilter    Boolean    true    For a real-time search, indicates if the indexer prefilters events.
	RealTimeIndexFilter *bool `json:"rt_indexfilter"`

	//rt_maxblocksecs    Number    60    For a real-time search with rt_blocking set to true, the maximum time to block.
	//Specify 0 to indicate no limit.
	RealTimeMaxBlockSeconds *uint `json:"rt_maxblocksecs"`

	//rt_queue_size    Number    10000 events    For a real-time search, the queue size (in events) that the indexer should use for this search.
	RealTimeQueueSize *uint `json:"rt_queue_size"`

	//search required    String        The search language string to execute, taking results from the local and remote servers.
	//Examples:
	//"search *"
	//"search * | outputcsv"
	Search string `json:"search"`

	//search_listener    String        Registers a search state listener with the search.
	//Use the format:
	//search_state;results_condition;http_method;uri;
	//For example:
	//search_listener=onResults;true;POST;/servicesNS/admin/search/saved/search/foobar/notify;
	SearchListener string `json:"search_listener"`

	//search_mode    Enum    normal    Valid values: (normal | realtime)
	//If set to realtime, search runs over live data. A real-time search may also be indicated by earliest_time and latest_time variables starting with 'rt' even if the search_mode is set to normal or is unset. For a real-time search, if both earliest_time and latest_time are both exactly 'rt', the search represents all appropriate live data received since the start of the search.
	//Additionally, if earliest_time and/or latest_time are 'rt' followed by a relative time specifiers then a sliding window is used where the time bounds of the window are determined by the relative time specifiers and are continuously updated based on the wall-clock time.
	SearchMode string `json:"search_mode"`
	//todo should we support realtime?

	//spawn_process    Boolean    true    Specifies whether the search should run in a separate spawned process. Default is true.
	//Searches against indexes must run in a separate process.
	SpawnProcess *bool `json:"spawn_process"`

	//status_buckets    Number    0    The most status buckets to generate.
	//0 indicates to not generate timeline information.
	StatusBuckets *uint `json:"status_buckets"`

	//sync_bundle_replication    Boolean        Specifies whether this search should wait for bundle replication to complete.
	SyncBundleReplication *bool `json:"sync_bundle_replication"`

	//time_format    String     %FT%T.%Q%:z    Used to convert a formatted time string from {start,end}_time into UTC seconds. The default value is the ISO-8601 format.
	TimeFormat string `json:"time_format"`

	//timeout    Number    86400    The number of seconds to keep this search after processing has stopped.
	Timeout *uint `json:"timeout"`

	// config json string Specified the datasets and rules used by this search.
	Config string `json:"config"`
}
