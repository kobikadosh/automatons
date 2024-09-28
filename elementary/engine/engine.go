package engine

import "math"

type Game interface {
	Rows() int
	Columns() int
	Get(row, col int) bool
	Tick() (bool, Game)
}

type game struct {
	rule    uint8
	board   [][]bool
	currRow int
}

func New(rule uint8, width, height int, initial []bool) Game {
	if len(initial) != width {
		panic("initial state should be the same length as width")
	}

	board := make([][]bool, height)
	board[0] = initial
	for row := 1; row < height; row++ {
		board[row] = make([]bool, width)
	}

	return game{
		rule:  rule,
		board: board,
	}
}

func (g game) Rows() int {
	return len(g.board)
}

func (g game) Columns() int {
	if g.Rows() == 0 {
		return 0
	}
	return len(g.board[0])
}

func (g game) Get(row, col int) bool {
	return g.board[row][col]
}

func (g game) Tick() (eof bool, game Game) {
	if g.currRow >= g.Rows()-1 {
		return true, g
	}

	prevRow := g.board[g.currRow]
	g.currRow++
	for col := range prevRow {
		g.board[g.currRow][col] = g.calculateCell(prevRow, col)
	}

	return false, g
}

func (g game) calculateCell(row []bool, col int) bool {
	length := len(row)
	prevCol := (length + col - 1) % length
	nextCol := (col + 1) % length
	pattern := getPatternValue(row[prevCol], row[col], row[nextCol])
	return g.rule&pattern == pattern
}

func getPatternValue(cells ...bool) uint8 {
	v := 0
	for i := range cells {
		if cells[i] {
			v |= 1 << i
		}
	}
	return uint8(math.Pow(2, float64(v)))
}
