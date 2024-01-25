// +---------------------------------------------------------------------------------------------------------------+
// | Memory Representation of 0x12345678 in Big-Endian with MSb 0 numbering                                        |
// +---------------------------------------------------------------------------------------------------------------+
// |           Byte 1 (MSB)           | Byte 2                 | Byte 3                  | Byte 4 (LSB)            |
// |           0x12                   | 0x34                   | 0x56                    | 0x78                    |
// +---------------------------------------------------------------------------------------------------------------+
// | Binary:   0  0  0  1  0  0  1  0 | 0  0  1  1  0  1  0  0 | 0  1  0  1  0  1  1  0  | 1  1  1  1  0  0  0  0  |
// +---------------------------------------------------------------------------------------------------------------+
// | Position: 0  1  2  3  4  5  6  7 | 8  9  10 11 12 13 14 15| 16 17 18 19 20 21 22 23 | 24 25 26 27 28 29 30 31 |
// +---------------------------------------------------------------------------------------------------------------+

package bitfield

import (
	"errors"
	"reflect"
	"testing"
)

// Test case structs

// Mocks

type MockBitManipulatorBE struct {
	*bigEndian        // Embed the default Big-Endian bit manipulator so we only have to override methods we care about
	SetBitFunc        func(bf *BitField, pos uint64) error
	ClearBitFunc      func(bf *BitField, pos uint64) error
	ToggleBitFunc     func(bf *BitField, pos uint64) error
	TestBitFunc       func(bf *BitField, pos uint64) (bool, error)
	InsertUint64Func  func(bf *BitField, offset, size, value uint64) error
	ExtractUint64Func func(bf *BitField, offset, size uint64) (uint64, error)
}

// Compile-time check to ensure MockBitManipulator implements BitManipulator
var _ BitManipulator = &MockBitManipulatorBE{}

func (m *MockBitManipulatorBE) SetBit(bf *BitField, pos uint64) error {
	if m.SetBitFunc != nil {
		return m.SetBitFunc(bf, pos)
	} else {
		return m.bigEndian.SetBit(bf, pos)
	}
}

func (m *MockBitManipulatorBE) ClearBit(bf *BitField, pos uint64) error {
	if m.ClearBitFunc != nil {
		return m.ClearBitFunc(bf, pos)
	} else {
		return m.bigEndian.ClearBit(bf, pos)
	}
}

func (m *MockBitManipulatorBE) ToggleBit(bf *BitField, pos uint64) error {
	if m.ToggleBitFunc != nil {
		return m.ToggleBitFunc(bf, pos)
	} else {
		return m.bigEndian.ToggleBit(bf, pos)
	}
}
func (m *MockBitManipulatorBE) TestBit(bf *BitField, pos uint64) (bool, error) {
	if m.TestBitFunc != nil {
		return m.TestBitFunc(bf, pos)
	} else {
		return m.bigEndian.TestBit(bf, pos)
	}
}

func (m *MockBitManipulatorBE) InsertUint64(bf *BitField, offset, size, value uint64) error {
	if m.InsertUint64Func != nil {
		return m.InsertUint64Func(bf, offset, size, value)
	} else {
		return m.bigEndian.InsertUint64(bf, offset, size, value)
	}
}
func (m *MockBitManipulatorBE) ExtractUint64(bf *BitField, offset, size uint64) (uint64, error) {
	if m.ExtractUint64Func != nil {
		return m.ExtractUint64Func(bf, offset, size)
	} else {
		return m.bigEndian.ExtractUint64(bf, offset, size)
	}
}

// Test cases

var setBitTestCasesBE = []SetBitTestCase{
	{
		name: "Set first bit to true in 2-byte field",
		bf: &BitField{
			data:        []byte{0b00000000, 0b00000000},
			size:        16,
			manipulator: BigEndian,
		},
		pos:          0,
		expectedBits: []byte{0b10000000, 0b00000000},
	},
	{
		name: "Set second bit in second byte",
		bf: &BitField{
			data:        []byte{0b00000000, 0b00000000},
			size:        16,
			manipulator: BigEndian,
		},
		pos:          9,
		expectedBits: []byte{0b00000000, 0b01000000},
	},
	{
		name: "Set bit out of range",
		bf: &BitField{
			data:        []byte{0b10101010},
			size:        8,
			manipulator: BigEndian,
		},
		pos:         100,
		expectError: true,
	},
}

