# Bitfield Manipulation Library

This is a Go library for manipulating bitfields. It provides functions for manipulating individual bits, as well as inserting and extracting unsigned integers of arbitrary size (currently, up to 64-bit) at arbitrary bit positions within a byte slice.

## Usage

To use this library, import it in your Go code:

```go
import "go.loafoe.dev/bitfield/v2"
```

Then, create a BitField and use its methods to manipulate bitfields:

```go
// Create a new BitField with a size of 128 bits
bf := bitfield.BigEndian.New(128)

// Insert a 4-bit value at the 48th bit
bf.InsertUint64(48, 4, 0b0111)

// Extract a 4-bit value from the 48th bit
val, _ := bf.ExtractUint64(48, 4)

// val == 7
```

## Endianness

There are two built-in implementations for manipulating bitfields:

1. Little-endian with LSb 0 numbering

   ```text
   +-----------------------------------------------------------+
   | 0x1234 in Little-Endian with LSb 0 numbering              |
   +-----------------------------------------------------------+
   |           Byte 1 (LSB)           | Byte 2 (MSB)           |
   |           0x34                   | 0x12                   |
   +-----------------------------------------------------------+
   | Binary:   0  0  1  1  0  1  0  0 | 0  0  0  1  0  0  1  0 |
   +-----------------------------------------------------------+
   | Position: 7  6  5  4  3  2  1  0 | 15 14 13 12 11 10 9  8 |
   +-----------------------------------------------------------+
   ```

2. Big-endian with MSb 0 numbering

   ```text
   +------------------------------------------------------------+
   | 0x1234 in Big-Endian with MSb 0 numbering                  |
   +------------------------------------------------------------+
   |           Byte 1 (MSB)           | Byte 2 (LSB)            |
   |           0x12                   | 0x34                    |
   +------------------------------------------------------------+
   | Binary:   0  0  0  1  0  0  1  0 | 0  0  1  1  0  1  0  0  |
   +------------------------------------------------------------+
   | Position: 0  1  2  3  4  5  6  7 | 8  9  10 11 12 13 14 15 |
   +------------------------------------------------------------+
   ```

## Testing

To run the unit tests, use the go test command:

```sh
go test ./...
```

## License

This project is licensed under the terms of the license provided in the LICENSE file.
