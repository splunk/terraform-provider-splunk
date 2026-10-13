package models

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// FlexString unmarshals a Splunk JSON value that is either a string or a
// number. Index fields such as bucketRebuildMemoryHint ("auto", "16MB") and
// repFactor ("auto", "0") are strings in indexes.conf, and Splunk sometimes
// returns those same fields as JSON numbers. null leaves the value unchanged.
type FlexString string

func (f *FlexString) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	switch v := raw.(type) {
	case nil:
	case string:
		*f = FlexString(v)
	case float64:
		s, err := wholeNumberString(v)
		if err != nil {
			return err
		}
		*f = FlexString(s)
	default:
		return fmt.Errorf("cannot unmarshal %T into string", raw)
	}
	return nil
}

// FlexHotBuckets unmarshals maxHotBuckets.
// Splunk 8.1 and earlier return a JSON number. Splunk 8.2 and later return a
// string such as "3" or "auto". indexes.conf defines "auto" as 3, so this
// stays an integer and existing Terraform config such as max_hot_buckets = 6
// keeps working.
// See https://github.com/splunk/terraform-provider-splunk/issues/124
type FlexHotBuckets int

func (f *FlexHotBuckets) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	switch v := raw.(type) {
	case nil:
	case float64:
		n, err := wholeInt(v)
		if err != nil {
			return err
		}
		*f = FlexHotBuckets(n)
	case string:
		s := strings.TrimSpace(v)
		switch {
		case strings.EqualFold(s, "auto"):
			*f = 3
		case s == "":
			*f = 0
		default:
			n, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("cannot parse %q as an integer", v)
			}
			*f = FlexHotBuckets(n)
		}
	default:
		return fmt.Errorf("cannot unmarshal %T into an integer", raw)
	}
	return nil
}

// UnmarshalJSON decodes an index and, when a flexible field rejects the JSON,
// names that REST field. encoding/json leaves custom UnmarshalJSON errors
// unlabeled.
func (o *IndexObject) UnmarshalJSON(data []byte) error {
	type plain IndexObject
	var decoded plain
	if err := json.Unmarshal(data, &decoded); err != nil {
		return labelIndexFieldError(data, err)
	}
	*o = IndexObject(decoded)
	return nil
}

func labelIndexFieldError(data []byte, err error) error {
	var raw map[string]json.RawMessage
	if json.Unmarshal(data, &raw) != nil {
		return err
	}
	checks := []struct {
		field string
		dest  json.Unmarshaler
	}{
		{field: "bucketRebuildMemoryHint", dest: new(FlexString)},
		{field: "maxHotBuckets", dest: new(FlexHotBuckets)},
		{field: "repFactor", dest: new(FlexString)},
	}
	for _, check := range checks {
		msg, ok := raw[check.field]
		if !ok || len(msg) == 0 || string(msg) == "null" {
			continue
		}
		if fieldErr := check.dest.UnmarshalJSON(msg); fieldErr != nil {
			return fmt.Errorf("%s: %w", check.field, fieldErr)
		}
	}
	return err
}

func wholeNumberString(v float64) (string, error) {
	n, err := wholeInt64(v)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(n, 10), nil
}

func wholeInt(v float64) (int, error) {
	n, err := wholeInt64(v)
	if err != nil {
		return 0, err
	}
	if n < math.MinInt || n > math.MaxInt {
		return 0, fmt.Errorf("cannot parse %v as an integer", v)
	}
	return int(n), nil
}

func wholeInt64(v float64) (int64, error) {
	// 1<<63 is the first magnitude float64 cannot store as an int64.
	// math.MaxInt64 itself rounds to that value, so compare against it directly.
	const maxInt64AsFloat = 9223372036854775808.0
	if math.IsNaN(v) || math.IsInf(v, 0) || math.Trunc(v) != v || v < -maxInt64AsFloat || v >= maxInt64AsFloat {
		return 0, fmt.Errorf("cannot parse %v as an integer", v)
	}
	return int64(v), nil
}
