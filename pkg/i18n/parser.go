package i18n

// Parser is an interface primarily used by the Catalog for loading keyValue data sets
type Parser interface {
	Parse(func(locale string, keyValue KeyValue)) error
}

// KeyValue is an interface used throughout the I18N package as a generic key and value storage object
type KeyValue interface {
	Key() string
	Value() string
}
