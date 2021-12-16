package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Paths struct {
	edges map[string][]string
}

func NewPaths() Paths {
	return Paths{make(map[string][]string)}
}

func (p *Paths) AddCorridor(p1, p2 string) {
	if p1 != "end" {
		if p2 != "start" {
			p.edges[p1] = append(p.edges[p1], p2)
		}
	}
	if p2 != "end" {
		if p1 != "start" {
			p.edges[p2] = append(p.edges[p2], p1)
		}
	}
}

func (p *Paths) FromStart() [][]string {
	visited := make(map[string]struct{})
	visited["start"] = struct{}{}
	return p.From([]string{"start"}, visited, false)
}
func (p *Paths) From(path []string, smallRoomsVisited map[string]struct{}, mulligan bool) [][]string {
	currentCave := path[len(path) - 1]
	if currentCave == "end" {
		return [][]string{path}
	} else {
		var paths [][]string
		for _, next := range p.edges[currentCave] {
			//fmt.Printf(" - checking %s-%s\n", currentCave, next)
			_, visited := smallRoomsVisited[next]

			if (! visited) || (! mulligan) {
				var newPath = append(path, next)
				//fmt.Printf("computing %v\n", newPath)
				thisPathVisited := make(map[string]struct{})
				for k := range smallRoomsVisited {
					thisPathVisited[k] = struct{}{}
				}
				if next != "end" && unicode.IsLower(rune(next[0])) {
					thisPathVisited[next] = struct{}{}
				}
				paths = append(paths, p.From(newPath, thisPathVisited, visited || mulligan)...)
			}
		}
		return paths
	}
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(data))
	paths := NewPaths()
	for _, pathInput := range strings.Split(string(data), "\n") {
		path := strings.Split(pathInput, "-")
		paths.AddCorridor(path[0], path[1])
	}
	for source, e := range paths.edges {
		fmt.Printf("%s - %v\n", source, e)
	}
	validPaths := paths.FromStart()
	for _, p := range validPaths {
		fmt.Printf("path %v\n", p)
	}
	fmt.Printf("%d paths\n", len(validPaths))
}
