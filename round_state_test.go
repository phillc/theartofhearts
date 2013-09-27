package main

import (
	"testing"
	"./lib/AgentVsAgent"
)

func TestPlay(t *testing.T) {
	roundState := createRoundState()
	position := (Position)("south")

	played := Cards{}
	trickState := TrickState{ leader: position, played: played }
	trickStates := []*TrickState{ &trickState }
	roundState.trickStates = trickStates
	card := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_TWO }

	if len(roundState.currentTrick().played) > 0 {
		t.Error("there should be no played cards")
	}
	if roundState.playerState(position).actions[card].played == true {
		t.Error("the card should not be played")
	}

	newRoundState := roundState.clone()
	newRoundState.play(card)

	if len(roundState.currentTrick().played) > 0 {
		t.Error("there should still be no played cards in the original")
	}
	if roundState.playerState(position).actions[card].played == true {
		t.Error("the card should still not be played in the original")
	}

	if len(newRoundState.currentTrick().played) != 1 {
		t.Error("newRoundState should have the played card")
	}
	if newRoundState.playerState(position).actions[card].played != true {
		t.Error("newRoundState should have the card marked as played")
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

	roundState := createRoundState()

	newRoundState := roundState.clone()
	newRoundState.pass(position, passedCards)

	for _, passedCard := range passedCards {
		if !newRoundState.south.actions[*passedCard].passed {
			t.Error("card should have been marked as passed")
		}
	}

	for _, keptCard := range keptCards {
		if newRoundState.south.actions[*keptCard].passed {
			t.Error("kept card should not have been marked as passed")
		}
	}
}

