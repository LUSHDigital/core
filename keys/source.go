package keys

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Source represents one or a chain of sources
type Source interface {
	Get(ctx context.Context) ([]byte, error)
}

// Sources defines a chain of sources
type Sources []Source

// Get iterates sources and returns the first successfully resolved
func (sources Sources) Get(ctx context.Context) ([]byte, error) {
	for _, source := range sources {
		bts, err := source.Get(ctx)
		if err == nil {
			return bts, nil
		}
	}
	return nil, ErrGetKeySource{"no sources could be resolved"}
}

// HTTPSource defines a source with a URL to resolve over HTTP
type HTTPSource string

// Get retrieves data from the URL over HTTP
func (source HTTPSource) Get(ctx context.Context) ([]byte, error) {
	if source == "" {
		return nil, ErrGetKeySource{"url cannot be empty"}
	}
	req, err := http.NewRequest(http.MethodGet, string(source), nil)
	if err != nil {
		return nil, ErrGetKeySource{err}
	}
	req = req.WithContext(ctx)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, ErrGetKeySource{err}
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, ErrGetKeySource{fmt.Sprintf("got status code %d got %d", http.StatusOK, res.StatusCode)}
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, ErrReadResponse{err}
	}
	return body, nil
}

// FileSource defines a path to a file on disk
type FileSource string

// Get retrieves data from the path to a file on disk
func (source FileSource) Get(ctx context.Context) ([]byte, error) {
	if source == "" {
		return nil, ErrGetKeySource{"file path cannot be empty"}
	}
	f, err := os.Open(string(source))
	if err != nil {
		return nil, ErrGetKeySource{err}
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, ErrReadResponse{err}
	}
	return content, nil
}

// StringSource defines the source as a string
type StringSource string

// Get converts the string to a byte slice
func (source StringSource) Get(ctx context.Context) ([]byte, error) {
	return []byte(source), nil
}
