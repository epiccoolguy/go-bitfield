package bitfield

import (
	"errors"
)

type bigEndian struct{}

// BigEndian is a BitManipulator implementation that follows the Most Significant Byte First order
// with Most Significant Bit 0 numbering.
var BigEndian BitManipulator = &bigEndian{}

// New creates a new BitField with the specified number of bits.
// It calculates the number of bytes needed to store 'n' bits and initializes the BitField
// with an empty byte slice of that size.
func (bm *bigEndian) New(n uint64) *BitField {
	byteSize := (n + 7) / 8

	return &BitField{
		data:        make([]byte, byteSize),
		size:        n,
		manipulator: BigEndian,
	}
}

// FromBytes creates a new BitField from a byte slice.
// It takes the byte slice and the size of the BitField in bits as parameters
// and returns a pointer to the created BitField.
func (bm *bigEndian) FromBytes(bytes []byte, size uint64) *BitField {
	data := make([]byte, len(bytes))
	copy(data, bytes)

	return &BitField{
		data:        data,
		size:        size,
		manipulator: BigEndian,
	}
}

func calcBitPosBE(bf *BitField, pos uint64) (bytePos, bitPos uint64, err error) {
	if pos >= bf.size {
		return 0, 0, errors.New("bit position out of range")
	}
	bytePos = pos / 8
	bitPos = 7 - (pos % 8)
	return
}

func (bm *bigEndian) SetBit(bf *BitField, pos uint64) error {
	if bytePos, bitPos, err := calcBitPosBE(bf, pos); err == nil {
		bf.data[bytePos] |= 1 << bitPos
		return nil
	} else {
		return err
	}
}

func (bm *bigEndian) ClearBit(bf *BitField, pos uint64) error {
	if bytePos, bitPos, err := calcBitPosBE(bf, pos); err == nil {
		bf.data[bytePos] &^= 1 << bitPos
		return nil
	} else {
		return err
	}
}

func (bm *bigEndian) ToggleBit(bf *BitField, pos uint64) error {
	if bytePos, bitPos, err := calcBitPosBE(bf, pos); err == nil {
		bf.data[bytePos] ^= 1 << bitPos
		return nil
	} else {
		return err
	}
}

func (bm *bigEndian) TestBit(bf *BitField, pos uint64) (bool, error) {
	if bytePos, bitPos, err := calcBitPosBE(bf, pos); err == nil {
		value := bf.data[bytePos] & (1 << bitPos)
		return value > 0, nil
	} else {
		return false, err
	}
}

func (bm *bigEndian) InsertUint64(bf *BitField, offset, size uint64, value uint64) error {
	if offset+size > bf.size || size > 64 {
		return errors.New("operation out of bounds or size is invalid")
	}

	for i := size; i > 0; i-- {
		pos := offset + i - 1
		if (value>>(size-i))&1 == 1 {
			if err := bf.manipulator.SetBit(bf, pos); err != nil {
				return err
			}
		} else {
			if err := bf.manipulator.ClearBit(bf, pos); err != nil {
				return err
			}
		}
	}

	return nil
}

func (bm *bigEndian) ExtractUint64(bf *BitField, offset, size uint64) (uint64, error) {
	if offset+size > bf.size || size > 64 {
		return 0, errors.New("operation out of bounds or size is invalid")
	}

	var group uint64
	for i := size; i > 0; i-- {
		bit, err := bf.manipulator.TestBit(bf, offset+i-1)
		if err != nil {
			return 0, err
		}
		if bit {
			group |= (1 << (size - i))
		}
	}
	return group, nil
}
