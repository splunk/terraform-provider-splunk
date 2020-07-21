package model

// OutputMode represents the out_mode param
type OutputMode string

const (
	// JSON fetch results in json
	JSON OutputMode = "json"
)

// EventFetchConfig specifies params for fetching /search/jobs/{sid}/events
type EventFetchConfig struct {
	Pagination       PaginationParams
	EarliestTime     string     `key:"earliest_time"`
	Fields           []string   `key:"f"`
	LatestTime       string     `key:"latest_time"`
	MaxLines         *uint      `key:"max_lines"`
	OutputMode       OutputMode `key:"output_mode"`
	TimeFormat       string     `key:"time_format"`
	OutputTimeFormat string     `key:"output_time_format"`
	Search           string     `key:"search"`
	TruncationMode   string     `key:"truncation_mode"`
	Segmentation     string     `key:"segmentation"`
}

// NewDefaultEventFetchConfig creates parameters according to Splunk Enterprise defaults
func NewDefaultEventFetchConfig() *EventFetchConfig {
	return &EventFetchConfig{
		OutputMode: JSON,
	}
}

// NewDefaultResultFetchConfig creates parameters according to Splunk Enterprise defaults
func NewDefaultResultFetchConfig() *ResultFetchConfig {
	return &ResultFetchConfig{
		OutputMode: JSON,
	}
}

// ResultFetchConfig specifies params for fetching /search/jobs/{sid}/results
type ResultFetchConfig struct {
	Pagination           PaginationParams
	AddSummaryToMetadata bool       `key:"add_summary_to_metadata"`
	Fields               []string   `key:"f"`
	OutputMode           OutputMode `key:"output_mode"`
	Search               string     `key:"search"`
}

// SortDir specifies ascending or descending sort order
type SortDir string

// Valid values for SortDir
const (
	ASC  SortDir = "asc"
	DESC SortDir = "desc"
)

// SortMode specifies how response content is sorted
type SortMode string

// Valid values for sort_mode
const (
	AUTO      SortMode = "auto"
	ALPHA     SortMode = "alpha"
	ALPHACASE SortMode = "alpha_case"
	NUM       SortMode = "num"
)

// FilteringParams specifies filtering parameters for certain supported requests
type FilteringParams struct {
	Search    string   `key:"search"`
	SortDir   SortDir  `key:"sort_dir"`
	SortKey   string   `key:"sort_key"`
	SortMode  SortMode `key:"sort_mode"`
	Summarize bool     `key:"summarize"`
}

// PaginationParams specifies pagination parameters for certain supported requests
type PaginationParams struct {
	Count  uint `key:"count"`
	Offset uint `key:"offset"`
}

// NewDefaultPaginationParams creates search pagination parameters according to Splunk Enterprise defaults
func NewDefaultPaginationParams() *PaginationParams {
	return &PaginationParams{
		Count:  30,
		Offset: 0,
	}
}

// NewDefaultFilteringParams creates search pagination parameters according to Splunk Enterprise defaults
func NewDefaultFilteringParams() *FilteringParams {
	return &FilteringParams{
		SortDir:   ASC,
		SortKey:   "name",
		SortMode:  AUTO,
		Summarize: false,
	}
}
