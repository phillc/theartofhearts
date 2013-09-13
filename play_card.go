package main

import (
	"time"
	"./lib/AgentVsAgent"
  // "fmt"
)

type PickEvaluation struct {
	card Card
	score int8
}

func pickCard(trick *Trick) *AgentVsAgent.Card {
	timeout := time.After(800 * time.Millisecond)
	cards := playableCards(trick)
	evalCh := make(chan PickEvaluation)
	evaluations := make(map[Card]PickEvaluation)

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
	var pick PickEvaluation
	for _, evaluation := range evaluations {
		trick.log("eval:", evaluation.card, evaluation.score)
		if evaluation.score >= pick.score {
			trick.log("winning evaluation")
			pick = evaluation
		}
		trick.log("current pick:", pick)
	}

	return pick.card.Card
}

func evaluateTrick(cards []*AgentVsAgent.Card, trick *Trick, evalCh chan PickEvaluation) {
	for _, card := range cards {
		go func(card Card) {
			evalCh <- evaluatePick(card, *trick)
		} (Card{card})
	}
}

func evaluatePick(card Card, trick Trick) PickEvaluation {
	trick.log("evaluating play of", card)
	return PickEvaluation{card, card.order()}
}
