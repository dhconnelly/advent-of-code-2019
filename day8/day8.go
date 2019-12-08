package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type image struct {
	width, height int
	layers        [][]byte
}

func readImage(path string, width, height int) image {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img := image{width, height, [][]byte{{}}}
	b := make([]byte, 1)
	for {
		n, err := f.Read(b)
		if n == 1 && b[0] != '\n' {
			if len(img.layers[len(img.layers)-1]) == width*height {
				img.layers = append(img.layers, []byte{})
			}
			i := len(img.layers) - 1
			img.layers[i] = append(img.layers[i], b[0]-'0')
		} else if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
	return img
}

func zeroes(pixels []byte) int {
	n := 0
	for _, p := range pixels {
		if p == 0 {
			n++
		}
	}
	return n
}

func layerChecksum(pixels []byte) int {
	ones, twos := 0, 0
	for _, p := range pixels {
		switch p {
		case 1:
			ones++
		case 2:
			twos++
		}
	}
	return ones * twos
}

func checksum(img image) int {
	fewestLayer := img.layers[0]
	fewestZeroes := zeroes(fewestLayer)
	x := layerChecksum(fewestLayer)
	for _, layer := range img.layers[1:] {
		if z := zeroes(layer); z < fewestZeroes {
			fewestZeroes = z
			fewestLayer = layer
			x = layerChecksum(layer)
		}
	}
	return x
}

func apply(base, layer []byte) {
	for i := 0; i < len(base); i++ {
		if layerPix := layer[i]; layerPix != 2 {
			base[i] = layerPix
		}
	}
}

func decode(img image) []byte {
	b := make([]byte, img.width*img.height)
	for i := len(img.layers) - 1; i >= 0; i-- {
		apply(b, img.layers[i])
	}
	return b
}

func printImage(b []byte, width, height int) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			pix := b[i*width+j]
			if pix == 0 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("%d", b[i*width+j])
			}
		}
		fmt.Println()
	}
}

func main() {
	img := readImage(os.Args[1], 25, 6)
	fmt.Println(checksum(img))
	printImage(decode(img), 25, 6)
}
