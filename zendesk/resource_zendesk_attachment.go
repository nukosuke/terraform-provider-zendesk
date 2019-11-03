package zendesk

import (
	"context"
	"io"
	"os"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nukosuke/go-zendesk/zendesk"
)

type attachment struct {
	zendesk.Attachment
	FilePath string
	Hash     string
}

func resourceZendeskAttachment() *schema.Resource {
	return &schema.Resource{
		Create: func(data *schema.ResourceData, i interface{}) error {
			zd := i.(zendesk.AttachmentAPI)
			return createAttachment(data, zd)
		},
		Read: func(data *schema.ResourceData, i interface{}) error {
			zd := i.(zendesk.AttachmentAPI)
			return readAttachment(data, zd)
		},
		Delete: func(data *schema.ResourceData, i interface{}) error {
			zd := i.(zendesk.AttachmentAPI)
			return deleteAttachment(data, zd)
		},
		Update: func(data *schema.ResourceData, i interface{}) error {
			zd := i.(zendesk.AttachmentAPI)
			return readAttachment(data, zd)
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
				ForceNew: true,
			},
			"file_hash": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"token": {
				Type:     schema.TypeString,
				Computed: true,
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
						"content_url": {
							Type:     schema.TypeString,
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
	filePath := d.Get("file_path").(string)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileName := d.Get("file_name").(string)
	ctx := context.Background()
	w := zd.UploadAttachment(ctx, fileName, "")

	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}

	result, err := w.Close()
	if err != nil {
		return err
	}

	a := result.Attachment
	d.SetId(strconv.FormatInt(a.ID, 10))
	err = d.Set("token", result.Token)
	if err != nil {
		return err
	}

	out := attachment{
		Attachment: a,
		FilePath:   filePath,
		Hash:       d.Get("file_hash").(string),
	}
	return marshalAttachment(d, out)
}

func deleteAttachment(d identifiableGetterSetter, zd zendesk.AttachmentAPI) error {
	v, ok := d.GetOk("token")
	if !ok {
		return nil
	}

	ctx := context.Background()
	return zd.DeleteUpload(ctx, v.(string))
}

func readAttachment(d identifiableGetterSetter, zd zendesk.AttachmentAPI) error {
	out := attachment{}

	if v, ok := d.GetOk("file_path"); ok {
		out.FilePath = v.(string)
	}

	if v, ok := d.GetOk("file_hash"); ok {
		out.Hash = v.(string)
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	a, err := zd.GetAttachment(ctx, id)
	if err != nil {
		return err
	}

	out.Attachment = a

	return marshalAttachment(d, out)
}

func marshalAttachment(d identifiableGetterSetter, a attachment) error {
	m := map[string]interface{}{
		"file_path":    a.FilePath,
		"file_hash":    a.Hash,
		"file_name":    a.FileName,
		"content_url":  a.ContentURL,
		"content_type": a.ContentType,
		"size":         a.Size,
		"inline":       a.Inline,
	}

	thumbnails := make([]map[string]interface{}, 0)
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
