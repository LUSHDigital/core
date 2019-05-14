package auth_test

import (
	"testing"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/test"
	"github.com/gofrs/uuid"
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

func TestConsumer_HasAnyRole(t *testing.T) {
	c := &auth.Consumer{
		Roles: []string{
			"test.foo",
			"test.bar",
			"test.baz",
		},
	}
	t.Run("when using one role that exists", func(t *testing.T) {
		test.Equals(t, true, c.HasAnyRole("test.foo"))
	})
	t.Run("when using two roles where one does not exist", func(t *testing.T) {
		test.Equals(t, true, c.HasAnyRole("test.foo", "doesnot.exist"))
	})
	t.Run("when using one role that does not exist", func(t *testing.T) {
		test.Equals(t, false, c.HasAnyRole("doesnot.exist"))
	})
	t.Run("when using two roles that does not exist", func(t *testing.T) {
		test.Equals(t, false, c.HasAnyRole("doesnot.exist", "has.no.access"))
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

func TestConsumer_HasUUID(t *testing.T) {
	id1 := uuid.Must(uuid.NewV4()).String()
	id2 := uuid.Must(uuid.NewV4()).String()
	c := &auth.Consumer{
		UUID: id1,
	}
	t.Run("when its the same user", func(t *testing.T) {
		test.Equals(t, true, c.HasUUID(id1))
	})
	t.Run("when its not the same user", func(t *testing.T) {
		test.Equals(t, false, c.HasUUID(id2))
	})
}
