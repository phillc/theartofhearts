package main

type RoundState struct {
	number int
	north *PlayerState
	east *PlayerState
	south *PlayerState
	west *PlayerState
	trickStates []*TrickState
}

func (roundState *RoundState) playerState(position Position) *PlayerState {
	switch string(position) {
	case "north": return roundState.north
	case "east": return roundState.east
	case "south": return roundState.south
	case "west": return roundState.west
	}
	return &PlayerState{}
}

func (roundState *RoundState) currentTrick() *TrickState {
	var trickState *TrickState
	if len(roundState.trickStates) > 0 {
		trickState = roundState.trickStates[len(roundState.trickStates) - 1]
	}
	return trickState
}

func (roundState *RoundState) isHeartsBroken() bool {
	broken := false
	for _, trick := range roundState.trickStates {
		cards := trick.played
		for _, card := range cards {
			if card.Suit == HEARTS {
				broken = true
				break
			}
		}
	}
	return broken
}

func (roundState *RoundState) playableCards() Cards {
	position := roundState.currentTrick().positionsMissing()[0]
	held := roundState.playerState(position).definitelyHeld()
  return roundState.playableCardsOutOf(held)
}

func (roundState *RoundState) playableCardsOutOf(startingCards Cards) Cards {
	held := Cards{}
	validCards := Cards{}
	for _, card := range startingCards {
		held = append(held, card)
		validCards = append(validCards, card)
	}
	trick := roundState.currentTrick()

	if trick.number == 1 && trick.isLeading() {
		validCards = validCards.onlyTwoClubs()
	}

	if trick.number == 1 {
		validCards = validCards.noPoints()
	}

	if trick.isLeading() && !roundState.isHeartsBroken() && len(held.noHearts()) > 0 {
		validCards = validCards.noHearts()
	}

	if !trick.isLeading() {
		newValidCards := validCards.allOfSuit(trick.played[0].Suit)
		if len(newValidCards) > 0 {
			validCards = newValidCards
		}
	}
	return validCards
}

func (roundState *RoundState) playableCardProbabilities() map[Card]int {
	position := roundState.currentTrick().positionsMissing()[0]

	knowledge := Knowledge{}
	knowledge.buildFrom(roundState)
	possiblyHeldCards := knowledge.possiblyHeldCardsFor(position)

	playableProbabilities := make(map[Card]int)
	for _, card := range roundState.playableCardsOutOf(possiblyHeldCards) {
		playableProbabilities[*card] = 100
	}

	return playableProbabilities
}

func (roundState *RoundState) evaluate(position Position) int {
	evaluation := 0
	// evaluation = evaluation - roundState.scores[position]
	handScore := 0

	// Take the average of each suit?
	// or something that promotes lower cards (2 + K > 7)? or is it?
	// how about (sum / len) - (len * 3)

	for card, action := range roundState.playerState(position).actions {
		if action.isDefinitelyHeld() {
			// todo: two of clubs doesn't matter if we can just simulate past a couple tricks
			if card.Suit == CLUBS && card.Rank == TWO {
				handScore = handScore - 13
			} else {
				handScore = handScore - card.order()
			}
		}
	}
	evaluation = evaluation + handScore

	for scorePosition, score := range roundState.scores() {
		if scorePosition == position {
			evaluation = evaluation - (score * 10)
		} else {
			evaluation = evaluation + (score * 3)
		}
	}

	return evaluation
}

func (roundState *RoundState) scores() map[Position]int {
	scores := make(map[Position]int, 4)

	for _, trickState := range roundState.trickStates {
		if len(trickState.played) > 0 {
			position := trickState.winner()
			scores[position] = scores[position] + trickState.score()
		}
	}
	for position, score := range scores {
		if score == 26 {
			scores["north"] = 26
			scores["east"] = 26
			scores["south"] = 26
			scores["west"] = 26
			scores[position] = 0
			break
		}
	}

	return scores
}

func (roundState *RoundState) clone() *RoundState {
	var newTrickStates []*TrickState
	for _, trickState := range roundState.trickStates {
	  newTrickStates = append(newTrickStates, trickState.clone())
	}

	newRoundState := *roundState
	newRoundState.trickStates = newTrickStates
	newRoundState.north = roundState.north.clone()
	newRoundState.east = roundState.east.clone()
	newRoundState.south = roundState.south.clone()
	newRoundState.west = roundState.west.clone()
	return &newRoundState
}

func (roundState *RoundState) pass(position Position, cards Cards) {
	playerState := roundState.playerState(position)
	actions := playerState.actions
	for _, passedCard := range cards {
		action := actions[*passedCard]
		action.passed = true
		actions[*passedCard] = action
	}
}

func (roundState *RoundState) play(card Card) {
	currentTrick := roundState.currentTrick()
	position := currentTrick.positionsMissing()[0]

	roundState.playerState(position).played(card)
	currentTrick.played = append(currentTrick.played, &card)

	roundState.nextTrick()
}

func (roundState *RoundState) nextTrick() {
	if len(roundState.currentTrick().played) == 4 && len(roundState.trickStates) != 13 {
		leader := roundState.currentTrick().winner()
		newTrickState := TrickState{ number: len(roundState.trickStates) + 1, leader: leader, played: Cards{} }
		roundState.trickStates = append(roundState.trickStates, &newTrickState)
	}
}

