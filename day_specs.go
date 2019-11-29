package main

import "io"

type testCase struct {
	input, output string
}

type part struct {
	solve     func(io.Reader) string
	testCases []testCase
}

type day struct {
	name  string
	parts []part
}

var days = []day{{
	name: "day1",
	parts: []part{
		{solve: solveDay1Part1, testCases: []testCase{
			{"1", "1"},
			{"2", "2"}}},
		{solve: solveDay1Part2, testCases: []testCase{
			{"3", "3"},
			{"4", "4"}}},
	}},
}
