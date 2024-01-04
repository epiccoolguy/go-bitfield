package bitfield

import (
	"reflect"
	"testing"
)

// Test case structs
type NewTestCase struct {
	name        string // Name of the test case
	n           uint   // Input size in bits for the New function
	expectedLen int    // Expected length of the underlying byte slice
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
