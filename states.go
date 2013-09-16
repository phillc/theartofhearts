package main

import (
	"fmt"
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

type CardMetadata struct {
	played bool
	cantOwn bool
}

type PlayerState struct {
	held map[Card]CardMetadata
}

type RoundState struct {
	scores map[Position]int //must be here for moon shooting
	players map[Position]PlayerState
	trick TrickState
}

func (roundState *RoundState) evaluate(position Position) int {
	evaluation := 0
	// evaluation = evaluation - roundState.scores[position]
	handScore := 0
	for card, meta := range roundState.players[position].held {
		if !meta.played {
			handScore = handScore - card.order()
		}
	}
	evaluation = evaluation + handScore

	return evaluation
}

type GameState struct {
	scores map[Position]int
	round RoundState
}

func (gameState *GameState) evaluate(position Position) int {
	evaluation := 0
	//evaluation = evaluation - gameState.scores[position]
	evaluation = evaluation + gameState.round.evaluate(position)
	return evaluation
}

func (gameState *GameState) play(position Position, card Card) *GameState {
	fmt.Println("calculating game state")

	fmt.Println("1>>>>>>>?????", gameState.round.players[position].held[card])
	newMeta := new(CardMetadata)
	*newMeta = gameState.round.players[position].held[card]
	(*newMeta).played = true

	newHeld := make(map[Card]CardMetadata, 13)
	for heldCard, meta := range gameState.round.players[position].held{
		newHeld[heldCard] = meta
	}
	newHeld[card] = *newMeta

	newPlayerState := new(PlayerState)
	*newPlayerState = gameState.round.players[position]
	(*newPlayerState).held = newHeld

	newPlayers := make(map[Position]PlayerState)
	for playerPosition, playerState := range gameState.round.players {
		newPlayers[playerPosition] = playerState
	}
	newPlayers[position] = *newPlayerState

	newRound := new(RoundState)
	*newRound = gameState.round
	(*newRound).players = newPlayers

	newGameState := new(GameState)
	*newGameState = *gameState
	(*newGameState).round = *newRound

	fmt.Println("7>>>>>>>?????", gameState.round.players[position].held[card])

	return newGameState
}
