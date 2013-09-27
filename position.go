package main

type Position string

func allPositions() []Position {
	return []Position{"north", "east", "south", "west"}
}

func otherPositions(excludedPosition Position) []Position {
	positions := []Position{}
	for _, position := range allPositions() {
		if position != excludedPosition {
			positions = append(positions, position)
		}
	}
	return positions
}

