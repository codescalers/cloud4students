// Package internal for internal details
package internal

// Contains check if a slice contains an element
func Contains[T comparable](elements []T, element T) bool {
	for _, e := range elements {
		if element == e {
			return true
		}
	}
	return false
}
