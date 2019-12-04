package ints

import "testing"

func TestAbs(t *testing.T) {
	for _, tc := range []struct {
		x, abs int
	}{
		{0, 0},
		{-7, 7},
		{19, 19},
	} {
		if abs := Abs(tc.x); abs != tc.abs {
			t.Errorf("Abs(%d) = %d, want %d", tc.x, abs, tc.abs)
		}
	}
}
