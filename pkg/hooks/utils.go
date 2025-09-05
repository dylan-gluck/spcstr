package hooks

import "time"

// CurrentTimestamp returns the current time in RFC3339 format
func CurrentTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// AddUniqueString adds a string to a slice if it doesn't already exist
func AddUniqueString(slice []string, str string) []string {
	for _, s := range slice {
		if s == str {
			return slice
		}
	}
	return append(slice, str)
}