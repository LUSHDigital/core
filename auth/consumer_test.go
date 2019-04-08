package auth_test

import (
	"testing"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/test"
)

func TestConsumer_HasAnyGrant(t *testing.T) {
	c := &auth.Consumer{
		Grants: []string{
			"test.foo",
			"test.bar",
			"test.baz",
		},
	}
	t.Run("when using one grant that exists", func(t *testing.T) {
		test.Equals(t, true, c.HasAnyGrant("test.foo"))
	})
	t.Run("when using two grants where one does not exist", func(t *testing.T) {
		test.Equals(t, true, c.HasAnyGrant("test.foo", "doesnot.exist"))
	})
	t.Run("when using one grant that does not exist", func(t *testing.T) {
		test.Equals(t, false, c.HasAnyGrant("doesnot.exist"))
	})
	t.Run("when using two grants that does not exist", func(t *testing.T) {
		test.Equals(t, false, c.HasAnyGrant("doesnot.exist", "has.no.access"))
	})
}

func TestConsumer_IsUser(t *testing.T) {
	c := &auth.Consumer{
		ID: 1,
	}
	t.Run("when its the same user", func(t *testing.T) {
		test.Equals(t, true, c.IsUser(1))
	})
	t.Run("when its not the same user", func(t *testing.T) {
		test.Equals(t, false, c.IsUser(2))
	})
}
