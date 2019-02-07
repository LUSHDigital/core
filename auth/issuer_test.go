package auth_test

import (
	"testing"

	"github.com/LUSHDigital/microservice-core-golang/auth"
)

var (
	err error
)

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
	equals(t, consumer.ID, claims.Consumer.ID)
}
