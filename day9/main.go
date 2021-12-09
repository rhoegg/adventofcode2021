package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type MapPoint struct {
	x, y int
	height int
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	inputLines := strings.Split(string(data), "\n")
	var heightMap [][]int
	for _, inputLine := range inputLines {
		var row []int
		for _, inputPoint := range strings.Split(inputLine, "") {
			height, _ := strconv.Atoi(inputPoint)
			row = append(row, height)
		}
		heightMap = append(heightMap, row)
	}

	var lowPoints []MapPoint
	for y, row := range heightMap {
		for x, height := range row {
			if x > 0 && row[x-1] <= height {
				continue
			}
			if x < (len(row) - 1) && row[x+1] <= height {
				continue
			}
			if y > 0 && heightMap[y-1][x] <= height {
				continue
			}
			if y < (len(heightMap) - 1) && heightMap[y+1][x] <= height {
				continue
			}
			lowPoints = append(lowPoints, MapPoint{x: x, y: y, height: height})
			//fmt.Printf("Low point %v\n", lowPoints[len(lowPoints)- 1])

		}
	}
	risk := 0
	for _, lowPoint := range lowPoints {
		risk += lowPoint.height + 1
	}
	fmt.Printf("Total low point risk: %d\n", risk)

	basins := make(map[MapPoint][]MapPoint)

	for _, lowPoint := range lowPoints {
		basins[lowPoint] = basin([]MapPoint{lowPoint}, heightMap)
	}

	var sizes []int
	for b, points := range basins {
		fmt.Printf("Basin %v: %v\n", b, points)
		sizes = append(sizes, len(points))
	}
	sort.Ints(sizes)
	fmt.Println(sizes)
}

func basin(points []MapPoint, heightMap [][]int) []MapPoint {
	seen := make(map[MapPoint]struct{})
	for _, p := range points {
		seen[p] = struct{}{}
	}
	var expandedBasin []MapPoint
	expandedBasin = append(expandedBasin, points...)
	for _, p := range points {
		if p.x > 0 && heightMap[p.y][p.x - 1] >= p.height {
			leftPoint := MapPoint{x: p.x - 1, y: p.y, height: heightMap[p.y][p.x - 1]}
			_, visited := seen[leftPoint]
			if (! visited) && leftPoint.height < 9 {
//				fmt.Printf("expanding left %v\n", leftPoint)
				expandedBasin = append(expandedBasin, leftPoint)
				seen[leftPoint] = struct{}{}
			}
		}
		if p.x < (len(heightMap[0]) - 1) && heightMap[p.y][p.x + 1] >= p.height {
			rightPoint := MapPoint{x: p.x + 1, y: p.y, height: heightMap[p.y][p.x + 1]}
			_, visited := seen[rightPoint]
			if (! visited) && rightPoint.height < 9 {
//				fmt.Printf("expanding right %v\n", rightPoint)
				expandedBasin = append(expandedBasin, rightPoint)
				seen[rightPoint] = struct{}{}
			}
		}
		if p.y > 0 && heightMap[p.y - 1][p.x] >= p.height {
			upPoint := MapPoint{x: p.x, y: p.y - 1, height: heightMap[p.y - 1][p.x]}
			_, visited := seen[upPoint]
			if (! visited) && upPoint.height < 9 {
//				fmt.Printf("expanding up %v\n", upPoint)
				expandedBasin = append(expandedBasin, upPoint)
				seen[upPoint] = struct{}{}
			}
		}
		if p.y < (len(heightMap) - 1) && heightMap[p.y + 1][p.x] >= p.height {
			downPoint := MapPoint{x: p.x, y: p.y + 1, height: heightMap[p.y + 1][p.x]}
			_, visited := seen[downPoint]
			if (! visited) && downPoint.height < 9 {
//				fmt.Printf("expanding down %v\n", downPoint)
				expandedBasin = append(expandedBasin, downPoint)
				seen[downPoint] = struct{}{}
			}
		}
	}
	if len(expandedBasin) > len(points) {
//		fmt.Println("recurse")
		expandedBasin = basin(expandedBasin, heightMap) // recurse
	}
	return expandedBasin
}
