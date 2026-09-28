package splunk

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/splunk/terraform-provider-splunk/client/models"
)

type namespacedRESTImportID struct {
	Owner string
	App   string
	Name  string
}

func parseNamespacedRESTImportID(id string, resourceParts ...string) (*namespacedRESTImportID, bool, error) {
	rawPath, isPathImport := importIDPath(id)
	if !isPathImport {
		return nil, false, nil
	}
	pathParts := strings.Split(strings.Trim(rawPath, "/"), "/")

	servicesNSIndex := -1
	for i, part := range pathParts {
		if part == "servicesNS" {
			servicesNSIndex = i
			break
		}
	}
	if servicesNSIndex == -1 {
		return nil, false, nil
	}

	pathParts = pathParts[servicesNSIndex:]
	minParts := 3 + len(resourceParts) + 1
	if len(pathParts) != minParts {
		return nil, true, fmt.Errorf("import path must be /servicesNS/<owner>/<app>/%s/<name>", strings.Join(resourceParts, "/"))
	}

	owner, err := decodeImportPathSegment(pathParts[1])
	if err != nil {
		return nil, true, fmt.Errorf("unable to decode owner in import path: %w", err)
	}
	app, err := decodeImportPathSegment(pathParts[2])
	if err != nil {
		return nil, true, fmt.Errorf("unable to decode app in import path: %w", err)
	}

	for i, expected := range resourceParts {
		actual, err := decodeImportPathSegment(pathParts[3+i])
		if err != nil {
			return nil, true, fmt.Errorf("unable to decode resource path in import path: %w", err)
		}
		if actual != expected {
			return nil, true, fmt.Errorf("import path must be /servicesNS/<owner>/<app>/%s/<name>", strings.Join(resourceParts, "/"))
		}
	}

	name, err := decodeImportPathSegment(pathParts[3+len(resourceParts)])
	if err != nil {
		return nil, true, fmt.Errorf("unable to decode resource name in import path: %w", err)
	}
	if owner == "" || app == "" || name == "" {
		return nil, true, fmt.Errorf("import path must include non-empty owner, app, and name")
	}

	return &namespacedRESTImportID{Owner: owner, App: app, Name: name}, true, nil
}

func parseNamespacedConfigsConfImportID(id string) (*namespacedRESTImportID, bool, error) {
	rawPath, isPathImport := importIDPath(id)
	if !isPathImport {
		return nil, false, nil
	}
	pathParts := strings.Split(strings.Trim(rawPath, "/"), "/")

	servicesNSIndex := -1
	for i, part := range pathParts {
		if part == "servicesNS" {
			servicesNSIndex = i
			break
		}
	}
	if servicesNSIndex == -1 {
		return nil, false, nil
	}

	pathParts = pathParts[servicesNSIndex:]
	const expectedParts = 6
	const expectedPath = "/servicesNS/<owner>/<app>/configs/conf-<configuration>/<stanza>"
	if len(pathParts) != expectedParts {
		return nil, true, fmt.Errorf("import path must be %s", expectedPath)
	}

	owner, err := decodeImportPathSegment(pathParts[1])
	if err != nil {
		return nil, true, fmt.Errorf("unable to decode owner in import path: %w", err)
	}
	app, err := decodeImportPathSegment(pathParts[2])
	if err != nil {
		return nil, true, fmt.Errorf("unable to decode app in import path: %w", err)
	}
	resource, err := decodeImportPathSegment(pathParts[3])
	if err != nil {
		return nil, true, fmt.Errorf("unable to decode resource path in import path: %w", err)
	}
	confPart, err := decodeImportPathSegment(pathParts[4])
	if err != nil {
		return nil, true, fmt.Errorf("unable to decode configuration in import path: %w", err)
	}
	stanza, err := decodeImportPathSegment(pathParts[5])
	if err != nil {
		return nil, true, fmt.Errorf("unable to decode stanza in import path: %w", err)
	}

	if resource != "configs" || !strings.HasPrefix(confPart, "conf-") {
		return nil, true, fmt.Errorf("import path must be %s", expectedPath)
	}
	conf := strings.TrimPrefix(confPart, "conf-")
	if owner == "" || app == "" || conf == "" || stanza == "" {
		return nil, true, fmt.Errorf("import path must include non-empty owner, app, configuration, and stanza")
	}

	return &namespacedRESTImportID{
		Owner: owner,
		App:   app,
		Name:  conf + "/" + stanza,
	}, true, nil
}

func importIDPath(id string) (string, bool) {
	trimmed := strings.TrimSpace(id)
	parsed, err := url.Parse(trimmed)
	if err == nil && (parsed.Scheme != "" || parsed.Host != "") {
		if parsed.RawPath != "" {
			return parsed.RawPath, true
		}
		if escapedPath := parsed.EscapedPath(); escapedPath != "" {
			return escapedPath, true
		}
		return parsed.Path, true
	}

	if !strings.HasPrefix(trimmed, "/") {
		return "", false
	}

	if i := strings.IndexAny(trimmed, "?#"); i >= 0 {
		return trimmed[:i], true
	}
	return trimmed, true
}

func decodeImportPathSegment(segment string) (string, error) {
	return url.PathUnescape(segment)
}

func importNamespacedResourceState(d *schema.ResourceData, resourceParts ...string) (bool, error) {
	parsed, matched, err := parseNamespacedRESTImportID(d.Id(), resourceParts...)
	if err != nil || !matched {
		return matched, err
	}

	d.SetId(parsed.Name)
	if err := d.Set("name", parsed.Name); err != nil {
		return true, err
	}

	aclObject := &models.ACLObject{
		Owner:   parsed.Owner,
		App:     parsed.App,
		Sharing: inferredNamespacedImportSharing(parsed.Owner),
	}
	if err := d.Set("acl", flattenACL(aclObject)); err != nil {
		return true, err
	}

	return true, nil
}

func inferredNamespacedImportSharing(owner string) string {
	if owner == "nobody" {
		return "app"
	}
	return "user"
}
