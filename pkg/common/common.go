package common

import (
	"log"
	"os"
	"path"
)

// GetEnvOrDefault ...
// given an env var return it's value, else return a default
func GetEnvOrDefault(envName string, defaultValue string, required ...bool) string {
	output, found := os.LookupEnv(envName)
	if found {
		return output
	}
	if len(required) > 0 && required[0] && output == "" {
		log.Panicf("error: env '%v' is empty when expected to be set", envName)
	}
	return defaultValue
}

func GetAppPort() (output string) {
	return GetEnvOrDefault("APP_PORT", ":8080")
}

func GetWebFolder() string {
	output := path.Join(GetEnvOrDefault("KO_DATA_PATH", "./kodata"), "web")
	if _, err := os.Stat(output); os.IsNotExist(err) {
		log.Panicf("error: web folder set (%v) does not exist", output)
	}
	return output
}
