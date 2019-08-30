// Package env provides functionality for ensuring we retrieve an environment
// variable
package env

import (
	"log"
	"os"

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
// This WILL NOT OVERRIDE an env variable that already exists.
// Those set prior to this call will remain. For env variables set in multiple
// files passed to this function, the FIRST one will prevail.
func TryLoadDefault(paths ...string) {
	paths = append([]string{"infra/.env", "infra/local.env"}, paths...)
	for _, p := range paths {
		if err := godotenv.Load(p); err != nil {
			log.Printf("could not load environment file: %s: skipping...\n", p)
		}
	}
}

// TryOverloadDefault will attempt to load the default environment variables.
// This WILL OVERRIDE an env variable that already exists.
// Those set prior to this call will be replaced if set here. For env variables
// set in multiple files passed to this function, the LAST one will prevail.
func TryOverloadDefault(paths ...string) {
	paths = append([]string{"infra/.env", "infra/local.env"}, paths...)
	for _, p := range paths {
		if err := godotenv.Overload(p); err != nil {
			log.Printf("could not load environment file: %s: skipping...\n", p)
		}
	}
}
