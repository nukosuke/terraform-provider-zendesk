package zendesk

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	. "github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

const attachmentConfig = `
resource "zendesk_attachment" "file" {
  file_name = "street.jpg"
  file_path = "%s"
  file_hash = "%s"
}
`

type mockUploadWriter struct {
	io.Writer
	Response zendesk.Upload
	Error    error
}

func (w mockUploadWriter) Close() (zendesk.Upload, error) {
	return w.Response, w.Error
}

func newMockUploadWriter(u zendesk.Upload, err error) zendesk.UploadWriter {
	return mockUploadWriter{
		Writer:   ioutil.Discard,
		Response: u,
		Error:    err,
	}
}

func TestCreateZendeskAttachment(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	u := zendesk.Upload{
		Attachment: zendesk.Attachment{
			ID:          1234,
			FileName:    "foo",
			ContentURL:  "content",
			ContentType: "jpg",
			Size:        1,
			Inline:      false,
		},
	}
	w := newMockUploadWriter(u, nil)

	m.EXPECT().UploadAttachment(Any(), Any(), Any()).Return(w)

	d := &identifiableMapGetterSetter{
		mapGetterSetter: mapGetterSetter{
			"file_path": "testdata/street.jpg",
			"file_name": "street.jpg",
			"file_hash": "foo",
		},
	}

	diags := createAttachment(context.Background(), d, m)
	if len(diags) != 0 {
		t.Fatalf("Create attachment returned an error %v", diags)
	}

	if v := d.Id(); v != "1234" {
		t.Fatalf("Resource did not have expected id %s", v)
	}
}

func TestDeleteZendeskAttachmentCallsWhenTokenIsSet(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)

	m.EXPECT().DeleteUpload(Any(), Any()).Return(nil)

	d := &identifiableMapGetterSetter{
		mapGetterSetter: mapGetterSetter{
			"token": "foo",
		},
	}

	diags := deleteAttachment(context.Background(), d, m)
	if len(diags) != 0 {
		t.Fatalf("delete attachment returned an error %v", diags)
	}
}

func TestDeleteZendeskAttachmentDoesNotCallWhenTokenIsNotSet(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)

	d := newIdentifiableGetterSetter()

	diags := deleteAttachment(context.Background(), d, m)
	if len(diags) != 0 {
		t.Fatalf("delete attachment returned an error %v", diags)
	}
}

func testAttachmentDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(zendesk.AttachmentAPI)

	for _, r := range s.RootModule().Resources {
		if r.Type != "zendesk_attachment" {
			continue
		}

		id, err := atoi64(r.Primary.ID)
		if err != nil {
			return err
		}

		_, err = client.GetAttachment(context.Background(), id)
		if err == nil {
			return fmt.Errorf("did not get error from zendesk when trying to fetch the destroyed ticket attachment")
		}

		zd, ok := err.(zendesk.Error)
		if !ok {
			return fmt.Errorf("error %v cannot be asserted as a zendesk error", err)
		}

		if zd.Status() != http.StatusNotFound {
			return fmt.Errorf(`did not get a "not found error"" after destroy. error was %v`, zd)
		}

	}

	return nil
}

func TestAccZendeskAttachment(t *testing.T) {
	original, err := os.Open("testdata/street.jpg")
	if err != nil {
		t.Fatalf("could not open street file")
	}
	defer original.Close()

	h := sha1.New()
	tee := io.TeeReader(original, h)

	tmpfile, err := ioutil.TempFile("", "new-streets.jpg")
	if err != nil {
		t.Fatalf("could not create temp file")
	}

	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	_, err = io.Copy(tmpfile, tee)
	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}

	hashString := hex.EncodeToString(h.Sum(nil))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAttachmentDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(attachmentConfig, original.Name(), hashString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("zendesk_attachment.file", "content_url"),
					resource.TestCheckResourceAttr("zendesk_attachment.file", "file_path", original.Name()),
				),
			},
			{
				Config: fmt.Sprintf(attachmentConfig, tmpfile.Name(), hashString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("zendesk_attachment.file", "content_url"),
					resource.TestCheckResourceAttr("zendesk_attachment.file", "file_path", tmpfile.Name()),
				),
			},
		},
	})

}
