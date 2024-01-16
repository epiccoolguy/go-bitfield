package bitfield

import (
	"errors"
	"reflect"
	"testing"
)

// Test case structs
type NewTestCase struct {
	name        string // Name of the test case
	n           uint   // Input size in bits for the New function
	expectedLen int    // Expected length of the underlying byte slice
}

type BytesTestCase struct {
	name         string    // Name of the test case
	bf           *BitField // Initial BitField for the test
	expectedBits []byte    // Expected byte slice after copying the bits
}

type SetBitTestCase struct {
	name         string    // Name of the test case
	bf           *BitField // Initial BitField for the test
	pos          uint      // Position of the bit to set
	expectError  bool      // Whether an error is expected
	expectedBits []byte    // Expected byte slice after setting the bit
}

type ClearBitTestCase struct {
	name         string    // Name of the test case
	bf           *BitField // Initial BitField for the test
	pos          uint      // Position of the bit to clear
	expectError  bool      // Whether an error is expected
	expectedBits []byte    // Expected byte slice after clearing the bit
}

type ToggleBitTestCase struct {
	name         string    // Name of the test case
	bf           *BitField // Initial BitField for the test
	pos          uint      // Position of the bit to toggle
	expectError  bool      // Whether an error is expected
	expectedBits []byte    // Expected byte slice after toggling the bit
}

type TestBitTestCase struct {
	name          string    // Name of the test case
	bf            *BitField // Initial BitField for the test
	pos           uint      // Position of the bit to get
	expectError   bool      // Whether an error is expected
	expectedValue bool      // Expected value of the bit
}

type InsertUintTestCase struct {
	name         string
	bf           *BitField
	man          BitManipulator // Optional manipulator (mock or nil)
	offset       uint
	size         uint
	value        uint64
	expectError  bool
	expectedBits []byte
}

type ExtractUintTestCase struct {
	name          string
	bf            *BitField
	man           BitManipulator // Optional manipulator (mock or nil)
	offset        uint
	size          uint
	expectError   bool
	expectedValue uint64
}

// Mock BitManipulator
type MockBitManipulator struct {
	SetBitFunc    func(pos uint) error
	ClearBitFunc  func(pos uint) error
	ToggleBitFunc func(pos uint) error
	TestBitFunc   func(pos uint) (bool, error)
}

func (m *MockBitManipulator) SetBit(pos uint) error {
	return m.SetBitFunc(pos)
}

func (m *MockBitManipulator) ClearBit(pos uint) error {
	return m.ClearBitFunc(pos)
}

