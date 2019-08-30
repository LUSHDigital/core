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

func ExampleTryOverloadDefault() {
	env.TryOverloadDefault()
}
func TestMustGet(t *testing.T) {
	os.Setenv("TMPENV", "HELLO WORLD")
	tmpenv := env.MustGet("TMPENV")
	test.Equals(t, tmpenv, "HELLO WORLD")
}

func TestTryLoadDefault(t *testing.T) {
	os.Setenv("TMPENV", "HELLO WORLD")
	env.TryLoadDefault("testdata/one.env", "testdata/two.env")
	test.Equals(t, "HELLO WORLD", os.Getenv("TMPENV"))
	test.Equals(t, "ONE", os.Getenv("TMPONE"))
}

func TestTryOverloadDefault(t *testing.T) {
	os.Setenv("TMPENV", "HELLO WORLD")
	env.TryOverloadDefault("testdata/one.env", "testdata/two.env")
	test.Equals(t, "TWO", os.Getenv("TMPENV"))
	test.Equals(t, "TWO", os.Getenv("TMPONE"))
}
