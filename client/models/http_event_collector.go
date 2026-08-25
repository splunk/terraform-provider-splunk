package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// IntBool is a custom type that can unmarshal JSON values that may be
// represented as integers ("0"/"1"), booleans ("true"/"false"), or their
// unquoted equivalents. Splunk Cloud 8.2+ returns "true"/"false" for useACK
// instead of "0"/"1" which breaks the standard json:",string" tag.
type IntBool int

func (ib *IntBool) UnmarshalJSON(data []byte) error {
	// Remove surrounding quotes if present
	s := strings.Trim(string(data), "\"")

	// Try parsing as integer first
	if i, err := strconv.Atoi(s); err == nil {
		*ib = IntBool(i)
		return nil
	}

	// Try parsing as boolean
	if b, err := strconv.ParseBool(s); err == nil {
		if b {
			*ib = 1
		} else {
			*ib = 0
		}
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into IntBool", string(data))
}

func (ib IntBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(ib))
}

// HTTP Input Response Schema
type HECResponse struct {
	Entry    []HECEntry     `json:"entry"`
	Messages []ErrorMessage `json:"messages"`
}

type HECEntry struct {
	Name    string                   `json:"name"`
	ACL     ACLObject                `json:"acl"`
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
	UseACK     IntBool       `json:"useACK" url:"useACK"`
}
