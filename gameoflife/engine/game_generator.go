package engine

import "math/rand/v2"

func NewRandomGame(width, height int, chanceForLife float64) Game {
	g := make(game, height+invisibleMarginsSize*2)
	for row := range g {
		g[row] = make([]Cell, width+invisibleMarginsSize*2)
		for col := range g[row] {
			g[row][col] = randomCell(chanceForLife)
		}
	}
	return g
}

func randomCell(chanceForLife float64) Cell {
	if rand.Float64() < chanceForLife {
		return true
	}
	return false
}
