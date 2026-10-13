package splunk

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/zclconf/go-cty/cty"
)

const newHttpEventCollectorInput = `
resource "splunk_global_http_event_collector" "http" {
  disabled     = false
  enable_ssl   = true
}

resource "splunk_inputs_http_event_collector" "new-token" {
    name = "new-token"
    source = "new"
    disabled = false
    use_ack = 0

    acl {
      app = "launcher"
      sharing = "global"
    }

    depends_on = ["splunk_global_http_event_collector.http"]
}
`

const updateHttpEventCollectorInput = `
resource "splunk_global_http_event_collector" "http" {
  disabled     = false
  enable_ssl   = true
  port         = 8088
}

resource "splunk_inputs_http_event_collector" "new-token" {
    name = "new-token"
    index = "main"
    indexes = ["main", "history"]
    source = "new"
    sourcetype = "new"
    disabled = false
    use_ack = 1

    depends_on = ["splunk_global_http_event_collector.http"]
}
`

const newHttpEventCollectorInputWithToken = `
resource "splunk_global_http_event_collector" "http" {
  disabled     = false
  enable_ssl   = true
}

resource "splunk_inputs_http_event_collector" "new-token" {
    name = "new-token"
    token = "12345678-0000-0000-0000-123456780000"
    source = "new"
    disabled = false
    use_ack = 0

    acl {
      app = "launcher"
      sharing = "global"
    }

    depends_on = ["splunk_global_http_event_collector.http"]
}
`

// hecTokenListEndpointConfig creates a HEC token with no acl block (default nobody/splunk_httpinput).
const hecTokenListEndpointConfig = `
resource "splunk_global_http_event_collector" "http" {
  disabled   = false
  enable_ssl = true
}

resource "splunk_inputs_http_event_collector" "list_endpoint_test" {
  name       = "test-hec-list-endpoint"
  index      = "main"
  indexes    = ["main"]
  disabled   = false
  use_ack    = 0
  depends_on = ["splunk_global_http_event_collector.http"]
}
`

func TestAccSplunkHttpEventCollectorInput(t *testing.T) {
	resourceName := "splunk_inputs_http_event_collector.new-token"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkHttpEventCollectorInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newHttpEventCollectorInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source", "new"),
					resource.TestCheckResourceAttr(resourceName, "index", "default"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "use_ack", "0"),
					resource.TestCheckResourceAttr(resourceName, "acl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "launcher"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "nobody"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "global"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.0", "admin"),
				),
			},
			{
				Config: updateHttpEventCollectorInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source", "new"),
					resource.TestCheckResourceAttr(resourceName, "index", "main"),
					resource.TestCheckResourceAttr(resourceName, "sourcetype", "new"),
					resource.TestCheckResourceAttr(resourceName, "indexes.#", "2"),
					testCheckHECIndex(resourceName, "main"),
					testCheckHECIndex(resourceName, "history"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "use_ack", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "launcher"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "nobody"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "global"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.0", "admin"),
				),
			},
			{
				ResourceName:      "splunk_inputs_http_event_collector.new-token",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSplunkHttpEventCollectorInputWithToken(t *testing.T) {
	resourceName := "splunk_inputs_http_event_collector.new-token"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkHttpEventCollectorInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newHttpEventCollectorInputWithToken,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source", "new"),
					resource.TestCheckResourceAttr(resourceName, "index", "default"),
					resource.TestCheckResourceAttr(resourceName, "token", "12345678-0000-0000-0000-123456780000"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "use_ack", "0"),
					resource.TestCheckResourceAttr(resourceName, "acl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "launcher"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "nobody"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "global"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.0", "admin"),
				),
			},
			{
				ResourceName:      "splunk_inputs_http_event_collector.new-token",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccSplunkHttpEventCollectorInputListEndpoint creates a HEC token with default ACL (no acl block)
// and asserts it can be read back. After create, the provider runs Read(), which lists tokens via
// ReadAllHttpEventCollectorObject(), finds the entry by name, then reads the token with that entry’s
// owner and app. The test passes only if the list returns the token and the subsequent read succeeds.
func TestAccSplunkHttpEventCollectorInputListEndpoint(t *testing.T) {
	resourceName := "splunk_inputs_http_event_collector.list_endpoint_test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkHttpEventCollectorInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: hecTokenListEndpointConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-hec-list-endpoint"),
					resource.TestCheckResourceAttr(resourceName, "index", "main"),
					resource.TestCheckResourceAttr(resourceName, "indexes.#", "1"),
					testCheckHECIndex(resourceName, "main"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "use_ack", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "token"),
					resource.TestCheckResourceAttr(resourceName, "acl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "splunk_httpinput"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "nobody"),
				),
			},
		},
	})
}

