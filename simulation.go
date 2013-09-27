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
	if len(simulation.children) == 0 {
		if len(simulation.roundState.trickStates) == 0 {
			twoClubs := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_TWO }

			cardKnowledge := CardKnowledge{}
			cardKnowledge.buildFrom(simulation.roundState, twoClubs)
			for _, position := range allPositions() {
				if cardKnowledge.isPossiblyHeldBy(position) {
					newRoundState := simulation.roundState.clone()
					newTrickState := TrickState{ number: 1, leader: position, played: Cards{} }
					newRoundState.trickStates = []*TrickState{ &newTrickState }
					newRoundState.play(twoClubs)
					newSimulation := Simulation{ roundState: newRoundState, probability: 100 }
					simulation.children = append(simulation.children, &newSimulation)
				}
			}
		} else {
			if len(simulation.roundState.currentTrick().played) == 4 {
				if len(simulation.roundState.trickStates) == 13 {
					return // nothing left to simulate
				} else {
					roundState := simulation.roundState.clone()
					leader := roundState.currentTrick().winner()
					newTrickState := TrickState{ number: len(roundState.trickStates) + 1, leader: leader, played: Cards{} }
					roundState.trickStates = append(roundState.trickStates, &newTrickState)
				}
			}

			// todo: card probabilities != move probabilities... and move probabilities need to add up to 1 (100?)
			for card, probability := range simulation.roundState.playableCardProbabilities() {
				newRoundState := simulation.roundState.clone()
				newRoundState.play(card)
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

