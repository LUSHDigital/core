package keys

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Source represents one key source
type Source interface {
	GetKey(ctx context.Context) ([]byte, error)
}

// Sources represents multiple key sources
type Sources []Source

// GetKey iterates sources and returns the first successful response
func (sources Sources) GetKey(ctx context.Context) ([]byte, error) {
	for _, source := range sources {
		bts, err := source.GetKey(ctx)
		if err == nil {
			return bts, nil
		}
	}
	return nil, ErrGetKeySource{"no sources could be resolved"}
}

// HTTPSource is a key with a URL source over HTTP
type HTTPSource string

// GetKey retrieves the key from the source over HTTP
func (source HTTPSource) GetKey(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, string(source), nil)
	if err != nil {
		return nil, ErrGetKeySource{}
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
