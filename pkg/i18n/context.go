package i18n

type contextKey string

const (
	CatalogContextKey       = contextKey("i18n-catalog")
	CatalogReaderContextKey = contextKey("i18n-catalog-reader")
)
