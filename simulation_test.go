package main

import (
	"testing"
	"./lib/AgentVsAgent"
)

func TestSimulation(t *testing.T) {
	roundState := createRoundState()
	card1 := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_ACE }
	card2 := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_KING }
	card3 := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_QUEEN }
	card4 := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_JACK }
	card5 := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_TEN }
	card6 := Card{ suit: AgentVsAgent.Suit_CLUBS, rank: AgentVsAgent.Rank_NINE }
	heldCards := Cards{ &card1, &card2, &card3, &card4, &card5, &card6 }

	for _, card := range heldCards {
		action := roundState.south.actions[card1]
		action.dealt = true
		roundState.south.actions[*card] = action
	}

	rootSimulation := Simulation{ roundState: &roundState }

	simEvaluation := rootSimulation.evaluate("south")
	roundEvaluation := roundState.evaluate("south")

	if simEvaluation != roundEvaluation {
		t.Error("Unadvanced simulation should have same evaluation as original", simEvaluation, roundEvaluation)
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

	trick := rootSimulation.children[0].roundState.trickStates[0]
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
	trick = simulation.roundState.trickStates[0]
	if len(trick.played) != 4 {
		t.Error("Four advances should have filled the trick", trick.played)
	}

	t.Log("advance 5")
	rootSimulation.advance()

	simulation = rootSimulation.children[0].children[0].children[0].children[0].children[0]
	trick = simulation.roundState.currentTrick()
	if len(trick.played) != 1 || trick.number != 2 {
		t.Error("Fifth advance should have created another trick", trick.number, trick.played)
	}
}

