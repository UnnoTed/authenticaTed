package util

import (
	"errors"
	"log"
	"os"
	"strings"
)

// MustGetEnv crashes the server when a environment variable is not specified
func MustGetEnv(env string) string {
	data := GetEnv(env, "")

	if data == "" {
		err := errors.New("Error: Environment variables[" + env + "] cannot be empty!")
		log.Fatal(err)
	}

	return data
}

// GetEnv returns a value from a environment variable
// when the variable is not found it returns the
// value from the variable def
func GetEnv(env string, def string) string {
	data := os.Getenv(env)

	if data == "" {
		return def
	}

	logdata := data

	lenv := strings.ToLower(env)
	if strings.Contains(lenv, "token") || strings.Contains(lenv, "password") || strings.Contains(lenv, "key") {
		logdata = "SECRET"
	}

	log.Println("[ENV]: " + env + " = " + logdata)
	return data
}
