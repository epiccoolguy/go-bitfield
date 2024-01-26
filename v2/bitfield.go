package bitfield

// BitField represents a field of bits. It provides methods for manipulating bits within a byte slice.
type BitField struct {
	data        []byte         // The underlying byte slice that stores the bits.
	size        uint64         // The size of the bit field in bits.
	manipulator BitManipulator // An interface that provides methods for bit manipulation.
	err         error          // An error that is set when a bit manipulation method fails.
}

// BitManipulator is an interface that defines methods for manipulating bits in a BitField.
// This includes setting, clearing, toggling, and testing individual bits, as well as
// inserting and extracting multi-bit values.
// BigEndian and LittleEndian are the included implementations of this interface.
type BitManipulator interface {
	New(n uint64) *BitField
	FromBytes(bytes []byte, size uint64) *BitField
	SetBit(bf *BitField, pos uint64) error
	ClearBit(bf *BitField, pos uint64) error
	ToggleBit(bf *BitField, pos uint64) error
	TestBit(bf *BitField, pos uint64) (bool, error)
	InsertUint64(bf *BitField, offset, size, value uint64) error
	ExtractUint64(bf *BitField, offset, size uint64) (uint64, error)
}

// Bytes returns a copy of the underlying data as a byte slice.
func (bf *BitField) Bytes() []byte {
	copiedBytes := make([]byte, len(bf.data))
	copy(copiedBytes, bf.data)
	return copiedBytes
}

// Size returns the size of the BitField in number of bits.
func (bf *BitField) Size() uint64 {
	return bf.size
}

// Error returns the error set by the last failing bit manipulation method.
func (bf *BitField) Error() error {
	return bf.err
}

func (bf *BitField) SetBit(pos uint64) error {
	if bf.err == nil {
		bf.err = bf.manipulator.SetBit(bf, pos)
	}
	return bf.err
}

func (bf *BitField) ClearBit(pos uint64) error {
	if bf.err == nil {
		bf.err = bf.manipulator.ClearBit(bf, pos)
	}
	return bf.err
}

func (bf *BitField) ToggleBit(pos uint64) error {
	if bf.err == nil {
		bf.err = bf.manipulator.ToggleBit(bf, pos)
	}
	return bf.err
}

func (bf *BitField) TestBit(pos uint64) (bool, error) {
	return bf.manipulator.TestBit(bf, pos)
}

func (bf *BitField) InsertUint64(offset, size, value uint64) error {
	if bf.err == nil {
		bf.err = bf.manipulator.InsertUint64(bf, offset, size, value)
	}
	return bf.err
}

func (bf *BitField) ExtractUint64(offset, size uint64) (uint64, error) {
	return bf.manipulator.ExtractUint64(bf, offset, size)
}
