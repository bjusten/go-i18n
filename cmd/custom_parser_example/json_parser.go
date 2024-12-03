package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bjusten/go-i18n/pkg/i18n"
)

type JSONParser struct {
	jsonFile string
}

func NewJSONParser(jsonFile string) JSONParser {
	return JSONParser{jsonFile: jsonFile}
}

func (jp JSONParser) Parse(addEntryFunc func(locale string, keyValue i18n.KeyValue)) error {
	b, err := os.ReadFile(jp.jsonFile)
	if err != nil {
		return fmt.Errorf("failed to read file '%s': %v", jp.jsonFile, err)
	}

	var jsonLocales JSONLocales
	if err := json.Unmarshal(b, &jsonLocales); err != nil {
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}

	for locale, keyValues := range jsonLocales.Locales {
		for key, value := range keyValues {
			value.K = key
			addEntryFunc(locale, value)
		}
	}

	return nil
}
