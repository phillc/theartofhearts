package main

import (
	"fmt"
	"sort"
	"./lib/AgentVsAgent"
)

type Position string

type Simulation struct {
	gameState *GameState
	children []*Simulation
	probability int
}

func (simulation *Simulation) advance() {
	positions := []Position{"north", "east", "south", "west"}
	gs := simulation.gameState.clone()
	roundState := gs.currentRound()

	if len(simulation.children) == 0 {
		if len(roundState.trickStates) == 0 {
			twoClubs := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_TWO }

			probabilities := roundState.probabilities()
			for _, position := range positions {
				probability := probabilities[position][twoClubs]
				if probability > 0 {
					newGameState := gs.clone()
					newTrickState := TrickState{ number: 1, leader: position, played: Cards{} }
					newGameState.currentRound().trickStates = []TrickState{ newTrickState }
					newGameState = newGameState.play(twoClubs)
					newSimulation := Simulation{ gameState: newGameState, probability: probability }
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

			for card, probability := range roundState.playableCardProbabilities() {
				newGameState := gs.play(card)
				newSimulation := Simulation{ gameState: newGameState, probability: probability }
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
			childEval := child.gameState.evaluate(position)
			if childEval < evaluation {
				evaluation = childEval
			}
		}
	} else {
		evaluation = simulation.gameState.evaluate(position)
	}
	return evaluation
}

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
	return evaluation
}

func (trickState *TrickState) winner() Position {
	matchingSuit := trickState.played.allOfSuit(trickState.played[0].suit)
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

type Action struct {
	dealt bool
	played bool
	passed bool
	received bool
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
			if card.suit == AgentVsAgent.Suit_HEARTS {
				broken = true
				break
			}
		}
	}
	return broken
}

func (roundState *RoundState) playableCards() Cards {
	position := roundState.currentTrick().positionsMissing()[0]
	held := roundState.playerState(position).definitelyHeld()
  return roundState.playableCardsOutOf(held)
}

func (roundState *RoundState) playableCardsOutOf(startingCards Cards) Cards {
	held := Cards{}
	validCards := Cards{}
	for _, card := range startingCards {
		held = append(held, card)
		validCards = append(validCards, card)
	}
	trick := roundState.currentTrick()

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
		newValidCards := validCards.allOfSuit(trick.played[0].suit)
		if len(newValidCards) > 0 {
			validCards = newValidCards
		}
	}
	return validCards
}

func (roundState *RoundState) playableCardProbabilities() map[Card]int {
	position := roundState.currentTrick().positionsMissing()[0]
	probabilities := roundState.probabilities()[position]

	possiblyHeldCards := Cards{}
	for card, probability := range probabilities {
		if probability > 0 {
			aCard := card
			possiblyHeldCards = append(possiblyHeldCards, &aCard)
		}
	}

	playableProbabilities := make(map[Card]int)
	for _, card := range roundState.playableCardsOutOf(possiblyHeldCards) {
		playableProbabilities[*card] = probabilities[*card]
	}

	return playableProbabilities
}

func (roundState *RoundState) probabilities() map[Position]map[Card]int {
	positions := []Position{"north", "east", "south", "west"}
	probabilities := make(map[Position]map[Card]int, 4)

	for _, position := range positions {
		probabilities[position] = make(map[Card]int)
	}

	cards := allCards()
	for _, card := range cards {
		for _, position := range positions {
			playerState := roundState.playerState(position)
			actions := playerState.actions[*card]
			if actions.isDefinitelyHeld() {
				probabilities[position][*card] = 100
				for _, otherPosition := range positions {
					if otherPosition != position {
						probabilities[otherPosition][*card] = 0
					}
				}
				break
			} else if actions.played {
				for _, otherPosition := range positions {
					probabilities[otherPosition][*card] = 0
				}
				break
			} else if !playerState.root {
				// todo: if played off suit, then zero and chage the other guys
				probabilities[position][*card] = 33
			}
		}
	}

	return probabilities
}

func (roundState *RoundState) evaluate(position Position) int {
	evaluation := 0
	// evaluation = evaluation - roundState.scores[position]
	handScore := 0

	// Take the average of each suit?
	// or something that promotes lower cards (2 + K > 7)? or is it?
	// how about (sum / len) - (len * 3)

	for card, action := range roundState.playerState(position).actions {
		if action.isDefinitelyHeld() {
			// todo: two of clubs doesn't matter if we can just simulate past a couple tricks
			if card.suit == AgentVsAgent.Suit_CLUBS && card.rank == AgentVsAgent.Rank_TWO {
				handScore = handScore - 13
			} else {
				handScore = handScore - card.order()
			}
		}
	}
	evaluation = evaluation + handScore

	// todo: full round score

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
	actions := playerState.actions
	for _, passedCard := range cards {
		action := actions[*passedCard]
		action.passed = true
		actions[*passedCard] = action
	}
	return newGameState
}

func (gameState *GameState) play(card Card) *GameState {
	newGameState := gameState.clone()
	currentRound := newGameState.currentRound()
	position := currentRound.currentTrick().positionsMissing()[0]

	playerState := currentRound.playerState(position)
	actions := playerState.actions
	action := actions[card]
	action.played = true
	actions[card] = action
	currentRound.playerState(position).actions[card] = action

	currentTrick := currentRound.currentTrick()
	currentTrick.played = append(currentTrick.played, &card)

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
