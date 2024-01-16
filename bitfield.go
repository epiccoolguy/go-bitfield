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

type BitManipulator interface {
	SetBit(pos uint) error
	ClearBit(pos uint) error
	ToggleBit(pos uint) error
	TestBit(pos uint) (bool, error)
}

// Compile-time check to ensure BitField implements BitManipulator
var _ BitManipulator = &BitField{}

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

// Bytes creates a copy of the underlying bits slice.
// Returns the newly created copy
func (bf *BitField) Bytes() []byte {
	copiedBytes := make([]byte, len(bf.bits))
	copy(copiedBytes, bf.bits)
	return copiedBytes
}

// FromBytes creates a new BitField from a slice of bytes.
// Returns a pointer to the newly created BitField.
func FromBytes(bytes []byte) *BitField {
	bf := New(uint(len(bytes) * 8))
	copy(bf.bits, bytes)
	return bf
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

// InsertUint sets a group of bits starting at a specified offset to the given value.
// The value is interpreted as an LSB-first bit sequence.
// Parameters:
// - offset uint: The starting position for setting the bits, in bits (0-based index).
// - size uint: The number of bits to set.
// - value uint64: The value to set, interpreted as LSB-first. Only the lowest 'size' bits are used.
// Returns an error if the operation goes beyond the bounds of the BitField or if size is invalid.
func (bf *BitField) InsertUint(man BitManipulator, offset, size uint, value uint64) error {
	// Validate that the operation is within bounds
	if offset+size > bf.sz || size > 64 {
		return errors.New("operation out of bounds or size is invalid")
	}

	for i := uint(0); i < size; i++ {
		// Calculate the position of the bit to be modified
		pos := offset + i

		// Determine whether to set or clear the bit based on the corresponding bit in 'value'
		if (value>>i)&1 == 1 {
			// Set the bit
			if err := man.SetBit(pos); err != nil {
				return err
			}
		} else {
			// Clear the bit
			if err := man.ClearBit(pos); err != nil {
				return err
			}
		}
	}

	return nil
}

// ExtractUint retrieves a group of bits starting at a specified offset within the BitField.
// It returns the bits as a uint64 value with LSB-first order.
// Parameters:
// - offset uint: The starting position for retrieving the bits, in bits (0-based index).
// - size uint: The number of bits to retrieve.
// Returns the retrieved bits as a uint64 value and an error if the operation goes beyond the bounds of the BitField.
func (bf *BitField) ExtractUint(man BitManipulator, offset, size uint) (uint64, error) {
	if offset+size > bf.sz {
		return 0, fmt.Errorf("range [%d, %d] out of bounds", offset, offset+size)
	}

	var group uint64
	for i := uint(0); i < size; i++ {
		bit, err := man.TestBit(offset + i)
		if err != nil {
			return 0, err
		}
		if bit {
			group |= (1 << i)
		}
	}
	return group, nil
}
