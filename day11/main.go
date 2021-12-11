package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func step(grid [][]int) ([][]int, [][]bool) {
	var flashes [][]bool
	for _, row := range grid {
		var flashRow []bool
		for x, level := range row {
			row[x] = level + 1
			flashRow = append(flashRow, false)
		}
		flashes = append(flashes, flashRow)
	}
	grid, flashes = flash(grid, flashes)
	// zero flashes
	for y, row := range flashes {
		for x, flash := range row {
			if flash {
				grid[y][x] = 0
			}
		}
	}
	return grid, flashes
}

func flash(grid [][]int, flashes [][]bool) ([][]int, [][]bool) {
	// new flashes
	var newFlashes [][]bool
	for _, row := range flashes {
		var newRow []bool
		for _ = range row {
			newRow = append(newRow, false)
		}
		newFlashes = append(newFlashes, newRow)
	}

	for y, row := range grid {
		for x, level := range row {
			if ! flashes[y][x] { // can't flash again
				if level > 9 {
					newFlashes[y][x] = true
				}
			}
		}
	}

	// excite neighbors
	flashCount := 0
	for y, row := range newFlashes {
		for x, flash := range row {
			if flash {
				flashCount++
				for gridx := x - 1; gridx <= x + 1; gridx++ {
					for gridy := y - 1; gridy <= y + 1; gridy++ {
						if inBounds(gridx, gridy) {
							grid[gridy][gridx] += 1
						}
					}
				}
			}
		}
	}

	flashes = matrixOr(flashes, newFlashes)

	if flashCount > 0 {
		// recurse
		grid, flashes = flash(grid, flashes)
	}
	return grid, flashes
}

func inBounds(x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}
	if x > 9 || y > 9 {
		return false
	}
	return true
}

func matrixOr(left, right [][]bool) [][]bool {
	for y, row := range left {
		for x := range row {
			row[x] = left[y][x] || right[y][x]
		}
	}
	return left
}


func makeGrid(data string) [][]int {
	var octopusGrid [][]int
	rows := strings.Split(data, "\n")
	for _, row := range rows {
		var levels []int
		for _, c := range row {
			level, _ := strconv.Atoi(string(c))
			levels = append(levels, level)
		}
		octopusGrid = append(octopusGrid, levels)
	}
	return octopusGrid
}

func countFlashes(flashes [][]bool) int {
	count := 0
	for _, row := range flashes {
		for _, flash := range row {
			if flash {
				count++
			}
		}
	}
	return count
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	octopusGrid := makeGrid(string(data))

	totalFlashes := 0

	for i := 0; i < 100; i++ {
		grid, flashes := step(octopusGrid)
		octopusGrid = grid
		totalFlashes += countFlashes(flashes)
	}

	fmt.Println(octopusGrid)
	fmt.Printf("Total flashes: %d\n", totalFlashes)

	// puzzle 2
	octopusGrid = makeGrid(string(data))

	flashesThisStep := 0
	for stepNum := 1 	; flashesThisStep < 100; stepNum++ {
		grid, flashes := step(octopusGrid)
		octopusGrid = grid
		flashesThisStep = countFlashes(flashes)
		fmt.Printf("Step %d: flashes = %d\n", stepNum, flashesThisStep)
	}
}
