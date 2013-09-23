package main

import (
	"testing"
	"./lib/AgentVsAgent"
)

func createGameState() GameState {
	roundState := RoundState{}
	roundState.north = PlayerState{ actions: make(map[Card]Action) }
	roundState.east = PlayerState{ actions: make(map[Card]Action) }
	roundState.south = PlayerState{ actions: make(map[Card]Action) }
	roundState.west = PlayerState{ actions: make(map[Card]Action) }
	roundStates := []RoundState{ roundState }
	gameState := GameState{ roundStates: roundStates }
	return gameState
}

func TestPlay(t *testing.T) {
	gameState := createGameState()
	position := (Position)("south")

	played := Cards{}
	trickState := TrickState{ leader: position, played: played }
	trickStates := []TrickState{ trickState }
	gameState.currentRound().trickStates = trickStates
	card := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_TWO }

	if len(gameState.currentRound().currentTrick().played) > 0 {
		t.Error("there should be no played cards")
	}
	if gameState.currentRound().playerState(position).actions[card].played == true {
		t.Error("the card should not be played")
	}

	newGameState := gameState.play(position, card)

	if len(gameState.currentRound().currentTrick().played) > 0 {
		t.Error("there should still be no played cards in the original")
	}
	if gameState.currentRound().playerState(position).actions[card].played == true {
		t.Error("the card should still not be played in the original")
	}

	if len(newGameState.currentRound().currentTrick().played) != 1 {
		t.Error("newGameState should have the played card")
	}
	if newGameState.currentRound().playerState(position).actions[card].played != true {
		t.Error("newGameState should have the card marked as played")
	}
}

func TestPass(t *testing.T) {
	position := (Position)("south")
	card1 := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_TWO }
	card2 := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_THREE }
	card3 := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_FOUR }
	card4 := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_FIVE }
	card5 := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_SIX }
	card6 := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_SEVEN }

	dealtCards := Cards{&card1, &card2, &card3, &card4, &card5, &card6}
	passedCards := dealtCards[0:3]
	keptCards := dealtCards[3:3]

	gameState := createGameState()

	newGameState := gameState.pass(position, passedCards)

	for _, passedCard := range passedCards {
		if !newGameState.currentRound().south.actions[*passedCard].passed {
			t.Error("card should have been marked as passed")
		}
	}

	for _, keptCard := range keptCards {
		if newGameState.currentRound().south.actions[*keptCard].passed {
			t.Error("kept card should not have been marked as passed")
		}
	}
}

func TestProbabilities(t *testing.T) {
	card1 := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_TWO }

	gameState := createGameState()
	actions := gameState.currentRound().south.actions[card1]
	actions.received = true
	gameState.currentRound().south.actions[card1] = actions

	probabilities := gameState.currentRound().probabilities()

	cardProbability := probabilities["south"][card1]
	if cardProbability != 100 {
		t.Error("Card should be there", cardProbability)
	}
	if probabilities["north"][card1] != 0 {
		t.Error("Card shouldn't be elsewhere")
	}
}


