package env_test

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/LUSHDigital/core/env"
)

func ExampleMustGet() {
	var dbURL = env.MustGet("DATABASE_URL")
	log.Println(dbURL)
}

func TestMustGet(t *testing.T) {
	os.Setenv("TMPENV", "HELLO WORLD")
	tmpenv := env.MustGet("TMPENV")
	equals(t, tmpenv, "HELLO WORLD")
}

func equals(tb testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		tb.Fatalf("\n\texp: %#[1]v (%[1]T)\n\tgot: %#[2]v (%[2]T)\n", expected, actual)
	}
}
