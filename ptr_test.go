package oki

import "testing"

func TestPTR(t *testing.T) {
	type sample struct {
		Name string
		Age  int
	}

	original := sample{Name: "Riley", Age: 27}
	ptr := PTR(original)

	if ptr == nil {
		t.Fatal("PTR returned nil pointer")
	}

	if *ptr != original {
		t.Fatalf("PTR(%+v) dereferenced to %+v", original, *ptr)
	}

	ptr.Name = "Jordan"
	if original.Name != "Riley" {
		t.Fatalf("PTR should return a pointer to a copy; original mutated to %q", original.Name)
	}
}
