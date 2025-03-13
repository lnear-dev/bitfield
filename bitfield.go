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

// storageType is a constraint that permits container types for storing bit fields.
type storageType interface {
	~uint | ~uint32 | ~uint64
}

// BitField represents a field of bits within a larger unsigned integer.
// T represents the type of values that can be stored in the field.
// U represents the container type where the bit field will be stored.
type BitField[T Unsigned, U storageType] struct {
	Shift uint // Position of the least significant bit of the field
	Size  uint // Number of bits in the field
	Mask  U    // Mask with 1s in the field position
}

// New creates a new BitField with the given shift and size.
// shift determines the position of the least significant bit of the field.
// size determines how many bits the field will occupy.
// The function creates a mask with 1s in the positions of the field.
// Note: This function doesn't perform validation, use Safe for validated creation.
func New[T Unsigned, U storageType](shift, size uint) BitField[T, U] {
	return BitField[T, U]{
		Shift: shift,
		Size:  size,
		Mask:  (U(1) << (shift + size)) - (U(1) << shift),
	}
}

// Safe creates a new BitField with the given shift and size, after validating the parameters.
// Returns an error if:
// - shift is greater than or equal to the bit size of type T
// - size is greater than or equal to the bit size of type T
// - shift + size exceeds the bit size of type T
// - size is less than or equal to 0
func Safe[T Unsigned, U storageType](shift, size uint) (BitField[T, U], error) {
	var bf BitField[T, U]
	switch bSize := unsignedSizeOf[T](); {
	case shift >= bSize:
		return bf, fmt.Errorf("invalid shift parameter")
	case size >= bSize:
		return bf, fmt.Errorf("invalid size parameter")
	case shift+size > bSize:
		return bf, fmt.Errorf("invalid shift/size parameters")
	case size <= 0:
		return bf, fmt.Errorf("invalid size parameter")
	}
	return New[T, U](shift, size), nil
}

// SafeNext creates a new BitField that starts after an existing BitField, with validation.
// It takes an existing BitField and creates a new one of the specified size
// that starts immediately after the end of the existing field.
// Returns an error if the new field would exceed the bounds of type T.
func SafeNext[
	T Unsigned,
	U storageType,
	Old Unsigned,
](
	bf BitField[Old, U],
	size uint,
) (BitField[T, U], error) {
	if bf.Shift+bf.Size+size > unsignedSizeOf[T]() {
		return BitField[T, U]{}, fmt.Errorf("field would exceed type bounds")
	}
	return Safe[T, U](bf.Shift+bf.Size, size)
}

// Next creates a new BitField that starts after an existing BitField.
// It takes an existing BitField and creates a new one of the specified size
// that starts immediately after the end of the existing field.
// Note: This function doesn't perform validation, use SafeNext for validated creation.
func Next[T Unsigned, U storageType, Old Unsigned](bf BitField[Old, U], size uint) BitField[T, U] {
	return New[T, U](bf.Shift+bf.Size, size)
}

// IsValid checks if the value fits within the bit field.
// Returns true if the value can be represented using the field's size.
func (bf BitField[T, U]) IsValid(value T) bool {
	return U(value) < U(1)<<bf.Size
}

// Encode encodes a value into the bit field.
// It shifts the value to the appropriate position.
// Panics if the value is too large for the field.
func (bf BitField[T, U]) Encode(value T) U {
	if !bf.IsValid(value) {
		panic(fmt.Sprintf("value %v out of range, max %v", value, U(1)<<bf.Size-1))
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
// Note: This is a method version of the Next function.
func (bf BitField[T, U]) NextBitField(size uint) BitField[T, U] {
	if totalBits := unsignedSizeOf[T](); bf.Shift+bf.Size+size > totalBits {
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

// unsignedSizeOf returns the size in bits of the unsigned type T.
func unsignedSizeOf[T Unsigned]() uint {
	return uint(unsafe.Sizeof(T(0)) * 8)
}
