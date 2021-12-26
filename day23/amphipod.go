package main

import "math"

type Amphipod struct {
	Type string
	Home Location
}

func NewAmphipod(amphipodType string, home Location) *Amphipod {
	return &Amphipod{
		Type: amphipodType,
		Home: home,
	}
}

func (a Amphipod) String() string {
	return string(a.Type[0])
}

func (a Amphipod) MoveEnergy() int64 {
	switch a.Type {
	case "Amber":
		return 1
	case "Bronze":
		return 10
	case "Copper":
		return 100
	case "Desert":
		return 1000
	default:
		return math.MaxInt64
	}
}
