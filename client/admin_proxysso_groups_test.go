package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/splunk/terraform-provider-splunk/client/models"
)

func TestAdminProxyssoGroupsRequests(t *testing.T) {
	t.Setenv(envVarHTTPScheme, "http")

	const groupName = "A/B+C D"
	obj := &models.AdminProxyssoGroupsObject{
		Name:  groupName,
		Roles: []string{"power", "admin"},
	}

	tests := []struct {
		name       string
		wantMethod string
		wantPath   string
		wantName   string
		wantRoles  []string
		call       func(*Client) error
	}{
		{
			name:       "create posts the collection",
			wantMethod: http.MethodPost,
			wantPath:   "/services/admin/ProxySSO-groups",
			wantName:   groupName,
			wantRoles:  []string{"power", "admin"},
			call: func(c *Client) error {
				return c.CreateAdminProxyssoGroups(groupName, obj)
			},
		},
		{
			name:       "update posts the escaped entity and omits name",
			wantMethod: http.MethodPost,
			wantPath:   "/services/admin/ProxySSO-groups/A%2FB+C%20D",
			wantRoles:  []string{"power", "admin"},
			call: func(c *Client) error {
				return c.UpdateAdminProxyssoGroups(groupName, obj)
			},
		},
		{
			name:       "delete uses the escaped entity",
			wantMethod: http.MethodDelete,
			wantPath:   "/services/admin/ProxySSO-groups/A%2FB+C%20D",
			call: func(c *Client) error {
				resp, err := c.DeleteAdminProxyssoGroups(groupName)
				if resp != nil {
					resp.Body.Close()
				}
				return err
			},
		},
		{
			name:       "read uses the escaped entity",
			wantMethod: http.MethodGet,
			wantPath:   "/services/admin/ProxySSO-groups/A%2FB+C%20D",
			call: func(c *Client) error {
				resp, err := c.ReadAdminProxyssoGroups(groupName)
				if resp != nil {
					resp.Body.Close()
				}
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotMethod, gotPath, body string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod = r.Method
				gotPath = r.URL.EscapedPath()
				payload, _ := io.ReadAll(r.Body)
				body = string(payload)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, `{"entry":[]}`)
			}))
			defer server.Close()

			c := testProxyssoClient(t, server)
			if err := tt.call(c); err != nil {
				t.Fatal(err)
			}
			if gotMethod != tt.wantMethod {
				t.Errorf("method = %s, want %s", gotMethod, tt.wantMethod)
			}
			if gotPath != tt.wantPath {
				t.Errorf("path = %s, want %s", gotPath, tt.wantPath)
			}
			if tt.wantMethod == http.MethodPost {
				values, err := url.ParseQuery(body)
				if err != nil {
					t.Fatal(err)
				}
				if got := values.Get("name"); got != tt.wantName {
					t.Errorf("name = %q, want %q", got, tt.wantName)
				}
				if strings.Join(values["roles"], ",") != strings.Join(tt.wantRoles, ",") {
					t.Errorf("roles = %#v, want %#v", values["roles"], tt.wantRoles)
				}
			}
		})
	}
}

func testProxyssoClient(t *testing.T, server *httptest.Server) *Client {
	t.Helper()
	backend, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewSplunkdClient("", [2]string{}, backend.Host, "", server.Client())
	if err != nil {
		t.Fatal(err)
	}
	return c
}
