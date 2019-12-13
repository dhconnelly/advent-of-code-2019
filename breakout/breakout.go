package breakout

import (
	"time"

	"github.com/dhconnelly/advent-of-code-2019/geom"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"github.com/gdamore/tcell"
)

type TileId int

const (
	EMPTY  TileId = 0
	WALL   TileId = 1
	BLOCK  TileId = 2
	PADDLE TileId = 3
	BALL   TileId = 4
)

var tileToRune = map[TileId]rune{
	EMPTY:  ' ',
	WALL:   '@',
	BLOCK:  'X',
	PADDLE: '-',
	BALL:   'o',
}

type JoystickPos int

const (
	NEUTRAL JoystickPos = 0
	LEFT    JoystickPos = -1
	RIGHT   JoystickPos = 1
)

type ScreenTiles map[geom.Pt2]TileId

type GameState struct {
	Joystick JoystickPos
	Score    int
	Tiles    ScreenTiles
}

func draw(screen tcell.Screen, x, y int, tile TileId) {
	screen.SetContent(x, y, tileToRune[tile], nil, 0)
	screen.Show()
}

func readEvents(screen tcell.Screen) chan *tcell.EventKey {
	ch := make(chan *tcell.EventKey)
	go func() {
		for {
			event := screen.PollEvent()
			switch e := event.(type) {
			case *tcell.EventKey:
				ch <- e
			}
		}
	}()
	return ch
}

func Play(
	data []int64,
	screen tcell.Screen,
	frameDelay time.Duration,
	joystickInit JoystickPos,
) (GameState, error) {
	in := make(chan int64)
	defer close(in)
	out := intcode.RunProgram(data, in)
	var events chan *tcell.EventKey
	if screen != nil {
		screen.Clear()
		events = readEvents(screen)
	}
	state := GameState{
		Tiles:    ScreenTiles(make(map[geom.Pt2]TileId)),
		Joystick: joystickInit,
	}
	tick := time.Tick(frameDelay)

loop:
	for {
		select {
		case e := <-events:
			switch e.Key() {
			case tcell.KeyCtrlC:
				break loop
			case tcell.KeyLeft:
				state.Joystick = LEFT
			case tcell.KeyRight:
				state.Joystick = RIGHT
			}

		case <-tick:
			select {
			case in <- int64(state.Joystick):
			default:
				state.Joystick = joystickInit
				continue
			}

		case x, ok := <-out:
			if !ok {
				break loop
			}

			y, z := <-out, <-out
			if x == -1 && y == 0 {
				if z > 0 {
					state.Score = int(z)
				}
			} else {
				tile := TileId(z)
				state.Tiles[geom.Pt2{int(x), int(y)}] = tile
				if screen != nil {
					draw(screen, int(x), int(y), tile)
				}
			}
		}
	}

	return state, nil
}