func (m *MockBitManipulator) ToggleBit(pos uint) error {
	return m.ToggleBitFunc(pos)
}
func (m *MockBitManipulator) TestBit(pos uint) (bool, error) {
	return m.TestBitFunc(pos)
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

var bytesTestCases = []BytesTestCase{
	{
		name:         "Empty BitField",
		bf:           New(0),
		expectedBits: []byte{},
	},
	{
		name: "BitField with set bits",
		bf: &BitField{
			bits: []byte{0b11111111, 0b11111111},
			sz:   16,
		},
		expectedBits: []byte{0b11111111, 0b11111111},
	},
}

var setBitTestCases = []SetBitTestCase{
	{
		name: "Set first bit to true in 2-byte field",
		bf: &BitField{
			bits: []byte{0b00000000, 0b00000000},
			sz:   16,
		},
		pos:          0,
		expectedBits: []byte{0b00000001, 0b00000000},
	},
	{
		name: "Set second bit in second byte",
		bf: &BitField{
			bits: []byte{0b00000000, 0b00000000},
			sz:   16,
		},
		pos:          9,
		expectedBits: []byte{0b00000000, 0b00000010},
	},
	{
		name: "Set bit out of range",
		bf: &BitField{
			bits: []byte{0b10101010},
			sz:   8,
		},
		pos:         100,
		expectError: true,
	},
}

var clearBitTestCases = []ClearBitTestCase{
	{
		name: "Clear first bit in 2-byte field",
		bf: &BitField{
			bits: []byte{0b11111111, 0b11111111},
			sz:   16,
		},
		pos:          0,
		expectedBits: []byte{0b11111110, 0b11111111},
	},
	{
		name: "Clear second bit in second byte",
		bf: &BitField{
			bits: []byte{0b11111111, 0b11111111},
			sz:   16,
		},
		pos:          9,
		expectedBits: []byte{0b11111111, 0b11111101},
	},
	{
		name: "Clear bit out of range",
		bf: &BitField{
			bits: []byte{0b10101010},
			sz:   8,
		},
		pos:         100,
		expectError: true,
	},
}

var toggleBitTestCases = []ToggleBitTestCase{
	{
		name: "Toggle first bit in 2-byte field",
		bf: &BitField{
			bits: []byte{0b00000000, 0b11111111},
			sz:   16,
		},
		pos:          0,
		expectedBits: []byte{0b00000001, 0b11111111},
	},
	{
		name: "Toggle second bit in second byte",
		bf: &BitField{
			bits: []byte{0b00000000, 0b11111111},
			sz:   16,
		},
		pos:          9,
		expectedBits: []byte{0b00000000, 0b11111101},
	},
	{
		name: "Toggle bit out of range",
		bf: &BitField{
			bits: []byte{0b10101010},
			sz:   8,
		},
		pos:         100,
		expectError: true,
	},
}

var testBitTestCases = []TestBitTestCase{
	{
		name: "Test first bit set in 2-byte field",
		bf: &BitField{
			bits: []byte{0b00000001, 0b00000000},
			sz:   16,
		},
		pos:           0,
		expectedValue: true,
	},
	{
		name: "Test second bit not set in second byte",
		bf: &BitField{
			bits: []byte{0b10000000, 0b00000000},
			sz:   16,
		},
		pos:           9,
		expectedValue: false,
	},
	{
		name: "Test bit out of range",
		bf: &BitField{
			bits: []byte{0b10101010},
			sz:   8,
		},
		pos:         100,
		expectError: true,
	},
}

var insertUintTestCases = []InsertUintTestCase{
	{
		name: "Insert within range",
		bf: &BitField{
			bits: []byte{0b00000000, 0b00000000},
			sz:   16,
		},
		offset:       0,
		size:         8,
		value:        0b10101010,
		expectedBits: []byte{0b10101010, 0b00000000},
	},
	{
		name: "Insert with overflow",
		bf: &BitField{
			bits: []byte{0b00000000, 0b00000000},
			sz:   16,
		},
		offset:      8,
		size:        10, // This goes beyond the size of BitField
		value:       0b1111111111,
		expectError: true,
	},
	{
		name: "Insert zero size",
		bf: &BitField{
			bits: []byte{0b00000000, 0b00000000},
			sz:   16,
		},
		offset:       5,
		size:         0,
		value:        0b1,
		expectedBits: []byte{0b00000000, 0b00000000},
	},
	{
		name: "Insert at offset",
		bf: &BitField{
			bits: []byte{0b00000000, 0b00000000},
			sz:   16,
		},
		offset:       4,
		size:         4,
		value:        0b1111,
		expectedBits: []byte{0b11110000, 0b00000000}, // The value 0b1111 starts at the 5th bit (offset 4), LSB-first
	},
	{
		name: "Invalid size greater than 64",
		bf: &BitField{
			bits: make([]byte, 16),
			sz:   128,
		},
		offset:      0,
		size:        65,
		value:       0b1111111111111111111111111111111111111111111111111111111111111111,
		expectError: true,
	},
	{
		name: "Insert spanning multiple bytes",
		bf: &BitField{
			bits: []byte{0b00000000, 0b00000000, 0b00000000, 0b00000000}, // Initial state with 4 bytes
			sz:   32,
		},
		offset:       4,
		size:         16,                 // 16-bit value
		value:        0b1010101010101010, // Spanning across multiple bytes
		expectedBits: []byte{0b10100000, 0b10101010, 0b00001010, 0b00000000},
	},
	{
		name: "Mock error with SetBit",
		bf: &BitField{
			bits: make([]byte, 1),
			sz:   8,
		},
		man: &MockBitManipulator{
			SetBitFunc: func(pos uint) error {
				return errors.New("mock error")
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
			bits: make([]byte, 1),
			sz:   8,
		},
		man: &MockBitManipulator{
			ClearBitFunc: func(pos uint) error {
				return errors.New("mock error")
			},
		},
		offset:      0,
		size:        8,
		value:       0b10101010,
		expectError: true,
	},
}

var extractUintTestCases = []ExtractUintTestCase{
	{
		name: "Extract within range",
		bf: &BitField{
			bits: []byte{0b10101010, 0b00000000},
			sz:   16,
		},
		offset:        0,
		size:          8,
		expectedValue: 0b10101010,
	},
	{
		name: "Extract with overflow",
		bf: &BitField{
			bits: []byte{0b00000000, 0b00000000},
			sz:   16,
		},
		offset:      8,
		size:        10, // This goes beyond the size of BitField
		expectError: true,
	},
	{
		name: "Extract zero size",
		bf: &BitField{
			bits: []byte{0b10101010, 0b00000000},
			sz:   16,
		},
		offset:        5,
		size:          0,
		expectedValue: 0,
	},
	{
		name: "Extract at offset",
		bf: &BitField{
			bits: []byte{0b11110000, 0b00000000},
			sz:   16,
		},
		offset:        4,
		size:          4,
		expectedValue: 0b1111,
	},
	{
		name: "Extract spanning multiple bytes",
		bf: &BitField{
			bits: []byte{0b10100000, 0b10101010, 0b00001010, 0b00000000},
			sz:   32,
		},
		offset:        4,
		size:          16,
		expectedValue: 0b1010101010101010,
	},
	{
		name: "Mock error with TestBit",
		bf: &BitField{
			bits: make([]byte, 1),
			sz:   8,
		},
		man: &MockBitManipulator{
			TestBitFunc: func(pos uint) (bool, error) {
				return false, errors.New("mock error")
			},
		},
		offset:      0,
		size:        8,
		expectError: true,
	},
}

func TestNew(t *testing.T) {
	for _, tc := range newTestCases {
		t.Run(tc.name, func(t *testing.T) {
			bf := New(tc.n)
			if len(bf.bits) != tc.expectedLen {
				t.Errorf("%s: expected byte length %d, got %d", tc.name, tc.expectedLen, len(bf.bits))
			}
			if bf.sz != tc.n {
				t.Errorf("%s: expected size %d, got %d", tc.name, tc.n, bf.sz)
			}
		})
	}
}

func TestBytes(t *testing.T) {
	for _, tc := range bytesTestCases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.bf.Bytes()
			// Compare the resulting byte slice with the expected slice, if no error is expected
			if !reflect.DeepEqual(c, tc.expectedBits) {
				t.Errorf("SetBit() got %v, want %v", tc.bf.bits, tc.expectedBits)
			}
		})
	}
}

