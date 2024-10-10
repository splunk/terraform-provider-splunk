package splunk

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/splunk/terraform-provider-splunk/client/models"
	"io"
)

func lookupTableFile() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"app": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The parent app to the lookup.",
			},
			"owner": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The owner of the lookup.",
			},
			"file_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "A file name for the lookup.",
			},
			"file_contents": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				Description: "The contents of the lookup.",
			},
		},
		Create: lookupTableFileCreate,
		Read:   lookupTableFileRead,
		Update: lookupTableFileUpdate,
		Delete: lookupTableFileDelete,
	}
}

func lookupTableFileCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	lookupTableFile := getLookupTableFile(d)

	err := (*provider.Client).CreateLookupTableFile(lookupTableFile.FileName, lookupTableFile.Owner, lookupTableFile.App, lookupTableFile.FileContents)
	if err != nil {
		return err
	}

	d.SetId(lookupTableFile.FileName)
	return lookupTableFileRead(d, meta)
}

func lookupTableFileRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	lookupTableFile := getLookupTableFile(d)

	resp, err := (*provider.Client).ReadLookupTableFile(lookupTableFile.FileName, lookupTableFile.Owner, lookupTableFile.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var fileContents [][]string
	if err := json.Unmarshal(bodyBytes, &fileContents); err != nil {
		return err
	}

	if err = d.Set("file_contents", fileContents); err != nil {
		return err
	}

	return nil
}

func lookupTableFileUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	lookupTableFile := getLookupTableFile(d)

	err := (*provider.Client).UpdateLookupTableFile(lookupTableFile.FileName, lookupTableFile.Owner, lookupTableFile.App, lookupTableFile.FileContents)
	if err != nil {
		return err
	}

	return lookupTableFileRead(d, meta)
}

func lookupTableFileDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	lookupTableFile := getLookupTableFile(d)

	resp, err := (*provider.Client).DeleteLookupTableFile(lookupTableFile.FileName, lookupTableFile.Owner, lookupTableFile.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.InputsUDPResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

func getLookupTableFile(d *schema.ResourceData) (lookupTableFile *models.LookupTableFile) {
	fileContents, _ := json.Marshal(d.Get("file_contents"))
	lookupTableFile = &models.LookupTableFile{
		App:          d.Get("app").(string),
		Owner:        d.Get("owner").(string),
		FileName:     d.Get("file_name").(string),
		FileContents: string(fileContents),
	}
	return lookupTableFile
}
