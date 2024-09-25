package engine

const invisibleMarginsSize = 10

type Game interface {
	Rows() int
	Columns() int
	GetCell(row, col int) Cell
	Tick()
}

type Cell bool

type game [][]Cell

func (g game) Rows() int {
	return len(g) - invisibleMarginsSize*2
}

func (g game) Columns() int {
	if g.Rows() == 0 {
		return 0
	}
	return len(g[0]) - invisibleMarginsSize*2
}

func (g game) GetCell(row, col int) Cell {
	return g[row+invisibleMarginsSize][col+invisibleMarginsSize]
}

func (g game) Tick() {
	changes := make([]func(), 0)
	for row := range g {
		for col := range g[row] {
			if change := g.calculateCellChange(row, col); change != nil {
				changes = append(changes, change)
			}
		}
	}

	// make changes to the game only after all cells were calculated
	for _, change := range changes {
		change()
	}
}

func (g game) calculateCellChange(row, col int) func() {
	c := g[row][col]
	pop := g.getNeighborhoodPopulation(row, col)

	if c && (pop < 2 || pop > 3) {
		// under/overpopulation, kill cell
		return g.setCellFunc(row, col, false)
	}
	if !c && pop == 3 {
		// reproduction, revive cell
		return g.setCellFunc(row, col, true)
	}
	return nil
}

// create a cell modification func, for future execution
func (g game) setCellFunc(row, col int, val Cell) func() {
	return func() {
		g[row][col] = val
	}
}

// get the number of live cells in the neighborhood
func (g game) getNeighborhoodPopulation(row, col int) int {
	population := 0
	for _, f := range neighborhoodFuncs {
		if f(g, row, col) {
			population++
		}
	}

	return population
}