var clearBitTestCasesBE = []ClearBitTestCase{
	{
		name: "Clear first bit in 2-byte field",
		bf: &BitField{
			data:        []byte{0b11111111, 0b11111111},
			size:        16,
			manipulator: BigEndian,
		},
		pos:          0,
		expectedBits: []byte{0b01111111, 0b11111111},
	},
	{
		name: "Clear second bit in second byte",
		bf: &BitField{
			data:        []byte{0b11111111, 0b11111111},
			size:        16,
			manipulator: BigEndian,
		},
		pos:          9,
		expectedBits: []byte{0b11111111, 0b10111111},
	},
	{
		name: "Clear bit out of range",
		bf: &BitField{
			data:        []byte{0b10101010},
			size:        8,
			manipulator: BigEndian,
		},
		pos:         100,
		expectError: true,
	},
}

var toggleBitTestCasesBE = []ToggleBitTestCase{
	{
		name: "Toggle first bit in 2-byte field",
		bf: &BitField{
			data:        []byte{0b00000000, 0b11111111},
			size:        16,
			manipulator: BigEndian,
		},
		pos:          0,
		expectedBits: []byte{0b10000000, 0b11111111},
	},
	{
		name: "Toggle second bit in second byte",
		bf: &BitField{
			data:        []byte{0b00000000, 0b11111111},
			size:        16,
			manipulator: BigEndian,
		},
		pos:          9,
		expectedBits: []byte{0b00000000, 0b10111111},
	},
	{
		name: "Toggle bit out of range",
		bf: &BitField{
			data:        []byte{0b10101010},
			size:        8,
			manipulator: BigEndian,
		},
		pos:         100,
		expectError: true,
	},
}

var testBitTestCasesBE = []TestBitTestCase{
	{
		name: "Test first bit set in 2-byte field",
		bf: &BitField{
			data:        []byte{0b10000000, 0b00000000},
			size:        16,
			manipulator: BigEndian,
		},
		pos:           0,
		expectedValue: true,
	},
	{
		name: "Test second bit not set in second byte",
		bf: &BitField{
			data:        []byte{0b11111111, 0b10111111},
			size:        16,
			manipulator: BigEndian,
		},
		pos:           9,
		expectedValue: false,
	},
	{
		name: "Test bit out of range",
		bf: &BitField{
			data:        []byte{0b10101010},
			size:        8,
			manipulator: BigEndian,
		},
		pos:         100,
		expectError: true,
	},
}

