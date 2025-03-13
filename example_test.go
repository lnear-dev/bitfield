package bitfield

import (
	"fmt"
	"testing"
)

// Define enum-like types for use with bitfields
type Color uint8

const (
	Red Color = iota
	Green
	Blue
	Yellow
	Purple
	Orange
	Cyan
	Magenta
)

type Priority uint8

const (
	Low Priority = iota
	Medium
	High
	Critical
)

type Status uint8

const (
	Inactive Status = iota
	Active
	Suspended
	Pending
	Error
)

func ExampleBitField_basic() {
	// Create a bit field for a 4-bit value starting at position 0 in a uint32
	bf := New[uint8, uint32](0, 4)

	// Encode a value into the bit field
	encoded := bf.Encode(5)
	fmt.Printf("5 encoded into field: 0x%08X\n", encoded)

	// Decode a value from the bit field
	value := bf.Decode(encoded)
	fmt.Printf("Decoded value: %d\n", value)

	// Update a value in an existing container
	container := uint32(0xF0F0F0F0)
	fmt.Printf("Container before update: 0x%08X\n", container)

	updated := bf.Update(container, 7)
	fmt.Printf("Container after update: 0x%08X\n", updated)

	// Check if container has specific value
	hasValue := bf.HasValue(updated, 7)
	fmt.Printf("Container has value 7: %t\n", hasValue)

	// Clear the bit field
	cleared := bf.Clear(updated)
	fmt.Printf("Container after clearing field: 0x%08X\n", cleared)

	// Output:
	// 5 encoded into field: 0x00000005
	// Decoded value: 5
	// Container before update: 0xF0F0F0F0
	// Container after update: 0xF0F0F0F7
	// Container has value 7: true
	// Container after clearing field: 0xF0F0F0F0
}

func ExampleBitField_multiple() {
	// Create multiple bit fields within a uint32
	flags := uint32(0)

	// Define our bit fields
	isActive := New[uint8, uint32](0, 1)   // 1-bit boolean at position 0
	priority := New[uint8, uint32](1, 3)   // 3-bit priority field (0-7) at position 1
	category := New[uint8, uint32](4, 4)   // 4-bit category ID (0-15) at position 4
	errorCode := New[uint16, uint32](8, 8) // 8-bit error code at position 8

	// Update the container with values for each field
	flags = isActive.Update(flags, 1)   // Set active flag
	flags = priority.Update(flags, 3)   // Set priority to 3 (medium)
	flags = category.Update(flags, 5)   // Set category to 5
	flags = errorCode.Update(flags, 42) // Set error code to 42

	fmt.Printf("Combined flags: 0x%08X\n", flags)

	// Extract values from each field
	fmt.Printf("Active: %d\n", isActive.Decode(flags))
	fmt.Printf("Priority: %d\n", priority.Decode(flags))
	fmt.Printf("Category: %d\n", category.Decode(flags))
	fmt.Printf("Error code: %d\n", errorCode.Decode(flags))

	// Output:
	// Combined flags: 0x00002A57
	// Active: 1
	// Priority: 3
	// Category: 5
	// Error code: 42
}

func ExampleBitField_nextBitField() {
	// Create an initial bit field
	field1 := New[uint16, uint32](0, 4) // 4-bit field at position 0

	// Create a second field that starts immediately after the first
	field2 := field1.NextBitField(6) // 6-bit field at position 4

	// Create a third field that starts after the second
	field3 := field2.NextBitField(6) // 6-bit field at position 10

	// Verify field positions and sizes
	fmt.Printf("Field 1: shift=%d, size=%d, mask=0x%08X\n", field1.Shift, field1.Size, field1.Mask)
	fmt.Printf("Field 2: shift=%d, size=%d, mask=0x%08X\n", field2.Shift, field2.Size, field2.Mask)
	fmt.Printf("Field 3: shift=%d, size=%d, mask=0x%08X\n", field3.Shift, field3.Size, field3.Mask)

	// Use all three fields together
	var container uint32 = 0
	container = field1.Update(container, 9)  // Set field 1 to 9
	container = field2.Update(container, 33) // Set field 2 to 33
	container = field3.Update(container, 42) // Set field 3 to 42

	fmt.Printf("Container with all fields: 0x%08X\n", container)
	fmt.Printf("Field 1 value: %d\n", field1.Decode(container))
	fmt.Printf("Field 2 value: %d\n", field2.Decode(container))
	fmt.Printf("Field 3 value: %d\n", field3.Decode(container))

	// Output:
	// Field 1: shift=0, size=4, mask=0x0000000F
	// Field 2: shift=4, size=6, mask=0x000003F0
	// Field 3: shift=10, size=6, mask=0x0000FC00
	// Container with all fields: 0x0000AA19
	// Field 1 value: 9
	// Field 2 value: 33
	// Field 3 value: 42
}

