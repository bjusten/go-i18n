# About

This repo contains a general purpose internationalization (I18N) implementation, known as the *i18n* package, for the GO programming language.

For a more detailed description about the package components and internals, go [here](https://github.com/bjusten/go-i18n/tree/main/detailed.md).

## Getting Started

### Clone the repo

```
git clone https://github.com/bjusten/go-i18n.git
```

### Example code

```
package main

import (
	"context"
	"fmt"

   	"github.com/bjusten/go-i18n/pkg/i18n"
)

const defaultLocale = "en-US"

func main() {
	// Create new catalog and initialize it using defaults
	catalog, err := i18n.NewCatalog().Initialize()
	if err != nil {
		panic(fmt.Sprintf("failed to load locales: %v", err))
	}

	// Create new catalog reader, specifying the locale and catalog, returning a new context
	_, ctx := i18n.NewCatalogReader().WithLocale(defaultLocale).WithCatalog(catalog).WithNewContext()

	printKey(ctx, "key-1")
}

func printKey(ctx context.Context, key string) {
	// Get the catalog reader from the context
	catalogReader := i18n.CatalogReaderFromContext(ctx)

	fmt.Printf("Key value: %s\n", catalogReader.Get(key).Value())
}
```

## Examples

Example go programs can be found under the [cmd](https://github.com/bjusten/go-i18n/tree/main/cmd) directory.

## Licensing

This project is under the [MIT License](https://github.com/bjusten/go-i18n/tree/main/LICENSE).