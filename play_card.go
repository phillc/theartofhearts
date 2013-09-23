package main

import (
	"time"
	"./lib/AgentVsAgent"
  "fmt"
)

type PlayEvaluation struct {
	card Card
	value int
}

func playCard(trick *Trick) *AgentVsAgent.Card {
	timeout := time.After(800 * time.Millisecond)
	evalCh := make(chan PlayEvaluation)
	evaluations := make(map[Card]PlayEvaluation)
	game := trick.round.game
	position := (Position)(game.info.Position)
	gameState := buildGameState(game)

	numEvals := evaluatePlays(gameState, position, evalCh)

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

	return play.card.toAvA()
}

func evaluatePlays(gameState *GameState, position Position, evalCh chan PlayEvaluation) int {
	cards := gameState.currentRound().playableCards()
	fmt.Println(">>> PLAYABLE CARDS:", cards)
	for _, card := range cards {
		go func(card Card) {
			evalCh <- PlayEvaluation{card, evaluatePlay(gameState, position, card)}
		} (*card)
	}
	return len(cards)
}

func evaluatePlay(gameState *GameState, position Position, card Card) int {
	fmt.Println(">>>>>>>>>>evaluating play of", card)
	newGameState := gameState.play(position, card)
	return newGameState.evaluate(position)
}

