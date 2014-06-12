package main

func (card Card) order() int {
	rank := card.Rank
	switch rank {
	case TWO: return 1
	case THREE: return 2
	case FOUR: return 3
	case FIVE: return 4
	case SIX: return 5
	case SEVEN: return 6
	case EIGHT: return 7
	case NINE: return 8
	case TEN: return 9
	case JACK: return 10
	case QUEEN: return 11
	case KING: return 12
	case ACE: return 13
	}

	log("*********Rank not found********")
	return 0
}

func (card Card) String() string {
	return "::" + card.Rank + " of " + card.Suit + "::"
}

func (card Card) score() int {
	value := 0
	if card.Suit == HEARTS {
		value = 1
	} else if card.Suit == SPADES && card.Rank == QUEEN {
		value = 13
	}
	return value
}

type Cards []*Card

func (s Cards) Len() int { return len(s) }
func (s Cards) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (cards *Cards) allOfSuit(suit string) Cards {
	newCards := Cards{}
	for _, card := range *cards {
		if card.Suit == suit {
			newCards = append(newCards, card)
		}
	}
	return newCards
}

func (cards *Cards) onlyTwoClubs() Cards {
	matchedCards := Cards{}
	for _, card := range *cards {
		if card.Suit == CLUBS && card.Rank == TWO {
			matchedCards = append(matchedCards, card)
		}
	}
	return matchedCards
}

func (cards *Cards) noHearts() Cards {
	matchedCards := Cards{}
	for _, card := range *cards {
		if card.Suit != HEARTS {
			matchedCards = append(matchedCards, card)
		}
	}
	return matchedCards
}

func (cards *Cards) noPoints() Cards {
	matchedCards := Cards{}
	for _, card := range cards.noHearts() {
		if !(card.Suit == SPADES && card.Rank == QUEEN) {
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

	suits := []string{ CLUBS, DIAMONDS, SPADES, HEARTS }
	ranks := []string{
		TWO,
		THREE,
		FOUR,
		FIVE,
		SIX,
		SEVEN,
		EIGHT,
		NINE,
		TEN,
		JACK,
		QUEEN,
		KING,
		ACE,
	}

	for _, suit := range suits {
		for _, rank := range ranks {
			cards = append(cards, &Card{ Suit: suit, Rank: rank })
		}
	}

	return cards
}
