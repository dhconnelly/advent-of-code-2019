package main

import (
	"strings"
	"testing"
)

func TestDayTestCases(t *testing.T) {
	for _, d := range days {
		for i, part := range d.parts {
			for _, tc := range part.testCases {
				r := strings.NewReader(tc.input)
				got := part.solve(r)
				if got != tc.output {
					t.Errorf("%s part%d: got %s, want %s", d.name, i+1, got, tc.output)
				}
			}
		}
	}
}
