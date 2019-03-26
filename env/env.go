// Package env provides functionality for ensuring we retrieve an environment variable
package env

import (
	"log"
	"os"
)

// MustGet returns an environment variable by name
// If the requested variable does not exist, we throw a fatal error
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
