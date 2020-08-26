package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"terraform-provider-splunk/client/models"
)

func confStanza() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"variables": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: `A map of key value pairs for a stanza.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `A '/' separated string consisting of {conf_file_name}/{stanza_name} ex. props/custom_stanza`,
			},
			"acl": aclSchema(),
		},
		Read:   confStanzaRead,
		Create: confStanzaCreate,
		Delete: confStanzaDelete,
		Update: confStanzaUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func confStanzaCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name:= d.Get("name").(string)
	confStanzaConfigObj := getConfStanzaConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.Owner = "nobody"
		aclObject.App = "search"
	}
	err := (*provider.Client).CreateConfStanzaObject(name, aclObject.Owner, aclObject.App, confStanzaConfigObj)
	if err != nil {
		return err
	}
	if _, ok := d.GetOk("acl"); ok {
		split_name := strings.Split(name, "/")
		conf_name := split_name[0]
		stanza_name := split_name[1]
		err := (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, stanza_name, aclObject, "configs", "conf-" + conf_name)
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return confStanzaRead(d, meta)
}

func confStanzaRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	split_name := strings.Split(name, "/")
	stanza_name := split_name[1]

	// We first get list of stanzas in a conf file to get owner and app name for the specific stanza
	resp, err := (*provider.Client).ReadAllConfStanzaObject(name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getConfStanzaConfigByName(stanza_name, resp)
	if err != nil {
		return err
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadConfStanzaObject(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contentResp, err := (*provider.Client).ReadConfStanzaObject(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	b, _ := ioutil.ReadAll(contentResp.Body)

	json.Unmarshal(b, &result)
	content := result["entry"].([]interface{})[0].(map[string]interface{})["content"].(map[string]interface{})

	for key, _ := range content {
		result, _ := regexp.MatchString(`eai:.*`, key)

		if result {
			delete(content, key)
		}
	}

	delete(content, "disabled")

	entry, err = getConfStanzaConfigByName(stanza_name, resp)
	if err != nil {
		return err
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}

	if err = d.Set("variables", content); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func confStanzaUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	confStanzaConfigObj := getConfStanzaConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	// Update will create a new resource with private `user` permissions if resource had shared permissions set
	var owner string
	if aclObject.Sharing != "user" {
		owner = "nobody"
	} else {
		owner = aclObject.Owner
	}
	name := d.Id()
	split_name := strings.Split(name, "/")
	conf_name := split_name[0]
	stanza_name := split_name[1]
	err := (*provider.Client).UpdateConfStanzaObject(d.Id(), owner, aclObject.App, confStanzaConfigObj)
	if err != nil {
		return err
	}

	//ACL update
	err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, stanza_name, aclObject, "configs", "conf-" + conf_name)
	if err != nil {
		return err
	}

	return confStanzaRead(d, meta)
}

func confStanzaDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))

	resp, err := (*provider.Client).DeleteConfStanzaObject(d.Id(), aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.ConfStanzaResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getConfStanzaConfig(d *schema.ResourceData) (confStanzaConfigObject *models.ConfStanzaObject) {
	confStanzaConfigObject = &models.ConfStanzaObject{}
	mapInterface := d.Get("variables").(map[string]interface {})
	mapString := make(map[string]string)
	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		mapString[strKey] = strValue
	}
	confStanzaConfigObject.Variables = mapString

	return confStanzaConfigObject
}

func getConfStanzaConfigByName(name string, httpResponse *http.Response) (confStanzaEntry *models.ConfStanzaEntry, err error) {
	response := &models.ConfStanzaResponse{}

	switch httpResponse.StatusCode {
	case 200, 201:
		_ =  json.NewDecoder(httpResponse.Body).Decode(&response)
		re := regexp.MustCompile(`(.*)`)

		for _, entry := range response.Entry {
			if name == re.FindStringSubmatch(entry.Name)[1] {
				return &entry, nil
			}
		}

	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return confStanzaEntry, err
	}

	return confStanzaEntry, nil
}