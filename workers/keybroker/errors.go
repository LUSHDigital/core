package keybroker

import (
	"fmt"
)

// ErrGetKeySource represents an error when failing to get the source
type ErrGetKeySource struct {
	msg interface{}
}

func (e ErrGetKeySource) Error() string {
	return fmt.Sprintf("failed to retrieve the key from source: %v", e.msg)
}

// ErrReadResponse represents an error when failing to read the source data
type ErrReadResponse struct {
	msg interface{}
}

func (e ErrReadResponse) Error() string {
	return fmt.Sprintf("failed to read the key response: %v", e.msg)
}

// ErrNoSourcesResolved represents an error for when no sources could be resolved at all
type ErrNoSourcesResolved struct {
	N int
}

func (e ErrNoSourcesResolved) Error() string {
	return fmt.Sprintf("no sources could be resolved: %d sources", e.N)
}

var (
	// ErrEmptyURL represents an error for when an expected url is an empty string
	ErrEmptyURL = ErrGetKeySource{"url cannot be empty"}

	// ErrEmptyFilePath represents an error for when an expected file path is an empty string
	ErrEmptyFilePath = ErrGetKeySource{"file path cannot be empty"}

	// ErrEmptyString represents an error for when an expected string should contain a public key
	ErrEmptyString = ErrGetKeySource{"string cannot be empty"}
)
