package env

import (
	"log"
	"os"
)

// MustGet - Return an environment variable by name.
// If the requested variable does not exist, throw a fatal error.
func MustGet(name string) string {
	envVar, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("environment variable (%s) has not been set", name)
	}
	if envVar == "" {
		log.Fatalf("environment variable (%s) is empty", name)
	}
	return envVar
}
