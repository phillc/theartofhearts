package main

import (
	"testing"
	"./lib/AgentVsAgent"
)

func TestCardKnowledgeWhenRootHasCard(t *testing.T) {
	card1 := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_TWO }
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
	twoClubs := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_TWO }
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
	card1 := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_TWO }
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
	card1 := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_TWO }
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
	card1 := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_TWO }
	roundState := createRoundState()

	cardKnowledge := CardKnowledge{}
	cardKnowledge.buildFrom(&roundState, card1)

	if !cardKnowledge.west {
		t.Error("Test assumes west could have the card at start")
	}

	roundState.west.discardedOn(card1.suit)

	cardKnowledge.buildFrom(&roundState, card1)

	if cardKnowledge.west {
		t.Error("Can't have any cards of that suit if discarded on it")
	}
}

