package provider

import "strings"

// isNotFoundError checks if an error indicates a resource was not found.
// It performs a case-insensitive check for "not found" in the error message.
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "not found")
}
