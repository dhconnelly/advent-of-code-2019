package main

import (
	"github.com/dhconnelly/advent-of-code-2019/breakout"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"github.com/gdamore/tcell"
	"log"
	"os"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err = screen.Init(); err != nil {
		log.Fatal(err)
	}
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	data[0] = 2 // play for free
	_, err = breakout.Play(data, screen)
	if err != nil {
		log.Fatal(err)
	}
	screen.Fini()
}
