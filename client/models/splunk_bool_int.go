package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// SplunkBoolInt handles Splunk Cloud API responses that return boolean strings
// ("true"/"false") for numeric fields instead of integer strings ("0"/"1").
// This inconsistency exists between Splunk Enterprise and Splunk Cloud APIs.
type SplunkBoolInt int

func (b *SplunkBoolInt) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		switch strings.ToLower(strings.TrimSpace(s)) {
		case "true", "1":
			*b = 1
		case "false", "0":
			*b = 0
		default:
			return fmt.Errorf("cannot parse %q as SplunkBoolInt", s)
		}
		return nil
	}
	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return fmt.Errorf("cannot parse %s as SplunkBoolInt: %w", data, err)
	}
	*b = SplunkBoolInt(i)
	return nil
}

func (b SplunkBoolInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.Itoa(int(b)))
}
