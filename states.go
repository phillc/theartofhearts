package main

import (
	"fmt"
	"./lib/AgentVsAgent"
)

type Position string

type Simulation struct {
	roundState *RoundState
	children []*Simulation
	probability int
}

func (simulation *Simulation) advance() {
	positions := []Position{"north", "east", "south", "west"}
	roundState := simulation.roundState.clone()

	if len(simulation.children) == 0 {
		if len(roundState.trickStates) == 0 {
			twoClubs := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_TWO }

			probabilities := roundState.probabilities()
			for _, position := range positions {
				probability := probabilities[position][twoClubs]
				if probability > 0 {
					newRoundState := roundState.clone()
					newTrickState := TrickState{ number: 1, leader: position, played: Cards{} }
					newRoundState.trickStates = []*TrickState{ &newTrickState }
					newRoundState = newRoundState.play(twoClubs)
					newSimulation := Simulation{ roundState: newRoundState, probability: probability }
					simulation.children = append(simulation.children, &newSimulation)
				}
			}
		} else {
			if len(roundState.currentTrick().played) == 4 {
				if len(roundState.trickStates) == 13 {
					return // nothing left to simulate
				} else {
					leader := roundState.currentTrick().winner()
					newTrickState := TrickState{ number: len(roundState.trickStates) + 1, leader: leader, played: Cards{} }
					roundState.trickStates = append(roundState.trickStates, &newTrickState)
				}
			}

			// todo: card probabilities != move probabilities
			for card, probability := range roundState.playableCardProbabilities() {
				newRoundState := roundState.play(card)
				newSimulation := Simulation{ roundState: newRoundState, probability: probability }
				simulation.children = append(simulation.children, &newSimulation)
			}
		}
	} else {
		/*simulation.children[0].advance()*/
		if len(simulation.children) > 20 {
			fmt.Println("# of children advancing:", len(simulation.children))
		}
		for _, child := range simulation.children {
			child.advance()
		}
	}
}

func (simulation *Simulation) evaluate(position Position) int {
	evaluation := 1000000 // or rather, infinity
	if len(simulation.children) > 0 {
		for _, child := range simulation.children {
			childEval := child.roundState.evaluate(position)
			if childEval < evaluation {
				evaluation = childEval
			}
		}
	} else {
		evaluation = simulation.roundState.evaluate(position)
	}
	return evaluation
}

func buildRoundState(round *Round) *RoundState {
	var trickStates []*TrickState
	for _, trick := range round.tricks {
		trickStates = append(trickStates, buildTrickState(trick))
	}

	rootPosition := (Position)(round.game.info.Position)
	players := make(map[Position]*PlayerState, 4)
	players["north"] = &PlayerState{}
	players["east"] = &PlayerState{}
	players["south"] = &PlayerState{}
	players["west"] = &PlayerState{}
	rootPlayer := players[rootPosition]
	rootPlayer.root = true

	passingTo := rootPlayer
	receivedFrom := rootPlayer
	positions := []Position{"north", "east", "south", "west"}
	rootIndex := -1
	for i, position := range positions {
		if position == rootPosition {
			rootIndex = i
			break
		}
	}
	positionsFromRoot := append(positions[rootIndex:4], positions[0:rootIndex]...)
	switch (round.number - 1) % 4 {
	case 0:
		// left
		passingTo = players[positionsFromRoot[3]]
		receivedFrom = players[positionsFromRoot[1]]
	case 1:
		// right
		passingTo = players[positionsFromRoot[1]]
		receivedFrom = players[positionsFromRoot[3]]
	case 2:
		// across
		passingTo = players[positionsFromRoot[2]]
		receivedFrom = players[positionsFromRoot[2]]
	}

	for _, aCard := range round.dealt {
		card := Card{ suit: aCard.Suit, rank: aCard.Rank }
		rootPlayer.dealt(card)
	}
	for _, aCard := range round.passed {
		card := Card{ suit: aCard.Suit, rank: aCard.Rank }
		rootPlayer.passed(card)
		passingTo.received(card)
	}
	for _, aCard := range round.received {
		card := Card{ suit: aCard.Suit, rank: aCard.Rank }
		rootPlayer.received(card)
		receivedFrom.passed(card)
	}

	for _, trickState := range trickStates {
		for index, position := range trickState.positionsFromLeader()[0:len(trickState.played)] {
			players[position].played(*trickState.played[index])
		}
	}

	return &RoundState{
		number: round.number,
		trickStates: trickStates,
		north: players["north"],
		east: players["east"],
		south: players["south"],
		west: players["west"],
	}
}

func buildTrickState(trick *Trick) *TrickState {
	var playedCards Cards
	for _, aCard := range trick.played {
		playedCards = append(playedCards, &Card{ suit: aCard.Suit, rank: aCard.Rank })
	}
	return &TrickState{ number: trick.number, leader: (Position)(trick.leader), played: playedCards }
}
