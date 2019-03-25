package zendesk

import (
	"io"
	"io/ioutil"
	"testing"

	. "github.com/golang/mock/gomock"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

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
