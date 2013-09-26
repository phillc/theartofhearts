package main

import (
	"fmt"
	"./lib/AgentVsAgent"
)

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
		if len(simulation.children) > 25 {
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

