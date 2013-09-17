package main

import (
	"fmt"
	"sort"
	"./lib/AgentVsAgent"
)

type Position string

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
			if card.Suit == AgentVsAgent.Suit_CLUBS && card.Rank == AgentVsAgent.Rank_TWO {
				handScore = handScore - 13
			} else {
				handScore = handScore - card.order()
			}
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

	// fmt.Println("1>>>>>>>?????", gameState.round.players[position].held[card])
	newMeta := new(CardMetadata)
	*newMeta = gameState.round.players[position].held[card]
	(*newMeta).played = true

	newHeld := make(map[Card]CardMetadata, 13)
	for heldCard, meta := range gameState.round.players[position].held {
		newHeld[heldCard] = meta
	}
	newHeld[card] = *newMeta

	newPlayed := new(Cards)
	*newPlayed = append(gameState.round.trick.played, &card)

	newTrick := new(TrickState)
	*newTrick = gameState.round.trick
	(*newTrick).played = *newPlayed

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
	(*newRound).trick = *newTrick

	newGameState := new(GameState)
	*newGameState = *gameState
	(*newGameState).round = *newRound

	// fmt.Println("7>>>>>>>?????", gameState.round.players[position].held[card])

	return newGameState
}

func buildGameState(game *Game) GameState {
	var scores map[Position]int
	roundState := buildRoundState(game.rounds[len(game.rounds) - 1])
	return GameState{ round: roundState, scores: scores }
}

func buildRoundState(round *Round) RoundState {
	players := buildPlayerStates(round)
	trickState := buildTrickState(round.tricks[len(round.tricks) - 1])
	return RoundState{ trick: trickState, players: players }
}

func buildPlayerStates(round *Round) map[Position]PlayerState {
	rootPosition := (Position)(round.game.info.Position)
	players := make(map[Position]PlayerState, 4)
	cards := make(map[Card]CardMetadata, 13)

	for _, aCard := range round.held {
		cards[Card{aCard}] = CardMetadata{ played: false, cantOwn: false }
	}
	rootPlayerState := PlayerState{ held: cards }

	players[rootPosition] = rootPlayerState
	return players
}

func buildTrickState(trick *Trick) TrickState {
	var playedCards Cards
	for _, aCard := range trick.played {
		playedCards = append(playedCards, &Card{aCard})
	}
	return TrickState{ leader: (Position)(trick.leader), played: playedCards }
}
