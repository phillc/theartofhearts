package main

import (
	"testing"
)

func TestCardProbabilities(t *testing.T) {
	cardKnowledge := CardKnowledge{}

	probabilities := cardKnowledge.probabilities()
	if probabilities["north"] != 0 || probabilities["east"] != 0 || probabilities["south"] != 0 || probabilities["west"] != 0 {
		t.Error("Card should not be anywhere", probabilities)
	}

	cardKnowledge.north = true
	probabilities = cardKnowledge.probabilities()
	if probabilities["north"] != 100 || probabilities["east"] != 0 || probabilities["south"] != 0 || probabilities["west"] != 0 {
		t.Error("Card should be at north", probabilities)
	}

	cardKnowledge.east = true
	probabilities = cardKnowledge.probabilities()
	if probabilities["north"] != 50 || probabilities["east"] != 50 || probabilities["south"] != 0 || probabilities["west"] != 0 {
		t.Error("Card should be at north or east", probabilities)
	}

	cardKnowledge.south = true
	probabilities = cardKnowledge.probabilities()
	if probabilities["north"] != 33 || probabilities["east"] != 33 || probabilities["south"] != 33 || probabilities["west"] != 0 {
		t.Error("Card should be at north, east, or south", probabilities)
	}

	cardKnowledge.west = true
	probabilities = cardKnowledge.probabilities()
	if probabilities["north"] != 25 || probabilities["east"] != 25 || probabilities["south"] != 25 || probabilities["west"] != 25 {
		t.Error("I guess the card could be anywhere", probabilities)
	}
}

func TestCardKnowledgeWhenRootHasCard(t *testing.T) {
	card1 := Card{ Suit: HEARTS, Rank: TWO }
	roundState := createRoundState()
	roundState.south.received(card1)

	cardKnowledge := CardKnowledge{}
	cardKnowledge.buildFrom(&roundState, card1)

	if !cardKnowledge.south {
		t.Error("Card should be there", cardKnowledge)
	}
	if cardKnowledge.north || cardKnowledge.east || cardKnowledge.west {
		t.Error("Card couldn't be elsewhere", cardKnowledge)
	}
}

func TestCardKnowledgeWhenRootDoesNotHaveCard(t *testing.T) {
	twoClubs := Card{ Suit: CLUBS, Rank: TWO }
	roundState := createRoundState()

	cardKnowledge := CardKnowledge{}
	cardKnowledge.buildFrom(&roundState, twoClubs)

	if cardKnowledge.south {
		t.Error("If root player doesn't see it, he can't have it", cardKnowledge)
	}
	if !cardKnowledge.north || !cardKnowledge.east || !cardKnowledge.west {
		t.Error("Well if the root player doesn't have it, it must be elsewhere")
	}
}

func TestCardKnowledgeWhenCardIsPlayed(t *testing.T) {
	card1 := Card{ Suit: HEARTS, Rank: TWO }
	roundState := createRoundState()

	cardKnowledge := CardKnowledge{}
	cardKnowledge.buildFrom(&roundState, card1)

	if !cardKnowledge.north || !cardKnowledge.east || !cardKnowledge.west || cardKnowledge.south {
		t.Error("Test assumes someone else has the card")
	}

	roundState.north.played(card1)

	cardKnowledge.buildFrom(&roundState, card1)

	if cardKnowledge.north || cardKnowledge.east || cardKnowledge.west || cardKnowledge.south {
		t.Error("If the card was played, no one should have the card", cardKnowledge)
	}
}

func TestCardKnowledgeWhenCardIsPlayedByRoot(t *testing.T) {
	card1 := Card{ Suit: HEARTS, Rank: TWO }
	roundState := createRoundState()
	roundState.south.received(card1)

	cardKnowledge := CardKnowledge{}
	cardKnowledge.buildFrom(&roundState, card1)

	if !cardKnowledge.south || cardKnowledge.north || cardKnowledge.east || cardKnowledge.west {
		t.Error("Test assumes root starts with card")
	}

	roundState.south.played(card1)
	cardKnowledge.buildFrom(&roundState, card1)

	if cardKnowledge.north || cardKnowledge.east || cardKnowledge.west || cardKnowledge.south {
		t.Error("If the card was played, no one should have the card", cardKnowledge)
	}
}

func TestCardKnowledgeWhenAPlayerDoesNotHaveASuit(t *testing.T) {
	card1 := Card{ Suit: HEARTS, Rank: TWO }
	roundState := createRoundState()

	cardKnowledge := CardKnowledge{}
	cardKnowledge.buildFrom(&roundState, card1)

	if !cardKnowledge.west {
		t.Error("Test assumes west could have the card at start")
	}

	roundState.west.discardedOn(card1.Suit)

	cardKnowledge.buildFrom(&roundState, card1)

	if cardKnowledge.west {
		t.Error("Can't have any cards of that suit if discarded on it")
	}
}

