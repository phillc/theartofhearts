package main

import (
	"time"
	"./lib/AgentVsAgent"
	"fmt"
	"sort"
)

type PassEvaluation struct {
	number int
	cards Cards
	value int
}

func passCards(round Round) []*AgentVsAgent.Card {
	fmt.Println("passing cards")
	timeout := time.After(800 * time.Millisecond)
	evalCh := make(chan PassEvaluation)
	evaluations := make(map[int]PassEvaluation)
	game := round.game
	position := (Position)(game.info.Position)
	gameState := buildGameState(game)

	numEvals := evaluatePasses(gameState, position, evalCh)

	for i := 0; i < numEvals; i++ {
		round.log("Waiting for a pass evaluation")
		select {
		case passEval := <-evalCh:
			round.log("Cards", passEval.cards, "evaluated at", passEval.value)
			evaluations[passEval.number] = passEval
		case <- timeout:
			round.log("*****Timeout*****")
			round.log("*****Timeout*****")
			round.log("*****Timeout*****")
			break
		}
	}

	round.log("Number of evaluations:", len(evaluations), evaluations)
	var pass *PassEvaluation
	for _, evaluation := range evaluations {
		round.log("eval:", evaluation.cards, evaluation.value)
		if pass == nil || evaluation.value >= pass.value {
			pass = new(PassEvaluation)
			*pass = evaluation
		}
	}

	var cardsToPass []*AgentVsAgent.Card
	for _, card := range pass.cards {
		cardsToPass = append(cardsToPass, card.toAvA())
	}

	return cardsToPass
}

func evaluatePasses(gameState *GameState, position Position, evalCh chan PassEvaluation) int {
	heldCards := gameState.currentRound().playerState(position).definitelyHeld()

	// too many combos right now, filter some out
	sort.Sort(sort.Reverse(ByOrder{heldCards}))
	heldCards = heldCards[0:10]

	var combinations []*Cards
	length := len(heldCards)
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			for k := 0; k < length; k++ {
				if i != j && j != k && i != k {
					combinations = append(combinations, &Cards{heldCards[i], heldCards[j], heldCards[k]})
				}
			}
		}
	}

	for number, cards := range combinations {
		go func(number int, cards Cards) {
			evalCh <- PassEvaluation{number, cards, evaluatePass(gameState, position, cards)}
		} (number, *cards)
	}
	return len(combinations)
}

func evaluatePass(gameState *GameState, position Position, cards Cards) int {
	fmt.Println(">>>>>>>>>>evaluating pass of", cards)
	newGameState := gameState.pass(position, cards)
	return newGameState.evaluate(position)
}
