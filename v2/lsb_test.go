// +---------------------------------------------------------------------------------------------------------------+
// | Memory Representation of 0x12345678 in Little-Endian with LSb 0 numbering                                     |
// +---------------------------------------------------------------------------------------------------------------+
// |           Byte 1 (LSB)           | Byte 2                 | Byte 3                  | Byte 4 (MSB)            |
// |           0x78                   | 0x56                   | 0x34                    | 0x12                    |
// +---------------------------------------------------------------------------------------------------------------+
// | Binary:   1  1  1  1  0  0  0  0 | 0  1  0  1  0  1  1  0 | 0  0  1  1  0  1  0  0  | 0  0  0  1  0  0  1  0  |
// +---------------------------------------------------------------------------------------------------------------+
// | Position: 7  6  5  4  3  2  1  0 | 15 14 13 12 11 10 9  8 | 23 22 21 20 19 18 17 16 | 31 30 29 28 27 26 25 24 |
// +---------------------------------------------------------------------------------------------------------------+

package bitfield

import (
	"errors"
	"reflect"
	"testing"
)

// Test case structs

// Mocks

type MockBitManipulatorLE struct {
	*littleEndian     // Embed the default Little-Endian bit manipulator so we only have to override methods we care about
	SetBitFunc        func(bf *BitField, pos uint64) error
	ClearBitFunc      func(bf *BitField, pos uint64) error
	ToggleBitFunc     func(bf *BitField, pos uint64) error
	TestBitFunc       func(bf *BitField, pos uint64) (bool, error)
	InsertUint64Func  func(bf *BitField, offset, size, value uint64) error
	ExtractUint64Func func(bf *BitField, offset, size uint64) (uint64, error)
}

// Compile-time check to ensure MockBitManipulator implements BitManipulator
var _ BitManipulator = &MockBitManipulatorLE{}

func (m *MockBitManipulatorLE) SetBit(bf *BitField, pos uint64) error {
	if m.SetBitFunc != nil {
		return m.SetBitFunc(bf, pos)
	} else {
		return m.littleEndian.SetBit(bf, pos)
	}
}

func (m *MockBitManipulatorLE) ClearBit(bf *BitField, pos uint64) error {
	if m.ClearBitFunc != nil {
		return m.ClearBitFunc(bf, pos)
	} else {
		return m.littleEndian.ClearBit(bf, pos)
	}
}

func (m *MockBitManipulatorLE) ToggleBit(bf *BitField, pos uint64) error {
	if m.ToggleBitFunc != nil {
		return m.ToggleBitFunc(bf, pos)
	} else {
		return m.littleEndian.ToggleBit(bf, pos)
	}
}
func (m *MockBitManipulatorLE) TestBit(bf *BitField, pos uint64) (bool, error) {
	if m.TestBitFunc != nil {
		return m.TestBitFunc(bf, pos)
	} else {
		return m.littleEndian.TestBit(bf, pos)
	}
}

func (m *MockBitManipulatorLE) InsertUint64(bf *BitField, offset, size, value uint64) error {
	if m.InsertUint64Func != nil {
		return m.InsertUint64Func(bf, offset, size, value)
	} else {
		return m.littleEndian.InsertUint64(bf, offset, size, value)
	}
}
func (m *MockBitManipulatorLE) ExtractUint64(bf *BitField, offset, size uint64) (uint64, error) {
	if m.ExtractUint64Func != nil {
		return m.ExtractUint64Func(bf, offset, size)
	} else {
		return m.littleEndian.ExtractUint64(bf, offset, size)
	}
}

// Test cases

var setBitTestCasesLE = []SetBitTestCase{
	{
		name: "Set first bit to true in 2-byte field",
		bf: &BitField{
			data:        []byte{0b00000000, 0b00000000},
			size:        16,
			manipulator: LittleEndian,
		},
		pos:          0,
		expectedBits: []byte{0b00000001, 0b00000000},
	},
	{
		name: "Set second bit in second byte",
		bf: &BitField{
			data:        []byte{0b00000000, 0b00000000},
			size:        16,
			manipulator: LittleEndian,
		},
		pos:          9,
		expectedBits: []byte{0b00000000, 0b00000010},
	},
	{
		name: "Set bit out of range",
		bf: &BitField{
			data:        []byte{0b10101010},
			size:        8,
			manipulator: LittleEndian,
		},
		pos:         100,
		expectError: true,
	},
}

var clearBitTestCasesLE = []ClearBitTestCase{
	{
		name: "Clear first bit in 2-byte field",
		bf: &BitField{
			data:        []byte{0b11111111, 0b11111111},
			size:        16,
			manipulator: LittleEndian,
		},
		pos:          0,
		expectedBits: []byte{0b11111110, 0b11111111},
	},
	{
		name: "Clear second bit in second byte",
		bf: &BitField{
			data:        []byte{0b11111111, 0b11111111},
			size:        16,
			manipulator: LittleEndian,
		},
		pos:          9,
		expectedBits: []byte{0b11111111, 0b11111101},
	},
	{
		name: "Clear bit out of range",
		bf: &BitField{
			data:        []byte{0b10101010},
			size:        8,
			manipulator: LittleEndian,
		},
		pos:         100,
		expectError: true,
	},
}

