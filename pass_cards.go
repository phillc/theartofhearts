package main

import (
	// "time"
	"./lib/AgentVsAgent"
	"fmt"
	"sort"
)


// evaluate permutations of three cards

func passCards(round Round) []*AgentVsAgent.Card {
	game := round.game
	gameState := buildGameState(game)
	position := (Position)(game.info.Position)

	var heldCards Cards
	for heldCard, _ := range gameState.currentRound().playerStates[position].held {
		heldCards = append(heldCards, &heldCard)
	}
	fmt.Println("HELD CARDS", heldCards)

	var combinations []*Cards
	for i := 0; i < 13; i++ {
		for j := 1; i < 13; i++ {
			for k := 2; i < 13; i++ {
				combinations = append(combinations, &Cards{heldCards[i], heldCards[j], heldCards[k]})
			}
		}
	}
	// Whatever modification to the state after passing needs to create the first trick, with two of clubs

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

