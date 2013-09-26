package main

import (
	"testing"
	"./lib/AgentVsAgent"
)

func TestProbabilities(t *testing.T) {
	card1 := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_TWO }

	roundState := createRoundState()
	actions1 := roundState.south.actions[card1]
	actions1.received = true
	roundState.south.actions[card1] = actions1

	probabilities := roundState.probabilities()

	if probabilities["south"][card1] != 100 {
		t.Error("Card should be there", probabilities["south"][card1], card1)
	}
	if probabilities["north"][card1] != 0 || probabilities["west"][card1] != 0 || probabilities["east"][card1] != 0 {
		t.Error("Card shouldn't be elsewhere", card1)
	}

	actions1 = roundState.south.actions[card1]
	actions1.played = true
	roundState.south.actions[card1] = actions1

	probabilities = roundState.probabilities()
	if probabilities["south"][card1] != 0 {
		t.Error("Card was played", probabilities["south"][card1], card1)
	}

	twoClubs := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_TWO }
	if probabilities["south"][twoClubs] != 0 {
		t.Error("If root player doesn't see it, he can't have it")
	}
	if probabilities["north"][twoClubs] == 0 || probabilities["west"][twoClubs] == 0 || probabilities["east"][twoClubs] == 0 {
		t.Error("Well if the root player doesn't have it, it must be elsewhere")
	}
}


func TestPlay(t *testing.T) {
	roundState := createRoundState()
	position := (Position)("south")

	played := Cards{}
	trickState := TrickState{ leader: position, played: played }
	trickStates := []TrickState{ trickState }
	roundState.trickStates = trickStates
	card := Card{ suit: AgentVsAgent.Suit_HEARTS, rank: AgentVsAgent.Rank_TWO }

	if len(roundState.currentTrick().played) > 0 {
		t.Error("there should be no played cards")
	}
	if roundState.playerState(position).actions[card].played == true {
		t.Error("the card should not be played")
	}

	newRoundState := roundState.play(card)

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

	newRoundState := roundState.pass(position, passedCards)

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

