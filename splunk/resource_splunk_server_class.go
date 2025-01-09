package splunk

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/splunk/terraform-provider-splunk/client/models"
)

func splunkServerClass() *schema.Resource {
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
		Read:   serverClassRead,
		Create: serverClassCreate,
		Delete: serverClassDelete,
		Update: serverClassUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func serverClassCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	serverClassConfigObj := getServerClassConfig(d)
	err := (*provider.Client).CreateServerClassObject(name, serverClassConfigObj)
	if err != nil {
		return err
	}

	d.SetId(name)
	return serverClassRead(d, meta)
}

func serverClassRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	resp, err := (*provider.Client).ReadServerClassObject(name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse the response and set the resource data fields
	// ...

	return nil
}

func serverClassUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	serverClassConfigObj := getServerClassConfig(d)
	err := (*provider.Client).UpdateServerClassObject(name, serverClassConfigObj)
	if err != nil {
		return err
	}

	return serverClassRead(d, meta)
}

func serverClassDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	err := (*provider.Client).DeleteServerClassObject(name)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func getServerClassConfig(d *schema.ResourceData) *models.ServerClassObject {
	return &models.ServerClassObject{
		RestartSplunkWeb: d.Get("restart_splunk_web").(bool),
		RestartSplunkd:   d.Get("restart_splunkd").(bool),
		Whitelist:        convertInterfaceToStringList(d.Get("whitelist").([]interface{})),
		Blacklist:        convertInterfaceToStringList(d.Get("blacklist").([]interface{})),
		Apps:             convertInterfaceToStringList(d.Get("apps").([]interface{})),
		Description:      d.Get("description").(string),
	}
}

func convertInterfaceToStringList(input []interface{}) []string {
	output := make([]string, len(input))
	for i, v := range input {
		output[i] = v.(string)
	}
	return output
}
