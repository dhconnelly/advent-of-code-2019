package main

import (
	"io"
	"io/ioutil"
)

func solveDay2Part1(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func solveDay2Part2(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(b)
}

var day2 = day{
	name: "day2",
	parts: []part{
		{solve: solveDay2Part1, testCases: []testCase{
			{"1", "1"},
			{"2", "2"},
		}},
		{solve: solveDay2Part2, testCases: []testCase{
			{"3", "3"},
			{"4", "4"},
		}},
	},
}
