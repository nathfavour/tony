package memory

import (
	"reflect"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Lock pins the memory of the given byte slice to RAM, preventing it from being swapped.
func Lock(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	return unix.Mlock(b)
}

// Unlock releases the memory lock on the given byte slice.
func Unlock(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	return unix.Munlock(b)
}

// Scrub overwrites the given byte slice with zeros.
func Scrub(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// ScrubString attempts to zero out the memory of a string.
// WARNING: Strings in Go are immutable; this uses unsafe pointers to bypass that restriction.
// Use with extreme caution and only on strings that are guaranteed not to be in read-only memory.
func ScrubString(s *string) {
	if s == nil || *s == "" {
		return
	}
	hdr := (*reflect.StringHeader)(unsafe.Pointer(s))
	b := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
	Scrub(b)
}
