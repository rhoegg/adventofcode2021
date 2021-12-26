package main

import (
	"container/heap"
	"fmt"
	"os"
	"strings"
)

func main() {
	data, err := os.ReadFile("input2.txt")
	if err != nil {
		panic(err)
	}

	burrow := buildBurrow(string(data))
	plan := CheapestPlan(burrow)
	fmt.Printf("Cheapest plan costs %d\n", plan.Spent)
}

func CheapestPlan(burrow []*Space) Plan {
	begin := Plan{Burrow: burrow, Spent: int64(0)}
	planQueue := make(PlanPriorityQueue, 1)
	planQueue[0] = &begin
	heap.Init(&planQueue)
	visited := make(map[string]struct{})
	for len(planQueue) > 0 {
		tryNext := heap.Pop(&planQueue).(*Plan)
		fmt.Printf("Checking plan costing %d %v\n", tryNext.Spent, tryNext)
		if tryNext.Success() {
			fmt.Printf("Found successful plan with cost %d\n", tryNext.Spent)
			return *tryNext
		}
		if _, ok := visited[tryNext.String()]; !ok { // possible this is unnecessary since Moves is finite? -- NOPE used 50GB RAM and killed it
			visited[tryNext.String()] = struct{}{}
			for _, move := range tryNext.Moves() {
				m := move
				heap.Push(&planQueue, &m)
			}
		}
	}
	return begin
}

func buildBurrow(input string) []*Space {
	var burrow []*Space
	lines := strings.Split(input, "\n")
	// hallway
	for _, c := range lines[1] {
		var space *Space
		switch c {
		case '#':
			space = &Space{Location: Wall, Id: len(burrow)}
			break
		case '.':
			space = &Space{Location: Hallway, Id: len(burrow)}
			break
		default:
			space = &Space{Location: Outside, Id: len(burrow)}
		}
		if len(burrow) > 0 {
			priorSpace := burrow[len(burrow)-1]
			space.Adjacents = []int{priorSpace.Id}
			priorSpace.Adjacents = append(priorSpace.Adjacents, space.Id)
		}
		burrow = append(burrow, space)
	}
	// rooms
	rooms := []Location{RoomA, RoomB, RoomC, RoomD}
	for i, c := range lines[2] {
		var space *Space
		switch c {
		case ' ':
			space = &Space{Location: Outside, Id: len(burrow)}
			break
		case '#':
			space = &Space{Location: Wall, Id: len(burrow)}
			break
		case 'A': // amphipod
			space = &Space{Location: rooms[0], Id: len(burrow), Occupier: NewAmphipod("Amber", RoomA)}
			rooms = rooms[1:]
			burrow[i].Location = Doorway
			space.Adjacents = []int{burrow[i].Id}
			break
		case 'B': // amphipod
			space = &Space{Location: rooms[0], Id: len(burrow), Occupier: NewAmphipod("Bronze", RoomB)}
			rooms = rooms[1:]
			burrow[i].Location = Doorway
			space.Adjacents = []int{burrow[i].Id}
			break
		case 'C': // amphipod
			space = &Space{Location: rooms[0], Id: len(burrow), Occupier: NewAmphipod("Copper", RoomC)}
			rooms = rooms[1:]
			burrow[i].Location = Doorway
			space.Adjacents = []int{burrow[i].Id}
			break
		case 'D': // amphipod
			space = &Space{Location: rooms[0], Id: len(burrow), Occupier: NewAmphipod("Desert", RoomD)}
			rooms = rooms[1:]
			burrow[i].Location = Doorway
			space.Adjacents = []int{burrow[i].Id}
			break
		default:
			panic("can't parse line 3: unrecognized " + string(c))
		}

		burrow[i].Adjacents = append(burrow[i].Adjacents, space.Id)
		burrow = append(burrow, space)
	}

	parseLineNum := 3
	for lines[parseLineNum] != "  #########" {
		rooms = []Location{RoomA, RoomB, RoomC, RoomD}
		for i, c := range lines[parseLineNum] {
			var space *Space
			switch c {
			case ' ':
				space = &Space{Location: Outside, Id: len(burrow)}
				break
			case '#':
				space = &Space{Location: Wall, Id: len(burrow)}
				break
			case 'A':
				space = &Space{Location: rooms[0], Id: len(burrow), Occupier: NewAmphipod("Amber", RoomA)}
				rooms = rooms[1:]
				doorway := burrow[i]

				room := burrow[doorway.Adjacents[len(doorway.Adjacents)-1]]
				for len(room.Adjacents) > 1 {
					room = burrow[room.Adjacents[len(room.Adjacents)-1]]
				}
				room.Adjacents = append(room.Adjacents, space.Id)
				space.Adjacents = []int{room.Id}
			case 'B':
				space = &Space{Location: rooms[0], Id: len(burrow), Occupier: NewAmphipod("Bronze", RoomB)}
				rooms = rooms[1:]
				doorway := burrow[i]

				room := burrow[doorway.Adjacents[len(doorway.Adjacents)-1]]
				for len(room.Adjacents) > 1 {
					room = burrow[room.Adjacents[len(room.Adjacents)-1]]
				}
				room.Adjacents = append(room.Adjacents, space.Id)
				space.Adjacents = []int{room.Id}
			case 'C':
				space = &Space{Location: rooms[0], Id: len(burrow), Occupier: NewAmphipod("Copper", RoomC)}
				rooms = rooms[1:]
				doorway := burrow[i]

				room := burrow[doorway.Adjacents[len(doorway.Adjacents)-1]]
				for len(room.Adjacents) > 1 {
					room = burrow[room.Adjacents[len(room.Adjacents)-1]]
				}
				room.Adjacents = append(room.Adjacents, space.Id)
				space.Adjacents = []int{room.Id}
			case 'D':
				space = &Space{Location: rooms[0], Id: len(burrow), Occupier: NewAmphipod("Desert", RoomD)}
				rooms = rooms[1:]
				doorway := burrow[i]

				room := burrow[doorway.Adjacents[len(doorway.Adjacents)-1]]
				for len(room.Adjacents) > 1 {
					room = burrow[room.Adjacents[len(room.Adjacents)-1]]
				}
				room.Adjacents = append(room.Adjacents, space.Id)
				space.Adjacents = []int{room.Id}
			}

			burrow = append(burrow, space)
		}
		parseLineNum++
	}

	return burrow
}
