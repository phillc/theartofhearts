package main

import (
	"testing"
	"./lib/AgentVsAgent"
)

func newGameState() GameState {
	roundState := RoundState{}
	roundState.north = PlayerState{ held: make(map[Card]CardMetadata) }
	roundState.east = PlayerState{ held: make(map[Card]CardMetadata) }
	roundState.south = PlayerState{ held: make(map[Card]CardMetadata) }
	roundState.west = PlayerState{ held: make(map[Card]CardMetadata) }
	roundStates := []RoundState{ roundState }
	gameState := GameState{ roundStates: roundStates }
	return gameState
}

func TestPlay(t *testing.T) {
	gameState := newGameState()
	position := (Position)("south")

	played := Cards{}
	trickState := TrickState{ leader: position, played: played }
	trickStates := []TrickState{ trickState }
	gameState.currentRound().trickStates = trickStates
	card := Card{ &AgentVsAgent.Card{ Suit: AgentVsAgent.Suit_HEARTS, Rank: AgentVsAgent.Rank_TWO } }


	if len(gameState.currentRound().currentTrick().played) > 0 {
		t.Error("there should be no played cards")
	}
	if gameState.currentRound().playerState(position).held[card].played == true {
		t.Error("the card should not be played")
	}

	newGameState := gameState.play(position, card)

	if len(gameState.currentRound().currentTrick().played) > 0 {
		t.Error("there should still be no played cards in the original")
	}
	if gameState.currentRound().playerState(position).held[card].played == true {
		t.Error("the card should still not be played in the original")
	}

	if len(newGameState.currentRound().currentTrick().played) != 1 {
		t.Error("newGameState should have the played card")
	}
	if newGameState.currentRound().playerState(position).held[card].played != true {
		t.Error("newGameState should have the card marked as played")
	}
}

func TestPass(t *testing.T) {
	position := (Position)("south")
	/*card := Card{ &AgentVsAgent.Card{ Suit: AgentVsAgent.Suit_HEARTS, Rank: AgentVsAgent.Rank_TWO } }*/
	cards := Cards{}

	gameState := newGameState()

	gameState.pass(position, cards)

	/*if len(gameState.currentRound().currentTrick().played) > 0 {*/
	/*	t.Error("there should be no played cards")*/
	/*}*/
	/*if gameState.currentRound().playerStates[position].held[card].played == true {*/
	/*	t.Error("the card should not be played")*/
	/*}*/

	/*newGameState := gameState.play(position, card)*/

	/*if len(gameState.currentRound().currentTrick().played) > 0 {*/
	/*	t.Error("there should still be no played cards in the original")*/
	/*}*/
	/*if gameState.currentRound().playerStates[position].held[card].played == true {*/
	/*	t.Error("the card should still not be played in the original")*/
	/*}*/

	/*if len(newGameState.currentRound().currentTrick().played) != 1 {*/
	/*	t.Error("newGameState should have the played card")*/
	/*}*/
	/*if newGameState.currentRound().playerStates[position].held[card].played != true {*/
	/*	t.Error("newGameState should have the card marked as played")*/
	/*}*/
}

