package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

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
	UseACK     int           `json:"useACK,omitempty" url:"useACK"`
}

func (a *HttpEventCollectorObject) UnmarshalJSON(data []byte) error {
	var content map[string]interface{}
	err := json.Unmarshal(data, &content)

	if err != nil {
		return err
	}

	if host, ok := content["host"]; ok {
		a.Host = host.(string)
	}

	if indexes, ok := content["indexes"]; ok {
		a.Indexes = indexes.([]interface{})
	}

	if index, ok := content["index"]; ok {
		a.Index = index.(string)
	}

	if source, ok := content["source"]; ok {
		a.Source = source.(string)
	}

	if sourcetype, ok := content["sourcetype"]; ok {
		a.SourceType = sourcetype.(string)
	}

	if token, ok := content["token"]; ok {
		a.Token = token.(string)
	}

	if disabled, ok := content["disabled"]; ok {
		a.Disabled = disabled.(bool)
	}

	if useack, ok := content["useACK"]; ok {
		a.UseACK, err = unmarshalUseAck(useack)
	}

	return err
}

func unmarshalUseAck(data interface{}) (int, error) {
	if i, ok := data.(int); ok {
		return i, nil
	}

	if i, ok := data.(float64); ok {
		return int(i), nil
	}

	if b, ok := data.(bool); ok {
		if b {
			return 1, nil
		}
		return 0, nil
	}

	if s, ok := data.(string); ok {
		if s == "true" {
			return 1, nil
		}

		if s == "false" {
			return 0, nil
		}

		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return int(val), nil
	}

	return 0, fmt.Errorf(`Could not parse "%v" as UseAck`, data)
}
