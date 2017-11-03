package microservicecore

import (
	"log"
	"os"
)

// GetEnvOrFail - Get an environment variable by name.
//
// Or if the requested variable does not exist, throw a fatal error.
//
// Params:
//     name string - The name of the environment variable to get.
//
// Return:
//     string - The value of the requested environment variable.
func GetEnvOrFail(name string) string {
	envVar, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("environment variable (%s) has not been set", name)
	}
	if envVar == "" {
		log.Fatalf("environment variable (%s) is empty", name)
	}
	return envVar
}
