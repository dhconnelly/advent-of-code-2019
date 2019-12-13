package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/breakout"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
)

func countTiles(
	state breakout.GameState,
	which breakout.TileId,
) int {
	n := 0
	for _, tile := range state.Tiles {
		if tile == which {
			n++
		}
	}
	return n
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	state, err := breakout.Play(data, nil, 1, breakout.NEUTRAL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(countTiles(state, breakout.BLOCK))

	data[0] = 2 // play for free
	state, err = breakout.Play(data, nil, 1, breakout.LEFT)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(state.Score)
}
