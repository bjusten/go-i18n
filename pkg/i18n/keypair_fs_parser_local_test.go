package i18n

import "testing"

const (
	missingLocalePath = `/file.ext`
	notAPath          = `not-a-valid-path`
	unixPath          = `/example/path/locale/file.ext`
	windowsPath       = `C:\example\path\locale\file.ext`
)

func TestExtractLocaleFromPath(t *testing.T) {
	locale, err := extractLocaleFromPath(unixPath)
	switch {
	case err != nil:
		t.Errorf("failed to regex match unix path: %v", err)
	case locale != "locale":
		t.Errorf("failed to match locale from unix path: got '%s'", locale)
	}

	locale, err = extractLocaleFromPath(windowsPath)
	switch {
	case err != nil:
		t.Errorf("failed to regex match windows path: %v", err)
	case locale != "locale":
		t.Errorf("failed to match locale from windows path: got '%s'", locale)
	}
}

func TestFailedExtractLocaleFromP(t *testing.T) {
	_, err := extractLocaleFromPath(missingLocalePath)
	if err == nil {
		t.Error("expected error from missing locale path")
	}

	_, err = extractLocaleFromPath(notAPath)
	if err == nil {
		t.Error("expected error from invalid path")
	}
}
