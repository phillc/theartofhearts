package main

import (
	"./lib/AgentVsAgent"
	"fmt"
)

type Card struct {
	*AgentVsAgent.Card
}

func (card Card) order() int {
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

func (card Card) score() int {
	value := 0
	if card.Suit == AgentVsAgent.Suit_HEARTS {
		value = 1
	} else if card.Suit == AgentVsAgent.Suit_SPADES && card.Rank == AgentVsAgent.Rank_QUEEN {
		value = 13
	}
	return value
}

type Cards []*Card

func (s Cards) Len() int { return len(s) }
func (s Cards) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s Cards) allOfSuit(suit AgentVsAgent.Suit) Cards {
	newCards := Cards{}
	for _, card := range s {
		if card.Suit == suit {
			newCards = append(newCards, card)
		}
	}
	return newCards
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

