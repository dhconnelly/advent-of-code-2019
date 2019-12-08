package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type image struct {
	width, height int
	layers        [][]byte
}

func readImage(path string, width, height int) image {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	txt := strings.TrimSpace(string(b))
	numLayers := len(txt) / (width * height)
	img := image{width, height, make([][]byte, numLayers)}
	for i := 0; i < numLayers; i++ {
		layerBase := i * width * height
		for row := 0; row < height; row++ {
			for col := 0; col < width; col++ {
				pixel := txt[layerBase+row*width+col] - '0'
				img.layers[i] = append(img.layers[i], pixel)
			}
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
