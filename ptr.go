package oki

// PTR returns a pointer to the given value.
func PTR[T any](v T) *T {
	return &v
}
