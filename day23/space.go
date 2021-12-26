package main

type Location int16

const (
	Outside Location = iota
	Wall
	Hallway
	Doorway
	RoomA
	RoomB
	RoomC
	RoomD
)
func (l Location) String() string {
	switch l {
	case Outside:
		return "Outside"
	case Wall:
		return "Wall"
	case Hallway:
		return "Hallway"
	case Doorway:
		return "Doorway"
	case RoomA:
		return "Room A"
	case RoomB:
		return "Room B"
	case RoomC:
		return "Room C"
	case RoomD:
		return "Room D"
	default:
		return ""
	}
}

func (l Location) Icon() string {
	switch l {
	case Outside:
		return " "
	case Wall:
		return "#"
	case Hallway:
		return "."
	case Doorway:
		return "."
	case RoomA:
		return "."
	case RoomB:
		return "."
	case RoomC:
		return "."
	case RoomD:
		return "."
	default:
		return ""
	}
}

type Space struct {
	Location
	Id        int
	Adjacents []int
	Occupier  *Amphipod
}

func (s Space) String() string {
	if s.Occupier != nil {
		return s.Occupier.String()
	}
	return s.Location.Icon()
}
