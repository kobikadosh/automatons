package engine

type neighborFunc func(g game, row, col int) Cell

var neighborhoodFuncs = []neighborFunc{
	newNeighborFunc(-1, -1),
	newNeighborFunc(-1, 0),
	newNeighborFunc(-1, 1),
	newNeighborFunc(0, -1),
	newNeighborFunc(0, 1),
	newNeighborFunc(1, -1),
	newNeighborFunc(1, 0),
	newNeighborFunc(1, 1),
}

func newNeighborFunc(rowDiff, colDiff int) neighborFunc {
	return func(g game, row, col int) Cell {
		row, col = row+rowDiff, col+colDiff
		if row >= 0 && row < len(g) && col >= 0 && col < len(g[row]) {
			return g[row][col]
		}
		return false
	}
}
