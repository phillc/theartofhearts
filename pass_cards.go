package main

import (
	// "time"
	"./lib/AgentVsAgent"
  // "fmt"
	"sort"
)


// evaluate permutations of three cards

func passCards(round Round) []*AgentVsAgent.Card {
	var cards Cards
	for _, card := range round.dealt {
		cards = append(cards, &Card{card})
	}
	sort.Sort(sort.Reverse(ByOrder{cards}))
	var cardsToPass []*AgentVsAgent.Card
	for _, card := range cards[0:3] {
		cardsToPass = append(cardsToPass, card.Card)
	}
	return cardsToPass
}

