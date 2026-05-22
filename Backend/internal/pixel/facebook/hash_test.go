package facebook

import "testing"

func TestHashPII(t *testing.T) {
	h := HashPII("Test@Example.com")
	if h == "" {
		t.Fatal("expected hash")
	}
	if HashPII("Test@Example.com") != h {
		t.Fatal("hash must be deterministic")
	}
}
