package sessions

import (
	"net/http"
	"time"
)

func ByteToHex(b byte) [2]byte {
	const hexvalues = "0123456789abcdef"
	return [2]byte{hexvalues[(b >> 4)], hexvalues[b&0xf]}
}

func SetAccessibilityCookie(w http.ResponseWriter, data string) {
	ac := &http.Cookie{
		Name:    "accessibility",
		Path:    "/",
		Value:   data,
		Expires: time.Now().Add(17520 * time.Hour),
	}

	http.SetCookie(w, ac)

}

// accessibility datatype
type AccessibilityData struct {
	EnableDyslexia   bool
	EnableHiContrast bool
	EnableNoImage    bool
}

type AccessibilityBitset uint8

// Bitset fields
// DO NOT MODIFY POSISITONS
const (
	EnableDyslexia AccessibilityBitset = 1 << iota
	EnableHiContrast
	EnableNoImage
)

func (A *AccessibilityBitset) Set(f AccessibilityBitset, b bool) {
	if b {
		(*A) |= f
	} else {
		(*A) &= ^f
	}
}

func (A AccessibilityBitset) Get(f AccessibilityBitset) bool {
	return A&f == f
}
