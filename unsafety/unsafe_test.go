package unsafety

import (
	"testing"
)

func assert(t *testing.T, b bool, s string) {

	if !b {
		t.Error("Assertion failure:", s)
	}

}

func TestBto8(t *testing.T) {

	assert(t, BtoU8(true) == 1, "BtoU8(true) == 1")
	assert(t, BtoU8(false) == 0, "BtoU8(false) == 0")
	assert(t, BtoS8(true) == 1, "BtoS8(true) == 1")
	assert(t, BtoS8(false) == 0, "BtoS8(false) == 0")

}

func TestReturnStringTrue(t *testing.T) {

	assert(t, ReturnStringTrue("asdf", false) == "", "ReturnStringTrue(\"asdf\", false) == \"\"")
	assert(t, ReturnStringTrue("asdf", true) == "asdf", "ReturnStringTrue(\"asdf\", true) == \"asdf\"")

}
