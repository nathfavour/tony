package memory

import (
	"testing"
)

func TestLockUnlock(t *testing.T) {
	data := []byte("secret password")
	if err := Lock(data); err != nil {
		t.Logf("Mlock failed (expected if not root or limit reached): %v", err)
	}
	if err := Unlock(data); err != nil {
		t.Errorf("Munlock failed: %v", err)
	}
}

func TestScrub(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	Scrub(data)
	for i, v := range data {
		if v != 0 {
			t.Errorf("byte at index %d was not scrubbed: %v", i, v)
		}
	}
}

func TestScrubString(t *testing.T) {
	s := "volatile secret"
	ScrubString(&s)
	for i := 0; i < len(s); i++ {
		if s[i] != 0 {
			t.Errorf("string char at index %d was not scrubbed: %v", i, s[i])
		}
	}
}
