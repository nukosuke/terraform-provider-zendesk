package zendesk

type getter interface {
	Get(string) interface{}
}

type setter interface {
	Set(string, interface{}) error
}

type identifiable interface {
	Id() string
}

type identifiableGetterSetter interface {
	identifiable
	getter
	setter
}

type mapGetterSetter map[string]interface{}

func (m mapGetterSetter) Get(k string) interface{} {
	v, ok := m[k]
	if !ok {
		return nil
	}

	return v
}

func (m mapGetterSetter) Set(k string, v interface{}) error {
	m[k] = v
	return nil
}

type identifiableMapGetterSetter struct {
	mapGetterSetter
	id string
}

func (i identifiableMapGetterSetter) Id() string {
	return i.id
}

func setSchemaFields(d setter, m map[string]interface{}) error {
	for k, v := range m {
		err := d.Set(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
