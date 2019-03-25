package zendesk

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/resource"
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

	m.EXPECT().UploadAttachment(Any(), Any()).Return(w)

	d := &identifiableMapGetterSetter{
		mapGetterSetter: mapGetterSetter{
			"file_path": "testdata/street.jpg",
			"file_name": "street.jpg",
			"file_hash": "foo",
		},
	}

	err := createAttachment(d, m)
	if err != nil {
		t.Fatalf("Create attachment returned an error %v", err)
	}

	if v := d.Id(); v != "1234" {
		t.Fatalf("Resource did not have expected id %s", v)
	}
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
		Providers: testAccProviders,
		//TODO: check destroyed
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(attachmentConfig, original.Name(), hashString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("zendesk_attachment.file", "content_url"),
				),
			},
			{
				Config: fmt.Sprintf(attachmentConfig, tmpfile.Name(), hashString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("zendesk_attachment.file", "content_url"),
				),
			},
		},
	})

}
