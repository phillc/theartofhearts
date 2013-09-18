package main

import (
	"testing"
	"./lib/AgentVsAgent"
)

func TestPlay(t *testing.T) {
	card := Card{ &AgentVsAgent.Card{ Suit: AgentVsAgent.Suit_HEARTS, Rank: AgentVsAgent.Rank_TWO } }
	played := Cards{}
	trickState := TrickState{ leader: (Position)("south"), played: played }
	roundState := RoundState{ trick: trickState }
	gameState := GameState{ round: roundState }

	if len(gameState.round.trick.played) > 0 {
		t.Error("there should be no played cards")
	}

	newGameState := gameState.play((Position)("west"), card)

	if len(gameState.round.trick.played) > 0 {
		t.Error("there should still be no played cards")
	}

	if len(newGameState.round.trick.played) != 1 {
		t.Error("newGameState should have the played card")
	}
}

