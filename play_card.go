package main

import (
	"time"
)

type PlayEvaluation struct {
	card Card
	value int
}

func playCard(trick *Trick) *Card {
	timeout := time.After(800 * time.Millisecond)
	evalCh := make(chan PlayEvaluation)
	evaluations := make(map[Card]PlayEvaluation)
	round := trick.round
	game := round.game
	position := (Position)(game.info["position"].(string))
	roundState := buildRoundState(round)

	numEvals := evaluatePlays(roundState, position, evalCh)

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

	trick.log("Number of evaluations:", len(evaluations), evaluations)
	var play *PlayEvaluation
	for _, evaluation := range evaluations {
		trick.log("eval:", evaluation.card, evaluation.value)
		if play == nil || evaluation.value >= play.value {
			play = new(PlayEvaluation)
			*play = evaluation
		}
	}

	return &play.card
}

func evaluatePlays(roundState *RoundState, position Position, evalCh chan PlayEvaluation) int {
	cards := roundState.playableCards()
	log(">>> PLAYABLE CARDS:", cards)
	for _, card := range cards {
		go func(card Card) {
			evalCh <- PlayEvaluation{card, evaluatePlay(roundState, position, card)}
		} (*card)
	}
	return len(cards)
}

func evaluatePlay(roundState *RoundState, position Position, card Card) int {
	log(">>>>>>>>>>evaluating play of", card)
	newRoundState := roundState.clone()
	newRoundState.play(card)
	simulation := Simulation{ roundState: newRoundState }

	for i := 0; i < 2; i++ {
	// for i := 0; i < (8 - len(newRoundState.currentTrick().played)); i++ {
		log("ADVANCE", i)
		simulation.advance()
	}
	log("simulating", simulation)
	evaluation := newRoundState.evaluate(position)
	simulatedEvaluation := simulation.evaluate(position)
	log("aaaaaaaaaaaaa!!done. eval:", evaluation, "simulated:", simulatedEvaluation)
	log("Scores: ", newRoundState.scores())
	return simulatedEvaluation
}

