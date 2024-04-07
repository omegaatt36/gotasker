package util

// Pointer
func Pointer[T any](value T) *T {
	return &value
}