func testCheckHECIndex(resourceName, index string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(
		resourceName,
		fmt.Sprintf("indexes.%d", schema.HashString(index)),
		index,
	)
}

func TestHttpEventCollectorIndexesSchema(t *testing.T) {
	res := inputsHttpEventCollector()
	indexes := res.Schema["indexes"]
	if indexes.Type != schema.TypeSet || !indexes.Optional || !indexes.Computed {
		t.Fatalf("indexes schema = type %v optional %v computed %v", indexes.Type, indexes.Optional, indexes.Computed)
	}
	if indexes.Set == nil || indexes.Set("main") != schema.HashString("main") {
		t.Fatal("indexes must use schema.HashString so acceptance checks find indexes.<hash>")
	}
	if res.SchemaVersion != 1 || len(res.StateUpgraders) != 1 || res.StateUpgraders[0].Version != 0 {
		t.Fatalf("schema version = %d upgraders = %#v", res.SchemaVersion, res.StateUpgraders)
	}
	if !res.Schema["token"].Sensitive {
		t.Fatal("token must stay Sensitive")
	}

	prior := inputsHttpEventCollectorResourceV0()
	if prior.Schema["indexes"].Type != schema.TypeList || prior.Schema["indexes"].Set != nil {
		t.Fatalf("v0 indexes type = %v set = %v", prior.Schema["indexes"].Type, prior.Schema["indexes"].Set)
	}
	if prior.Schema["indexes"] == res.Schema["indexes"] {
		t.Fatal("v0 indexes schema aliases the current schema")
	}
}

func TestGetHttpEventCollectorConfigIndexesSet(t *testing.T) {
	d := schema.TestResourceDataRaw(t, inputsHttpEventCollector().Schema, map[string]interface{}{
		"name":    "hec-token",
		"index":   "main",
		"indexes": []interface{}{"history", "main", "main"},
	})
	got := getHttpEventCollectorConfig(d)
	if !sameIndexes(got.Indexes, []string{"history", "main"}) {
		t.Fatalf("indexes = %#v", got.Indexes)
	}

	unset := schema.TestResourceDataRaw(t, inputsHttpEventCollector().Schema, map[string]interface{}{
		"name": "hec-token",
	})
	if indexes := getHttpEventCollectorConfig(unset).Indexes; len(indexes) != 0 {
		t.Fatalf("unset indexes = %#v", indexes)
	}
}

