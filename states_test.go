package main

import (
	"testing"
	"./lib/AgentVsAgent"
)

func createGameState() GameState {
	roundState := RoundState{}
	roundState.north = PlayerState{ actions: make(map[Card]Action) }
	roundState.east = PlayerState{ actions: make(map[Card]Action) }
	roundState.south = PlayerState{ actions: make(map[Card]Action), root: true }
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

	newGameState := gameState.play(card)

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
	actions1 := gameState.currentRound().south.actions[card1]
	actions1.received = true
	gameState.currentRound().south.actions[card1] = actions1

	probabilities := gameState.currentRound().probabilities()

	if probabilities["south"][card1] != 100 {
		t.Error("Card should be there", probabilities["south"][card1], card1)
	}
	if probabilities["north"][card1] != 0 || probabilities["west"][card1] != 0 || probabilities["east"][card1] != 0 {
		t.Error("Card shouldn't be elsewhere", card1)
	}

	actions1 = gameState.currentRound().south.actions[card1]
	actions1.played = true
	gameState.currentRound().south.actions[card1] = actions1

	probabilities = gameState.currentRound().probabilities()
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

func TestSimulation(t *testing.T) {
	gameState := createGameState()
	card1 := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_ACE }
	card2 := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_KING }
	card3 := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_QUEEN }
	card4 := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_JACK }
	card5 := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_TEN }
	card6 := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_NINE }
	heldCards := Cards{ &card1, &card2, &card3, &card4, &card5, &card6 }

	for _, card := range heldCards {
		action := gameState.currentRound().south.actions[card1]
		action.dealt = true
		gameState.currentRound().south.actions[*card] = action
	}

	rootSimulation := Simulation{ gameState: &gameState }

	simEvaluation := rootSimulation.evaluate("south")
	gameEvaluation := gameState.evaluate("south")

	if simEvaluation != gameEvaluation {
		t.Error("Unadvanced simulation should have same evaluation as original game state", simEvaluation, gameEvaluation)
	}

	if len(rootSimulation.children) > 0 {
		t.Error("Shouldn't have children until advanced")
	}

	t.Log("advance 1")
	rootSimulation.advance()

	if len(rootSimulation.children) != 3 {
		t.Error("Should guess the two of clubs being with other players")
		for _, child := range rootSimulation.children {
			t.Log("child:", child)
		}
	}

	trick := rootSimulation.children[0].gameState.currentRound().currentTrick()
	if len(trick.played) != 1 || (trick.played[0].suit != AgentVsAgent.Suit_CLUBS || trick.played[0].rank != AgentVsAgent.Rank_TWO) {
		t.Error("Should have played just the two of clubs", trick)
	}

	t.Log("advance 2")
	rootSimulation.advance()

	t.Log("advance 3")
	rootSimulation.advance()

	t.Log("advance 4")
	rootSimulation.advance()

	simulation := rootSimulation.children[0].children[0].children[0].children[0]
	trick = simulation.gameState.currentRound().currentTrick()
	if len(trick.played) != 4 {
		t.Error("Four advances should have filled the trick", trick.played)
	}

	t.Log("advance 5")
	rootSimulation.advance()

	simulation = rootSimulation.children[0].children[0].children[0].children[0].children[0]
	trick = simulation.gameState.currentRound().currentTrick()
	if len(trick.played) != 1 || trick.number != 2 {
		t.Error("Fifth advance should have created another trick", trick.number, trick.played)
	}
}

