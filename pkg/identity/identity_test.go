package identity

import (
	"bytes"
	"testing"
)

func TestDerivation(t *testing.T) {
	masterSeed := [32]byte{0xAA, 0xBB, 0xCC}
	manager := NewManager(masterSeed)

	id1, err := manager.DerivePersona("m/agent-1/task-1")
	if err != nil {
		t.Fatalf("failed to derive id1: %v", err)
	}

	id2, err := manager.DerivePersona("m/agent-1/task-1")
	if err != nil {
		t.Fatalf("failed to derive id2: %v", err)
	}

	if !bytes.Equal(id1.Seed[:], id2.Seed[:]) {
		t.Errorf("deterministic derivation failed: seeds do not match for same path")
	}

	id3, err := manager.DerivePersona("m/agent-1/task-2")
	if err != nil {
		t.Fatalf("failed to derive id3: %v", err)
	}

	if bytes.Equal(id1.Seed[:], id3.Seed[:]) {
		t.Errorf("unlinked derivation failed: seeds match for different paths")
	}
}
