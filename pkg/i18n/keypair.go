package i18n

import (
	"fmt"
	"strings"
)

type KeyPair struct {
	key   string
	value string
}

const UnknownKeyFormat = "[unknown key:%s]"

// NewKeyPair returns a new KeyPair object based on the specified key and value
func NewKeyPair(key string, value string) KeyPair {
	return KeyPair{key: key, value: value}
}

// Key returns the key associated with this KeyPair object
func (kp KeyPair) Key() string {
	return kp.key
}

// Value returns the value associated with this KeyPair object
func (kp KeyPair) Value() string {
	return kp.value
}

// NewKeyPairFromString returns a KeyPair object based on the specified keypair string (using the 'key=value' format)
func NewKeyPairFromString(keypair string) (KeyPair, error) {
	tokens := strings.Split(keypair, "=")
	switch {
	case len(tokens) != 2:
		return KeyPair{}, fmt.Errorf("failed to parse keypair from string; found %d token(s)", len(tokens))
	default:
		return NewKeyPair(tokens[0], tokens[1]), nil
	}
}

// NewUnknownKeyPair returns a KeyPair object that represents a non-existant key
func NewUnknownKeyPair(key string) KeyPair {
	return NewKeyPair(key, fmt.Sprintf(UnknownKeyFormat, key))
}
