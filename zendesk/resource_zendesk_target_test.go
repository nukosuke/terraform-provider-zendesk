package zendesk

import (
	"strconv"
	"testing"

	. "github.com/golang/mock/gomock"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

func TestMarshalTarget(t *testing.T) {
	expectedURL := "https://example.com"
	expectedType := "email_target"
	expectedTitle := "target :: email :: john.doe@example.com"
	expectedEmail := "john.doe@example.com"
	expectedSubject := "New ticket created"

	g := zendesk.Target{
		URL:  expectedURL,
		Type: expectedType,
		Title: expectedTitle,
		Email: expectedEmail,
		Subject: expectedSubject,
	}
	m := &identifiableMapGetterSetter{
		mapGetterSetter: mapGetterSetter{},
	}

	err := marshalTarget(g, m)
	if err != nil {
		t.Fatalf("Could marshal map %v", err)
	}

	v, ok := m.GetOk("url")
	if !ok {
		t.Fatalf("Failed to get url value")
	}
	if v != expectedURL {
		t.Fatalf("target had incorrect url value %v. should have been %v", v, expectedURL)
	}

	v, ok = m.GetOk("type")
	if !ok {
		t.Fatalf("Failed to get type value")
	}
	if v != expectedType {
		t.Fatalf("target had incorrect type value %v. should have been %v", v, expectedType)
	}

	v, ok = m.GetOk("title")
	if !ok {
		t.Fatalf("Failed to get title value")
	}
	if v != expectedTitle {
		t.Fatalf("target had incorrect title value %v. should have been %v", v, expectedTitle)
	}

	v, ok = m.GetOk("email")
	if !ok {
		t.Fatalf("Failed to get email value")
	}
	if v != expectedEmail {
		t.Fatalf("target had incorrect email value %v. should have been %v", v, expectedEmail)
	}

	v, ok = m.GetOk("subject")
	if !ok {
		t.Fatalf("Failed to get subject value")
	}
	if v != expectedSubject {
		t.Fatalf("target had incorrect subject value %v. should have been %v", v, expectedSubject)
	}
}

func TestUnmarshalTarget(t *testing.T) {
	m := &identifiableMapGetterSetter{
		id: "1234",
		mapGetterSetter: mapGetterSetter{
			"url":  "https://example.zendesk.com/api/v2/targets/360011737434.json",
			"type": "email_target",
			"title": "target :: email :: john.doe@example.com",
			"email": "john.doe@example.com",
			"subject": "New ticket created",
		},
	}

	g, err := unmarshalTarget(m)
	if err != nil {
		t.Fatalf("Could marshal map %v", err)
	}

	if v := m.Get("url"); g.URL != v {
		t.Fatalf("target had url value %v. should have been %v", g.URL, v)
	}

	if v := m.Get("type"); g.Type != v {
		t.Fatalf("target had type value %v. should have been %v", g.Type, v)
	}

	if v := m.Get("title"); g.Title != v {
		t.Fatalf("target had title value %v. should have been %v", g.Title, v)
	}

	if v := m.Get("email"); g.Email != v {
		t.Fatalf("target had email value %v. should have been %v", g.Email, v)
	}

	if v := m.Get("subject"); g.Subject != v {
		t.Fatalf("target had subject value %v. should have been %v", g.Subject, v)
	}
}

func TestReadTarget(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	id := 1234
	gs := &identifiableMapGetterSetter{
		mapGetterSetter: make(mapGetterSetter),
		id:              strconv.Itoa(id),
	}

	field := zendesk.Target{
		ID:   int64(id),
		URL:  "foo",
		Type: "email_target",
		Title: "target :: email :: john.doe@example.com",
		Email: "john.doe@example.com",
		Subject: "New ticket created",
	}

	m.EXPECT().GetTarget(Any()).Return(field, nil)
	if err := readTarget(gs, m); err != nil {
		t.Fatalf("readTarget returned an error: %v", err)
	}

	if v := gs.mapGetterSetter["url"]; v != field.URL {
		t.Fatalf("url field %v does not have expected value %v", v, field.URL)
	}

	if v := gs.mapGetterSetter["type"]; v != field.Type {
		t.Fatalf("type field %v does not have expected value %v", v, field.Type)
	}

	if v := gs.mapGetterSetter["title"]; v != field.Title {
		t.Fatalf("title field %v does not have expected value %v", v, field.Title)
	}

	if v := gs.mapGetterSetter["email"]; v != field.Email {
		t.Fatalf("email field %v does not have expected value %v", v, field.Email)
	}

	if v := gs.mapGetterSetter["subject"]; v != field.Subject {
		t.Fatalf("subject field %v does not have expected value %v", v, field.Subject)
	}
}

func TestCreateTarget(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		mapGetterSetter: make(mapGetterSetter),
	}

	out := zendesk.Target{
		ID: 12345,
	}

	m.EXPECT().CreateTarget(Any()).Return(out, nil)
	if err := createTarget(i, m); err != nil {
		t.Fatalf("create target returned an error: %v", err)
	}

	if v := i.Id(); v != "12345" {
		t.Fatalf("Create did not set resource id. Id was %s", v)
	}
}

func TestUpdateTarget(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		id:              "12345",
		mapGetterSetter: make(mapGetterSetter),
	}

	m.EXPECT().UpdateTarget(Eq(int64(12345)), Any()).Return(zendesk.Target{}, nil)
	if err := updateTarget(i, m); err != nil {
		t.Fatalf("updateTarget returned an error: %v", err)
	}
}

func TestDeleteTarget(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		id: "12345",
	}

	m.EXPECT().DeleteTarget(Eq(int64(12345))).Return(nil)
	if err := deleteTarget(i, m); err != nil {
		t.Fatalf("deleteTarget returned an error: %v", err)
	}
}
