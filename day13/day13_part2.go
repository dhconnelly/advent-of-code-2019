package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/breakout"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"github.com/gdamore/tcell"
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
	state, err := breakout.Play(data, screen)
	if err != nil {
		log.Fatal(err)
	}
	screen.Fini()
	fmt.Println(state.Score)
}
