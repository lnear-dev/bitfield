// Package bitfield provides functionality for working with bit fields in Go.
// It allows for type-safe manipulation of bit fields within unsigned integer types.
package bitfield

import (
	"fmt"
	"unsafe"
)

// Unsigned is a constraint that permits any unsigned integer type.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// BitField represents a field of bits within a larger unsigned integer.
// T represents the type of values that can be stored in the field.
// U represents the container type where the bit field will be stored.
type BitField[T Unsigned, U ~uint | ~uint32 | ~uint64] struct {
	Shift int // Position of the least significant bit of the field
	Size  int // Number of bits in the field
	Mask  U   // Mask with 1s in the field position
}

// New creates a new BitField with the given shift and size.
// shift determines the position of the least significant bit of the field.
// size determines how many bits the field will occupy.
// It panics if:
// - shift is greater than or equal to the bit size of type T
// - size is greater than or equal to the bit size of type T
// - shift + size exceeds the bit size of type T
// - size is less than or equal to 0
func New[T Unsigned, U ~uint | ~uint32 | ~uint64](shift, size int) BitField[T, U] {
	var zero T
	switch bSize := int(unsafe.Sizeof(zero) * 8); {
	case shift >= bSize:
		panic("invalid shift parameter")
	case size >= bSize:
		panic("invalid size parameter")
	case shift+size > bSize:
		panic("invalid shift/size parameters")
	case size <= 0:
		panic("invalid size parameter")
	}
	return BitField[T, U]{
		Shift: shift,
		Size:  size,
		Mask:  (U(1) << (shift + size)) - (U(1) << shift),
	}
}

// IsValid checks if the value fits within the bit field.
// Returns true if the value can be represented using the field's size.
func (bf BitField[T, U]) IsValid(value T) bool {
	return value < T(1<<bf.Size)
}

// Encode encodes a value into the bit field.
// It shifts the value to the appropriate position.
// Panics if the value is too large for the field.
func (bf BitField[T, U]) Encode(value T) U {
	if !bf.IsValid(value) {
		panic(fmt.Sprintf("value %v out of range", value))
	}
	return U(value) << bf.Shift
}

// Update updates the bit field within an existing value.
// It clears the existing bits in the field and sets them to the new value.
// Panics if the new value is too large for the field.
func (bf BitField[T, U]) Update(previous U, value T) U {
	return (previous &^ bf.Mask) | bf.Encode(value)
}

// Decode extracts the bit field from a value.
// It masks out all other bits and shifts the field down to position 0.
func (bf BitField[T, U]) Decode(value U) T {
	return T((value & bf.Mask) >> bf.Shift)
}

// NextBitField returns a new BitField starting from the end of the current one.
// The new field will have the specified size.
// Panics if the new field would exceed the bounds of type T.
func (bf BitField[T, U]) NextBitField(size int) BitField[T, U] {
	var zero T
	if totalBits := int(unsafe.Sizeof(zero) * 8); bf.Shift+bf.Size+size > totalBits {
		panic("field would exceed type bounds")
	}
	return New[T, U](bf.Shift+bf.Size, size)
}

// Clear zeroes out the bits in this field while preserving all other bits.
// Returns the modified value with this field cleared.
func (bf BitField[T, U]) Clear(value U) U {
	return value &^ bf.Mask
}

// HasValue checks if the field contains a specific value.
// Returns true if the value in the container matches the provided value.
func (bf BitField[T, U]) HasValue(container U, value T) bool {
	return bf.Decode(container) == value
}
