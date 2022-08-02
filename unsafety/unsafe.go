package unsafety

import (
	"unsafe"
)

// Returns a binary representation of a bool as if it were a uint8
func BtoU8(b bool) uint8 {
	return *(*uint8)(unsafe.Pointer(&b))
}

// Returns a binary representation of a bool as if it were a int8
func BtoS8(b bool) int8 {
	return *(*int8)(unsafe.Pointer(&b))
}

// Internal representation of a string
type stringHeader struct {
	Data unsafe.Pointer
	Len  int
}

// A branchless method to return a string if true, else return a string with len(0)
func ReturnStringTrue(s string, b bool) string {
	((*stringHeader)(unsafe.Pointer(&s))).Len *= int(BtoU8(b))
	return s
}
