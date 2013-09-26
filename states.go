package main

type Position string

func buildRoundState(round *Round) *RoundState {
	var trickStates []*TrickState
	for _, trick := range round.tricks {
		trickStates = append(trickStates, buildTrickState(trick))
	}

	rootPosition := (Position)(round.game.info.Position)
	players := make(map[Position]*PlayerState, 4)
	players["north"] = &PlayerState{ actions: make(map[Card]Action, 13) }
	players["east"] = &PlayerState{ actions: make(map[Card]Action, 13) }
	players["south"] = &PlayerState{ actions: make(map[Card]Action, 13) }
	players["west"] = &PlayerState{ actions: make(map[Card]Action, 13) }
	rootPlayer := players[rootPosition]
	rootPlayer.root = true

	passingTo := rootPlayer
	receivedFrom := rootPlayer
	positions := []Position{"north", "east", "south", "west"}
	rootIndex := -1
	for i, position := range positions {
		if position == rootPosition {
			rootIndex = i
			break
		}
	}
	positionsFromRoot := append(positions[rootIndex:4], positions[0:rootIndex]...)
	switch (round.number - 1) % 4 {
	case 0:
		// left
		passingTo = players[positionsFromRoot[3]]
		receivedFrom = players[positionsFromRoot[1]]
	case 1:
		// right
		passingTo = players[positionsFromRoot[1]]
		receivedFrom = players[positionsFromRoot[3]]
	case 2:
		// across
		passingTo = players[positionsFromRoot[2]]
		receivedFrom = players[positionsFromRoot[2]]
	}

	for _, aCard := range round.dealt {
		card := Card{ suit: aCard.Suit, rank: aCard.Rank }
		rootPlayer.dealt(card)
	}
	for _, aCard := range round.passed {
		card := Card{ suit: aCard.Suit, rank: aCard.Rank }
		rootPlayer.passed(card)
		passingTo.received(card)
	}
	for _, aCard := range round.received {
		card := Card{ suit: aCard.Suit, rank: aCard.Rank }
		rootPlayer.received(card)
		receivedFrom.passed(card)
	}

	for _, trickState := range trickStates {
		for index, position := range trickState.positionsFromLeader()[0:len(trickState.played)] {
			playedCard := *trickState.played[index]

			leadSuit := trickState.played[0].suit
			if playedCard.suit != leadSuit {
				players[position].discardedOn(leadSuit)
			}

			players[position].played(playedCard)
		}
	}

	return &RoundState{
		number: round.number,
		trickStates: trickStates,
		north: players["north"],
		east: players["east"],
		south: players["south"],
		west: players["west"],
	}
}

func buildTrickState(trick *Trick) *TrickState {
	var playedCards Cards
	for _, aCard := range trick.played {
		playedCards = append(playedCards, &Card{ suit: aCard.Suit, rank: aCard.Rank })
	}
	return &TrickState{ number: trick.number, leader: (Position)(trick.leader), played: playedCards }
}
