package main

import (
	"testing"
	"./lib/AgentVsAgent"
)

func TestPlay(t *testing.T) {
	card := Card{ &AgentVsAgent.Card{ Suit: AgentVsAgent.Suit_HEARTS, Rank: AgentVsAgent.Rank_TWO } }
	played := Cards{}
	position := (Position)("south")

	trickState := TrickState{ leader: position, played: played }
	trickStates := []TrickState{ trickState }
	held := make(map[Card]CardMetadata)
	playerState := PlayerState{ held: held }
	playerStates := make(map[Position]PlayerState)
	playerStates[position] = playerState
	roundState := RoundState{ trickStates: trickStates, playerStates: playerStates }
	roundStates := []RoundState{ roundState }
	gameState := GameState{ roundStates: roundStates }

	if len(gameState.currentRound().currentTrick().played) > 0 {
		t.Error("there should be no played cards")
	}
	if gameState.currentRound().playerStates[position].held[card].played == true {
		t.Error("the card should not be played")
	}

	newGameState := gameState.play(position, card)

	if len(gameState.currentRound().currentTrick().played) > 0 {
		t.Error("there should still be no played cards in the original")
	}
	if gameState.currentRound().playerStates[position].held[card].played == true {
		t.Error("the card should still not be played in the original")
	}

	if len(newGameState.currentRound().currentTrick().played) != 1 {
		t.Error("newGameState should have the played card")
	}
	if newGameState.currentRound().playerStates[position].held[card].played != true {
		t.Error("newGameState should have the card marked as played")
	}
}