func TestHttpEventCollectorStateUpgradeV0(t *testing.T) {
	res := inputsHttpEventCollector()
	upgraded := upgradeHECFlatmapState(t, map[string]string{
		"id":                     "hec-token",
		"name":                   "hec-token",
		"token":                  "test-token",
		"index":                  "main",
		"host":                   "localhost",
		"source":                 "new",
		"sourcetype":             "new",
		"disabled":               "false",
		"use_ack":                "1",
		"indexes.#":              "3",
		"indexes.0":              "history",
		"indexes.1":              "main",
		"indexes.2":              "main",
		"acl.#":                  "1",
		"acl.0.app":              "launcher",
		"acl.0.owner":            "nobody",
		"acl.0.sharing":          "global",
		"acl.0.read.#":           "1",
		"acl.0.read.0":           "admin",
		"acl.0.write.#":          "1",
		"acl.0.write.0":          "admin",
		"acl.0.can_change_perms": "true",
		"acl.0.can_share_app":    "true",
		"acl.0.can_share_global": "true",
		"acl.0.can_share_user":   "true",
		"acl.0.can_write":        "true",
		"acl.0.removable":        "false",
	})

	if upgraded.Attributes["name"] != "hec-token" || upgraded.Attributes["token"] != "test-token" || upgraded.Attributes["index"] != "main" {
		t.Fatalf("upgraded identity attrs = %#v", upgraded.Attributes)
	}
	if upgraded.Attributes["disabled"] != "false" || upgraded.Attributes["use_ack"] != "1" {
		t.Fatalf("upgraded flags = disabled %q use_ack %q", upgraded.Attributes["disabled"], upgraded.Attributes["use_ack"])
	}
	if upgraded.Attributes["acl.0.app"] != "launcher" || upgraded.Attributes["acl.0.owner"] != "nobody" || upgraded.Attributes["acl.0.read.0"] != "admin" {
		t.Fatalf("upgraded acl = %#v", upgraded.Attributes)
	}
	if upgraded.Attributes["indexes.#"] != "2" {
		t.Fatalf("indexes.# = %q, want 2 after dropping the duplicate", upgraded.Attributes["indexes.#"])
	}
	for _, index := range []string{"history", "main"} {
		key := fmt.Sprintf("indexes.%d", schema.HashString(index))
		if upgraded.Attributes[key] != index {
			t.Fatalf("state[%s] = %q", key, upgraded.Attributes[key])
		}
	}
	if _, ok := upgraded.Attributes["indexes.0"]; ok {
		t.Fatal("upgraded state still has list key indexes.0")
	}

	absent := upgradeHECFlatmapState(t, map[string]string{
		"id":    "hec-token",
		"name":  "hec-token",
		"token": "test-token",
		"index": "main",
	})
	if absent.Attributes["name"] != "hec-token" || absent.Attributes["token"] != "test-token" {
		t.Fatalf("state without indexes dropped other attributes: %#v", absent.Attributes)
	}
	if absent.Attributes["indexes.#"] != "" && absent.Attributes["indexes.#"] != "0" {
		t.Fatalf("missing indexes became %#v", absent.Attributes["indexes.#"])
	}

	current := roundTripHECState(t, res, res.CoreConfigSchema().ImpliedType(), map[string]string{
		"id":        "hec-token",
		"name":      "hec-token",
		"indexes.#": "2",
		fmt.Sprintf("indexes.%d", schema.HashString("summary")): "summary",
		fmt.Sprintf("indexes.%d", schema.HashString("main")):    "main",
	})
	if current.Attributes["indexes.#"] != "2" || current.Attributes[fmt.Sprintf("indexes.%d", schema.HashString("summary"))] != "summary" {
		t.Fatalf("current set state = %#v", current.Attributes)
	}
}

