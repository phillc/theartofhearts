package main

import (
	// "time"
	"./lib/AgentVsAgent"
  // "fmt"
)

func passCard(round Round) []*AgentVsAgent.Card {
	return round.dealt[0:3]
}

