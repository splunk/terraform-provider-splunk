package splunk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/splunk/terraform-provider-splunk/client/models"
)

// resourceSplunkServerClass defines the schema and CRUD operations for the Splunk server class resource
func serverClass() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the server class.",
			},
			"restart_splunk_web": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether to restart Splunk Web.",
			},
			"restart_splunkd": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether to restart Splunkd.",
			},
			"whitelist": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of whitelist entries.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"blacklist": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of blacklist entries.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"apps": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of apps associated with the server class.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the server class.",
			},
		},
		Create: serverClassCreate,
		Read:   serverClassRead,
		Update: serverClassUpdate,
		Delete: serverClassDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// resourceSplunkServerClassCreate handles the creation of a Splunk server class
func serverClassCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*SplunkProvider)
	name := d.Get("name").(string)

	serverClass := getServerClassConfig(d)

	body, err := json.Marshal(serverClass)
	if err != nil {
		return fmt.Errorf("error marshaling server class: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/services/deployment/server/serverclasses/%s", provider.Client.BaseURL, name), strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := provider.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating server class: %s", resp.Status)
	}

	d.SetId(name)
	return serverClassRead(d, m)
}

// resourceSplunkServerClassRead handles reading a Splunk server class
func serverClassRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(*SplunkProvider)
	name := d.Id()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/services/deployment/server/serverclasses/%s", provider.Client.BaseURL, name), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}

	resp, err := provider.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error reading server class: %s", resp.Status)
	}

	var serverClass models.ServerClassObject
	if err := json.NewDecoder(resp.Body).Decode(&serverClass); err != nil {
		return fmt.Errorf("error decoding response: %s", err)
	}

	d.Set("name", serverClass.Name)
	d.Set("description", serverClass.Description)
	d.Set("whitelist", serverClass.Whitelist)
	d.Set("blacklist", serverClass.Blacklist)
	d.Set("apps", serverClass.Apps)
	d.Set("restart_splunk_web", serverClass.RestartSplunkWeb)
	d.Set("restart_splunkd", serverClass.RestartSplunkd)

	return nil
}

// resourceSplunkServerClassUpdate handles updating a Splunk server class
func serverClassUpdate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*SplunkProvider)
	name := d.Id()

	serverClass := getServerClassConfig(d)

	body, err := json.Marshal(serverClass)
	if err != nil {
		return fmt.Errorf("error marshaling server class: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/services/deployment/server/serverclasses/%s", provider.Client.BaseURL, name), strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := provider.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error updating server class: %s", resp.Status)
	}

	return serverClassRead(d, m)
}

// resourceSplunkServerClassDelete handles deleting a Splunk server class
func serverClassDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(*SplunkProvider)
	name := d.Id()

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/services/deployment/server/serverclasses/%s", provider.Client.BaseURL, name), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}

	resp, err := provider.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting server class: %s", resp.Status)
	}

	d.SetId("")
	return nil
}

// getServerClassConfig extracts the server class configuration from the resource data
func getServerClassConfig(d *schema.ResourceData) *models.ServerClassObject {
	return &models.ServerClassObject{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Whitelist:        convertInterfaceToStringList(d.Get("whitelist").([]interface{})),
		Blacklist:        convertInterfaceToStringList(d.Get("blacklist").([]interface{})),
		Apps:             convertInterfaceToStringList(d.Get("apps").([]interface{})),
		RestartSplunkWeb: d.Get("restart_splunk_web").(bool),
		RestartSplunkd:   d.Get("restart_splunkd").(bool),
	}
}

// convertInterfaceToStringList converts a list of interfaces to a list of strings
func convertInterfaceToStringList(input []interface{}) []string {
	output := make([]string, len(input))
	for i, v := range input {
		output[i] = v.(string)
	}
	return output
}
