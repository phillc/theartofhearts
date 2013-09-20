package main

import (
	"fmt"
	"sort"
	"./lib/AgentVsAgent"
)

type Position string

type TrickState struct {
	number int
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
	fmt.Println("trick eval", evaluation)
	return evaluation
}

func (trickState *TrickState) winner() Position {
	matchingSuit := trickState.played.allOfSuit(trickState.played[0].Suit)
	sort.Sort(sort.Reverse(ByOrder{matchingSuit}))
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

func (trickState *TrickState) positionsMissing() []Position {
	return trickState.positionsFromLeader()[len(trickState.played):4]
}

func (trickState *TrickState) isLeading() bool {
	return len(trickState.played) == 0
}

func (trickState *TrickState) clone() *TrickState {
	newTrickState := *trickState
	return &newTrickState
}

type CardMetadata struct {
	//dealt bool
	played bool
	passed bool
	received bool
}

// Uhh... is this needed? why not go staight to the map?
type PlayerState struct {
	held map[Card]CardMetadata
}

func (playerState *PlayerState) clone() *PlayerState {
	newHeld := make(map[Card]CardMetadata)
	for card, meta := range playerState.held {
		newHeld[card] = meta
	}

	newPlayerState := *playerState
	newPlayerState.held = newHeld
	return &newPlayerState
}

type RoundState struct {
	number int
	north PlayerState
	east PlayerState
	south PlayerState
	west PlayerState
	trickStates []TrickState
}

func (roundState *RoundState) playerState(position Position) *PlayerState {
	switch string(position) {
	case "north": return &roundState.north
	case "east": return &roundState.east
	case "south": return &roundState.south
	case "west": return &roundState.west
	}
	return &PlayerState{}
}

func (roundState *RoundState) currentTrick() *TrickState {
	var trickState *TrickState
	if len(roundState.trickStates) > 0 {
		trickState = &roundState.trickStates[len(roundState.trickStates) - 1]
	}
	return trickState
}

func (roundState *RoundState) isHeartsBroken() bool {
	broken := false
	for _, trick := range roundState.trickStates {
		cards := trick.played
		for _, card := range cards {
			if card.Suit == AgentVsAgent.Suit_HEARTS {
				broken = true
				break
			}
		}
	}
	return broken
}

func (roundState *RoundState) playableCards() Cards {
	position := roundState.currentTrick().positionsMissing()[0]
	held := Cards{}
	for card, _ := range roundState.playerState(position).held {
		newCard := card
		held = append(held, &newCard)
	}

	trick := roundState.currentTrick()

	validCards := held

	if trick.number == 1 && trick.isLeading() {
		validCards = validCards.onlyTwoClubs()
	}

	if trick.number == 1 {
		validCards = validCards.noPoints()
	}

	if trick.isLeading() && !roundState.isHeartsBroken() && len(held.noHearts()) > 0 {
		validCards = validCards.noHearts()
	}

	if !trick.isLeading() {
		newValidCards := validCards.allOfSuit(trick.played[0].Suit)
		if len(newValidCards) > 0 {
			validCards = newValidCards
		}
	}

	fmt.Println(">>Valid cards:", validCards)
	return validCards
}

func (roundState *RoundState) evaluate(position Position) int {
	evaluation := 0
	// evaluation = evaluation - roundState.scores[position]
	handScore := 0

	// Take the average of each suit?
	// or something that promotes lower cards (2 + K > 7)? or is it?
	// how about (sum / len) - (len * 3)

	for card, meta := range roundState.playerState(position).held {
		if !meta.played {
			if card.Suit == AgentVsAgent.Suit_CLUBS && card.Rank == AgentVsAgent.Rank_TWO {
				handScore = handScore - 13
			} else {
				handScore = handScore - card.order()
			}
		}
	}
	evaluation = evaluation + handScore

	if roundState.currentTrick() != nil {
		evaluation = evaluation + roundState.currentTrick().evaluate(position)
	}

	return evaluation
}

func (roundState *RoundState) clone() *RoundState {
	var newTrickStates []TrickState
	for _, trickState := range roundState.trickStates {
	  newTrickStates = append(newTrickStates, *trickState.clone())
	}

	newRoundState := *roundState
	newRoundState.trickStates = newTrickStates
	newRoundState.north = *roundState.north.clone()
	newRoundState.east = *roundState.east.clone()
	newRoundState.south = *roundState.south.clone()
	newRoundState.west = *roundState.west.clone()
	return &newRoundState
}

type GameState struct {
	roundStates []RoundState
}

func (gameState *GameState) currentRound() *RoundState {
	return &gameState.roundStates[len(gameState.roundStates) - 1]
}

func (gameState *GameState) evaluate(position Position) int {
	evaluation := 0
	//evaluation = evaluation - gameState.scores[position]
	evaluation = evaluation + gameState.currentRound().evaluate(position)
	return evaluation
}

func (gameState *GameState) clone() *GameState {
	newGameState := *gameState
	var newRoundStates []RoundState
	for _, roundState := range gameState.roundStates {
	  newRoundStates = append(newRoundStates, *roundState.clone())
	}
	newGameState.roundStates = newRoundStates
	return &newGameState
}

func (gameState *GameState) pass(position Position, cards Cards) *GameState {
	newGameState := gameState.clone()
	currentRound := newGameState.currentRound()
	playerState := currentRound.playerState(position)
	held := playerState.held
	for _, passedCard := range cards {
		meta := held[*passedCard]
		meta.passed = true
		held[*passedCard] = meta
	}
	return newGameState
}

func (gameState *GameState) play(position Position, card Card) *GameState {
	fmt.Println("creating game state from play")
	// fmt.Println("before>>>>>>>?????", gameState.currentRound().currentTrick().played)

	newGameState := gameState.clone()
	currentRound := newGameState.currentRound()

	playerState := currentRound.playerState(position)
	held := playerState.held
	meta := held[card]
	meta.played = true
	held[card] = meta
	currentRound.playerState(position).held[card] = meta

	currentTrick := currentRound.currentTrick()
	currentTrick.played = append(currentTrick.played, &card)

	// fmt.Println("after>>>>>>>?????", gameState.currentRound().currentTrick().played)
	return newGameState
}

func buildGameState(game *Game) *GameState {
	var roundStates []RoundState
	for _, round := range game.rounds {
		roundStates = append(roundStates, *buildRoundState(round))
	}
	return &GameState{ roundStates: roundStates }
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
	cards := make(map[Card]CardMetadata, 13)

	for _, aCard := range round.held {
		cards[Card{aCard}] = CardMetadata{}
	}
	rootPlayerState := PlayerState{ held: cards }

	players[rootPosition] = rootPlayerState
	return players
}

func buildTrickState(trick *Trick) *TrickState {
	var playedCards Cards
	for _, aCard := range trick.played {
		playedCards = append(playedCards, &Card{aCard})
	}
	return &TrickState{ number: trick.number, leader: (Position)(trick.leader), played: playedCards }
}
