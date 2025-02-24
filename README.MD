# BitField

A Go package for type-safe bit field manipulation with generic support.

## Installation

```bash
go get github.com/lnear-dev/bitfield
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/lnear-dev/bitfield"
)

func main() {
    // Create a 3-bit field starting at position 2
    field := bitfield.NewBitField[uint8, uint32](2, 3)

    // Encode a value
    encoded := field.Encode(5)
    fmt.Printf("Encoded: %032b\n", encoded)

    // Update an existing value
    container := uint32(0xFFFFFFFF)
    updated := field.Update(container, 3)
    fmt.Printf("Updated: %032b\n", updated)

    // Decode a value
    decoded := field.Decode(encoded)
    fmt.Printf("Decoded: %d\n", decoded)
}
```

## Features

- Generic support for different unsigned integer types
- Type-safe bit field manipulation
- Support for:
  - Field creation with bounds checking
  - Value encoding and decoding
  - Field updates
  - Value validation
  - Field clearing
  - Adjacent field creation

## API Documentation

[Go to full documentation on pkg.go.dev](https://pkg.go.dev/github.com/lnear-dev/bitfield)

## License

MIT License - see LICENSE file
