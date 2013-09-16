package main

import (
)

type Position string

type Player interface {
}

type MyTrick struct {
	*Trick
}

type TrickState struct {
	leader Position
	played Cards
}

func (trickState *TrickState) score() {
}

type CardProbability struct {
	played bool
	cantOwn bool
}

type PlayerState struct {
	held map[Card]CardProbability
}

type RoundState struct {
	scores map[Position]int //must be here for moon shooting
	players map[Position]PlayerState
	trick *TrickState
}

func (roundState *RoundState) evaluate(position Position) int {
	evaluation := 0
	// evaluation = evaluation - roundState.scores[position]
	handScore := 0
	for card, _ := range roundState.players[position].held {
		handScore = handScore - card.order()
	}
	evaluation = evaluation + handScore

	return evaluation
}

type GameState struct {
	scores map[Position]int
	round *RoundState
}

func (gameState *GameState) evaluate(position Position) int {
	evaluation := 0
	//evaluation = evaluation - gameState.scores[position]
	evaluation = evaluation + gameState.round.evaluate(position)
	return evaluation
}
