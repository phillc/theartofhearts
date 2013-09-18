package main

import (
	"./lib/AgentVsAgent"
)

func doPassCards(round Round) []*AgentVsAgent.Card {
	cardsToPass := passCards(round)
	round.log("Passing cards", cardsToPass)

	return cardsToPass
}

func doPlayCard(trick Trick) *AgentVsAgent.Card {
	trick.log("Current trick:", trick.number, &trick.round, trick.leader, trick.played)
	cardToPlay := playCard(&trick)
	trick.log("Playing card:", cardToPlay)
	return cardToPlay
}

func main() {
	play(doPassCards, doPlayCard)
}