func ExampleBitField_withEnums() {
	// Create bit fields for different enum types
	var config uint32 = 0

	// Create fields for our enum types
	colorField := New[Color, uint32](0, 3)         // 3 bits for Color (up to 8 values)
	priorityField := Next[Priority](colorField, 2) // 2 bits for Priority (up to 4 values)
	statusField := Next[Status](priorityField, 3)  // 3 bits for Status (up to 8 values)

	// Set values using enum constants
	config = colorField.Update(config, Blue)
	config = priorityField.Update(config, High)
	config = statusField.Update(config, Active)

	fmt.Printf("Config with enums: 0x%08X\n", config)

	// Decode values as enums
	color := colorField.Decode(config)
	priority := priorityField.Decode(config)
	status := statusField.Decode(config)

	// Print decoded enum values
	fmt.Printf("Color: %d\n", color)
	fmt.Printf("Priority: %d\n", priority)
	fmt.Printf("Status: %d\n", status)

	// Check for specific enum values
	isBlue := colorField.HasValue(config, Blue)
	isHigh := priorityField.HasValue(config, High)
	isActive := statusField.HasValue(config, Active)

	fmt.Printf("Is Blue? %t\n", isBlue)
	fmt.Printf("Is High Priority? %t\n", isHigh)
	fmt.Printf("Is Active? %t\n", isActive)

	// Change enum values
	config = colorField.Update(config, Purple)
	config = statusField.Update(config, Pending)

	fmt.Printf("Updated config: 0x%08X\n", config)
	fmt.Printf("New Color: %d\n", colorField.Decode(config))
	fmt.Printf("New Status: %d\n", statusField.Decode(config))

	// Output:
	// Config with enums: 0x00000032
	// Color: 2
	// Priority: 2
	// Status: 1
	// Is Blue? true
	// Is High Priority? true
	// Is Active? true
	// Updated config: 0x00000074
	// New Color: 4
	// New Status: 3
}

func ExampleBitField_enumSafety() {
	// Demonstrate how the package handles invalid enum values

	// Create a field for Color enum (3 bits, values 0-7)
	colorField := New[Color, uint32](0, 3)

	// All defined colors work fine
	fmt.Println("Setting valid colors:")
	fmt.Printf("Red: 0x%08X\n", colorField.Encode(Red))
	fmt.Printf("Green: 0x%08X\n", colorField.Encode(Green))
	fmt.Printf("Magenta: 0x%08X\n", colorField.Encode(Magenta))

	// Even undefined enum values work if they fit in the field
	var undefinedColor Color = 6 // This is Cyan in our enum, but imagine it's undefined
	fmt.Printf("Undefined color (6): 0x%08X\n", colorField.Encode(undefinedColor))

	// Output:
	// Setting valid colors:
	// Red: 0x00000000
	// Green: 0x00000001
	// Magenta: 0x00000007
	// Undefined color (6): 0x00000006
}

func TestBitField_PanicsOnInvalidValue(t *testing.T) {
	// Create a 2-bit field
	bf := New[uint8, uint32](0, 2)

	// This should work fine
	_ = bf.Encode(3)

	// This should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected Encode to panic on invalid value, but it didn't")
		}
	}()

	_ = bf.Encode(4) // Should panic
}

