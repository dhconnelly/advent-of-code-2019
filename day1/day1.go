package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func fuelForMass(mass int) int {
	return (mass / 3) - 2
}

func fuelSum(mass int) int {
	sum := 0
	for {
		fuel := fuelForMass(mass)
		if fuel < 0 {
			break
		}
		sum += fuel
		mass = fuel
	}
	return sum
}

func solvePart1(masses <-chan int) int {
	sum := 0
	for mass := range masses {
		sum += fuelForMass(mass)
	}
	return sum
}

func solvePart2(masses <-chan int) int {
	sum := 0
	for mass := range masses {
		sum += fuelSum(mass)
	}
	return sum
}

func scanLines(r io.Reader) <-chan int {
	ch := make(chan int)
	go func() {
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			i, err := strconv.Atoi(scan.Text())
			if err != nil {
				log.Fatal(err)
			}
			ch <- i
		}
		if err := scan.Err(); err != nil {
			log.Fatal(err)
		}
		close(ch)
	}()
	return ch
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Println(solvePart1(scanLines(f)))
	f.Seek(0, 0)
	fmt.Println(solvePart2(scanLines(f)))
}
