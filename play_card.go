package main

import (
	"time"
	"./lib/AgentVsAgent"
  // "fmt"
)

type PlayEvaluation struct {
	card Card
	score int8
}

func playCard(trick *Trick) *AgentVsAgent.Card {
	timeout := time.After(800 * time.Millisecond)
	cards := playableCards(trick)
	evalCh := make(chan PlayEvaluation)
	evaluations := make(map[Card]PlayEvaluation)

	go evaluateTrick(cards, trick, evalCh)

	for i := 0; i < len(cards); i++ {
		trick.log("Waiting for an evaluation")
		select {
		case cardEval := <-evalCh:
			trick.log("Card", cardEval.card, "evaluated at", cardEval.score)
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
		trick.log("eval:", evaluation.card, evaluation.score)
		if evaluation.score >= play.score {
			trick.log("winning evaluation")
			play = evaluation
		}
		trick.log("current play:", play)
	}

	return play.card.Card
}

// maybe be evaluate round... to determine the value of a round given the state of everything.
// that way, the pass cards logic can use evaluate round to determine position
func evaluateTrick(cards []*AgentVsAgent.Card, trick *Trick, evalCh chan PlayEvaluation) {
	for _, card := range cards {
		go func(card Card) {
			evalCh <- evaluatePlay(card, *trick)
		} (Card{card})
	}
}

func evaluatePlay(card Card, trick Trick) PlayEvaluation {
	trick.log("evaluating play of", card)
	return PlayEvaluation{card, card.order()}
}
