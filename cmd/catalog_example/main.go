package main

import (
	"fmt"

	"github.com/bjusten/go-i18n/pkg/i18n"
)

const Locales = "en" // Filter locales and only load english

func main() {
	// Create a catalog, filter it to only populate the english locale and initialize
	catalog, err := i18n.NewCatalog().WithLocales(Locales).Initialize()
	if err != nil {
		panic(fmt.Sprintf("failed to load locales: %v", err))
	}

	// Print english key value, good result
	fmt.Printf("english: %s\n", catalog.Get("en", "key-1").Value())

	// Print french key value, bad result (fr is not a loaded locale)
	fmt.Printf("french: %s\n", catalog.Get("fr", "key-1").Value())

	// Set the catalog default locale to english
	catalog = catalog.WithDefaultLocale("en")

	// Print french key value, good result (fr is not loaded, fallback to the default locale)
	fmt.Printf("french: %s\n", catalog.Get("fr", "key-1").Value())
}
