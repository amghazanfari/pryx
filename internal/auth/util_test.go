package auth

import "testing"

func TestGenerateAndHash(t *testing.T) {
	plain, prefix, hash, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("GenerateAPIKey: %v", err)
	}
	if len(plain) < 20 || prefix == "" || hash == "" {
		t.Fatalf("bad key parts: plain=%q prefix=%q hash-len=%d", plain, prefix, len(hash))
	}
	if got := HashAPIKey(plain); got != hash {
		t.Fatalf("hash mismatch: %s vs %s", got, hash)
	}
	if prefix != plain[:8] {
		t.Fatalf("prefix mismatch: %q vs %q", prefix, plain[:8])
	}
}

func TestHasScope(t *testing.T) {
	if !HasScope("a,b , c", "b") {
		t.Fatal("expected to find scope b")
	}
	if HasScope("a, b, c", "d") {
		t.Fatal("did not expect to find scope d")
	}
}
