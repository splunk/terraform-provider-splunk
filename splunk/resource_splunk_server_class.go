package splunk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/splunk/terraform-provider-splunk/client/models"
)

func serverClass() *schema.Resource {
	return &schema.Resource{
		Create: resourceSplunkServerClassCreate,
		Read:   resourceSplunkServerClassRead,
		Update: resourceSplunkServerClassUpdate,
		Delete: resourceSplunkServerClassDelete,

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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceSplunkServerClassCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)

	serverClassConfigObj := getServerClassConfig(d)
	err := provider.Client.CreateServerClassObject(name, serverClassConfigObj)
	if err != nil {
		return fmt.Errorf("error creating Splunk server class: %s", err)
	}

	d.SetId(name)
	return resourceSplunkServerClassRead(d, meta)
}

func resourceSplunkServerClassRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()

	resp, err := provider.Client.ReadServerClassObject(name)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading Splunk server class: %s", err)
	}
	defer resp.Body.Close()

	var serverClassResponse models.ServerClassResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverClassResponse); err != nil {
		return fmt.Errorf("error decoding response: %s", err)
	}

	if len(serverClassResponse.Entry) == 0 {
		d.SetId("")
		return nil
	}

	serverClass := serverClassResponse.Entry[0].Content

	d.Set("name", serverClass.Name)
	d.Set("description", serverClass.Description)
	d.Set("whitelist", serverClass.Whitelist)
	d.Set("blacklist", serverClass.Blacklist)
	d.Set("apps", serverClass.Apps)
	d.Set("restart_splunk_web", serverClass.RestartSplunkWeb)
	d.Set("restart_splunkd", serverClass.RestartSplunkd)

	return nil
}

func resourceSplunkServerClassUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()

	serverClassConfigObj := getServerClassConfig(d)
	err := provider.Client.UpdateServerClassObject(name, serverClassConfigObj)
	if err != nil {
		return fmt.Errorf("error updating Splunk server class: %s", err)
	}

	return resourceSplunkServerClassRead(d, meta)
}

func resourceSplunkServerClassDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()

	err := provider.Client.DeleteServerClassObject(name)
	if err != nil {
		return fmt.Errorf("error deleting Splunk server class: %s", err)
	}

	d.SetId("")
	return nil
}

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

func convertInterfaceToStringList(input []interface{}) []string {
	output := make([]string, len(input))
	for i, v := range input {
		output[i] = v.(string)
	}
	return output
}
