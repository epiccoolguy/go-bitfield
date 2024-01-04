package bitfield

import (
	"reflect"
	"testing"
)

// Test case struct for creating new BitFields
type NewTestCase struct {
	name        string // Name of the test case
	n           uint   // Input size in bits for the New function
	expectedLen int    // Expected length of the underlying byte slice
}

// Test cases for creating new BitFields
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

// TestNew tests the New function of the bitfield package.
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
