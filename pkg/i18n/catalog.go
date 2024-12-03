// The i18n package is a general purpose I18N implementation.
package i18n

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Catalog struct {
	parser Parser

	locales map[string]map[string]KeyValue

	defaultLocale string
	localeFilters map[string]bool

	stats CatalogStats

	lock sync.RWMutex
}

type CatalogStats struct {
	Locales int
	Keys    int
}

const (
	DefaultLocaleEnvironment = "DEFAULT_LOCALE"
	FilterLocalesEnvironment = "FILTER_LOCALES"
)

// NewCatalog returns a new I18N Catalog
func NewCatalog() *Catalog {
	c := &Catalog{
		parser:        NewKeyPairFSParser([]string{"./locales"}),
		locales:       make(map[string]map[string]KeyValue),
		localeFilters: make(map[string]bool),
	}

	env, exists := os.LookupEnv(DefaultLocaleEnvironment)
	if exists {
		c.WithDefaultLocale(env)
	}

	env, exists = os.LookupEnv(FilterLocalesEnvironment)
	if exists {
		tokens := strings.Split(env, ",")
		c.WithLocales(tokens...)
	}

	return c
}

// WithParser sets the catalog parsing interface to the specified parser
func (c *Catalog) WithParser(parser Parser) *Catalog {
	if c != nil {
		c.parser = parser
	}
	return c
}

// WithDefaultLocale sets the default (fallback) locale to use if a key is not found with the specified locale
func (c *Catalog) WithDefaultLocale(locale string) *Catalog {
	if c != nil {
		c.defaultLocale = locale
	}
	return c
}

// WithLocales will take the specified locales and apply them as a filter when loading catalog entries
func (c *Catalog) WithLocales(locales ...string) *Catalog {
	if c != nil {
		c.lock.Lock()
		defer c.lock.Unlock()

		for _, locale := range locales {
			c.localeFilters[locale] = true
		}
	}
	return c
}

var (
	errNoCatalog = errors.New("catalog is nil")
	errNoParser  = errors.New("catalog parser is nil")
)

// Initialize loads keyValues using the catalog parser
func (c *Catalog) Initialize() (*Catalog, error) {
	switch {
	case c == nil:
		return nil, errNoCatalog
	case c.parser == nil:
		return c, errNoParser
	default:
		return c, c.parser.Parse(c.AddKeyValue)
	}
}

// Initialize loads keyValues using the catalog parser and returns a new context
func (c *Catalog) InitializeWithContext() (context.Context, error) {
	switch {
	case c == nil:
		return nil, errNoCatalog
	default:
		if _, err := c.Initialize(); err != nil {
			return nil, err
		}

		_, ctx := c.WithNewContext()
		return ctx, nil
	}
}

// Get returns the KeyValue for the specified key for the specified locale
func (c *Catalog) Get(locale string, key string) KeyValue {
	if c == nil {
		return NewUnknownKeyPair(key)
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	localeEntry, exists := c.locales[locale]
	switch {
	case !exists && len(c.defaultLocale) > 0 && c.defaultLocale != locale:
		return c.Get(c.defaultLocale, key)
	case !exists:
		return NewUnknownKeyPair(key)
	}

	value, exists := localeEntry[key]
	switch {
	case !exists && len(c.defaultLocale) > 0 && c.defaultLocale != locale:
		return c.Get(c.defaultLocale, key)
	case !exists:
		return NewUnknownKeyPair(key)
	}

	return value
}

// Locales returns a copy of the catalog locale filters
func (c *Catalog) Locales() []string {
	switch {
	case c == nil:
		return []string{}
	default:
		c.lock.RLock()
		defer c.lock.RUnlock()

		locales := make([]string, 0)
		for locale := range c.localeFilters {
			locales = append(locales, locale)
		}

		return locales
	}
}

// Stats returns a copy of the catalog stats
func (c *Catalog) Stats() CatalogStats {
	switch {
	case c == nil:
		return CatalogStats{}
	default:
		return CatalogStats{
			Locales: c.stats.Locales,
			Keys:    c.stats.Keys,
		}
	}
}

// WithContext will add, and return, the catalog to the specified context
func (c *Catalog) WithContext(ctx context.Context) (*Catalog, context.Context) {
	switch {
	case c == nil:
		return nil, ctx
	default:
		ctx = context.WithValue(ctx, CatalogContextKey, c)
		return c, ctx
	}
}

// WithNewContext will add, and return, the catalog to a new context
func (c *Catalog) WithNewContext() (*Catalog, context.Context) {
	return c.WithContext(context.Background())
}

// AddKeyValue will add the specified keyValue to the catalog under the specified locale
func (c *Catalog) AddKeyValue(locale string, keyValue KeyValue) {
	if c == nil {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	filtered, exists := c.localeFilters[locale]
	if len(c.localeFilters) > 0 && (!exists || !filtered) {
		return
	}

	localeEntry, exists := c.locales[locale]
	if !exists {
		localeEntry = make(map[string]KeyValue)
		c.locales[locale] = localeEntry
		c.stats.Locales++
	}

	_, exists = localeEntry[keyValue.Key()]
	if !exists {
		c.stats.Keys++
	}

	localeEntry[keyValue.Key()] = keyValue
}

// CatalogFromContext will return the catalog contained in the specified context
func CatalogFromContext(ctx context.Context) *Catalog {
	switch {
	case ctx == nil:
		return nil
	default:
		return ctx.Value(CatalogContextKey).(*Catalog)
	}
}

func (c *Catalog) PrintAll() {
	for k := range c.locales {
		fmt.Println(k)
		for k, v := range c.locales[k] {
			fmt.Printf("\t%s: %s\n", k, v)
		}
	}
}
