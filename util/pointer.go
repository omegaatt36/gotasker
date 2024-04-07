package util

// Pointer returns a pointer to a value
func Pointer[T any](value T) *T {
	return &value
}
