package zendesk

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

func TestMarshalTrigger(t *testing.T) {
	expected := zendesk.Trigger{
		Title:       "title",
		Description: "blabla",
		Active:      true,
	}
	m := &identifiableMapGetterSetter{
		mapGetterSetter: mapGetterSetter{},
	}

	err := marshalTrigger(expected, m)
	if err != nil {
		t.Fatalf("Failed to marshal map %v", err)
	}

	v, ok := m.GetOk("title")
	if !ok {
		t.Fatal("Failed to get title value")
	}
	if v != expected.Title {
		t.Fatalf("trigger had incorrect title value %v. should have been %v", v, expected.Title)
	}

	v, ok = m.GetOk("description")
	if !ok {
		t.Fatal("Failed to get description value")
	}
	if v != expected.Description {
		t.Fatalf("trigger had incorrect description value %v. should have been %v", v, expected.Description)
	}

	v, ok = m.GetOk("active")
	if !ok {
		t.Fatal("Failed to get active value")
	}
	if v != expected.Active {
		t.Fatalf("trigger had incorrect active value %v. should have been %v", v, expected.Active)
	}
}

func TestUnmarshalTrigger(t *testing.T) {
	m := &identifiableMapGetterSetter{
		id: "100",
		mapGetterSetter: mapGetterSetter{
			"title":       "Auto reply",
			"description": "reply automatically",
			"active":      true,
		},
	}

	trg, err := unmarshalTrigger(m)
	if err != nil {
		t.Fatalf("unmarshal returned an error: %v", err)
	}

	if v := m.Get("title"); trg.Title != v {
		t.Fatalf("trigger had title value %v. should have been %v", trg.Title, v)
	}
}

func TestDeleteTrigger(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := mock.NewClient(ctrl)
	d := newIdentifiableGetterSetter()

	d.SetId("1234")

	c.EXPECT().DeleteTrigger(gomock.Eq(int64(1234))).Return(nil)
	err := deleteTrigger(d, c)
	if err != nil {
		t.Fatalf("Got error from resource delete: %v", err)
	}
}
