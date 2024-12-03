package i18n_test

import (
	"context"
	"strings"
	"testing"

	"github.com/bjusten/go-i18n/pkg/i18n"
)

func TestReaderExistingContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), testContextKey1, testContextKey1Value)

	catalogReader := i18n.NewCatalogReader()
	_, ctx = catalogReader.WithContext(ctx)

	v := ctx.Value(testContextKey1).(string)
	if v != testContextKey1Value {
		t.Errorf("failed to get proper context key value; expected '%s' but got '%s'", testContextKey1Value, v)
	}

	if catalogReader = i18n.CatalogReaderFromContext(ctx); catalogReader == nil {
		t.Error("could not find catalog reader in context")
	}
}

func TestReaderNewContext(t *testing.T) {
	catalog, _, err := newTestCatalogAndContext()
	if err != nil {
		t.Fatalf("unexpected error creating test catalog and context; %v", err)
	}

	_, ctx := i18n.NewCatalogReader().WithCatalog(catalog).WithNewContext()

	catalogReader := i18n.CatalogReaderFromContext(ctx)
	if catalogReader == nil {
		t.Error("failed to get catalog reader from context")
	}

	v := catalogReader.GetWithLocale("test", testKey)
	switch {
	case strings.HasPrefix(v.Value(), "[unknown key:"):
		t.Error("failed to get value for key")
	case v.Value() != testValue:
		t.Errorf("failed to get proper value for key; expected '%s' but got '%s'", testValue, v.Value())
	}
}

func TestReaderCatalogFromContext(t *testing.T) {
	_, ctx, err := newTestCatalogAndContext()
	if err != nil {
		t.Fatalf("unexpected error creating test catalog and context; %v", err)
	}

	catalogReader := i18n.NewCatalogReader().WithCatalogFromContext(ctx)

	v := catalogReader.GetWithLocale("test", testKey)
	switch {
	case strings.HasPrefix(v.Value(), "[unknown key:"):
		t.Error("failed to get value for key")
	case v.Value() != testValue:
		t.Errorf("failed to get proper value for key; expected '%s' but got '%s'", testValue, v.Value())
	}
}

func TestReaderLocaleFromEnvironment(t *testing.T) {
	t.Setenv(i18n.ReaderLocaleEnvironment, "test")

	catalog, _, err := newTestCatalogAndContext()
	if err != nil {
		t.Fatalf("unexpected error creating test catalog and context; %v", err)
	}

	catalogReader := i18n.NewCatalogReader().WithCatalog(catalog)

	v := catalogReader.Get(testKey)
	if v.Value() != testValue {
		t.Errorf("failed to get proper value for key; expected '%s' but got '%s'", testValue, v)
	}
}

func TestReaderNoCatalog(t *testing.T) {
	catalogReader := i18n.NewCatalogReader()

	v := catalogReader.Get(testKey)
	if !strings.HasPrefix(v.Value(), "[unknown key:") {
		t.Errorf("expected unknown key but got '%s'", v)
	}
}
