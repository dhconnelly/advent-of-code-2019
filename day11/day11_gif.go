package main

import (
	"fmt"
	"github.com/dhconnelly/advent-of-code-2019/geom"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"github.com/dhconnelly/advent-of-code-2019/ints"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"os"
)

type tileColor int64

const (
	BLACK tileColor = iota
	WHITE
)

type direction int64

const (
	TURN_LEFT direction = iota
	TURN_RIGHT
)

type orientation rune

const (
	LEFT orientation = iota
	UP
	RIGHT
	DOWN
)

func turn(cur orientation, dir direction) orientation {
	switch dir {
	case TURN_LEFT:
		if cur == LEFT {
			return DOWN
		}
		return cur - 1
	case TURN_RIGHT:
		return (cur + 1) % 4
	}
	log.Fatalf("bad dir: %d", dir)
	return -1
}

var diffs = map[orientation]geom.Pt2{
	LEFT:  geom.Pt2{-1, 0},
	UP:    geom.Pt2{0, 1},
	RIGHT: geom.Pt2{1, 0},
	DOWN:  geom.Pt2{0, -1},
}

func move(cur geom.Pt2, o orientation) geom.Pt2 {
	return cur.Add(diffs[o])
}

const (
	tileHeight = 10
	tileWidth  = 10
)

type grid map[geom.Pt2]tileColor

func (g grid) ColorModel() color.Model {
	return color.RGBAModel
}

func (g grid) Bounds() image.Rectangle {
	minX, minY := math.MaxInt64, math.MaxInt64
	maxX, maxY := math.MinInt64, math.MinInt64
	for p, _ := range g {
		minX, minY = ints.Min(minX, p.X), ints.Min(minY, p.Y)
		maxX, maxY = ints.Max(maxX, p.X), ints.Max(maxY, p.Y)
	}
	return image.Rect(
		0, 0,
		int((maxX-minX+1)*tileWidth),
		int((maxY-minY+1)*tileHeight),
	)
}

var toColor = map[tileColor]color.Color{
	BLACK: color.Black,
	WHITE: color.White,
}

func (g grid) At(x, y int) color.Color {
	p := geom.Pt2{(x / tileWidth), -(y / tileHeight)}
	return toColor[g[p]]
}

func run(data []int64, initial tileColor) grid {
	in := make(chan int64)
	out := intcode.RunProgram(data, in)
	g := grid(make(map[geom.Pt2]tileColor))
	p := geom.Zero2
	g[p] = initial
	o := UP
loop:
	for {
		select {
		case c, ok := <-out:
			if !ok {
				break loop
			}
			g[p] = tileColor(c)
			dir := direction(<-out)
			o = turn(o, dir)
			p = move(p, o)
		case in <- int64(g[p]):
		}
	}
	return g
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	out, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	g := run(data, BLACK)
	fmt.Println(len(g))
	g = run(data, WHITE)
	if err = gif.Encode(out, g, nil); err != nil {
		log.Fatal(err)
	}
}
