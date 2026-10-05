package uuid

import (
	"strings"
	"testing"
)

func TestGenerateRandomBytesNegative(t *testing.T) {
	b, err := GenerateRandomBytes(-1)
	if err == nil {
		t.Fatal("expected error for negative size")
	}
	if b != nil {
		t.Fatalf("expected nil bytes, got %v", b)
	}
	if !strings.Contains(err.Error(), "non-negative") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGenerateRandomBytesZero(t *testing.T) {
	b, err := GenerateRandomBytes(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 0 {
		t.Fatalf("len = %d, want 0", len(b))
	}
}
