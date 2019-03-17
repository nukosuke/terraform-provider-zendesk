package zendesk

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/pkg/errors"
)

type attachment struct {
	zendesk.Attachment
	FilePath string
	Hash     []byte
}

func resourceZendeskAttachment() *schema.Resource {
	return &schema.Resource{
		Create: func(data *schema.ResourceData, i interface{}) error {
			zd := i.(zendesk.AttachmentAPI)
			return createAttachment(data, zd)
		},
		Read: func(data *schema.ResourceData, i interface{}) error {
			return nil
		},
		Delete: func(data *schema.ResourceData, i interface{}) error {
			return nil
		},
		Update: func(data *schema.ResourceData, i interface{}) error {
			return errors.New("Update attachment not supported")
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"file_path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: isValidFile(),
			},
			"file_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_hash": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"content_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"inline": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"thumbnails": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"file_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"content_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func createAttachment(d identifiableGetterSetter, zd zendesk.AttachmentAPI) error {
	path := d.Get("file_path").(string)
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	fileName := d.Get("file_name").(string)
	w := zd.UploadAttachment(fileName, "")
	tee := io.TeeReader(file, w)

	h := sha1.New()
	_, err = io.Copy(h, tee)
	if err != nil {
		return err
	}

	result, err := w.Close()
	if err != nil {
		return err
	}

	a := result.Attachment
	d.SetId(strconv.FormatInt(a.ID, 10))
	return marshalAttachment(d, attachment{
		Attachment: a,
		FilePath:   path,
		Hash:       h.Sum(nil),
	})
}

func marshalAttachment(d identifiableGetterSetter, a attachment) error {
	m := map[string]interface{}{
		"file_path":    a.FilePath,
		"file_hash":    hex.EncodeToString(a.Hash),
		"file_name":    a.FileName,
		"content_url":  a.ContentURL,
		"content_type": a.ContentType,
		"size":         a.Size,
		"inline":       a.Inline,
	}

	thumbnails := make([]map[string]interface{}, len(a.Thumbnails))
	for _, v := range a.Thumbnails {
		thumb := map[string]interface{}{
			"id":           v.ID,
			"file_name":    v.FileName,
			"content_url":  v.ContentURL,
			"content_type": v.ContentType,
			"size":         v.Size,
		}
		thumbnails = append(thumbnails, thumb)
	}

	m["thumbnails"] = thumbnails
	return setSchemaFields(d, m)
}
