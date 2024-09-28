package engine

import (
	"math/rand/v2"
	"slices"
)

func NewRandom(rule uint8, width, height int, chance float64) Game {
	if chance <= 0 {
		return NewBlank(rule, width, height)
	} else if chance >= 1 {
		return NewFull(rule, width, height)
	}

	initial := make([]bool, width)
	for i := range initial {
		initial[i] = randomBool(chance)
	}
	return New(rule, width, height, initial)
}

func randomBool(chance float64) bool {
	if rand.Float64() < chance {
		return true
	}
	return false
}

func NewBlank(rule uint8, width, height int) Game {
	return New(rule, width, height, make([]bool, width))
}

func NewFull(rule uint8, width, height int) Game {
	return New(rule, width, height, slices.Repeat([]bool{true}, width))
}

func NewCenterBlock(rule uint8, width, height int) Game {
	initial := make([]bool, width)
	initial[width/2] = true
	return New(rule, width, height, initial)
}
