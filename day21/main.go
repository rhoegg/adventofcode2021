package main

import "fmt"

type GameState struct {
	Positions [2]int
	Scores [2] int
	Universes int64
}

type DiracSimulator struct {
	WinningScore int
	rollFrequency map[int]int64
}

func NewDiracSimulator(winningScore int) DiracSimulator {
	rollFrequency := make(map[int]int64)

	for d1 := 1; d1 <= 3; d1++ {
		for d2 := 1; d2 <= 3; d2++ {
			for d3 := 1; d3 <= 3; d3++ {
				rollFrequency[d1 + d2 + d3] += 1
			}
		}
	}
	return DiracSimulator{WinningScore: winningScore, rollFrequency: rollFrequency}
}

func main() {

	sim := NewDiracSimulator(21)

	now := GameState{
		//Positions: [2]int{4, 8},
		Positions: [2]int{8, 3},
		Scores:    [2]int{0, 0},
		Universes: 1,
	}
	counts := sim.WinCounts([]GameState{now})
	fmt.Printf("Win counts %d / %d", counts[0], counts[1])
}

func MovePawn(position, roll int) int {
	return ((position - 1 + roll) % 10) + 1
}

func  (s *DiracSimulator) Possibilities(states []GameState) []GameState {
	var newStates []GameState
	for _, state := range states {
		for roll, universeCount := range s.rollFrequency {
			newPosition := MovePawn(state.Positions[0], roll)
			newStates = append(newStates, GameState{
				Positions: [2]int{state.Positions[1], newPosition},
				Scores: [2]int{state.Scores[1], state.Scores[0] + newPosition},
				Universes: state.Universes * universeCount,
			})
		}
	}
	return newStates
}

func (s *DiracSimulator) WinCounts(states []GameState) [2]int64 {
	outcomes := [2]int64{0, 0}
	var indeterminate []GameState
	for _, p := range s.Possibilities(states) {
		if p.Scores[0] >= s.WinningScore {
			outcomes[0] += p.Universes
		} else if p.Scores[1] >= s.WinningScore {
			outcomes[1] += p.Universes
		} else {
			indeterminate = append(indeterminate, p)
		}
	}
	if len(indeterminate) > 0 {
		nextTurnOutcomes := s.WinCounts(indeterminate)
		for _, i := range []int{0, 1} {
			outcomes[i] += nextTurnOutcomes[(i + 1) % 2]
		}
	}
	return outcomes
}