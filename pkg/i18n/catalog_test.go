package i18n_test

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/bjusten/go-i18n/pkg/i18n"
)

type testContextKey string

const (
	testDirectory = "./test_data"

	testDataLocale = "locale2"
	testDataKey    = "locale-2-key-1"
	testDataValue  = "Test 3"

	testContextKey1      = testContextKey("test1")
	testContextKey1Value = "test1value"

	testKey     = "test-key"
	testValue   = "test-value"
	testKeyPair = testKey + "=" + testValue

	testLocales = 2
	testKeys    = 4
)

func newTestCatalogFromTestData() (*i18n.Catalog, error) {
	parser := i18n.NewKeyPairFSParser([]string{"./test_data"})
	catalog, err := i18n.NewCatalog().WithParser(parser).Initialize()
	if err != nil {
		return nil, fmt.Errorf("failed to load locales from test directory: %w", err)
	}

	return catalog, nil
}

func newTestCatalogAndContext() (*i18n.Catalog, context.Context, error) {
	kp, err := i18n.NewKeyPairFromString(testKeyPair)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create key pair from string: %w", err)
	}

	catalog, ctx := i18n.NewCatalog().WithNewContext()
	catalog.AddKeyValue("test", kp)

	return catalog, ctx, nil
}

func TestStats(t *testing.T) {
	catalog, err := newTestCatalogFromTestData()
	if err != nil {
		t.Errorf("failed to load test catalog; %v", err)
	}

	stats := catalog.Stats()
	if stats.Locales != testLocales {
		t.Errorf("expected %d locales but found %d", testLocales, stats.Locales)
	}

	if stats.Keys != testKeys {
		t.Errorf("expected %d keys but found %d", testKeys, stats.Keys)
	}
}

func TestExistingContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), testContextKey1, testContextKey1Value)

	catalog := i18n.NewCatalog()
	_, ctx = catalog.WithContext(ctx)

	v := ctx.Value(testContextKey1).(string)
	if v != testContextKey1Value {
		t.Errorf("failed to get proper context key value; expected '%s' but got '%s'", testContextKey1Value, v)
	}

	if catalog = i18n.CatalogFromContext(ctx); catalog == nil {
		t.Error("could not find catalog in context")
	}
}

func TestNewContext(t *testing.T) {
	_, ctx, err := newTestCatalogAndContext()
	if err != nil {
		t.Fatalf("unexpected error creating test catalog and context; %v", err)
	}

	catalog := ctx.Value(i18n.CatalogContextKey).(*i18n.Catalog)
	if catalog == nil {
		t.Error("failed to get catalog from context")
	}

	v := catalog.Get("test", testKey)
	switch {
	case strings.HasPrefix(v.Value(), "[unknown key:"):
		t.Error("failed to get value for key")
	case v.Value() != testValue:
		t.Errorf("failed to get proper value for key; expected '%s' but got '%s'", testValue, v.Value())
	}
}

func TestInitializeAndNewContext(t *testing.T) {
	parser := i18n.NewKeyPairFSParser([]string{"./test_data"})
	ctx, err := i18n.NewCatalog().WithParser(parser).InitializeWithContext()
	if err != nil {
		t.Fatalf("unexpected error creating test catalog and context; %v", err)
	}

	catalog := i18n.CatalogFromContext(ctx)

	kp := catalog.Get(testDataLocale, testDataKey)
	switch {
	case kp.Key() != testDataKey:
		t.Errorf("failed to get proper key; expected '%s' but got '%s'", testDataKey, kp.Key())
	case strings.HasPrefix(kp.Value(), "[unknown key:"):
		t.Error("failed to get value for key")
	case kp.Value() != testDataValue:
		t.Errorf("failed to get proper value for key; expected '%s' but got '%s'", testDataValue, kp.Value())
	}
}

func TestDefaultLocale(t *testing.T) {
	catalog, _, err := newTestCatalogAndContext()
	if err != nil {
		t.Fatalf("unexpected error creating test catalog and context; %v", err)
	}

	v := catalog.Get("test", testKey)
	switch {
	case strings.HasPrefix(v.Value(), "[unknown key:"):
		t.Error("failed to get value for key")
	case v.Value() != testValue:
		t.Errorf("failed to get proper value for key; expected '%s' but got '%s'", testValue, v.Value())
	}

	v = catalog.Get("invalid", testKey)
	if !strings.HasPrefix(v.Value(), "[unknown key:") {
		t.Errorf("expected key to be unknown but got '%s'", v.Value())
	}

	catalog = catalog.WithDefaultLocale("test")

	v = catalog.Get("invalid", testKey)
	if v.Value() != testValue {
		t.Errorf("failed to get proper value for key; expected '%s' but got '%s'", testValue, v.Value())
	}
}

func TestDefaultLocaleFromEnvironment(t *testing.T) {
	t.Setenv(i18n.DefaultLocaleEnvironment, "test")

	catalog, _, err := newTestCatalogAndContext()
	if err != nil {
		t.Fatalf("unexpected error creating test catalog and context; %v", err)
	}

	v := catalog.Get("invalid", testKey)
	if v.Value() != testValue {
		t.Errorf("failed to get proper value for key; expected '%s' but got '%s'", testValue, v.Value())
	}
}

func TestLocales(t *testing.T) {
	catalog := i18n.NewCatalog().WithLocales([]string{"locale-1", "locale-2"}...)
	locales := catalog.Locales()
	sort.Strings(locales)

	switch {
	case len(locales) != 2:
		t.Errorf("expected 2 locales, but got %d", len(locales))
	case locales[0] != "locale-1":
		t.Errorf("expected 'locale-1', but got '%s'", locales[0])
	case locales[1] != "locale-2":
		t.Errorf("expected 'locale-2', but got '%s'", locales[1])
	}
}

func TestLocalesFromEnvironment(t *testing.T) {
	t.Setenv(i18n.FilterLocalesEnvironment, "locale-1,locale-2")

	catalog := i18n.NewCatalog()
	locales := catalog.Locales()
	sort.Strings(locales)

	switch {
	case len(locales) != 2:
		t.Errorf("expected 2 locales, but got %d", len(locales))
	case locales[0] != "locale-1":
		t.Errorf("expected 'locale-1', but got '%s'", locales[0])
	case locales[1] != "locale-2":
		t.Errorf("expected 'locale-2', but got '%s'", locales[1])
	}
}

func TestUnknownKey(t *testing.T) {
	catalog := i18n.NewCatalog()

	v := catalog.Get("invalid", "invalid")
	switch {
	case v.Key() != "invalid":
		t.Errorf("failed to get proper key; expected '%s' but got '%s'", "invalid", v.Key())
	case !strings.HasPrefix(v.Value(), "[unknown key:"):
		t.Errorf("expected unknown key, but got '%s'", v.Value())
	}
}