var insertUint64TestCasesBE = []InsertUint64TestCase{
	{
		name: "Insert within range",
		bf: &BitField{
			data:        []byte{0b00000000, 0b00000000},
			size:        16,
			manipulator: BigEndian,
		},
		offset:       0,
		size:         8,
		value:        0b10101010,
		expectedBits: []byte{0b10101010, 0b00000000},
	},
	{
		name: "Insert with overflow",
		bf: &BitField{
			data:        []byte{0b00000000, 0b00000000},
			size:        16,
			manipulator: BigEndian,
		},
		offset:      8,
		size:        10, // This goes beyond the size of BitField
		value:       0b1111111111,
		expectError: true,
	},
	{
		name: "Insert zero size",
		bf: &BitField{
			data:        []byte{0b00000000, 0b00000000},
			size:        16,
			manipulator: BigEndian,
		},
		offset:       5,
		size:         0,
		value:        0b1,
		expectedBits: []byte{0b00000000, 0b00000000},
	},
	{
		name: "Insert at offset",
		bf: &BitField{
			data:        []byte{0b00000000, 0b00000000},
			size:        16,
			manipulator: BigEndian,
		},
		offset:       4,
		size:         4,
		value:        0b1111,
		expectedBits: []byte{0b00001111, 0b00000000}, // The value 0b1111 starts at the 5th bit (offset 4), BE-first
	},
	{
		name: "Invalid size greater than 64",
		bf: &BitField{
			data:        make([]byte, 16),
			size:        128,
			manipulator: BigEndian,
		},
		offset:      0,
		size:        65,
		value:       0b1111111111111111111111111111111111111111111111111111111111111111,
		expectError: true,
	},
	{
		name: "Insert spanning multiple bytes",
		bf: &BitField{
			data:        []byte{0b00000000, 0b00000000, 0b00000000, 0b00000000}, // Initial state with 4 bytes
			size:        32,
			manipulator: BigEndian,
		},
		offset:       4,
		size:         16,                 // 16-bit value
		value:        0b1010101010101010, // Spanning across multiple bytes
		expectedBits: []byte{0b00001010, 0b10101010, 0b10100000, 0b00000000},
	},
	{
		name: "Mock error with SetBit",
		bf: &BitField{
			data: make([]byte, 1),
			size: 8,
			manipulator: &MockBitManipulatorBE{
				SetBitFunc: func(bf *BitField, pos uint64) error {
					return errors.New("mock error")
				},
			},
		},
		offset:      0,
		size:        8,
		value:       0b01010101,
		expectError: true,
	},
	{
		name: "Mock error with ClearBit",
		bf: &BitField{
			data: make([]byte, 1),
			size: 8,
			manipulator: &MockBitManipulatorBE{
				ClearBitFunc: func(bf *BitField, pos uint64) error {
					return errors.New("mock error")
				},
			},
		},
		offset:      0,
		size:        8,
		value:       0b10101010,
		expectError: true,
	},
}

var extractUint64TestCasesBE = []ExtractUint64TestCase{
	{
		name: "Extract within range",
		bf: &BitField{
			data:        []byte{0b10101010, 0b00000000},
			size:        16,
			manipulator: BigEndian,
		},
		offset:        0,
		size:          8,
		expectedValue: 0b10101010,
	},
	{
		name: "Extract with overflow",
		bf: &BitField{
			data:        []byte{0b00000000, 0b00000000},
			size:        16,
			manipulator: BigEndian,
		},
		offset:      8,
		size:        10, // This goes beyond the size of BitField
		expectError: true,
	},
	{
		name: "Extract zero size",
		bf: &BitField{
			data:        []byte{0b10101010, 0b00000000},
			size:        16,
			manipulator: BigEndian,
		},
		offset:        5,
		size:          0,
		expectedValue: 0,
	},
	{
		name: "Extract at offset",
		bf: &BitField{
			data:        []byte{0b00001111, 0b00000000},
			size:        16,
			manipulator: BigEndian,
		},
		offset:        4,
		size:          4,
		expectedValue: 0b1111,
	},
	{
		name: "Extract spanning multiple bytes",
		bf: &BitField{
			data:        []byte{0b10100000, 0b10101010, 0b00001010, 0b00000000},
			size:        32,
			manipulator: BigEndian,
		},
		offset:        4,
		size:          16,
		expectedValue: 0b0000101010100000,
	},
	{
		name: "Mock error with TestBit",
		bf: &BitField{
			data: make([]byte, 1),
			size: 8,
			manipulator: &MockBitManipulatorBE{
				TestBitFunc: func(bf *BitField, pos uint64) (bool, error) {
					return false, errors.New("mock error")
				},
			},
		},
		offset:      0,
		size:        8,
		expectError: true,
	},
}

// Test functions

func TestNewBE(t *testing.T) {
	for _, tc := range newTestCases {
		t.Run(tc.name, func(t *testing.T) {
			bf := BigEndian.New(tc.n)
			if len(bf.data) != tc.expectedLen {
				t.Errorf("%s: expected byte length %d, got %d", tc.name, tc.expectedLen, len(bf.data))
			}
			if bf.size != tc.n {
				t.Errorf("%s: expected size %d, got %d", tc.name, tc.n, bf.size)
			}
			if !reflect.DeepEqual(bf.manipulator, BigEndian) {
				t.Errorf("%s: expected manipulator %v, got %v", tc.name, BigEndian, bf.manipulator)
			}
		})
	}
}

