package security

import "testing"

func TestPasswordHashRoundTrip(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatal(err)
	}

	if !CheckPassword("secret", hash) {
		t.Fatal("expected password to match")
	}
	if CheckPassword("wrong", hash) {
		t.Fatal("expected wrong password to fail")
	}
}
