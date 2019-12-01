package main

import (
	"fmt"
	"io"
	"log"
	"os"
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

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sum, recSum := 0, 0
	var mass int
	for {
		if _, err := fmt.Fscanf(f, "%d\n", &mass); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		sum += fuelForMass(mass)
		recSum += fuelSum(mass)
	}
	fmt.Println(sum)
	fmt.Println(recSum)
}
