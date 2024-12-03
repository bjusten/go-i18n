package i18n

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type KeyPairFSParser struct {
	directories []string
}

// NewKeyPairFSParser returns a new KeyPairFSParser that will process all of the specified directories when parsing
func NewKeyPairFSParser(directories []string) KeyPairFSParser {
	return KeyPairFSParser{directories: directories}
}

func (p KeyPairFSParser) Parse(addEntryFunc func(locale string, keyValue KeyValue)) error {
	for _, directory := range p.directories {
		if err := p.FromDirectory(addEntryFunc, directory); err != nil {
			return err
		}
	}

	return nil
}

var errInvalidScanner = errors.New("scanner is nil")

// FromScanner will attempt to read keyValues from the specified BufIO scanner for the specified locale
func (p KeyPairFSParser) FromScanner(addEntryFunc func(locale string, keyValue KeyValue), locale string, scanner *bufio.Scanner) error {
	if scanner == nil {
		return errInvalidScanner
	}

	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "#") || strings.HasPrefix(text, "//") {
			continue
		}

		kp, err := NewKeyPairFromString(text)
		if err != nil {
			return fmt.Errorf("failed to create keypair from string: %v", err)
		}

		addEntryFunc(locale, kp)
	}

	return nil
}

// FromFile will attempt to load keyValues from a file located at the specified path
func (p KeyPairFSParser) FromFile(addEntryFunc func(locale string, keyValue KeyValue), path string) error {
	locale, err := extractLocaleFromPath(path)
	if err != nil {
		return fmt.Errorf("failed to extract locale from path '%s': %w", path, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open '%s': %w", path, err)
	}
	defer f.Close()

	return p.FromScanner(addEntryFunc, locale, bufio.NewScanner(f))
}

// FromDirectory will attempt to load keyValues from all files located in the specified directory
func (p KeyPairFSParser) FromDirectory(addEntryFunc func(locale string, keyValue KeyValue), directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory '%s': %w", directory, err)
	}

	for _, entry := range entries {
		path := directory + "/" + entry.Name()
		switch {
		case entry.IsDir():
			if err := p.FromDirectory(addEntryFunc, path); err != nil {
				return fmt.Errorf("failed to load catalog from directory '%s': %w", path, err)
			}
		default:
			if err := p.FromFile(addEntryFunc, path); err != nil {
				return fmt.Errorf("failed to load catalog from file '%s': %w", path, err)
			}
		}
	}

	return nil
}

var _extractRegex = regexp.MustCompile(`.*[\/\\]([^\/\\]+)[\/\\]`)

func extractLocaleFromPath(path string) (string, error) {
	matches := _extractRegex.FindStringSubmatch(path)
	switch {
	case len(matches) != 2:
		return "", fmt.Errorf("expected 2 matches; found %d", len(matches))
	default:
		return matches[1], nil
	}
}
