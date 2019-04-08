package auth_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/test"
)

var (
	err error
)

func ExampleIssuer_Issue() {
	consumer := &auth.Consumer{
		ID:        999,
		FirstName: "Testy",
		LastName:  "McTest",
		Grants: []string{
			"testing.read",
			"testing.create",
		},
	}
	raw, err := issuer.Issue(consumer)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(raw)
}

func TestIssuer_Issue(t *testing.T) {
	consumer := &auth.Consumer{
		ID:        999,
		FirstName: "Testy",
		LastName:  "McTest",
		Grants: []string{
			"testing.read",
			"testing.create",
		},
	}
	raw, err := issuer.Issue(consumer)
	if err != nil {
		t.Error(err)
	}
	claims, err := parser.Claims(raw)
	if err != nil {
		t.Error(err)
	}
	test.Equals(t, consumer.ID, claims.Consumer.ID)
}