func TestHttpEventCollectorCreateReadUpdateIndexes(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	const tokenName = "hec-token"
	const entityPath = "/servicesNS/nobody/splunk_httpinput/data/inputs/http/hec-token"
	const listPath = "/services/data/inputs/http"
	const aclPath = entityPath + "/acl"

	stored := []string{"history", "main"}
	expect := []string{"main", "history"}
	readBody := func() string {
		quoted := make([]string, len(stored))
		for i, index := range stored {
			quoted[i] = `"` + index + `"`
		}
		return `{"entry":[{"name":"http://hec-token","acl":{"owner":"nobody","app":"splunk_httpinput","sharing":"app","perms":{"read":["admin"],"write":["admin"]}},"content":{"index":"main","indexes":[` +
			strings.Join(quoted, ",") +
			`],"source":"new","token":"test-token","disabled":false,"useACK":0}}]}`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == entityPath:
			values, err := url.ParseQuery(string(body))
			if err != nil {
				t.Errorf("parse form: %v", err)
			}
			if !sameStringSet(values["indexes"], expect) {
				t.Errorf("posted indexes = %#v, want %v", values["indexes"], expect)
			}
			stored = append([]string(nil), values["indexes"]...)
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, readBody())
		case r.Method == http.MethodPost && r.URL.Path == aclPath:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{}`)
		case r.Method == http.MethodGet && (r.URL.Path == listPath || r.URL.Path == entityPath):
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, readBody())
		default:
			http.Error(w, "unexpected "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	provider := proxyssoTestProvider(t, server)
	d := schema.TestResourceDataRaw(t, inputsHttpEventCollector().Schema, map[string]interface{}{
		"name":    tokenName,
		"index":   "main",
		"indexes": []interface{}{"main", "history"},
		"source":  "new",
	})
	if err := httpEventCollectorInputCreate(d, provider); err != nil {
		t.Fatalf("create: %v", err)
	}
	assertIndexes(t, d, "history", "main")

	if err := d.Set("indexes", []interface{}{"summary", "main"}); err != nil {
		t.Fatal(err)
	}
	expect = []string{"summary", "main"}
	stored = []string{"main", "summary"}
	if err := httpEventCollectorInputUpdate(d, provider); err != nil {
		t.Fatalf("update: %v", err)
	}
	assertIndexes(t, d, "summary", "main")
}

func upgradeHECFlatmapState(t *testing.T, attrs map[string]string) *terraform.InstanceState {
	t.Helper()
	res := inputsHttpEventCollector()
	prior := inputsHttpEventCollectorResourceV0()
	priorType := prior.CoreConfigSchema().ImpliedType()
	val, err := schema.StateValueFromInstanceState(&terraform.InstanceState{
		ID:         attrs["id"],
		Attributes: attrs,
	}, priorType)
	if err != nil {
		t.Fatalf("decode v0 state: %v", err)
	}
	jsonMap, err := schema.StateValueToJSONMap(val, priorType)
	if err != nil {
		t.Fatalf("v0 json: %v", err)
	}
	upgraded, err := httpEventCollectorStateUpgradeV0(jsonMap, nil)
	if err != nil {
		t.Fatalf("upgrade: %v", err)
	}
	newVal, err := schema.JSONMapToStateValue(upgraded, res.CoreConfigSchema())
	if err != nil {
		t.Fatalf("encode v1 state: %v", err)
	}
	state, err := res.ShimInstanceStateFromValue(newVal)
	if err != nil {
		t.Fatalf("v1 flatmap: %v", err)
	}
	return state
}

func roundTripHECState(t *testing.T, res *schema.Resource, ty cty.Type, attrs map[string]string) *terraform.InstanceState {
	t.Helper()
	val, err := schema.StateValueFromInstanceState(&terraform.InstanceState{
		ID:         attrs["id"],
		Attributes: attrs,
	}, ty)
	if err != nil {
		t.Fatalf("decode state: %v", err)
	}
	state, err := res.ShimInstanceStateFromValue(val)
	if err != nil {
		t.Fatalf("encode state: %v", err)
	}
	return state
}

func assertIndexes(t *testing.T, d *schema.ResourceData, want ...string) {
	t.Helper()
	got := d.Get("indexes").(*schema.Set)
	if got.Len() != len(want) {
		t.Fatalf("indexes = %#v, want %v", got.List(), want)
	}
	for _, index := range want {
		if !got.Contains(index) {
			t.Fatalf("indexes = %#v, missing %s", got.List(), index)
		}
	}
}

func sameIndexes(got []interface{}, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	seen := map[string]int{}
	for _, item := range got {
		name, ok := item.(string)
		if !ok {
			return false
		}
		seen[name]++
	}
	for _, name := range want {
		seen[name]--
	}
	for _, n := range seen {
		if n != 0 {
			return false
		}
	}
	return true
}

func sameStringSet(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	seen := map[string]int{}
	for _, name := range got {
		seen[name]++
	}
	for _, name := range want {
		seen[name]--
	}
	for _, n := range seen {
		if n != 0 {
			return false
		}
	}
	return true
}

func testAccSplunkHttpEventCollectorInputDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_inputs_http_event_collector":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "splunk_httpinput", "data", "inputs", "http", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
