package main

import (
	"testing"
	// "./lib/AgentVsAgent"
)

func TestEquality(t *testing.T) {
	card1 := *allCards()[0]
	card2 := *allCards()[0]
	if card1 != card2 {
		t.Error("same card should be equal to itself", card1, card2)
	}
}
