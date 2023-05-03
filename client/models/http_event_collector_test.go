package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func randomString() string {
	n := rand.Intn(100)
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sb := strings.Builder{}
	sb.Grow(n)

	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}

	return sb.String()
}

func TestDecodeHttpEventCollectorObject(t *testing.T) {
	indexes := make([]string, rand.Intn(10))

	for i := range indexes {
		indexes[i] = randomString()
	}

	host := randomString()
	index := randomString()
	source := randomString()
	sourcetype := randomString()
	token := randomString()

	testCase := fmt.Sprintf(
		`{
			"host":"%s",
			"indexes":["%s"],
			"index":"%s",
			"source":"%s",
			"sourcetype":"%s",
			"token":"%s",
			"disabled":false,
			"useACK":"1"
		}`,
		host,
		strings.Join(indexes, `","`),
		index,
		source,
		sourcetype,
		token,
	)

	content := HttpEventCollectorObject{}
	err := json.NewDecoder(strings.NewReader(testCase)).Decode(&content)

	if err != nil {
		t.Errorf("Could not parse HttpEventCollectorObject: %v", err)
	}

	if content.Host != host {
		t.Errorf("Host did not parse %s as %s", content.Host, host)
	}

	for i := range indexes {
		if v := content.Indexes[i]; v != indexes[i] {
			t.Errorf("Indexes[%d] did not parse %v as %s", i, v, indexes[i])
		}
	}

	if content.Index != index {
		t.Errorf("Index did not parse %s as %s", content.Index, index)
	}

	if content.Source != source {
		t.Errorf("Source did not parse %s as %s", content.Source, source)
	}

	if content.SourceType != sourcetype {
		t.Errorf("SourceType did not parse %s as %s", content.SourceType, sourcetype)
	}

	if content.Token != token {
		t.Errorf("Token did not parse %s as %s", content.Token, token)
	}

	if content.Disabled {
		t.Error("Disabled should've been false")
	}

	if content.UseACK != 1 {
		t.Errorf("UseACK did not parse %d as 1", content.UseACK)
	}
}

func TestDecodeUseAckAsBool(t *testing.T) {
	testCases := map[string]int{"false": 0, "true": 1}

	for k, v := range testCases {
		t.Run(fmt.Sprintf("%s should be %d", k, v), func(t *testing.T) {
			body := fmt.Sprintf(`{"useACK":%s}`, k)
			content := HttpEventCollectorObject{}
			err := json.NewDecoder(strings.NewReader(body)).Decode(&content)

			if err != nil {
				t.Errorf("Could not parse useACK from bool: %v", err)
			}

			if int(content.UseACK) != v {
				t.Errorf("UseACK did not parse %s as %d", k, v)
			}
		})
	}
}

func TestDecodeUseAckAsBoolString(t *testing.T) {
	testCases := map[string]int{"false": 0, "true": 1}

	for k, v := range testCases {
		t.Run(fmt.Sprintf("%s should be %d", k, v), func(t *testing.T) {
			body := fmt.Sprintf(`{"useACK":"%s"}`, k)
			content := HttpEventCollectorObject{}
			err := json.NewDecoder(strings.NewReader(body)).Decode(&content)

			if err != nil {
				t.Errorf("Could not parse useACK from bool-string: %v", err)
			}

			if int(content.UseACK) != v {
				t.Errorf("UseACK did not parse %s as %d", k, v)
			}
		})
	}
}

func TestDecodeUseAckAsString(t *testing.T) {
	for v := range rand.Perm(3) {
		t.Run(fmt.Sprintf(`"%d" should be %d`, v, v), func(t *testing.T) {
			body := fmt.Sprintf(`{"useACK":"%d"}`, v)
			content := HttpEventCollectorObject{}
			err := json.NewDecoder(strings.NewReader(body)).Decode(&content)

			if err != nil {
				t.Errorf("Could not parse useACK string: %v", err)
			}

			if int(content.UseACK) != v {
				t.Errorf("UseACK did not parse %d", v)
			}
		})
	}
}

func TestDecodeUseAckAsInt(t *testing.T) {
	for v := range rand.Perm(3) {
		t.Run(fmt.Sprintf(`%d should be %d`, v, v), func(t *testing.T) {
			body := fmt.Sprintf(`{"useACK":%d}`, v)
			content := HttpEventCollectorObject{}
			err := json.NewDecoder(strings.NewReader(body)).Decode(&content)

			if err != nil {
				t.Errorf("Could not parse useACK from int: %v", err)
			}

			if int(content.UseACK) != v {
				t.Errorf("UseACK did not parse %d", v)
			}
		})
	}
}
