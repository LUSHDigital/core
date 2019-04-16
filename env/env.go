// Package env provides functionality for ensuring we retrieve an environment
// variable
package env

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
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

// TryLoadDefault will attempt to load the default environment variables.
func TryLoadDefault(paths ...string) {
	paths = append([]string{"infra/.env", "infra/local.env"}, paths...)
	if err := godotenv.Load(paths...); err != nil {
		log.Printf("could not load environment files: %s: skipping...\n", strings.Join(paths, " "))
	}
}
