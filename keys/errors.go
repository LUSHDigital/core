package keys

import "fmt"

// ErrGetKeySource represents an error when failing to get the source
type ErrGetKeySource struct {
	msg interface{}
}

func (e ErrGetKeySource) Error() string {
	return fmt.Sprintf("failed to get the key source: %v", e.msg)
}

// ErrReadResponse represents an error when failing to read the source data
type ErrReadResponse struct {
	msg interface{}
}

func (e ErrReadResponse) Error() string {
	return fmt.Sprintf("failed to read the key response: %v", e.msg)
}
