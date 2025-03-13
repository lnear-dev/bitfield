package bitfield

import (
	"fmt"
	"testing"
)

func TestNewBitField(t *testing.T) {
	tests := []struct {
		name    string
		shift   uint
		size    uint
		wantErr bool
	}{
		{"valid field", 0, 3, false},
		{"max valid field", 29, 3, false},
		{"invalid shift", 64, 1, true},
		{"invalid size", 0, 65, true},
		{"invalid combined", 62, 3, true},
		{"zero size", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf, err := Safe[uint64, uint64](tt.shift, tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("Safe(%v, %v): err = %v, want err = %v",
					tt.shift, tt.size, err, tt.wantErr)
			}
			if !tt.wantErr {
				if bf.Shift != tt.shift {
					t.Errorf("shift = %v, want %v", bf.Shift, tt.shift)
				}
				if bf.Size != tt.size {
					t.Errorf("size = %v, want %v", bf.Size, tt.size)
				}
			}
		})
	}
}

func TestBitField_IsValid(t *testing.T) {
	bf := New[uint8, uint32](0, 3)
	tests := []struct {
		value uint8
		want  bool
	}{
		{0, true},
		{7, true},
		{8, false},
		{255, false},
	}

	for _, tt := range tests {
		if got := bf.IsValid(tt.value); got != tt.want {
			t.Errorf("IsValid(%v) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestBitField_Encode(t *testing.T) {
	bf := New[uint8, uint32](2, 3)
	tests := []struct {
		value     uint8
		want      uint32
		wantPanic bool
	}{
		{0, 0, false},
		{7, 28, false},
		{8, 0, true},
	}

	for _, tt := range tests {
		func() {
			defer func() {
				panicked := recover() != nil
				if panicked != tt.wantPanic {
					t.Errorf("Encode(%v): panic = %v, want panic = %v",
						tt.value, panicked, tt.wantPanic)
				}
			}()
			if got := bf.Encode(tt.value); !tt.wantPanic && got != tt.want {
				t.Errorf("Encode(%v) = %v, want %v", tt.value, got, tt.want)
			}
		}()
	}
}

func TestBitField_Update(t *testing.T) {
	bf := New[uint8, uint32](2, 3)
	tests := []struct {
		previous uint32
		value    uint8
		want     uint32
	}{
		{0xFFFFFFFF, 0, 0xFFFFFFE3},
		{0, 7, 28},
		{0xFFFFFFFF, 3, 0xFFFFFFEF},
	}

	for _, tt := range tests {
		if got := bf.Update(tt.previous, tt.value); got != tt.want {
			t.Errorf("Update(%v, %v) = %v, want %v",
				tt.previous, tt.value, got, tt.want)
		}
	}
}

func TestBitField_Decode(t *testing.T) {
	bf := New[uint8, uint32](2, 3)
	tests := []struct {
		value uint32
		want  uint8
	}{
		{0, 0},
		{28, 7},
		{0xFFFFFFFF, 7},
	}

	for _, tt := range tests {
		if got := bf.Decode(tt.value); got != tt.want {
			t.Errorf("Decode(%v) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestBitField_NextBitField(t *testing.T) {
	bf := New[uint8, uint32](0, 3)
	tests := []struct {
		size      uint
		wantShift uint
		wantPanic bool
	}{
		{3, 3, false},
		{5, 3, false},
		{6, 0, true}, // Would exceed uint8
	}

	for _, tt := range tests {
		t.Run(
			// "size="+string(tt.size),
			fmt.Sprintf("size=%d", tt.size),
			func(t *testing.T) {
				defer func() {
					panicked := recover() != nil
					if panicked != tt.wantPanic {
						t.Errorf("NextBitField(%v): panic = %v, want panic = %v",
							tt.size, panicked, tt.wantPanic)
					}
				}()
				next := bf.NextBitField(tt.size)
				if !tt.wantPanic && next.Shift != tt.wantShift {
					t.Errorf("NextBitField(%v).Shift = %v, want %v",
						tt.size, next.Shift, tt.wantShift)
				}
			})
	}
}

func TestBitField_Clear(t *testing.T) {
	bf := New[uint8, uint32](2, 3)
	tests := []struct {
		value uint32
		want  uint32
	}{
		{0xFFFFFFFF, 0xFFFFFFE3},
		{28, 0},
		{0, 0},
	}

	for _, tt := range tests {
		if got := bf.Clear(tt.value); got != tt.want {
			t.Errorf("Clear(%v) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestBitField_HasValue(t *testing.T) {
	bf := New[uint8, uint32](2, 3)
	tests := []struct {
		container uint32
		value     uint8
		want      bool
	}{
		{28, 7, true},
		{28, 6, false},
		{0, 0, true},
		{0xFFFFFFFF, 7, true},
	}

	for _, tt := range tests {
		if got := bf.HasValue(tt.container, tt.value); got != tt.want {
			t.Errorf("HasValue(%v, %v) = %v, want %v",
				tt.container, tt.value, got, tt.want)
		}
	}
}
