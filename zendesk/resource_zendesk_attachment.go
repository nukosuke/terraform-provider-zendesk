package zendesk

import (
	"context"
	"io"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nukosuke/go-zendesk/zendesk"
)

type attachment struct {
	zendesk.Attachment
	FilePath string
	Hash     string
}

func resourceZendeskAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Provides an attachment resource.",
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
				Description: "The name of the image file.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"file_hash": {
				Description: "SHA256 hash of the image file. Terraform built-in `filesha256()` is convenient to calculate it.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"token": {
				Description: "The token of the uploaded attachment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"content_url": {
				Description: "A full URL where the attachment image file can be downloaded. The file may be hosted externally so take care not to inadvertently send Zendesk authentication credentials.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"content_type": {
				Description: `The content type of the image. Example value: "image/png"`,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"size": {
				Description: "The size of the image file in bytes.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"inline": {
				Description: "If true, the attachment is excluded from the attachment list and the attachment's URL can be referenced within the comment of a ticket. Default is false.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"thumbnails": {
				Description: "A list of attachments.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Attachment id.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"file_name": {
							Description: "File name of the image",
							Type:        schema.TypeString,
							Required:    true,
						},
						"content_type": {
							Description: "Content-Type of the image",
							Type:        schema.TypeString,
							Required:    true,
						},
						"size": {
							Description: "File size of the image.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"content_url": {
							Description: "URL of the image.",
							Type:        schema.TypeString,
							Required:    true,
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
