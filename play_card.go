package main

import (
	"time"
	"./lib/AgentVsAgent"
  // "fmt"
)

type PlayEvaluation struct {
	card Card
	value int
}

func playCard(trick *Trick) *AgentVsAgent.Card {
	timeout := time.After(800 * time.Millisecond)
	evalCh := make(chan PlayEvaluation)
	evaluations := make(map[Card]PlayEvaluation)

	numEvals := evaluateTrick(trick, evalCh)

	for i := 0; i < numEvals; i++ {
		trick.log("Waiting for an evaluation")
		select {
		case cardEval := <-evalCh:
			trick.log("Card", cardEval.card, "evaluated at", cardEval.value)
			evaluations[cardEval.card] = cardEval
		case <- timeout:
			trick.log("*****Timeout*****")
			trick.log("*****Timeout*****")
			trick.log("*****Timeout*****")
			break
		}
	}

	trick.log("Number of evaluations:", len(evaluations))
	var play PlayEvaluation
	for _, evaluation := range evaluations {
		trick.log("eval:", evaluation.card, evaluation.value)
		if evaluation.value >= play.value {
			trick.log("winning evaluation")
			play = evaluation
		}
		trick.log("current play:", play)
	}

	return play.card.Card
}

// maybe be evaluate round... to determine the value of a round given the state of everything.
// that way, the pass cards logic can use evaluate round to determine position
func evaluateTrick(trick *Trick, evalCh chan PlayEvaluation) int {
	game := trick.round.game
	gameState := buildGameState(game)
	cards := playableCards(trick)
	for _, card := range cards {
		go func(card Card) {
			evalCh <- PlayEvaluation{card, evaluatePlay(*gameState, &card, trick)}
		} (Card{card})
	}
	return len(cards)
}

func evaluatePlay(gameState GameState, card *Card, trick *Trick) int {
	// position := (Position)(game.info.Position)
	trick.log("evaluating play of", card)
	return card.order()
}

func buildGameState(game *Game) *GameState {
	var scores map[Position]int
	roundState := buildRoundState(game.rounds[len(game.rounds) - 1])
	return &GameState{ round: roundState, scores: scores }
}

func buildRoundState(round *Round) *RoundState {
	trickState := buildTrickState(round.tricks[len(round.tricks) - 1])
	return &RoundState{ trick: trickState }
}

func buildTrickState(trick *Trick) *TrickState {
	var cards Cards
	for _, aCard := range trick.played {
		cards = append(cards, &Card{aCard})
	}
	return &TrickState{ leader: (Position)(trick.leader), played: cards }
}
