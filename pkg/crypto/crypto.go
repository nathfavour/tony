package crypto

import (
	"crypto/ed25519"
	"crypto/sha512"
	"errors"

	"github.com/nathfavour/tony/pkg/memory"
	"golang.org/x/crypto/curve25519"
)

// Identity represents a unified Curve25519 identity.
type Identity struct {
	Seed       [32]byte
	EdPublic   ed25519.PublicKey
	EdPrivate  ed25519.PrivateKey
	X255Public [32]byte
	X255Secret [32]byte
}

// NewIdentity creates a new identity from a 32-byte seed.
func NewIdentity(seed [32]byte) (*Identity, error) {
	edPriv := ed25519.NewKeyFromSeed(seed[:])
	edPub := edPriv.Public().(ed25519.PublicKey)

	// Convert Ed25519 private key to X25519 secret key.
	// The Ed25519 private key is the SHA-512 hash of the seed.
	// The first 32 bytes of this hash are used as the X25519 secret key after clamping.
	h := sha512.Sum512(seed[:])
	x255Secret := [32]byte{}
	copy(x255Secret[:], h[:32])
	
	// Scrub temporary hash
	memory.Scrub(h[:])
	
	x255Public, err := curve25519.X25519(x255Secret[:], curve25519.Basepoint)
	if err != nil {
		return nil, err
	}

	id := &Identity{
		Seed:       seed,
		EdPublic:   edPub,
		EdPrivate:  edPriv,
		X255Secret: x255Secret,
	}
	copy(id.X255Public[:], x255Public)

	return id, nil
}

// Destroy zeroes out the sensitive components of the identity.
func (id *Identity) Destroy() {
	memory.Scrub(id.Seed[:])
	memory.Scrub(id.EdPrivate)
	memory.Scrub(id.X255Secret[:])
}

// Sign signs a message using Ed25519.
func (id *Identity) Sign(message []byte) []byte {
	return ed25519.Sign(id.EdPrivate, message)
}

// Verify verifies an Ed25519 signature.
func Verify(pub ed25519.PublicKey, message, sig []byte) bool {
	return ed25519.Verify(pub, message, sig)
}

// Seal encrypts a message for a peer using X25519 (Diffie-Hellman + shared secret logic would go here).
// This is a placeholder for the X25519 key exchange.
func (id *Identity) SharedSecret(peerPublic [32]byte) ([32]byte, error) {
	secret, err := curve25519.X25519(id.X255Secret[:], peerPublic[:])
	if err != nil {
		return [32]byte{}, err
	}
	var res [32]byte
	copy(res[:], secret)
	return res, nil
}

var ErrInvalidSeed = errors.New("invalid seed length")
