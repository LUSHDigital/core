package test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

// BindJSON takes a json stream and binds it to a struct.
func BindJSON(r io.Reader, target interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}

// ToJSONBody turns a struct into a json body.
func ToJSONBody(tb testing.TB, i interface{}) io.Reader {
	j, err := json.Marshal(i)
	if err != nil {
		tb.Fatalf("error stringifying: %v", err)
	}

	return strings.NewReader(string(j))
}
