package main

type KeyValueWithDescription struct {
	K string `json:"-"`
	V string `json:"value"`
	D string `json:"description"`
}

func (kvwd KeyValueWithDescription) Key() string {
	return kvwd.K
}

func (kvwd KeyValueWithDescription) Value() string {
	return kvwd.V
}

func (kvwd KeyValueWithDescription) Description() string {
	return kvwd.D
}

type JSONLocales struct {
	Locales map[string]map[string]KeyValueWithDescription `json:"locales"`
}
