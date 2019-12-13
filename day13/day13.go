package main

import (
	"fmt"
	"github.com/dhconnelly/advent-of-code-2019/geom"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"log"
	"os"
)

type tileId int

const (
	EMPTY  tileId = 0
	WALL   tileId = 1
	BLOCK  tileId = 2
	PADDLE tileId = 3
	BALL   tileId = 4
)

type joystickPos int

const (
	NEUTRAL joystickPos = 0
	LEFT    joystickPos = -1
	RIGHT   joystickPos = 1
)

type screenTiles map[geom.Pt2]tileId

type gameState struct {
	score int
	tiles screenTiles
}

func run(data []int64) gameState {
	in := make(chan int64)
	out := intcode.RunProgram(data, in)

	state := gameState{
		tiles: screenTiles(make(map[geom.Pt2]tileId)),
	}
	joystick := NEUTRAL
loop:
	for {
		select {
		case x, ok := <-out:
			if !ok {
				break loop
			}
			y, z := <-out, <-out
			if x == -1 && y == 0 {
				state.score = int(z)
			} else {
				state.tiles[geom.Pt2{int(x), int(y)}] = tileId(z)
			}
		case in <- int64(joystick):
		}
	}
	close(in)
	return state
}

func countTiles(screen screenTiles, which tileId) int {
	n := 0
	for _, tile := range screen {
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

	state := run(data)
	fmt.Println(countTiles(state.tiles, BLOCK))

	state = run(data)
	fmt.Println(state.score)
}
