package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/intcode"
)

type game struct {
	in  chan<- int64
	out <-chan int64
	r   *bufio.Scanner
}

func NewGame(data []int64, r io.Reader) game {
	in := make(chan int64)
	out := intcode.RunProgram(data, in)
	return game{in, out, bufio.NewScanner(r)}
}

func (g game) readLine() (string, bool) {
	var b []byte
	var ok bool
	var c int64
	for c, ok = <-g.out; ok && c != '\n'; c = <-g.out {
		b = append(b, byte(c))
	}
	if len(b) > 0 {
		return string(b), true
	}
	return "", ok
}

func (g game) getCommand() (string, bool) {
	if g.r.Scan() {
		return g.r.Text(), true
	}
	if err := g.r.Err(); err != nil {
		log.Fatalf("failed to read command: %s", err)
	}
	return "", false
}

func (g game) writeLine(line string) {
	for _, c := range line {
		g.in <- int64(c)
	}
	g.in <- '\n'
}

func (g game) loop() {
	for {
		line, ok := g.readLine()
		if !ok {
			fmt.Println("machine halted; exiting")
			return
		}
		if line != prompt {
			fmt.Println(line)
			continue
		}
		fmt.Println(prompt)
		cmd, ok := g.getCommand()
		if !ok {
			fmt.Println("no more commands; exiting")
			return
		}
		g.writeLine(cmd)
	}
}

const (
	prompt = "Command?"
)

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	g := NewGame(data, os.Stdin)
	g.loop()
}