var toggleBitTestCasesLE = []ToggleBitTestCase{
	{
		name: "Toggle first bit in 2-byte field",
		bf: &BitField{
			data:        []byte{0b00000000, 0b11111111},
			size:        16,
			manipulator: LittleEndian,
		},
		pos:          0,
		expectedBits: []byte{0b00000001, 0b11111111},
	},
	{
		name: "Toggle second bit in second byte",
		bf: &BitField{
			data:        []byte{0b00000000, 0b11111111},
			size:        16,
			manipulator: LittleEndian,
		},
		pos:          9,
		expectedBits: []byte{0b00000000, 0b11111101},
	},
	{
		name: "Toggle bit out of range",
		bf: &BitField{
			data:        []byte{0b10101010},
			size:        8,
			manipulator: LittleEndian,
		},
		pos:         100,
		expectError: true,
	},
}

var testBitTestCasesLE = []TestBitTestCase{
	{
		name: "Test first bit set in 2-byte field",
		bf: &BitField{
			data:        []byte{0b00000001, 0b00000000},
			size:        16,
			manipulator: LittleEndian,
		},
		pos:           0,
		expectedValue: true,
	},
	{
		name: "Test second bit not set in second byte",
		bf: &BitField{
			data:        []byte{0b11111111, 0b11111101},
			size:        16,
			manipulator: LittleEndian,
		},
		pos:           9,
		expectedValue: false,
	},
	{
		name: "Test bit out of range",
		bf: &BitField{
			data:        []byte{0b10101010},
			size:        8,
			manipulator: LittleEndian,
		},
		pos:         100,
		expectError: true,
	},
}

var insertUint64TestCasesLE = []InsertUint64TestCase{
	{
		name: "Insert within range",
		bf: &BitField{
			data:        []byte{0b00000000, 0b00000000},
			size:        16,
			manipulator: LittleEndian,
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
			manipulator: LittleEndian,
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
			manipulator: LittleEndian,
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
			manipulator: LittleEndian,
		},
		offset:       4,
		size:         4,
		value:        0b1111,
		expectedBits: []byte{0b11110000, 0b00000000}, // The value 0b1111 starts at the 5th bit (offset 4), LE-first
	},
	{
		name: "Invalid size greater than 64",
		bf: &BitField{
			data:        make([]byte, 16),
			size:        128,
			manipulator: LittleEndian,
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
			manipulator: LittleEndian,
		},
		offset:       4,
		size:         16,                 // 16-bit value
		value:        0b1010101010101010, // Spanning across multiple bytes
		expectedBits: []byte{0b10100000, 0b10101010, 0b00001010, 0b00000000},
	},
	{
		name: "Mock error with SetBit",
		bf: &BitField{
			data: make([]byte, 1),
			size: 8,
			manipulator: &MockBitManipulatorLE{
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
			manipulator: &MockBitManipulatorLE{
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

var extractUint64TestCasesLE = []ExtractUint64TestCase{
	{
		name: "Extract within range",
		bf: &BitField{
			data:        []byte{0b10101010, 0b00000000},
			size:        16,
			manipulator: LittleEndian,
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
			manipulator: LittleEndian,
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
			manipulator: LittleEndian,
		},
		offset:        5,
		size:          0,
		expectedValue: 0,
	},
	{
		name: "Extract at offset",
		bf: &BitField{
			data:        []byte{0b11110000, 0b00000000},
			size:        16,
			manipulator: LittleEndian,
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
			manipulator: LittleEndian,
		},
		offset:        4,
		size:          16,
		expectedValue: 0b1010101010101010,
	},
	{
		name: "Mock error with TestBit",
		bf: &BitField{
			data: make([]byte, 1),
			size: 8,
			manipulator: &MockBitManipulatorLE{
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

func TestNewLE(t *testing.T) {
	for _, tc := range newTestCases {
		t.Run(tc.name, func(t *testing.T) {
			bf := LittleEndian.New(tc.n)
			if len(bf.data) != tc.expectedLen {
				t.Errorf("%s: expected byte length %d, got %d", tc.name, tc.expectedLen, len(bf.data))
			}
			if bf.size != tc.n {
				t.Errorf("%s: expected size %d, got %d", tc.name, tc.n, bf.size)
			}
			if !reflect.DeepEqual(bf.manipulator, LittleEndian) {
				t.Errorf("%s: expected manipulator %v, got %v", tc.name, LittleEndian, bf.manipulator)
			}
		})
	}
}

func TestFromBytesLE(t *testing.T) {
	for _, tc := range fromBytesTestCases {
		t.Run(tc.name, func(t *testing.T) {
			bf := LittleEndian.FromBytes(tc.bytes, tc.size)
			if bf.size != tc.size {
				t.Errorf("%s: expected size %d, got %d", tc.name, tc.size, bf.size)
			}
			if !reflect.DeepEqual(bf.data, tc.bytes) {
				t.Errorf("FromBytes() got %v, want %v", bf.data, tc.bytes)
			}
			if !reflect.DeepEqual(bf.manipulator, LittleEndian) {
				t.Errorf("%s: expected manipulator %v, got %v", tc.name, LittleEndian, bf.manipulator)
			}
		})
	}
}

func TestSetBitLE(t *testing.T) {
	for _, tc := range setBitTestCasesLE {
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

func TestClearBitLE(t *testing.T) {
	for _, tc := range clearBitTestCasesLE {
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

func TestToggleBitLE(t *testing.T) {
	for _, tc := range toggleBitTestCasesLE {
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

func TestTestBitLE(t *testing.T) {
	for _, tc := range testBitTestCasesLE {
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

func TestInsertUint64LE(t *testing.T) {
	for _, tc := range insertUint64TestCasesLE {
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

func TestExtractUint64LE(t *testing.T) {
	for _, tc := range extractUint64TestCasesLE {
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
