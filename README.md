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
err := bf.InsertUint64(48, 4, 0b0111)

// Check for error
if err != nil {
  panic(err)
}

// Extract a 4-bit value from the 48th bit
val, err := bf.ExtractUint64(48, 4)

// Check for error
if err != nil {
  panic(err)
}

println(val) // 7
```

## Error Handling

This library uses a fail-fast error handling strategy. If an error occurs during a _mutating_ method call, the error is stored (in addition to being returned) and subsequent _mutating_ method calls become no-ops that return the stored error. This allows you to perform a sequence of operations and then check the error once at the end.

```go
// Create a new BitField with a size of 32 bits
bf := bitfield.BigEndian.New(32)

// Insert many values
bf.InsertUint64(0, 4, 0b1010)
bf.InsertUint64(4, 4, 0b0000)
// ...
bf.InsertUint64(28, 4, 0b1010)

// Check for error once
if err := bf.Error(); err != nil {
  panic(err)
}
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
