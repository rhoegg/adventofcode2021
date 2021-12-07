package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func totalError(alignPos int, positions []int) int {
	var error = 0
	for _, p := range positions {
		error += int(math.Abs(float64(p - alignPos)))
	}
	return error
}

func totalTriangularFuel(alignPos int, positions []int) int {
	var fuel = 0
	for _, p := range positions {
		steps := int(math.Abs(float64(p - alignPos)))
		fuel += (steps * (steps + 1) / 2)
	}
	return fuel
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	inputPositions := strings.Split(string(data), ",")
	fmt.Printf("crab positions: %d\n", len(inputPositions))
	var positions []int
	minP, _ := strconv.Atoi(inputPositions[0])
	maxP, _ := strconv.Atoi(inputPositions[0])
	for _, inputPosition := range inputPositions {
		p, _ := strconv.Atoi(inputPosition)
		if p < minP {
			minP = p
		}
		if p > maxP {
			maxP = p
		}
		positions = append(positions, p)
	}
	startPos := int((maxP - minP) / 2)
	lastFuel := totalTriangularFuel(startPos, positions)
	fmt.Printf("Error for %d = %d\n", startPos, lastFuel)

	//step := 16
	direction := 1
	pos := startPos + direction // + step
	fuel := totalTriangularFuel(pos, positions)
	fmt.Printf("Error for %d = %d\n", pos, fuel)
	if (fuel > lastFuel) { // worse
		direction = -1
		fuel, lastFuel = lastFuel, fuel
		pos = pos - 1
	}

	for (fuel < lastFuel) {
		pos = pos + direction
		lastFuel, fuel = fuel, totalTriangularFuel(pos, positions)
	}
	fmt.Printf("Error for %d = %d\n", pos - direction, lastFuel)
}
