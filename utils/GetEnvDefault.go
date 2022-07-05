package utils

import (
	"log"
	"os"
	"strconv"
)

// Returns given defaultvalue if env var is nil
func GetEnvDefault(var_name string, defaultValue string) string {
	// Get Variable
	value := os.Getenv(var_name)
	if value == "" {
		value = defaultValue
	}

	// Return Parse Integer
	return value
}

// Returns default if integer is invalid or nil
func GetEnvDefaultInt(var_name string, defaultValue int) int {
	// Get Variable
	value := os.Getenv(var_name)
	if value == "" {
		return defaultValue
	}

	// Parse Integer
	parsedInt, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Failed to parse value for Variable \"%s\": %s", var_name, err.Error())
		return defaultValue
	}

	// Return Parsed Integer
	return parsedInt
}
