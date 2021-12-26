package main

import (
	"fmt"
	"math"
)

type Plan struct {
	Burrow []*Space
	Spent  int64
}
func (p Plan) String() string {
	text := ""
	for i, v := range p.Burrow {
		if i == 0 || i == 13 || i == 26 || i == 37 || i == 48 {
			text += "\n"
		}
		text += v.String()
	}
	return text
}

type PlanPriorityQueue []*Plan

func (pq PlanPriorityQueue) Len() int { return len(pq) }
func (pq PlanPriorityQueue) Less(i, j int) bool {
	return pq[i].Spent < pq[j].Spent
}
func (pq PlanPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PlanPriorityQueue) Push(plan interface{}) {
	*pq = append(*pq, plan.(*Plan))
}
func (pq *PlanPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func (p Plan) Success() bool {
	for _, space := range p.Burrow {
		if nil != space.Occupier && space.Occupier.Home != space.Location {
			//fmt.Printf("%s not home yet\n", space.Occupier.Type)
			return false
		}
	}
	return true
}

func (p *Plan) Moves() []Plan {
	var moves []Plan
	for _, space := range p.Burrow {
		if nil != space.Occupier && (! p.RoomReady(*space) || space.Occupier.Home != space.Location) {
			//fmt.Printf("Looking for moves for %s in space %d\n", space.Occupier.Type, space.Id)
			possibles := p.PossibleMoves(*space)
			//for _, possible := range possibles {
			//	fmt.Printf("Guessing (%d) %v\n", possible.Spent, possible)
			//}
			moves = append(moves, possibles...)
		}
	}
	return moves
}

func (p *Plan) PossibleMoves(from Space) []Plan {
	if nil == from.Occupier {
		fmt.Printf("No moves because source is empty %d\n", from.Id)
		return nil
	}
	if from.Occupier.Home == from.Location && p.RoomReady(from) {
		return nil
	}
	var moves []Plan
	for _, destination := range p.Burrow {
		if p.AllowedMove(from, *destination) && p.Reachable(from, *destination) {
			// don't stop one short of the back of the room
			ignore := false
			if destination.Location == from.Occupier.Home && len(destination.Adjacents) == 2 {
				// our room but not the back
				// assuming room size 2
				for _, neighbor := range destination.Adjacents {
					otherRoomLocation := p.Burrow[neighbor]
					if otherRoomLocation.Location == destination.Location &&
						len(otherRoomLocation.Adjacents) == 1 &&
						otherRoomLocation.Occupier == nil {
						ignore = true
					}
				}
			}
			if !ignore {
				moves = append(moves, p.Move(from, *destination))
			}
		}
	}
	return moves
}

func (p *Plan) Move(from, to Space) Plan {

	newBurrow := make([]*Space, len(p.Burrow))
	for _, s := range p.Burrow {
		c := *s
		newBurrow[c.Id] = &c
	}
	if nil == from.Occupier {
		panic(fmt.Sprintf("tried to move and we have no Occupier in %d", from.Id))
	}
	if nil != to.Occupier {
		panic(fmt.Sprintf("tried to move %s on top of %s in space %d", from.Occupier.Type, to.Occupier.Type, to.Id))
	}
	//fmt.Printf("Moving %s from %d to %d in burrow sized %d\n", from.Occupier.Type, from.Id, to.Id, len(newBurrow))
	amphipod := newBurrow[from.Id].Occupier
	newBurrow[from.Id].Occupier = nil
	newBurrow[to.Id].Occupier = amphipod

	return Plan{
		Burrow: newBurrow,
		Spent:  p.Spent + p.ComputeSteps(from, to) * from.Occupier.MoveEnergy(),
	}
}

func (p *Plan) ComputeSteps(from, to Space) int64 {
	type Step struct {
		space *Space
		steps int64
	}
	var cheapest int64 = math.MaxInt64
	nextSteps := []Step{{space: &from, steps: int64(0)}}
	visited := make(map[int]struct{})
	var step Step
	for len(nextSteps) > 0 {
		step, nextSteps = nextSteps[0], nextSteps[1:]
		if step.space.Id == to.Id && step.steps < cheapest {
			cheapest = step.steps
		} else if _, ok := visited[step.space.Id]; !ok {
			visited[step.space.Id] = struct{}{}
			for _, adj := range step.space.Adjacents {
				nextSteps = append(nextSteps, Step{space: p.Burrow[adj], steps: step.steps + 1})
			}
		}
	}
	return cheapest
}

func (p Plan) RoomReady(s Space) bool {
	visited := make(map[int]struct{})
	locationsToCheck := []Space{s}
	for len(locationsToCheck) > 0 {
		s, locationsToCheck = locationsToCheck[0], locationsToCheck[1:]
		if _, ok := visited[s.Id]; !ok {
			if nil != s.Occupier && s.Location != s.Occupier.Home {
				return false
			}
			visited[s.Id] = struct{}{}
			for _, adj := range s.Adjacents {
				if p.Burrow[adj].Location == s.Location {
					locationsToCheck = append(locationsToCheck, *p.Burrow[adj])
				}
			}
		}
	}

	return true
}

func (p Plan) AllowedMove(from, to Space) bool {
	if to.Location == Doorway || to.Location == Wall || to.Location == Outside || to.Occupier != nil {
		return false
	}

	if from.Occupier == nil {
		return false
	}

	if from.Location == Hallway {
		return to.Location == from.Occupier.Home && p.RoomReady(to)
	}
	return (to.Location == from.Occupier.Home && p.RoomReady(to)) || to.Location == Hallway
}


func (p Plan) Reachable(from, to Space) bool {
	visited := make(map[int]struct{})
	neighbors := []Space{from}
	for len(neighbors) > 0 {
		from, neighbors = neighbors[0], neighbors[1:]
		if _, ok := visited[from.Id]; !ok {
			visited[from.Id] = struct{}{}
			if from.Id == to.Id {
				return true
			}
			for _, neighborId := range from.Adjacents {
				neighbor := p.Burrow[neighborId]
				if neighbor.Location != Wall && neighbor.Location != Outside && neighbor.Occupier == nil {
					neighbors = append(neighbors, *neighbor)
				}
			}
		}
	}
	return false
}
