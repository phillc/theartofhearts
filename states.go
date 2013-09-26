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
					newRoundState.trickStates = []TrickState{ newTrickState }
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
					roundState.trickStates = append(roundState.trickStates, newTrickState)
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
		fmt.Println("# of children advancing:", len(simulation.children))
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

type Action struct {
	dealt bool
	played bool
	passed bool
	received bool
	// can't have
}

func (action Action) String() string {
	str := "Action:<"
	if action.dealt {
		str = str + " dealt"
	}
	if action.played {
		str = str + " played"
	}
	if action.passed {
		str = str + " passed"
	}
	if action.received {
		str = str + " received"
	}
	if action.isDefinitelyHeld() {
		str = str + " :definitely held:"
	}
	str = str + " >"
	return str
}

func (action *Action) isDefinitelyHeld() bool {
	return !action.played && ((action.dealt && !action.passed) || action.received)
}

type PlayerState struct {
	actions map[Card]Action
	root bool
}

func (playerState *PlayerState) definitelyHeld() Cards {
	cards := Cards{}
	for card, action := range playerState.actions {
		if action.isDefinitelyHeld() {
			aCard := card
			cards = append(cards, &aCard)
		}
	}
	return cards
}

func (playerState *PlayerState) clone() *PlayerState {
	newActions := make(map[Card]Action)
	for card, action := range playerState.actions {
		newActions[card] = action
	}

	newPlayerState := *playerState
	newPlayerState.actions = newActions
	return &newPlayerState
}

func buildRoundState(round *Round) *RoundState {
	players := buildPlayerStates(round)
	var trickStates []TrickState
	for _, trick := range round.tricks {
		trickStates = append(trickStates, *buildTrickState(trick))
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

func buildPlayerStates(round *Round) map[Position]PlayerState {
	rootPosition := (Position)(round.game.info.Position)
	players := make(map[Position]PlayerState, 4)
	cards := make(map[Card]Action, 13)

	for _, aCard := range round.dealt {
		card := Card{ suit: aCard.Suit, rank: aCard.Rank }
		actions := cards[card]
		actions.dealt = true
		actions.played = true
		cards[card] = actions
	}
	for _, aCard := range round.passed {
		//todo: mark received for player passed to
		card := Card{ suit: aCard.Suit, rank: aCard.Rank }
		actions := cards[card]
		actions.passed = true
		actions.played = false
		cards[card] = actions
	}
	for _, aCard := range round.received {
		//todo: mark dealt and passed for player received from
		card := Card{ suit: aCard.Suit, rank: aCard.Rank }
		actions := cards[card]
		actions.received = true
		actions.played = true
		cards[card] = actions
	}
	for _, aCard := range round.held {
		card := Card{ suit: aCard.Suit, rank: aCard.Rank }
		actions := cards[card]
		actions.played = false
		cards[card] = actions
	}
	rootPlayerState := PlayerState{ actions: cards, root: true }

	players[rootPosition] = rootPlayerState
	return players
}

func buildTrickState(trick *Trick) *TrickState {
	var playedCards Cards
	for _, aCard := range trick.played {
		playedCards = append(playedCards, &Card{ suit: aCard.Suit, rank: aCard.Rank })
	}
	return &TrickState{ number: trick.number, leader: (Position)(trick.leader), played: playedCards }
}
