package zendesk

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
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
			zd := i.(zendesk.AttachmentAPI)
			return readAttachment(data, zd)
		},
		Delete: func(data *schema.ResourceData, i interface{}) error {
			zd := i.(zendesk.AttachmentAPI)
			return deleteAttachment(data, zd)
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
				ForceNew: true,
			},
			"file_hash": {
				Type:     schema.TypeString,
				Computed: true,
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
		return fmt.Errorf("error reading file %v: %v", filePath, err)
	}
	defer file.Close()

	fileName := d.Get("file_name").(string)
	w := zd.UploadAttachment(fileName, "")
	tee := io.TeeReader(file, w)

	h := sha1.New()
	_, err = io.Copy(h, tee)
	if err != nil {
		return fmt.Errorf("error reading file %v: %v", fileName, err)
	}

	result, err := w.Close()
	if err != nil {
		return fmt.Errorf("getting response from zendesk: %v", err)
	}

	a := result.Attachment
	d.SetId(strconv.FormatInt(a.ID, 10))
	err = d.Set("token", result.Token)
	if err != nil {
		return err
	}

	return marshalAttachment(d, attachment{
		Attachment: a,
		FilePath:   filePath,
		Hash:       h.Sum(nil),
	})
}

func deleteAttachment(d identifiableGetterSetter, zd zendesk.AttachmentAPI) error {
	v, ok := d.GetOk("token")
	if !ok {
		return nil
	}

	return zd.DeleteUpload(v.(string))
}

func readAttachment(d identifiableGetterSetter, zd zendesk.AttachmentAPI) error {
	out := attachment{}
	if v, ok := d.GetOk("file_path"); ok {
		file, err := os.Open(v.(string))
		if err != nil {
			return fmt.Errorf("error opening file %v: %v", v, err)
		}

		defer file.Close()

		h := sha1.New()
		_, err = io.Copy(h, file)
		if err != nil {
			return fmt.Errorf("error reading file %v: %v", v, err)
		}

		out.FilePath = v.(string)
		out.Hash = h.Sum(nil)
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	a, err := zd.GetAttachment(id)
	if err != nil {
		return fmt.Errorf("getting response from zendesk: %v", err)
	}

	out.Attachment = a

	return marshalAttachment(d, out)
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

	//thumbnails := make([]map[string]interface{}, len(a.Thumbnails))
	//for _, v := range a.Thumbnails {
	//	thumb := map[string]interface{}{
	//		"id":           v.ID,
	//		"file_name":    v.FileName,
	//		"content_url":  v.ContentURL,
	//		"content_type": v.ContentType,
	//		"size":         v.Size,
	//	}
	//	thumbnails = append(thumbnails, thumb)
	//}
	//
	//m["thumbnails"] = thumbnails
	return setSchemaFields(d, m)
}
