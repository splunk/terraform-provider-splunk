package splunk

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
	"regexp"
	"terraform-provider-splunk/client/models"
)

func inputHttpEventCollector() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "HTTP Event Collector Token name",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HTTP Event Collector Token name",
			},
			"index": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Index to store generated events",
			},
			"indexes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Default host value for events with this token",
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Default source for events with this token",
			},
			"sourcetype": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Default sourcetype for events with this token",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Input disabled indicator: 0 = Input Not disabled, 1 = Input disabled",
			},
			"use_ack": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0",
				Description: "Indexer acknowledgement for this token",
			},
			"acl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sharing": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"read": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"write": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
		Read:   httpEventCollectorInputRead,
		Create: httpEventCollectorInputCreate,
		Delete: httpEventCollectorInputDelete,
		Update: httpEventCollectorInputUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func httpEventCollectorInputCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	httpInputConfigObj := getHttpEventCollectorConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).CreateHttpEventCollectorObject(name, aclObject.Owner, aclObject.App, httpInputConfigObj)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, ok := d.GetOk("acl"); ok {
		resp, err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, "http", name, aclObject)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}

	d.SetId(name)
	return httpEventCollectorInputRead(d, meta)
}

func httpEventCollectorInputRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	// We first get list of inputs to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadAllHttpEventCollectorObject()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getConfigByName(name, resp)
	if err != nil {
		return err
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadHttpEventCollectorObject(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getConfigByName(name, resp)
	if err != nil {
		return err
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}

	if err = d.Set("host", entry.Content.Host); err != nil {
		return err
	}

	if err = d.Set("token", entry.Content.Token); err != nil {
		return err
	}

	if err = d.Set("index", entry.Content.Index); err != nil {
		return err
	}

	if err = d.Set("indexes", entry.Content.Indexes); err != nil {
		return err
	}

	if err = d.Set("source", entry.Content.Source); err != nil {
		return err
	}

	if err = d.Set("sourcetype", entry.Content.SourceType); err != nil {
		return err
	}

	if err = d.Set("disabled", entry.Content.Disabled); err != nil {
		return err
	}

	if err = d.Set("use_ack", entry.Content.UseACK); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func httpEventCollectorInputUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	httpInputConfigObj := getHttpEventCollectorConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).UpdateHttpEventCollectorObject(d.Id(), aclObject.Owner, aclObject.App, httpInputConfigObj)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//ACL update
	resp, err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, "http", d.Id(), aclObject)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return httpEventCollectorInputRead(d, meta)
}

func httpEventCollectorInputDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).DeleteHttpEventCollectorObject(d.Id(), aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.HECResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getHttpEventCollectorConfig(d *schema.ResourceData) (httpInputConfigObject *models.HttpEventCollectorObject) {
	httpInputConfigObject = &models.HttpEventCollectorObject{}
	httpInputConfigObject.Index = d.Get("index").(string)
	httpInputConfigObject.Indexes = d.Get("indexes").([]interface{})
	httpInputConfigObject.Source = d.Get("source").(string)
	httpInputConfigObject.SourceType = d.Get("sourcetype").(string)
	httpInputConfigObject.UseACK = d.Get("use_ack").(string)
	httpInputConfigObject.Disabled = d.Get("disabled").(bool)
	return httpInputConfigObject
}

func getACLConfig(r []interface{}) (acl *models.ACLObject) {
	acl = &models.ACLObject{}
	for _, v := range r {
		a := v.(map[string]interface{})

		if a["app"] != "" {
			acl.App = a["app"].(string)
		} else {
			acl.App = "splunk_httpinput"
		}

		if a["owner"] != "" {
			acl.Owner = a["owner"].(string)
		} else {
			acl.Owner = "nobody"
		}

		if a["sharing"] != "" {
			acl.Sharing = a["sharing"].(string)
		} else {
			acl.Sharing = "global"
		}

		for _, v := range a["read"].([]interface{}) {
			acl.Perms.Read = append(acl.Perms.Read, v.(string))
		}

		for _, w := range a["write"].([]interface{}) {
			acl.Perms.Write = append(acl.Perms.Write, w.(string))
		}
	}

	return acl
}

func getConfigByName(name string, httpResponse *http.Response) (hecEntry *models.HECEntry, err error) {
	response := &models.HECResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:
		_ = json.NewDecoder(httpResponse.Body).Decode(&response)
		re := regexp.MustCompile(`http://(.*)`)
		for _, entry := range response.Entry {
			if name == re.FindStringSubmatch(entry.Name)[1] {
				return &entry, nil
			}
		}

	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return hecEntry, err
	}

	return hecEntry, nil
}

func flattenACL(acl *models.ACLObject) []interface{} {
	if acl == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{}
	m["app"] = acl.App
	m["owner"] = acl.Owner
	m["sharing"] = acl.Sharing
	m["read"] = acl.Perms.Read
	m["write"] = acl.Perms.Write
	return []interface{}{m}
}
