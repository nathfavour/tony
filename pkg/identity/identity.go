package identity

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"strings"

	"github.com/nathfavour/tony/pkg/crypto"
	"github.com/nathfavour/tony/pkg/memory"
)

// Manager handles hierarchical deterministic identity derivation.
type Manager struct {
	MasterSeed [32]byte
}

// NewManager creates a new identity manager with the given master seed.
func NewManager(seed [32]byte) *Manager {
	return &Manager{MasterSeed: seed}
}

// DerivePersona generates a specific identity based on a derivation path.
// Example path: m/agent-001/github-persona/commit-signing/0
func (m *Manager) DerivePersona(path string) (*crypto.Identity, error) {
	if !strings.HasPrefix(path, "m/") {
		return nil, fmt.Errorf("invalid derivation path: must start with 'm/'")
	}

	parts := strings.Split(path, "/")
	currentSeed := m.MasterSeed

	for _, part := range parts[1:] {
		nextSeed := deriveChild(currentSeed, part)
		// If we are not at the root, we might want to scrub currentSeed if it was intermediate,
		// but here it's passed by value and currentSeed is just a local copy.
		currentSeed = nextSeed
	}

	id, err := crypto.NewIdentity(currentSeed)
	// Scrub the leaf seed after identity is created
	memory.Scrub(currentSeed[:])
	return id, err
}

// Destroy zeroes out the master seed.
func (m *Manager) Destroy() {
	memory.Scrub(m.MasterSeed[:])
}

// deriveChild computes a child seed from a parent seed and a label.
// It uses HMAC-SHA512 to ensure independent, unlinked keys.
func deriveChild(parent [32]byte, label string) [32]byte {
	mac := hmac.New(sha512.New, parent[:])
	mac.Write([]byte(label))
	digest := mac.Sum(nil)

	var child [32]byte
	copy(child[:], digest[:32])
	return child
}
