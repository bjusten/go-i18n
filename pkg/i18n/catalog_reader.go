package i18n

import (
	"context"
	"os"
)

type CatalogReader struct {
	catalog *Catalog
	locale  string
}

const (
	ReaderLocaleEnvironment = "READER_LOCALE"
)

// NewCatalogReader returns a new I18N Catalog Reader
func NewCatalogReader() *CatalogReader {
	cr := &CatalogReader{
		locale: "default",
	}

	env, exists := os.LookupEnv(ReaderLocaleEnvironment)
	if exists {
		cr.WithLocale(env)
	}

	return cr
}

// Get returns the KeyValue associated with the specified key
func (cr *CatalogReader) Get(key string) KeyValue {
	if cr != nil {
		return cr.catalog.Get(cr.locale, key)
	}
	return NewUnknownKeyPair(key)
}

// GetWithLocale returns the KeyValue associated with the specified key using the specified locale
func (cr *CatalogReader) GetWithLocale(locale string, key string) KeyValue {
	if cr != nil {
		return cr.catalog.Get(locale, key)
	}
	return NewUnknownKeyPair(key)
}

// WithCatalog sets the catalog to use for the catalog reader to the specified catalog
func (cr *CatalogReader) WithCatalog(catalog *Catalog) *CatalogReader {
	if cr != nil {
		cr.catalog = catalog
	}
	return cr
}

// WithCatalogFromContext set the catalog to use for the catalog reader found in the specified context
func (cr *CatalogReader) WithCatalogFromContext(ctx context.Context) *CatalogReader {
	return cr.WithCatalog(CatalogFromContext(ctx))
}

// WithLocale sets the locale to use for the catalog reader to the specified locale
func (cr *CatalogReader) WithLocale(locale string) *CatalogReader {
	if cr != nil {
		cr.locale = locale
	}
	return cr
}

// WithContext will add, and return, the catalog reader to the specified context
func (cr *CatalogReader) WithContext(ctx context.Context) (*CatalogReader, context.Context) {
	switch {
	case cr == nil:
		return nil, ctx
	default:
		ctx = context.WithValue(ctx, CatalogReaderContextKey, cr)
		return cr, ctx
	}
}

// WithNewContext will add, and return, the catalog reader to a new context
func (cr *CatalogReader) WithNewContext() (*CatalogReader, context.Context) {
	return cr.WithContext(context.Background())
}

// CatalogReaderFromContext will return the catalog reader contained in the specified context
func CatalogReaderFromContext(ctx context.Context) *CatalogReader {
	switch {
	case ctx == nil:
		return nil
	default:
		return ctx.Value(CatalogReaderContextKey).(*CatalogReader)
	}
}
