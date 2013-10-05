package main

import (
)

type Knowledge struct {
	trackedCards map[Card]CardKnowledge
}

func (knowledge *Knowledge) buildFrom(roundState *RoundState) {
	trackedCards := make(map[Card]CardKnowledge)

	cards := allCards()
	for _, card := range cards {
		cardKnowledge := CardKnowledge{}
		cardKnowledge.buildFrom(roundState, *card)
		trackedCards[*card] = cardKnowledge
	}

	knowledge.trackedCards = trackedCards
}

func (knowledge *Knowledge) possiblyHeldCardsFor(position Position) Cards {
	cards := Cards{}
	for card, cardKnowledge := range knowledge.trackedCards {
		if cardKnowledge.isPossiblyHeldBy(position) {
			aCard := card
			cards = append(cards, &aCard)
		}
	}
	return cards
}


type CardKnowledge struct {
	north bool
	east bool
	south bool
	west bool
}

func (cardKnowledge *CardKnowledge) buildFrom(roundState *RoundState, card Card) {
	cardKnowledge.north = true
	cardKnowledge.east = true
	cardKnowledge.south = true
	cardKnowledge.west = true

	for _, position := range allPositions() {
		playerState := roundState.playerState(position)
		actions := playerState.actions[card]
		if actions.isDefinitelyHeld() {
			cardKnowledge.ruleOut(otherPositions(position)...)
			break
		} else if actions.played {
			cardKnowledge.ruleOut(allPositions()...)
			break
		} else if playerState.root || playerState.hasDiscardedOn(card.suit){
			cardKnowledge.ruleOut(position)
		}
	}
}

func (cardKnowledge *CardKnowledge) ruleOut(positions ...Position) {
	for _, position := range positions {
		switch position {
		case "north": cardKnowledge.north = false
		case "east": cardKnowledge.east = false
		case "south": cardKnowledge.south = false
		case "west": cardKnowledge.west = false
		}
	}
}

func (cardKnowledge *CardKnowledge) isPossiblyHeldBy(position Position) bool {
	possible := false
	switch position {
	case "north": possible = cardKnowledge.north
	case "east": possible = cardKnowledge.east
	case "south": possible = cardKnowledge.south
	case "west": possible = cardKnowledge.west
	}
	return possible
}

func (cardKnowledge *CardKnowledge) probabilities() map[Position]int {
	probabilities := make(map[Position]int, 4)
	count := 0
	for _, position := range allPositions() {
		if cardKnowledge.isPossiblyHeldBy(position) {
			count++
		}
	}

	for _, position := range allPositions() {
		if cardKnowledge.isPossiblyHeldBy(position) {
			probabilities[position] = (100 / count)
		}
	}
	return probabilities
}
