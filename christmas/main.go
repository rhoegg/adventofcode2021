package main

import (
	"fmt"
	"os"
	"strings"
)


func main() {
	seaFloor := parseInput("input.txt")
	steps := 0
	for seaFloor.Step() > 0 {
		steps++
		fmt.Printf("Step %d:\n%v\n", steps, seaFloor)
	}
	steps++
	fmt.Printf("Stopped moving after step %d:\n%v\n", steps, seaFloor)
}

type SeaCucumber struct {
	X, Y int
}

type SeaCucumberHerd struct {
	Cucumbers []*SeaCucumber
}

func (h *SeaCucumberHerd) AddCucumber(x, y int) {
	h.Cucumbers = append(h.Cucumbers, &SeaCucumber{X: x, Y: y})
}

type SeaFloor struct {
	Width, Height int
	EastFacing *SeaCucumberHerd
	SouthFacing *SeaCucumberHerd
}

func (f SeaFloor) Occupied(x, y int) bool {
	for _, herd := range []*SeaCucumberHerd{f.EastFacing, f.SouthFacing} {
		for _, c := range herd.Cucumbers {
			if c.X == x && c.Y == y {
				return true
			}
		}
	}
	return false
}

func (f SeaFloor) Step() int {
	moveCount := 0
	var willMove []*SeaCucumber
	for _, c := range f.EastFacing.Cucumbers {
		if ! f.Occupied((c.X + 1) % f.Width, c.Y) {
			willMove = append(willMove, c)
		}
	}
	for _, mover := range willMove {
		mover.X = (mover.X + 1) % f.Width
		moveCount++
	}
	willMove = nil
	for _, c := range f.SouthFacing.Cucumbers {
		if ! f.Occupied(c.X, (c.Y + 1) % f.Height) {
			willMove = append(willMove, c)
		}
	}
	for _, mover := range willMove {
		mover.Y = (mover.Y + 1) % f.Height
		moveCount++
	}
	return moveCount
}

func (f SeaFloor) String() string {
	var floor [][]string
	for y := 0; y < f.Height; y++ {
		floor = append(floor, []string{})
		for x := 0; x < f.Width; x++ {
			floor[len(floor)-1] = append(floor[len(floor)-1], ".")
		}
	}
	for _, eCucumber := range f.EastFacing.Cucumbers {
		floor[eCucumber.Y][eCucumber.X] = ">"
	}
	for _, sCucumber := range f.SouthFacing.Cucumbers {
		floor[sCucumber.Y][sCucumber.X] = "v"
	}
	var lines []string
	for _, floorRow := range floor {
		lines = append(lines, strings.Join(floorRow, ""))
	}
	return strings.Join(lines, "\n")
}

func parseInput(filename string) SeaFloor {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	inputLines := strings.Split(string(data), "\n")
	seaFloor := SeaFloor{
		Width: len(inputLines[0]),
		Height: len(inputLines),
		EastFacing: &SeaCucumberHerd{},
		SouthFacing: &SeaCucumberHerd{},
	}

	for y, inputLine := range inputLines {
		for x, inputCell := range inputLine {
			switch inputCell {
			case '>':
				seaFloor.EastFacing.AddCucumber(x, y)
			case 'v':
				seaFloor.SouthFacing.AddCucumber(x, y)
			}
		}
	}

	return seaFloor
}