func TestFromBytesBE(t *testing.T) {
	for _, tc := range fromBytesTestCases {
		t.Run(tc.name, func(t *testing.T) {
			bf := BigEndian.FromBytes(tc.bytes, tc.size)
			if bf.size != tc.size {
				t.Errorf("%s: expected size %d, got %d", tc.name, tc.size, bf.size)
			}
			if !reflect.DeepEqual(bf.data, tc.bytes) {
				t.Errorf("FromBytes() got %v, want %v", bf.data, tc.bytes)
			}
			if !reflect.DeepEqual(bf.manipulator, BigEndian) {
				t.Errorf("%s: expected manipulator %v, got %v", tc.name, BigEndian, bf.manipulator)
			}
		})
	}
}

func TestSetBitBE(t *testing.T) {
	for _, tc := range setBitTestCasesBE {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.bf.SetBit(tc.pos)

			// Check for error consistency
			if (err != nil) != tc.expectError {
				t.Errorf("SetBit() returned unexpected error: got %v, want %v", err, tc.expectError)
				return
			}

			// Compare the resulting byte slice with the expected slice, if no error is expected
			if !tc.expectError && !reflect.DeepEqual(tc.bf.data, tc.expectedBits) {
				t.Errorf("SetBit() got %v, want %v", tc.bf.data, tc.expectedBits)
			}
		})
	}
}

func TestClearBitBE(t *testing.T) {
	for _, tc := range clearBitTestCasesBE {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.bf.ClearBit(tc.pos)

			if (err != nil) != tc.expectError {
				t.Errorf("ClearBit() returned unexpected error: got %v, want %v", err, tc.expectError)
				return
			}

			if !tc.expectError && !reflect.DeepEqual(tc.bf.data, tc.expectedBits) {
				t.Errorf("ClearBit() got %v, want %v", tc.bf.data, tc.expectedBits)
			}
		})
	}
}

func TestToggleBitBE(t *testing.T) {
	for _, tc := range toggleBitTestCasesBE {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.bf.ToggleBit(tc.pos)

			if (err != nil) != tc.expectError {
				t.Errorf("ToggleBit() returned unexpected error: got %v, want %v", err, tc.expectError)
				return
			}

			if !tc.expectError && !reflect.DeepEqual(tc.bf.data, tc.expectedBits) {
				t.Errorf("ToggleBit() got %v, want %v", tc.bf.data, tc.expectedBits)
			}
		})
	}
}

func TestTestBitBE(t *testing.T) {
	for _, tc := range testBitTestCasesBE {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.bf.TestBit(tc.pos)

			if errOccurred := err != nil; tc.expectError != errOccurred {
				if tc.expectError {
					t.Errorf("%s: expected an error, but got none", tc.name)
				} else {
					t.Errorf("%s: did not expect an error, but got %v", tc.name, err)
				}
			} else if !tc.expectError && value != tc.expectedValue {
				t.Errorf("%s: expected value %t, got %t", tc.name, tc.expectedValue, value)
			}
		})
	}
}

func TestInsertUint64BE(t *testing.T) {
	for _, tc := range insertUint64TestCasesBE {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.bf.InsertUint64(tc.offset, tc.size, tc.value)

			if (err != nil) != tc.expectError {
				t.Errorf("InsertUint() returned unexpected error: got %v, want %v", err, tc.expectError)
				return
			}

			if !tc.expectError && !reflect.DeepEqual(tc.bf.data, tc.expectedBits) {
				t.Errorf("InsertUint() got %v, want %v", tc.bf.data, tc.expectedBits)
			}
		})
	}
}

func TestExtractUint64BE(t *testing.T) {
	for _, tc := range extractUint64TestCasesBE {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.bf.ExtractUint64(tc.offset, tc.size)

			if (err != nil) != tc.expectError {
				t.Errorf("ExtractUint() returned unexpected error: got %v, want %v", err, tc.expectError)
				return
			}

			if !tc.expectError && value != tc.expectedValue {
				t.Errorf("ExtractUint() got %v, want %v", value, tc.expectedValue)
			}
		})
	}
}
