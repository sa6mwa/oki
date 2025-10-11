package oki

import "testing"

func TestDEREF(t *testing.T) {
	t.Run("non-nil pointer", func(t *testing.T) {
		val := 42
		got := DEREF(&val)
		if got != val {
			t.Fatalf("DEREF(&%d) = %d, want %d", val, got, val)
		}
	})
	t.Run("nil pointer returns zero value", func(t *testing.T) {
		got := DEREF[int](nil)
		if got != 0 {
			t.Fatalf("DEREF[int](nil) = %d, want 0", got)
		}
	})
}

func TestDEREFWithDefault(t *testing.T) {
	t.Run("non-nil pointer returns value", func(t *testing.T) {
		val := "hello"
		got := DEREFWithDefault(&val, "fallback")
		if got != val {
			t.Fatalf("DEREFWithDefault(&%q, %q) = %q, want %q", val, "fallback", got, val)
		}
	})
	t.Run("nil pointer returns provided default", func(t *testing.T) {
		got := DEREFWithDefault(nil, "fallback")
		if got != "fallback" {
			t.Fatalf("DEREFWithDefault[string](nil, %q) = %q, want %q", "fallback", got, "fallback")
		}
	})
}
