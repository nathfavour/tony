package crypto

import (
	"bytes"
	"testing"
)

func TestNewIdentity(t *testing.T) {
	seed := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	id, err := NewIdentity(seed)
	if err != nil {
		t.Fatalf("failed to create identity: %v", err)
	}

	if !bytes.Equal(id.Seed[:], seed[:]) {
		t.Errorf("seed mismatch")
	}

	if len(id.EdPublic) != 32 {
		t.Errorf("invalid Ed25519 public key length: %d", len(id.EdPublic))
	}
}

func TestSignVerify(t *testing.T) {
	seed := [32]byte{42}
	id, _ := NewIdentity(seed)
	msg := []byte("hello world")
	sig := id.Sign(msg)

	if !Verify(id.EdPublic, msg, sig) {
		t.Errorf("signature verification failed")
	}
}

func TestSharedSecret(t *testing.T) {
	seed1 := [32]byte{1}
	id1, _ := NewIdentity(seed1)

	seed2 := [32]byte{2}
	id2, _ := NewIdentity(seed2)

	secret1, err := id1.SharedSecret(id2.X255Public)
	if err != nil {
		t.Fatalf("id1 failed shared secret: %v", err)
	}

	secret2, err := id2.SharedSecret(id1.X255Public)
	if err != nil {
		t.Fatalf("id2 failed shared secret: %v", err)
	}

	if !bytes.Equal(secret1[:], secret2[:]) {
		t.Errorf("shared secrets do not match")
	}
}
