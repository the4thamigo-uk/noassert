package testlib

import (
	"testing"
)

func testError(t *testing.T) {
	if false {
		t.Error("My Error")
	}
}

func TestF1(t *testing.T) {
	f()
}

func TestF2(t *testing.T) {
	f()
	if false {
		t.Error("My Error")
	}
}

func TestF3(t *testing.T) {
	f()
	testError(t)
}
