package main

import (
	"./lib/AgentVsAgent"
	"fmt"
)

type Card struct {
	suit AgentVsAgent.Suit
	rank AgentVsAgent.Rank
}

func (card Card) toAvA() *AgentVsAgent.Card {
	avaCard := AgentVsAgent.Card{ Suit: card.suit, Rank: card.rank }
	return &avaCard
}

func (card Card) order() int {
	rank := card.rank
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

func (card Card) score() int {
	value := 0
	if card.suit == AgentVsAgent.Suit_HEARTS {
		value = 1
	} else if card.suit == AgentVsAgent.Suit_SPADES && card.rank == AgentVsAgent.Rank_QUEEN {
		value = 13
	}
	return value
}

type Cards []*Card

func (s Cards) Len() int { return len(s) }
func (s Cards) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (cards *Cards) allOfSuit(suit AgentVsAgent.Suit) Cards {
	newCards := Cards{}
	for _, card := range *cards {
		if card.suit == suit {
			newCards = append(newCards, card)
		}
	}
	return newCards
}

func (cards *Cards) onlyTwoClubs() Cards {
	matchedCards := Cards{}
	for _, card := range *cards {
		if card.suit == AgentVsAgent.Suit_CLUBS && card.rank == AgentVsAgent.Rank_TWO {
			matchedCards = append(matchedCards, card)
		}
	}
	return matchedCards
}

func (cards *Cards) noHearts() Cards {
	matchedCards := Cards{}
	for _, card := range *cards {
		if card.suit != AgentVsAgent.Suit_HEARTS {
			matchedCards = append(matchedCards, card)
		}
	}
	return matchedCards
}

func (cards *Cards) noPoints() Cards {
	matchedCards := Cards{}
	for _, card := range cards.noHearts() {
		if !(card.suit == AgentVsAgent.Suit_SPADES && card.rank == AgentVsAgent.Rank_QUEEN) {
			matchedCards = append(matchedCards, card)
		}
	}
	return matchedCards
}

func (s Cards) indexOf(card *Card) int {
	for i, aCard := range s {
		if card == aCard {
			return i
		}
	}
	return -1
}

type ByOrder struct{ Cards }
func (s ByOrder) Less(i, j int) bool { return s.Cards[i].order() < s.Cards[j].order() }

func allCards() Cards {
	cards := Cards{}

	suits := []AgentVsAgent.Suit{ AgentVsAgent.Suit_CLUBS, AgentVsAgent.Suit_DIAMONDS, AgentVsAgent.Suit_SPADES, AgentVsAgent.Suit_HEARTS }
	ranks := []AgentVsAgent.Rank{
		AgentVsAgent.Rank_TWO,
		AgentVsAgent.Rank_THREE,
		AgentVsAgent.Rank_FOUR,
		AgentVsAgent.Rank_FIVE,
		AgentVsAgent.Rank_SIX,
		AgentVsAgent.Rank_SEVEN,
		AgentVsAgent.Rank_EIGHT,
		AgentVsAgent.Rank_NINE,
		AgentVsAgent.Rank_TEN,
		AgentVsAgent.Rank_JACK,
		AgentVsAgent.Rank_QUEEN,
		AgentVsAgent.Rank_KING,
		AgentVsAgent.Rank_ACE,
	}

	for _, suit := range suits {
		for _, rank := range ranks {
			cards = append(cards, &Card{ suit: suit, rank: rank })
		}
	}

	return cards
}
