package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/splunk/terraform-provider-splunk/client/models"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func appsLocal() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Required. Literal app name or path for the file to install, depending on the value of filename." +
					"filename = false indicates that name is the literal app name and that the app is created from Splunkbase using a template." +
					"filename = true indicates that name is the URL or path to the local .tar, .tgz or .spl file. " +
					"If name is the Splunkbase URL, set auth or session to authenticate the request. " +
					"The app folder name cannot include spaces or special characters.",
			},
			"auth": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Splunkbase session token for operations like install and update that require login. " +
					"Use auth or session when installing or updating an app through Splunkbase.",
			},
			"author": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "For apps posted to Splunkbase, use your Splunk account username. For internal apps, use your full name and contact information.",
			},
			"configured": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
				Description: "Custom setup completion indicator." +
					"true = Custom app setup complete." +
					"false = Custom app setup not complete",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Short app description also displayed below the app title in Splunk Web Launcher.",
			},
			"explicit_appname": {
				Type:        schema.TypeString,
				Description: "Custom app name. Required when installing an app from a file where filename is set to true.",
				Optional:    true,
			},
			"filename": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Indicates whether to use the name value as the app source location." +
					"true indicates that name is a path to a file to install." +
					"false indicates that name is the literal app name and that the app is created from Splunkbase using a template.",
				RequiredWith: []string{"explicit_appname"},
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "App name displayed in Splunk Web, from five to 80 characters and excluding the prefix \"Splunk For\".",
			},
			"session": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login session token for installing or updating an app on Splunkbase. Alternatively, use auth.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "App version.",
			},
			"update": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "File-based update indication: true specifies that filename should be used to update an existing app. " +
					"If not specified, update defaults to false, which indicates that filename should not be used to update an existing app.",
			},
			"visible": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: "Indicates whether the app is visible and navigable from Splunk Web." +
					"true = App is visible and navigable." +
					"false = App is not visible or navigable.",
			},
			"acl": aclSchema(),
		},
		Read:   appsLocalRead,
		Create: appsLocalCreate,
		Delete: appsLocalDelete,
		Update: appsLocalUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func appsLocalCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	appsLocalConfigObj := getAppsLocalConfig(d)
	err := (*provider.Client).CreateAppsLocalObject(name, appsLocalConfigObj)
	if err != nil {
		return err
	}
	if e, ok := d.GetOk("explicit_appname"); ok {
		name = e.(string)
	}
	if r, ok := d.GetOk("acl"); ok {
		aclObject := getACLConfig(r.([]interface{}))
		err = (*provider.Client).UpdateAcl("nobody", "system", name, aclObject, "apps", "local")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return appsLocalRead(d, meta)
}
func appsLocalRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	resp, err := (*provider.Client).ReadAllAppsLocalObject()
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	entry, err := getAppsLocalConfigByName(name, resp)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	resp, err = (*provider.Client).ReadAppsLocalObject(name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	entry, err = getAppsLocalConfigByName(name, resp)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	if err = d.Set("author", entry.Content.Author); err != nil {
		return err
	}

	if err = d.Set("configured", entry.Content.Configured); err != nil {
		return err
	}

	if err = d.Set("description", entry.Content.Description); err != nil {
		return err
	}

	if err = d.Set("label", entry.Content.Label); err != nil {
		return err
	}

	if err = d.Set("version", entry.Content.Version); err != nil {
		return err
	}

	if err = d.Set("visible", entry.Content.Visible); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}
func appsLocalUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	appsLocalObject := getAppsLocalConfig(d)
	err := (*provider.Client).UpdateAppsLocalObject(d.Id(), appsLocalObject)
	if r, ok := d.GetOk("acl"); ok {
		aclObject := getACLConfig(r.([]interface{}))
		err = (*provider.Client).UpdateAcl("nobody", "system", d.Id(), aclObject, "apps", "local")
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return appsLocalRead(d, meta)
}
func appsLocalDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).DeleteAppsLocalObject(d.Id())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200, 201:
		return nil
	default:
		errorResponse := &models.AppsLocalResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getAppsLocalConfig(d *schema.ResourceData) (appsLocalObject *models.AppsLocalObject) {
	appsLocalObject = &models.AppsLocalObject{}
	appsLocalObject.Auth = d.Get("auth").(string)
	appsLocalObject.Author = d.Get("author").(string)
	appsLocalObject.Configured = d.Get("configured").(bool)
	appsLocalObject.Description = d.Get("description").(string)
	appsLocalObject.ExplicitAppName = d.Get("explicit_appname").(string)
	appsLocalObject.Filename = d.Get("filename").(bool)
	appsLocalObject.Label = d.Get("label").(string)
	appsLocalObject.Name = d.Get("name").(string)
	appsLocalObject.Session = d.Get("session").(string)
	appsLocalObject.Version = d.Get("version").(string)
	appsLocalObject.Visible = d.Get("visible").(bool)
	appsLocalObject.Update = d.Get("update").(bool)
	return appsLocalObject
}

func getAppsLocalConfigByName(name string, httpResponse *http.Response) (appsLocalEntry *models.AppsLocalEntry, err error) {
	response := &models.AppsLocalResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:
		_ = json.NewDecoder(httpResponse.Body).Decode(&response)
		re := regexp.MustCompile(`(.*)`)
		for _, entry := range response.Entry {
			if name == re.FindStringSubmatch(entry.Name)[1] {
				return &entry, nil
			}
		}
	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return appsLocalEntry, err
	}
	return appsLocalEntry, nil
}
