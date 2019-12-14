package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"os"
)

type tileColor int

const (
	BLACK       tileColor = 0
	WHITE       tileColor = 1
	TRANSPARENT tileColor = 2
)

const (
	HEIGHT = 6
	WIDTH  = 25
)

type layer [HEIGHT * WIDTH]tileColor

func readImage(r io.Reader) []layer {
	br := bufio.NewReader(r)
	var layers []layer
	for {
		var l layer
		for i := 0; i < HEIGHT*WIDTH; i++ {
			b, err := br.ReadByte()
			if err == io.EOF {
				return layers
			}
			if err != nil {
				log.Fatalf("bad color at row,height = %d,%d: %s", err)
			}
			l[i] = tileColor(b - '0')
		}
		layers = append(layers, l)
	}
	return layers
}

const (
	pixelSize = 5
)

var (
	layerRect    = image.Rect(0, 0, WIDTH*pixelSize, HEIGHT*pixelSize)
	layerPalette = []color.Color{
		color.Black,
		color.White,
		color.Transparent,
	}
)

var toColor = map[tileColor]color.Color{
	BLACK:       color.Black,
	WHITE:       color.White,
	TRANSPARENT: color.Transparent,
}

func toPaletted(from *layer) *image.Paletted {
	p := image.NewPaletted(layerRect, layerPalette)
	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			p.Set(j, i, toColor[from[i*WIDTH+j]])
		}
	}
	return p
}

func apply(to, from *layer) {
	for i := 0; i < HEIGHT*WIDTH; i++ {
		if c := from[i]; c != TRANSPARENT {
			to[i] = c
		}
	}
}

func flattenLayers(ls []layer) {
	prev := &ls[len(ls)-1]
	for i := len(ls) - 2; i >= 0; i-- {
		next := &ls[i]
		apply(prev, next)
		next = prev
	}
}

func printLayer(l *layer) {
	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			switch c := l[i*WIDTH+j]; c {
			case BLACK:
				fmt.Print(" ")
			default:
				fmt.Printf("%d", c)
			}
		}
		fmt.Println()
	}
}

func main() {
	in, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	out, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	m := readImage(in)
	flattenLayers(m)
	final := &m[len(m)-1]
	if err := gif.Encode(out, toPaletted(final), nil); err != nil {
		log.Fatal(err)
	}
	printLayer(final)
	//g := toGif(img)
	//if err := gif.EncodeAll(out, g); err != nil {
	//log.Fatal(err)
	//}
}
