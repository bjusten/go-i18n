package main

import (
	"context"
	"fmt"

	"github.com/bjusten/go-i18n/pkg/i18n"
)

func main() {
	// Create a catalog, load the local locales and return it with a new context
	catalog, catalogCtx := makeCatalogWithContext()

	// Create a catalog reader, populate it with a catalog from a context and return a new context
	_, readerCtx := i18n.NewCatalogReader().WithCatalogFromContext(catalogCtx).WithNewContext()

	// Output key values using the reader context using data from the catalog context
	printKey(readerCtx, "en", "key-1")
	printKey(readerCtx, "fr", "key-1")

	// Output a key value from a locale that doesn't exist
	printKey(readerCtx, "invalid", "key-1")

	// Set a default locale for the catalog
	_ = catalog.WithDefaultLocale("fr")

	// Output a key value from a locale that doesn't exist and fallback to the french locale
	printKey(readerCtx, "invalid", "key-1")
}

func makeCatalogWithContext() (*i18n.Catalog, context.Context) {
	// Create a catalog with a new context
	catalog, ctx := i18n.NewCatalog().WithNewContext()

	// Initialize the catalog using the default './locales' directory
	if _, err := catalog.Initialize(); err != nil {
		panic("failed to load locales: " + err.Error())
	}

	return catalog, ctx
}

func printKey(ctx context.Context, locale string, key string) {
	// Get the catalog reader from the context
	catalogReader := i18n.CatalogReaderFromContext(ctx)

	// Print the key value
	fmt.Printf("%s: %s\n", locale, catalogReader.GetWithLocale(locale, key).Value())
}
