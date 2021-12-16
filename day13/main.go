package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	//"strings"
)

type Point struct {
	x, y int
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

func FoldPaperX(paper map[Point]struct{}, crease int) map[Point]struct{} {
	foldedPaper := make(map[Point]struct{})
	for p, _ := range paper {
		newPoint := Point{p.x, p.y}
		if p.x > crease {
			newPoint.x = 2 * crease - p.x
		}
		foldedPaper[newPoint] = struct{}{}
	}
	return foldedPaper
}

func FoldPaperY(paper map[Point]struct{}, crease int) map[Point]struct{} {
	foldedPaper := make(map[Point]struct{})
	for p, _ := range paper {
		newPoint := Point{p.x, p.y}
		if p.y > crease {
			newPoint.y = 2 * crease - p.y
		}
		foldedPaper[newPoint] = struct{}{}
	}
	return foldedPaper
}


func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(data))
	inputSections := strings.Split(string(data), "\n\n")

	originalDotsInput := strings.Split(inputSections[0], "\n")
	pointMap := make(map[Point]struct{})
	for _, dotInput := range originalDotsInput {
		coords := strings.Split(dotInput, ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		pointMap[Point{x, y}] = struct{}{}
	}


	fmt.Printf("original dots %d\n", len(pointMap))

	var paper = pointMap
	foldInstructions := strings.Split(inputSections[1], "\n")
	for _, instruction := range foldInstructions {
		tokens := strings.Split(instruction, " ")
		foldParams := strings.Split(tokens[2], "=")
		creasePosition, _ := strconv.Atoi(foldParams[1])
		if foldParams[0] == "x" {
			paper = FoldPaperX(paper, creasePosition)
		} else {
			paper = FoldPaperY(paper, creasePosition)
		}
	}

	var xMax, yMax int
	for p := range paper {
		if p.x > xMax {
			xMax = p.x
		}
		if p.y > yMax {
			yMax = p.y
		}
	}

	for y := 0; y <= yMax; y++ {
		for x := 0; x <= xMax; x++ {
			_, ok := paper[Point{x, y}]
			if ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
