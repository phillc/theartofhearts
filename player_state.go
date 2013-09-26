package main

import (
	"./lib/AgentVsAgent"
)

type Action struct {
	dealt bool
	played bool
	passed bool
	received bool
	// can't have
}

func (action Action) String() string {
	str := "Action:<"
	if action.dealt {
		str = str + " dealt"
	}
	if action.played {
		str = str + " played"
	}
	if action.passed {
		str = str + " passed"
	}
	if action.received {
		str = str + " received"
	}
	if action.isDefinitelyHeld() {
		str = str + " :definitely held:"
	}
	str = str + " >"
	return str
}

func (action *Action) isDefinitelyHeld() bool {
	return !action.played && ((action.dealt && !action.passed) || action.received)
}

type PlayerState struct {
	actions map[Card]Action
	root bool
	emptySuits []AgentVsAgent.Suit // uniqueness?
}

func (playerState *PlayerState) definitelyHeld() Cards {
	cards := Cards{}
	for card, action := range playerState.actions {
		if action.isDefinitelyHeld() {
			aCard := card
			cards = append(cards, &aCard)
		}
	}
	return cards
}

func (playerState *PlayerState) clone() *PlayerState {
	newActions := make(map[Card]Action)
	for card, action := range playerState.actions {
		newActions[card] = action
	}

	newPlayerState := *playerState
	newPlayerState.actions = newActions
	return &newPlayerState
}

func (playerState *PlayerState) received(card Card) {
	action := playerState.actions[card]
	action.received = true
	playerState.actions[card] = action
}

func (playerState *PlayerState) played(card Card) {
	action := playerState.actions[card]
	action.played = true
	playerState.actions[card] = action
}

func (playerState *PlayerState) dealt(card Card) {
	action := playerState.actions[card]
	action.dealt = true
	playerState.actions[card] = action
}

func (playerState *PlayerState) passed(card Card) {
	action := playerState.actions[card]
	action.passed = true
	playerState.actions[card] = action
}

func (playerState *PlayerState) discardedOn(suit AgentVsAgent.Suit) {
	playerState.emptySuits = append(playerState.emptySuits, suit)
}