func TestSetBit(t *testing.T) {
	for _, tc := range setBitTestCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.bf.SetBit(tc.pos)

			// Check for error consistency
			if (err != nil) != tc.expectError {
				t.Errorf("SetBit() error = %v, expectError %v", err, tc.expectError)
				return
			}

			// Compare the resulting byte slice with the expected slice, if no error is expected
			if !tc.expectError && !reflect.DeepEqual(tc.bf.bits, tc.expectedBits) {
				t.Errorf("SetBit() got %v, want %v", tc.bf.bits, tc.expectedBits)
			}
		})
	}
}

func TestClearBit(t *testing.T) {
	for _, tc := range clearBitTestCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.bf.ClearBit(tc.pos)

			// Check for error consistency
			if (err != nil) != tc.expectError {
				t.Errorf("ClearBit() error = %v, expectError %v", err, tc.expectError)
				return
			}

			// Compare the resulting byte slice with the expected slice, if no error is expected
			if !tc.expectError && !reflect.DeepEqual(tc.bf.bits, tc.expectedBits) {
				t.Errorf("ClearBit() got %v, want %v", tc.bf.bits, tc.expectedBits)
			}
		})
	}
}

func TestToggleBit(t *testing.T) {
	for _, tc := range toggleBitTestCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.bf.ToggleBit(tc.pos)

			// Check for error consistency
			if (err != nil) != tc.expectError {
				t.Errorf("ToggleBit() error = %v, expectError %v", err, tc.expectError)
				return
			}

			// Compare the resulting byte slice with the expected slice, if no error is expected
			if !tc.expectError && !reflect.DeepEqual(tc.bf.bits, tc.expectedBits) {
				t.Errorf("ToggleBit() got %v, want %v", tc.bf.bits, tc.expectedBits)
			}
		})
	}
}

func TestTestBit(t *testing.T) {
	for _, tc := range testBitTestCases {
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

func TestInsertUint(t *testing.T) {
	for _, tc := range insertUintTestCases {
		// Use the provided manipulator if it's not nil, else use the BitField itself
		man := tc.man
		if man == nil {
			man = tc.bf
		}

		t.Run(tc.name, func(t *testing.T) {
			err := tc.bf.InsertUint(man, tc.offset, tc.size, tc.value)

			if (err != nil) != tc.expectError {
				t.Errorf("InsertUint() error = %v, expectError %v", err, tc.expectError)
				return
			}

			if !tc.expectError && !reflect.DeepEqual(tc.bf.bits, tc.expectedBits) {
				t.Errorf("InsertUint() got %v, want %v", tc.bf.bits, tc.expectedBits)
			}
		})
	}
}

func TestExtractUint(t *testing.T) {
	for _, tc := range extractUintTestCases {
		// Use the provided manipulator if it's not nil, else use the BitField itself
		man := tc.man
		if man == nil {
			man = tc.bf
		}

		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.bf.ExtractUint(man, tc.offset, tc.size)

			if (err != nil) != tc.expectError {
				t.Errorf("ExtractUint() error = %v, expectError %v", err, tc.expectError)
				return
			}

			if !tc.expectError && value != tc.expectedValue {
				t.Errorf("ExtractUint() got %v, want %v", value, tc.expectedValue)
			}
		})
	}
}
