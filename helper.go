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
	envVar := os.Getenv(name)
	if envVar == "" {
		log.Fatalf("Environment variable (%s) has not been set.", name)
	}

	return envVar
}
