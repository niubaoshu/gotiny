package gotiny

import (
	"testing"
)

func TestUint32ToInt32(t *testing.T) {
	tests := []struct {
		name string
		u    uint32
		want int32
	}{
		{"max", 4294967295, -2147483648},
		{"9=-5", 9, -5},
		{"7=-4", 7, -4},
		{"5=-3", 5, -3},
		{"3=-2", 3, -2},
		{"1=-1", 1, -1},
		{"0=0", 0, 0},
		{"2=1", 2, 1},
		{"4=2", 4, 2},
		{"6=3", 6, 3},
		{"8=4", 8, 4},
		{"10=5", 10, 5},
		{"12=6", 12, 6},
		{"max", 4294967294, 2147483647},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uint32ToInt32(tt.u); got != tt.want {
				t.Errorf("uint32ToInt32(%v) = %v, want %v", tt.u, got, tt.want)
			}
		})
	}
}

func TestUint16ToInt16(t *testing.T) {
	tests := []struct {
		name string
		u    uint16
		want int16
	}{
		{"max", 65535, -32768},
		{"9=-5", 9, -5},
		{"7=-4", 7, -4},
		{"5=-3", 5, -3},
		{"3=-2", 3, -2},
		{"1=-1", 1, -1},
		{"0=0", 0, 0},
		{"2=1", 2, 1},
		{"4=2", 4, 2},
		{"6=3", 6, 3},
		{"8=4", 8, 4},
		{"10=5", 10, 5},
		{"12=6", 12, 6},
		{"max", 65530, 32765},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uint16ToInt16(tt.u); got != tt.want {
				t.Errorf("uint16ToInt16(%v) = %v, want %v", tt.u, got, tt.want)
			}
		})
	}
}
func TestReverse64Byte(t *testing.T) {
	tests := []struct {
		name string
		u    uint64
		want uint64
	}{
		{"all zeros", 0, 0},
		{"all ones", 0xffffffffffffffff, 0xffffffffffffffff},
		{"simple pattern", 0x1234567890abcdef, 0xefcdab9078563412},
		{"random pattern", 0x8765432109876543, 0x4365870921436587},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reverse64Byte(tt.u)
			if got != tt.want {
				t.Errorf("reverse64Byte(%x) = %x, want %x", tt.u, got, tt.want)
			}
		})
	}
}
func TestReverse32Byte(t *testing.T) {
	tests := []struct {
		name string
		u    uint32
		want uint32
	}{
		{"all zeros", 0, 0},
		{"all ones", 0xFFFFFFFF, 0xFFFFFFFF},
		{"simple pattern", 0x12345678, 0x78563412},
		{"random pattern", 0x87654321, 0x21436587},
		{"edge case 1", 0x80000000, 0x00000080},
		{"edge case 2", 0x7FFFFFFF, 0xFFFFFF7F},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reverse32Byte(tt.u)
			if got != tt.want {
				t.Errorf("reverse32Byte(%x) = %x, want %x", tt.u, got, tt.want)
			}
		})
	}
}

func TestInt64ToUint64(t *testing.T) {
	tests := []struct {
		name string
		want uint64
		v    int64
	}{
		{"max", 4294967295, -2147483648},
		{"9=-5", 9, -5},
		{"7=-4", 7, -4},
		{"5=-3", 5, -3},
		{"3=-2", 3, -2},
		{"1=-1", 1, -1},
		{"0=0", 0, 0},
		{"2=1", 2, 1},
		{"4=2", 4, 2},
		{"6=3", 6, 3},
		{"8=4", 8, 4},
		{"10=5", 10, 5},
		{"12=6", 12, 6},
		{"max", 4294967294, 2147483647},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := int64ToUint64(tt.v)
			if got != tt.want {
				t.Errorf("int64ToUint64(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}
