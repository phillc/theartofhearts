package main

import (
)

func createRoundState() RoundState {
	roundState := RoundState{}
	roundState.north = PlayerState{ actions: make(map[Card]Action) }
	roundState.east = PlayerState{ actions: make(map[Card]Action) }
	roundState.south = PlayerState{ actions: make(map[Card]Action), root: true }
	roundState.west = PlayerState{ actions: make(map[Card]Action) }
	return roundState
}
