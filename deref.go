package oki

// Deref returns the value pointed to by p, or empty if p is nil.
func DEREF[T any](p *T) T {
	var zero T
	if p != nil {
		return *p
	}
	return zero
}

// DEREFWithDefault returns the value pointed to by p, or def if p is nil.
func DEREFWithDefault[T any](p *T, def T) T {
	if p != nil {
		return *p
	}
	return def
}
