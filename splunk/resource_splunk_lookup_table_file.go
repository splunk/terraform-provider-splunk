package splunk

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/splunk/terraform-provider-splunk/client/models"
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
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"file_contents", "file_path"},
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				Description: "The contents of the lookup.",
			},
			"file_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"file_contents", "file_path"},
				Description:  "Path to a local CSV file. Use this for large files to avoid sending the contents through Terraform's provider RPC.",
			},
			"file_contents_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hash of the CSV contents when file_path is used.",
			},
		},
		CustomizeDiff: lookupTableFileCustomizeDiff,
		Create:        lookupTableFileCreate,
		Read:          lookupTableFileRead,
		Update:        lookupTableFileUpdate,
		Delete:        lookupTableFileDelete,
	}
}

func lookupTableFileCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	lookupTableFile := getLookupTableFile(d)
	contents, err := lookupTableFileContents(d)
	if err != nil {
		return err
	}

	err = (*provider.Client).CreateLookupTableFile(lookupTableFile.FileName, lookupTableFile.Owner, lookupTableFile.App, contents)
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
	if d.Get("file_path").(string) != "" {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var fileContents [][]string
		if err := json.Unmarshal(bodyBytes, &fileContents); err != nil {
			return err
		}
		contents, err := json.Marshal(fileContents)
		if err != nil {
			return err
		}
		hash := sha256.Sum256(contents)
		return d.Set("file_contents_hash", hex.EncodeToString(hash[:]))
	}

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
	contents, err := lookupTableFileContents(d)
	if err != nil {
		return err
	}

	err = (*provider.Client).UpdateLookupTableFile(lookupTableFile.FileName, lookupTableFile.Owner, lookupTableFile.App, contents)
	if err != nil {
		return err
	}

	return lookupTableFileRead(d, meta)
}

func lookupTableFileCustomizeDiff(d *schema.ResourceDiff, _ interface{}) error {
	filePath := d.Get("file_path").(string)
	if filePath == "" {
		return d.SetNew("file_contents_hash", "")
	}

	contents, err := lookupTableFileContentsFromPath(filePath)
	if err != nil {
		return fmt.Errorf("reading lookup CSV %q: %w", filePath, err)
	}
	hash := sha256.Sum256([]byte(contents))
	return d.SetNew("file_contents_hash", hex.EncodeToString(hash[:]))
}

func lookupTableFileContents(d *schema.ResourceData) (string, error) {
	if filePath := d.Get("file_path").(string); filePath != "" {
		return lookupTableFileContentsFromPath(filePath)
	}
	return getLookupTableFile(d).FileContents, nil
}

func lookupTableFileContentsFromPath(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("opening lookup CSV %q: %w", filePath, err)
	}
	defer file.Close()

	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return "", fmt.Errorf("reading lookup CSV %q: %w", filePath, err)
	}
	contents, err := json.Marshal(rows)
	if err != nil {
		return "", err
	}
	return string(contents), nil
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
	rawContents := d.Get("file_contents").([]interface{})
	contents := make([][]string, len(rawContents))
	for i, row := range rawContents {
		rawRow := row.([]interface{})
		contents[i] = make([]string, len(rawRow))
		for j, val := range rawRow {
			if val == nil {
				contents[i][j] = ""
			} else {
				contents[i][j] = val.(string)
			}
		}
	}
	fileContents, _ := json.Marshal(contents)
	lookupTableFile = &models.LookupTableFile{
		App:          d.Get("app").(string),
		Owner:        d.Get("owner").(string),
		FileName:     d.Get("file_name").(string),
		FileContents: string(fileContents),
	}
	return lookupTableFile
}
