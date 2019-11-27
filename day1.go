package main

import (
	"io"
	"io/ioutil"
)

func solveDay1Part1(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func solveDay1Part2(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(b)
}

var day1 = day{
	name: "day1",
	parts: []part{
		{solve: solveDay1Part1, testCases: []testCase{
			{"1", "1"},
			{"2", "2"},
		}},
		{solve: solveDay1Part2, testCases: []testCase{
			{"3", "3"},
			{"4", "4"},
		}},
	},
}
