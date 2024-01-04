package bitfield

import (
	"errors"
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

// calculateBitPosition calculates and returns the byte index and bit position within the byte
// for a specified bit position in the BitField.
// Parameters:
// - pos uint: The position of the bit, in bits (0-based index).
// Returns the byte index (int), the bit position within that byte (uint), and an error if the
// given position is out of the range of the BitField (i.e., if it exceeds the size of the BitField).
func (bf *BitField) calculateBitPosition(pos uint) (byteIndex int, bitPosition uint, err error) {
	if pos >= bf.sz {
		return 0, 0, errors.New("bit position out of range")
	}
	byteIndex = int(pos / 8) // Calculate the byte index.
	bitPosition = pos % 8    // Calculate the bit position within the byte.
	return byteIndex, bitPosition, nil
}

// SetBit sets a bit at a specified position within the BitField to 1.
// Parameters:
// - pos uint: The position of the bit to set, in bits (0-based index).
// Returns an error if the position is out of the range of the BitField.
func (bf *BitField) SetBit(pos uint) error {
	byteIndex, bitPosition, err := bf.calculateBitPosition(pos)
	if err != nil {
		return err
	}
	bf.bits[byteIndex] |= (1 << bitPosition) // Set the bit.
	return nil
}

// ClearBit clears a bit at a specified position within the BitField to 0.
// Parameters:
// - pos uint: The position of the bit to clear, in bits (0-based index).
// Returns an error if the position is out of the range of the BitField.
func (bf *BitField) ClearBit(pos uint) error {
	byteIndex, bitPosition, err := bf.calculateBitPosition(pos)
	if err != nil {
		return err
	}
	bf.bits[byteIndex] &= ^(1 << bitPosition) // Clear the bit.
	return nil
}

// ToggleBit toggles a bit at a specified position within the BitField.
// Parameters:
// - pos uint: The position of the bit to toggle, in bits (0-based index).
// Returns an error if the position is out of the range of the BitField.
func (bf *BitField) ToggleBit(pos uint) error {
	byteIndex, bitPosition, err := bf.calculateBitPosition(pos)
	if err != nil {
		return err
	}
	bf.bits[byteIndex] ^= (1 << bitPosition) // Toggle the bit.
	return nil
}

// TestBit retrieves the value of a bit at a specified position within the BitField.
// Parameters:
// - pos uint: The position of the bit to retrieve, in bits (0-based index).
// Returns the value of the bit (true for 1, false for 0) and an error if the position is out of range.
func (bf *BitField) TestBit(pos uint) (bool, error) {
	byteIndex, bitPosition, err := bf.calculateBitPosition(pos)
	if err != nil {
		return false, err
	}

	// Retrieve the value of the bit by checking if the bit at the bitPosition is set in the byte.
	return bf.bits[byteIndex]&(1<<bitPosition) != 0, nil
}
