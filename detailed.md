# i18n Package Details

This resource contains general information about the i18n go package.

## Components

The i18n go package contains the catalog and a catalog reader.

### Catalog

The catalog is the main object that handles storing and retrieving KeyValue objects.

It is primarily used as the backend in conjunction with catalog readers, but can be fully functional as a stand alone object.

### Catalog reader

A catalog reader retrieves KeyValue objects from the catalog.

The reader is meant to be a lightweight object stored and utilized by an end user with their own perferred locale, in applications such as apis or web UIs.

## Sub-Components

The i18n go package also contains the KeyValue interface, the Parser interface, and the keypair file system parser.

### KeyValue

KeyValue is a generic interface used by the catalog in order to support custom objects that may contain additional metadata.

### Parser

Parser is a generic interface used by the catalog in order to support custom data sources.

### KeyPair File System Parser

The KeyPair File System Parser is an implementation of the Parser interface.

It reads the configured directories and attempts to load all files it finds within.
It reads each file line-by-line and parses each line assuming the 'key=value' format.

An example english locale file, such as `./locales/en/entries.i18n`:

    key-1=value-1
    key-2=value-2

## Defaults

- Catalog uses the KeyPair File System Parser by default and its default directory for loading data from is './locales'.
- Catalog does not filter locales by default.
- Catalog does not have a default locale.

