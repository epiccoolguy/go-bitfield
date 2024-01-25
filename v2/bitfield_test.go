package bitfield

import (
	"reflect"
	"testing"
)

// Test case structs

type NewTestCase struct {
	name        string // Name of the test case
	n           uint64 // Input size in bits for the New function
	expectedLen int    // Expected length of the underlying byte slice
}

type FromBytesTestCase struct {
	name  string // Name of the test case
	bytes []byte // Expected byte slice after copying the bits
	size  uint64 // Expected length of the underlying byte slice
}

type BytesTestCase struct {
	name         string    // Name of the test case
	bf           *BitField // Initial BitField for the test
	expectedBits []byte    // Expected byte slice after copying the bits
}

type SizeTestCase struct {
	name         string    // Name of the test case
	bf           *BitField // Initial BitField for the test
	expectedSize uint64    // Expected size of the BitField
}

type SetBitTestCase struct {
	name         string    // Name of the test case
	bf           *BitField // Initial BitField for the test
	pos          uint64    // Position of the bit to set
	expectError  bool      // Whether an error is expected
	expectedBits []byte    // Expected byte slice after setting the bit
}

type ClearBitTestCase struct {
	name         string    // Name of the test case
	bf           *BitField // Initial BitField for the test
	pos          uint64    // Position of the bit to clear
	expectError  bool      // Whether an error is expected
	expectedBits []byte    // Expected byte slice after clearing the bit
}

type ToggleBitTestCase struct {
	name         string    // Name of the test case
	bf           *BitField // Initial BitField for the test
	pos          uint64    // Position of the bit to toggle
	expectError  bool      // Whether an error is expected
	expectedBits []byte    // Expected byte slice after toggling the bit
}

type TestBitTestCase struct {
	name          string    // Name of the test case
	bf            *BitField // Initial BitField for the test
	pos           uint64    // Position of the bit to get
	expectError   bool      // Whether an error is expected
	expectedValue bool      // Expected value of the bit
}

type InsertUint64TestCase struct {
	name         string    // The name of the test case.
	bf           *BitField // The BitField to insert the value into.
	offset       uint64    // The offset at which to insert the value.
	size         uint64    // The size of the value to insert.
	value        uint64    // The value to insert.
	expectError  bool      // Indicates whether an error is expected during insertion.
	expectedBits []byte    // The expected bits after insertion.
}

type ExtractUint64TestCase struct {
	name          string    // name of the test case
	bf            *BitField // BitField instance
	offset        uint64    // offset of the value in the BitField
	size          uint64    // size of the value in bits
	expectError   bool      // indicates if an error is expected
	expectedValue uint64    // expected extracted value
}

// Test cases

var newTestCases = []NewTestCase{
	{
		name:        "Zero size",
		n:           0,
		expectedLen: 0,
	},
	{
		name:        "Non-multiple of 8 size",
		n:           10,
		expectedLen: 2, // 10 bits require at least 2 bytes
	},
	{
		name:        "Exact byte size",
		n:           8,
		expectedLen: 1,
	},
	{
		name:        "Large size",
		n:           1024,
		expectedLen: 128, // 1024 bits = 128 bytes
	},
}

var fromBytesTestCases = []FromBytesTestCase{
	{
		name:  "Empty BitField",
		bytes: []byte{},
	},
	{
		name:  "Set bits",
		bytes: []byte{0b11111111, 0b11111111},
		size:  16,
	},
}

var bytesTestCases = []BytesTestCase{
	{
		name:         "Empty BitField",
		bf:           LittleEndian.New(0),
		expectedBits: []byte{},
	},
	{
		name: "BitField with set bits",
		bf: &BitField{
			data:        []byte{0b11111111, 0b11111111},
			size:        16,
			manipulator: LittleEndian,
		},
		expectedBits: []byte{0b11111111, 0b11111111},
	},
}

var sizeTestCases = []SizeTestCase{
	{
		name:         "Empty BitField",
		bf:           LittleEndian.New(0),
		expectedSize: 0,
	},
	{
		name: "BitField with set bits",
		bf: &BitField{
			data:        []byte{0b11111111, 0b11111111},
			size:        16,
			manipulator: LittleEndian,
		},
		expectedSize: 16,
	},
}

func TestBytes(t *testing.T) {
	for _, tc := range bytesTestCases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.bf.Bytes()
			// Compare the resulting byte slice with the expected slice
			if !reflect.DeepEqual(c, tc.expectedBits) {
				t.Errorf("Bytes() got %v, want %v", tc.bf.data, tc.expectedBits)
			}
		})
	}
}

func TestSize(t *testing.T) {
	for _, tc := range sizeTestCases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.bf.Size()
			// Compare the resulting size with the expected size
			if c != tc.expectedSize {
				t.Errorf("Size() got %v, want %v", c, tc.expectedSize)
			}
		})
	}
}
