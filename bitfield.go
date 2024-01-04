package bitfield

import (
	"errors"
	"fmt"
)

// BitField represents a field of bits and provides methods for manipulating bits.
type BitField struct {
	bits []byte // The underlying storage for the bit field, as a slice of bytes.
	sz   uint   // The size of the bit field in bits.
}

// New creates a new BitField with the specified number of bits.
// Parameters:
// - n uint: The size of the bit field in bits.
// Returns a pointer to the newly created BitField.
func New(n uint) *BitField {
	size := (n + 7) / 8 // Calculate the number of bytes needed to store 'n' bits.
	return &BitField{
		bits: make([]byte, size),
		sz:   n,
	}
}
