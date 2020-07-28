package models

// HTTP Input Response Schema
type HECResponse struct {
	Entry    []HECEntry     `json:"entry"`
	Messages []ErrorMessage `json:"messages"`
}

type HECEntry struct {
	Name    string                   `json:"name"`
	Content HttpEventCollectorObject `json:"content"`
}

type HttpEventCollectorObject struct {
	Host       string        `json:"host,omitempty" url:"host,omitempty"`
	Indexes    []interface{} `json:"indexes,omitempty" url:"indexes,omitempty"`
	Index      string        `json:"index,omitempty" url:"index,omitempty"`
	Source     string        `json:"source,omitempty" url:"source,omitempty"`
	SourceType string        `json:"sourcetype,omitempty" url:"sourcetype,omitempty"`
	Token      string        `json:"token,omitempty" url:"token,omitempty"`
	Disabled   bool          `json:"disabled,omitempty" url:"disabled"`
	UseACK     bool          `json:"useACK,omitempty" url:"useACK"`
}
