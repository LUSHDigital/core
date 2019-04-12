package env_test

import (
	"log"
	"os"
	"testing"

	"github.com/LUSHDigital/core/env"
	"github.com/LUSHDigital/core/test"
)

func ExampleMustGet() {
	var dbURL = env.MustGet("DATABASE_URL")
	log.Println(dbURL)
}

func ExampleTryLoadDefault() {
	env.TryLoadDefault()
}

func TestMustGet(t *testing.T) {
	os.Setenv("TMPENV", "HELLO WORLD")
	tmpenv := env.MustGet("TMPENV")
	test.Equals(t, tmpenv, "HELLO WORLD")
}
