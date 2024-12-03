package main

import (
	"fmt"

	"github.com/bjusten/go-i18n/pkg/i18n"
)

const JSONFile = "./locales.json"

func main() {
	// Create a new catalog using our custom JSON parser and initialize it
	catalog, err := i18n.NewCatalog().WithParser(NewJSONParser(JSONFile)).Initialize()
	if err != nil {
		panic(fmt.Sprintf("error loading catalog: %v", err))
	}

	printKey(catalog, "en", "key-1")
	printKey(catalog, "fr", "key-1")
}

func printKey(catalog *i18n.Catalog, locale string, key string) {
	// Pull out the keyValue and cast it to our custom type
	kvwd, ok := catalog.Get(locale, key).(KeyValueWithDescription)
	if !ok {
		panic("keyValue is not of type: KeyValueWithDescription")
	}

	// Print the key value along with our custom description
	fmt.Printf("%s: %s (%s)\n", locale, kvwd.Value(), kvwd.Description())
}
