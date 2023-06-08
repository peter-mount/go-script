package calculator

type KeyValue struct {
	key   string
	value interface{}
}

func NewKeyValue(key string, value interface{}) *KeyValue {
	return &KeyValue{key: key, value: value}
}

func GetKeyValue(v interface{}) (*KeyValue, bool) {
	if v == nil {
		return nil, false
	}

	kv, ok := v.(*KeyValue)
	return kv, ok
}

func (v *KeyValue) Key() string { return v.key }

func (v *KeyValue) Value() interface{} { return v.value }
