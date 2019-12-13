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
	PADDLE: '_',
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
	Score int
	Tiles ScreenTiles
}

/*
func copied(tiles ScreenTiles) ScreenTiles {
	copied := ScreenTiles(make(map[geom.Pt2]TileId))
	for k, v := range tiles {
		copied[k] = v
	}
	return copied
}

func updateScreen(frames chan ScreenTiles) {
	for {
		frame, ok := <-frames
		if !ok {
			return
		}
		fmt.Println("drawing frame:", frame)
	}
}
*/

func draw(screen tcell.Screen, x, y int, tile TileId) {
	screen.SetContent(x, y, tileToRune[tile], nil, 0)
	screen.Show()
}

func erase(screen tcell.Screen) {
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
) (GameState, error) {
	in := make(chan int64)
	out := intcode.RunProgram(data, in)

	state := GameState{
		Tiles: ScreenTiles(make(map[geom.Pt2]TileId)),
	}
	joystick := NEUTRAL

	if screen != nil {
		screen.Clear()
	}
loop:
	for _ = range time.Tick(33000000) {
		select {
		case x, ok := <-out:
			if !ok {
				break loop
			}
			y, z := <-out, <-out
			if x == -1 && y == 0 {
				state.Score = int(z)
			} else {
				tile := TileId(z)
				state.Tiles[geom.Pt2{int(x), int(y)}] = tile
				if screen != nil {
					draw(screen, int(x), int(y), tile)
				}
			}
		case e := <-readEvents(screen):
			switch e.Key() {
			case tcell.KeyCtrlC:
				break loop
			case tcell.KeyLeft:
				joystick = LEFT
			case tcell.KeyRight:
				joystick = RIGHT
			}
		case in <- int64(joystick):
		}
	}
	close(in)
	return state, nil
}
