package model

import (
	"encoding/json"
	"reflect"
	"strings"
)

// SearchJobEntry specifies the fields returned for a /search/jobs entry
type SearchJobEntry struct {
	Name      string                 `json:"name"`
	ID        string                 `json:"id"`
	Updated   string                 `json:"updated"`
	Links     map[string]interface{} `json:"links"`
	Published string                 `json:"published"`
	Author    string                 `json:"author"`
	Content   SearchJobContent       `json:"content"`
	ACL       map[string]interface{} `json:"acl"`
}

// PagingInfo captures fields returned for endpoints supporting paging
type PagingInfo struct {
	Total   float64 `json:"total"`
	PerPage float64 `json:"perPage"`
	Offset  float64 `json:"offset"`
}

// SearchJobsResponse represents a response that can be unmarshalled from /search/jobs
type SearchJobsResponse struct {
	Links     map[string]interface{} `json:"links"`
	Origin    string                 `json:"origin"`
	Updated   string                 `json:"updated"`
	Generator map[string]interface{} `json:"generator"`
	Entry     []SearchJobEntry       `json:"entry"`
	Paging    PagingInfo             `json:"paging"`
}

// SearchJobContent represents the content json object from /search/jobs response
type SearchJobContent struct {
	Sid              string                 `json:"sid"`
	EventCount       int                    `json:"eventCount"`
	DispatchState    string                 `json:"dispatchState"`
	DiskUsage        int64                  `json:"diskUsage"`
	IsFinalized      bool                   `json:"isFinalized"`
	OptimizedSearch  string                 `json:"optimizedSearch"`
	ScanCount        int64                  `json:"scanCount"`
	AdditionalFields map[string]interface{} `json:"-"`
}

type searchJobContent SearchJobContent

// UnmarshalJSON unmarshals named properties as fields and additional properties into a map
func (content *SearchJobContent) UnmarshalJSON(b []byte) error {
	tmp := searchJobContent{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &tmp.AdditionalFields); err != nil {
		return err
	}

	typ := reflect.TypeOf(tmp)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonTag != "" && jsonTag != "-" {
			delete(tmp.AdditionalFields, jsonTag)
		}
	}

	*content = SearchJobContent(tmp)
	return nil
}
