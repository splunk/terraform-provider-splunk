package splunk

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
	"terraform-provider-splunk/client/models"
)

func inputHttpEventCollector() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
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
				Default:     "main",
				Description: "Index to store generated events",
			},
			"indexes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default source for events with this token",
			},
			"sourcetype": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default sourcetype for events with this token",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Input disabled indicator: 0 = Input Not disabled, 1 = Input disabled",
			},
			"use_ack": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indexer acknowledgement for this token",
			},
			"acl": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "splunk_httpinput",
							ForceNew: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "nobody",
						},
						"sharing": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "app",
						},
						"read": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:    schema.TypeString,
								Default: "*",
							},
							Optional: true,
						},
						"write": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:    schema.TypeString,
								Default: "*",
							},
							Optional: true,
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
	httpInputConfigObj := createHttpInputConfigObject(d)
	aclObject := createACLObject(d)
	resp, err := (*provider.Client).CreateHttpEventCollectorObject(name, aclObject.Owner, aclObject.App, httpInputConfigObj)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := unmarshalHttpInputResponse(resp)
	d.SetId(name)
	err = d.Set("token", content.Token)
	if err != nil {
		return err
	}

	//ACL update
	resp, err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, "http", name, aclObject)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func httpEventCollectorInputRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	aclObject := createACLObject(d)
	resp, err := (*provider.Client).ReadHttpEventCollectorObject(name, aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := unmarshalHttpInputResponse(resp)
	if err != nil {
		return err
	}

	d.SetId(name)
	err = d.Set("token", content.Token)
	if err != nil {
		return err
	}

	// ACL Read and Set values
	resp, err = (*provider.Client).GetAcl(aclObject.Owner, aclObject.App, "http", name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	aclContent, err := unmarshalAclResponse(resp)
	err = d.Set("acl", flattenACL(aclContent))
	if err != nil {
		return err
	}

	return nil
}

func httpEventCollectorInputUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	httpInputConfigObj := createHttpInputConfigObject(d)
	aclObject := createACLObject(d)
	resp, err := (*provider.Client).UpdateHttpEventCollectorObject(name, aclObject.Owner, aclObject.App, httpInputConfigObj)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := unmarshalHttpInputResponse(resp)
	if err != nil {
		return err
	}

	d.SetId(name)
	err = d.Set("token", content.Token)
	if err != nil {
		return err
	}

	//ACL update
	resp, err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, "http", name, aclObject)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func httpEventCollectorInputDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	aclObject := createACLObject(d)
	resp, err := (*provider.Client).DeleteHttpEventCollectorObject(name, aclObject.Owner, aclObject.App)
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
func createHttpInputConfigObject(d *schema.ResourceData) (httpInputConfigObject *models.HttpEventCollectorObject) {
	httpInputConfigObject = &models.HttpEventCollectorObject{}
	httpInputConfigObject.Index = d.Get("index").(string)
	httpInputConfigObject.Indexes = d.Get("indexes").([]interface{})
	httpInputConfigObject.Source = d.Get("source").(string)
	httpInputConfigObject.SourceType = d.Get("sourcetype").(string)
	httpInputConfigObject.UseACK = d.Get("use_ack").(bool)
	httpInputConfigObject.Disabled = d.Get("disabled").(bool)
	return httpInputConfigObject
}

func createACLObject(d *schema.ResourceData) (aclObject *models.ACLObject) {
	aclObject = &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.Owner = "nobody"
		aclObject.App = "splunk_httpinput"
		aclObject.Sharing = "app"
	}
	return aclObject
}

func getACLConfig(r []interface{}) (acl *models.ACLObject) {
	for _, v := range r {
		a := v.(map[string]interface{})
		acl = &models.ACLObject{
			App:     a["app"].(string),
			Owner:   a["owner"].(string),
			Sharing: a["sharing"].(string),
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

func unmarshalHttpInputResponse(httpResponse *http.Response) (httpEventCollectorObj *models.HttpEventCollectorObject, err error) {
	response := &models.HECResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:
		_ = json.NewDecoder(httpResponse.Body).Decode(&response)
		return &response.Entry[0].Content, nil

	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return httpEventCollectorObj, err
	}
}

func unmarshalAclResponse(httpResponse *http.Response) (aclObj *models.ACLObject, err error) {
	response := &models.ACLResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:
		_ = json.NewDecoder(httpResponse.Body).Decode(&response)
		return &response.Entry[0].Content, nil

	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return aclObj, err
	}
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
