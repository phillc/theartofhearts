package main

import (
	"./lib/AgentVsAgent"
  "fmt"
)

func isLeadingTrick(trick *Trick) bool {
	return len(trick.played) == 0
}

func isHeartsBroken(trick *Trick) bool {
	broken := false
	for h := 0; h < len(trick.round.tricks); h++ {
		cards := trick.round.tricks[h].played
		for i := 0; i < len(cards); i++ {
			if cards[i].Suit == AgentVsAgent.Suit_HEARTS {
				broken = true
			}
		}
	}
	return broken
}

func onlyTwoClubs(cards []*AgentVsAgent.Card) []*AgentVsAgent.Card {
	var matchedCards []*AgentVsAgent.Card
	for i := 0; i < len(cards); i++ {
		if cards[i].Suit == AgentVsAgent.Suit_CLUBS && cards[i].Rank == AgentVsAgent.Rank_TWO {
			matchedCards = append(matchedCards, cards[i])
		}
	}
	return matchedCards
}

func noHearts(cards []*AgentVsAgent.Card) []*AgentVsAgent.Card {
	var matchedCards []*AgentVsAgent.Card
	for i := 0; i < len(cards); i++ {
		if cards[i].Suit != AgentVsAgent.Suit_HEARTS {
			matchedCards = append(matchedCards, cards[i])
		}
	}
	return matchedCards
}

func noPoints(allCards []*AgentVsAgent.Card) []*AgentVsAgent.Card {
	var matchedCards []*AgentVsAgent.Card
	cards := noHearts(allCards)
	for i := 0; i < len(cards); i++ {
		if !(cards[i].Suit == AgentVsAgent.Suit_SPADES && cards[i].Rank == AgentVsAgent.Rank_QUEEN) {
			matchedCards = append(matchedCards, cards[i])
		}
	}
	return matchedCards
}

func followSuit(cards []*AgentVsAgent.Card, trick *Trick) []*AgentVsAgent.Card {
	var matchedCards []*AgentVsAgent.Card
	suit := trick.played[0].Suit
	for i := 0; i < len(cards); i++ {
		if cards[i].Suit == suit {
			matchedCards = append(matchedCards, cards[i])
		}
	}
	if len(matchedCards) == 0 {
		matchedCards = cards
	}
	return matchedCards
}

func playableCards(trick *Trick) []*AgentVsAgent.Card {
	validCards := trick.round.held

	if trick.number == 1 && isLeadingTrick(trick) {
		validCards = onlyTwoClubs(validCards)
	}

	if trick.number == 1 {
		validCards = noPoints(validCards)
	}

	if isLeadingTrick(trick) && !isHeartsBroken(trick) && len(noHearts(trick.round.held)) > 0 {
		validCards = noHearts(validCards)
	}

	if !isLeadingTrick(trick) {
		validCards = followSuit(validCards, trick)
	}

	trick.log("Valid cards:", validCards)
	return validCards
}

func doPassCards(round Round) []*AgentVsAgent.Card {
	cardsToPass := passCard(round)
	round.log("Passing cards", cardsToPass)

	return cardsToPass
}

func doPlayCard(trick Trick) *AgentVsAgent.Card {
	trick.log("Current trick:", trick)
	cardToPlay := pickCard(&trick)
	trick.log("Playing card:", cardToPlay)
	return cardToPlay
}

func main() {
	play(doPassCards, doPlayCard)
}

type Card struct {
	*AgentVsAgent.Card
}

func (card Card) order() int8 {
	rank := card.Rank
	switch rank {
	case AgentVsAgent.Rank_TWO: return 1
	case AgentVsAgent.Rank_THREE: return 2
	case AgentVsAgent.Rank_FOUR: return 3
	case AgentVsAgent.Rank_FIVE: return 4
	case AgentVsAgent.Rank_SIX: return 5
	case AgentVsAgent.Rank_SEVEN: return 6
	case AgentVsAgent.Rank_EIGHT: return 7
	case AgentVsAgent.Rank_NINE: return 8
	case AgentVsAgent.Rank_TEN: return 9
	case AgentVsAgent.Rank_JACK: return 10
	case AgentVsAgent.Rank_QUEEN: return 11
	case AgentVsAgent.Rank_KING: return 12
	case AgentVsAgent.Rank_ACE: return 13
	}

	fmt.Println("*********Rank not found********")
	return 0
}

