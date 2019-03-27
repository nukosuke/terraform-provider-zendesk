package zendesk

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

func TestTriggerDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := mock.NewClient(ctrl)
	d := newIdentifiableGetterSetter()

	d.SetId("1234")

	c.EXPECT().DeleteTrigger(gomock.Eq(int64(1234))).Return(nil)
	err := resourceZendeskTriggerDelete(d, c)
	if err != nil {
		t.Fatalf("Got error from resource delete: %v", err)
	}
}