func TestBitField_EnumValidation(t *testing.T) {
	// Create a bitfield for the Priority enum (needs 2 bits)
	priorityField := New[Priority, uint32](0, 2)

	// Valid values: 0-3 (Low, Medium, High, Critical)
	validValues := []Priority{Low, Medium, High, Critical}
	for _, val := range validValues {
		if !priorityField.IsValid(val) {
			t.Errorf("Expected Priority %d to be valid for 2-bit field", val)
		}
	}

	// Create an out-of-range priority value
	var invalidPriority Priority = 4 // This is outside our enum range
	if priorityField.IsValid(invalidPriority) {
		t.Errorf("Expected Priority %d to be invalid for 2-bit field", invalidPriority)
	}
}

func ExampleBitField_typeSafety() {
	// Demonstrate type safety with different container and value types

	// 8-bit value in a uint32 container
	bf1 := New[uint8, uint32](0, 8)
	container32 := bf1.Encode(255)
	fmt.Printf("8-bit value in uint32: 0x%08X\n", container32)

	// 16-bit value in a uint64 container
	bf2 := New[uint16, uint64](16, 16)
	container64 := bf2.Encode(65535)
	fmt.Printf("16-bit value in uint64: 0x%016X\n", container64)

	// Create a bit field that uses uint for both parameters
	bf3 := New[uint, uint](4, 4)
	containerUint := bf3.Encode(10)
	fmt.Printf("4-bit value in uint: %d\n", containerUint)

	// Output:
	// 8-bit value in uint32: 0x000000FF
	// 16-bit value in uint64: 0x00000000FFFF0000
	// 4-bit value in uint: 160
}

func ExampleBitField_deviceRegister() {
	// Real-world example: Device control register using enums

	// Define device-related enum types
	type DeviceMode uint8
	const (
		Standby DeviceMode = iota
		LowPower
		Normal
		Performance
	)

	type ErrorHandling uint8
	const (
		Ignore ErrorHandling = iota
		Report
		RetryOnce
		RetryMultiple
		Abort
	)

	type InterruptMode uint8
	const (
		Disabled InterruptMode = iota
		Edge
		Level
		Both
	)

	// Create bit fields for a device control register
	var deviceCtrl uint32 = 0

	// Define fields for the register
	powerModeField := New[DeviceMode, uint32](0, 2)              // 2 bits for power mode
	errorHandlingField := Next[ErrorHandling](powerModeField, 3) // 3 bits for error handling
	intModeField := Next[InterruptMode](errorHandlingField, 2)   // 2 bits for interrupt mode
	enabledField := Next[uint8](intModeField, 1)                 // 1 bit for device enabled

	// Configure device
	deviceCtrl = powerModeField.Update(deviceCtrl, Normal)
	deviceCtrl = errorHandlingField.Update(deviceCtrl, RetryOnce)
	deviceCtrl = intModeField.Update(deviceCtrl, Edge)
	deviceCtrl = enabledField.Update(deviceCtrl, 1) // Enable device

	fmt.Printf("Device control register: 0x%08X\n", deviceCtrl)

	// Read current configuration
	powerMode := powerModeField.Decode(deviceCtrl)
	errorMode := errorHandlingField.Decode(deviceCtrl)
	intMode := intModeField.Decode(deviceCtrl)
	enabled := enabledField.Decode(deviceCtrl)

	fmt.Printf("Power mode: %d\n", powerMode)
	fmt.Printf("Error handling: %d\n", errorMode)
	fmt.Printf("Interrupt mode: %d\n", intMode)
	fmt.Printf("Enabled: %d\n", enabled)

	// Toggle device on/off while preserving other settings
	deviceCtrl = enabledField.Update(deviceCtrl, 0) // Disable
	fmt.Printf("After disable: 0x%08X\n", deviceCtrl)

	deviceCtrl = enabledField.Update(deviceCtrl, 1) // Enable again
	fmt.Printf("After re-enable: 0x%08X\n", deviceCtrl)

	// Change to low power mode
	deviceCtrl = powerModeField.Update(deviceCtrl, LowPower)
	fmt.Printf("After power mode change: 0x%08X\n", deviceCtrl)

	// Output:
	// Device control register: 0x000000AA
	// Power mode: 2
	// Error handling: 2
	// Interrupt mode: 1
	// Enabled: 1
	// After disable: 0x0000002A
	// After re-enable: 0x000000AA
	// After power mode change: 0x000000A9
}
