package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

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

var days = []day{
	day1,
}

func openOrDie(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func solveDay(d day, inputDir string) {
	fmt.Println(d.name)
	path := filepath.Join(inputDir, d.name)
	f := openOrDie(path)
	defer f.Close()
	for i, part := range d.parts {
		fmt.Printf("> part%d: %q\n", i+1, part.solve(f))
		f.Seek(0, 0)
	}
	fmt.Println()
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		log.Fatal("Usage: advent-of-code-2019 input_dir [day_name]")
	}

	fmt.Println("=================================")
	fmt.Println("🎄🎄🎄 Advent of Code 2019 🎄🎄🎄")
	fmt.Println("=================================")
	fmt.Println()

	inputDir := os.Args[1]
	fmt.Printf("Reading input from directory %s\n", inputDir)

	var whichDay string
	if len(os.Args) == 3 {
		whichDay = os.Args[2]
		fmt.Printf("Running %s\n", whichDay)
	} else {
		fmt.Println("Running all days")
	}
	fmt.Println()

	for _, d := range days {
		if whichDay == "" || whichDay == d.name {
			solveDay(d, inputDir)
		}
	}
}
