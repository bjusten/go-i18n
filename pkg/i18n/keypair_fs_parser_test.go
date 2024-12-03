package i18n_test

import (
	"strings"
	"testing"

	"github.com/bjusten/go-i18n/pkg/i18n"
)

const (
	testInvalidDirectory = "invalid-directory"

	testDataFilteredLocale = "locale1"
	testDataFilteredKey    = "locale-1-key-1"
	testDataFilteredValue  = "Test 1"
)

func TestFromDirectory(t *testing.T) {
	parser := i18n.NewKeyPairFSParser([]string{})
	catalog := i18n.NewCatalog().WithParser(parser)

	if err := parser.FromDirectory(catalog.AddKeyValue, testInvalidDirectory); err == nil {
		t.Error("expected error from FromDirectory but found none")
	}

	if err := parser.FromDirectory(catalog.AddKeyValue, testDirectory); err != nil {
		t.Errorf("failed to load test data; %v", err)
	}

	kp := catalog.Get(testDataLocale, testDataKey)
	switch {
	case kp.Key() != testDataKey:
		t.Errorf("failed to get proper key; expected '%s' but got '%s'", testDataKey, kp.Key())
	case strings.HasPrefix(kp.Value(), "[unknown key:"):
		t.Error("failed to get value for key")
	case kp.Value() != testDataValue:
		t.Errorf("failed to get proper value for key; expected '%s' but got '%s'", testDataValue, kp.Value())
	}

	kp = catalog.Get(testDataFilteredLocale, testDataFilteredKey)
	switch {
	case kp.Key() != testDataFilteredKey:
		t.Errorf("failed to get proper key; expected '%s' but got '%s'", testDataFilteredKey, kp.Key())
	case strings.HasPrefix(kp.Value(), "[unknown key:"):
		t.Error("failed to get value for key")
	case kp.Value() != testDataFilteredValue:
		t.Errorf("failed to get proper value for key; expected '%s' but got '%s'", testDataFilteredValue, kp.Value())
	}
}

func TestFilteredFromDirectory(t *testing.T) {
	parser := i18n.NewKeyPairFSParser([]string{"./locales"})
	catalog := i18n.NewCatalog().WithLocales(testDataLocale)

	if err := parser.FromDirectory(catalog.AddKeyValue, testInvalidDirectory); err == nil {
		t.Error("expected error from FromDirectory but found none")
	}

	if err := parser.FromDirectory(catalog.AddKeyValue, testDirectory); err != nil {
		t.Errorf("failed to load test data; %v", err)
	}

	kp := catalog.Get(testDataLocale, testDataKey)
	switch {
	case kp.Key() != testDataKey:
		t.Errorf("failed to get proper key; expected '%s' but got '%s'", testDataKey, kp.Key())
	case strings.HasPrefix(kp.Value(), "[unknown key:"):
		t.Error("failed to get value for key")
	case kp.Value() != testDataValue:
		t.Errorf("failed to get proper value for key; expected '%s' but got '%s'", testDataValue, kp.Value())
	}

	kp = catalog.Get(testDataFilteredLocale, testDataFilteredKey)
	switch {
	case kp.Key() != testDataFilteredKey:
		t.Errorf("failed to get proper key; expected '%s' but got '%s'", testDataFilteredKey, kp.Key())
	case !strings.HasPrefix(kp.Value(), "[unknown key:"):
		t.Errorf("found key when it should be filtered, got '%s'", kp.Value())
	}
}
