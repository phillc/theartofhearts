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

	trick.log("Number of evaluations:", len(evaluations), evaluations)
	var play *PlayEvaluation
	for _, evaluation := range evaluations {
		trick.log("eval:", evaluation.card, evaluation.value)
		if play == nil || evaluation.value >= play.value {
			play = new(PlayEvaluation)
			*play = evaluation
		}
	}

	return play.card.Card
}

func evaluateTrick(trick *Trick, evalCh chan PlayEvaluation) int {
	game := trick.round.game
	position := (Position)(game.info.Position)
	gameState := buildGameState(game)
	cards := playableCards(trick)
	for _, card := range cards {
		go func(card Card) {
			evalCh <- PlayEvaluation{card, evaluatePlay(gameState, position, card)}
		} (Card{card})
	}
	return len(cards)
}

func evaluatePlay(gameState GameState, position Position, card Card) int {
	fmt.Println(">>>>>>>>>>evaluating play of", card)
	// fmt.Println("was:", card.order())

	/*fmt.Println("pre ???", gameState.round.players[position].held[card])*/
	/*fmt.Println("pre ref:", &gameState.round)*/
	/*fmt.Println("pre gs:", gameState.round.players)*/
	// fmt.Println("pre now:", gameState.evaluate(position))

	newGameState := gameState.play(position, card)

	/*fmt.Println("???", gameState.round.players[position].held[card])*/
	/*fmt.Println("ref:", &gameState.round)*/
	/*fmt.Println("gs:", gameState.round.players)*/
	// fmt.Println("now:", gameState.evaluate(position))

	/*fmt.Println("new gs:", newGameState.round.players)*/
	/*fmt.Println("new ref:", &gameState.round)*/
	/*fmt.Println("new???", newGameState.round.players[position].held[card])*/
	// fmt.Println("new now:", newGameState.evaluate(position))
	return newGameState.evaluate(position)
}

