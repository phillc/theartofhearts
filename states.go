package main

import (
	"fmt"
	"sort"
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

func (trickState *TrickState) evaluate(position Position) int {
	var evaluation int
	if len(trickState.played) == 0 {
		evaluation = 0
	} else if trickState.winner() == position {
		evaluation = evaluation - (trickState.score() * 10)
	} else {
		evaluation = evaluation + (trickState.score() * 3)
	}
	fmt.Println("TRICK STATE EVALUATED AT", evaluation)
	return evaluation
}

func (trickState *TrickState) winner() Position {
	matchingSuit := trickState.played.allOfSuit(trickState.played[0].Suit)
	sort.Sort(ByOrder{matchingSuit})
	winningCard := matchingSuit[0]
	winningCardIndex := trickState.played.indexOf(winningCard)
	return trickState.positionsFromLeader()[winningCardIndex]
}

func (trickState *TrickState) score() int {
	score := 0
	for _, card := range trickState.played {
		score = score + card.score()
	}
	return score
}

func (trickState *TrickState) positionsFromLeader() []Position {
	positions := []Position{"north", "east", "south", "west"}
	leaderIndex := -1
	for i, position := range positions {
		if position == trickState.leader {
			leaderIndex = i
			break
		}
	}
	sortedPositions := append(positions[leaderIndex:4], positions[0:leaderIndex]...)
	return sortedPositions
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

	evaluation = evaluation + roundState.trick.evaluate(position)

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
