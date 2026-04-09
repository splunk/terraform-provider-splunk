package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/splunk/terraform-provider-splunk/client/models"
)

const ACLGetModeCloud = "cloud"

func (client *Client) getAclHTTP(endpoint url.URL) (*http.Response, error) {
	req, err := client.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

// https://docs.splunk.com/Documentation/Splunk/8.0.4/RESTUM/RESTusing#Access_Control_List
func (client *Client) GetAcl(owner, app, name, sharing string, resources ...string) (*http.Response, error) {
	resourcePath := []string{"servicesNS", owner, app}
	resourcePath = append(resourcePath, resources...)
	resourcePath = append(resourcePath, name, "acl")

	var q url.Values
	if strings.EqualFold(strings.TrimSpace(client.ACLGetMode), ACLGetModeCloud) {
		q = url.Values{}
		if owner != "" {
			q.Set("owner", owner)
		}
		if sharing != "" {
			q.Set("sharing", sharing)
		}
	}
	endpoint := client.BuildSplunkURL(q, resourcePath...)
	resp, err := client.getAclHTTP(endpoint)
	if err != nil {
		return nil, fmt.Errorf("GET failed for endpoint %s: %s", endpoint.Path, err)
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return resp, nil
	}
	body, readErr := io.ReadAll(resp.Body)
	resp.Body.Close()
	if readErr != nil {
		return nil, readErr
	}
	return nil, fmt.Errorf("%s: %s", resp.Status, string(body))
}

func (client *Client) ResourcesAndNameForPath(path string) (resources []string, name string, ok bool) {
	parts := strings.Split(path, "/")

	// an example path of apps/local/myapp would have parts:
	// * [apps, local] - parts[0:1]
	// * myapp         - parts[2]

	// 2 is the absolute minimum number of path parts that are valid to be parsed into resources and name
	if len(parts) < 2 {
		ok = false
		return
	}

	resources = parts[0 : len(parts)-1]
	name = parts[len(parts)-1]
	ok = true

	return
}

func (client *Client) UpdateAcl(owner, app, name string, acl *models.ACLObject, resources ...string) error {
	resourcePath := []string{"servicesNS", owner, app}
	resourcePath = append(resourcePath, resources...)
	resourcePath = append(resourcePath, name, "acl")
	isCloud := strings.EqualFold(strings.TrimSpace(client.ACLGetMode), ACLGetModeCloud)
	readPerms := strings.Join(acl.Perms.Read, ",")
	writePerms := strings.Join(acl.Perms.Write, ",")

	var (
		resp     *http.Response
		err      error
		endpoint url.URL
	)

	if isCloud {
		values := url.Values{}
		if owner != "" {
			values.Set("owner", owner)
		}
		if acl.Sharing != "" {
			values.Set("sharing", acl.Sharing)
		}
		values.Set("perms.read", readPerms)
		values.Set("perms.write", writePerms)

		buildPath := client.path
		for _, pathPart := range resourcePath {
			pathPart = strings.ReplaceAll(pathPart, " ", "+")
			buildPath = path.Join(buildPath, pathPart)
		}
		endpoint = url.URL{
			Scheme: getEnv(envVarHTTPScheme, defaultScheme),
			Host:   client.host,
			Path:   buildPath,
		}

		req, reqErr := client.NewRequest(http.MethodPost, endpoint.String(), strings.NewReader(values.Encode()))
		if reqErr != nil {
			return fmt.Errorf("POST failed for endpoint %s: %s", endpoint.Path, reqErr)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err = client.Do(req)
	} else {
		values, valuesErr := query.Values(&acl)
		if valuesErr != nil {
			return valuesErr
		}
		values.Del("app")
		values.Del("perms[read]")
		values.Del("perms[write]")
		values.Set("perms.read", readPerms)
		values.Set("perms.write", writePerms)

		endpoint = client.BuildSplunkURL(nil, resourcePath...)
		resp, err = client.Post(endpoint, values)
	}
	if err != nil {
		return fmt.Errorf("POST failed for endpoint %s: %s", endpoint.Path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("POST failed for endpoint %s: %s", endpoint.Path, resp.Status)
	}

	return nil
}

func (client *Client) Move(owner, app, name string, acl *models.ACLObject, resources ...string) error {
	values, err := query.Values(&acl)
	if err != nil {
		return err
	}
	// Adding resources
	resourcePath := []string{"servicesNS", owner, app}
	resourcePath = append(resourcePath, resources...)
	resourcePath = append(resourcePath, name, "move")
	endpoint := client.BuildSplunkURL(nil, resourcePath...)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
