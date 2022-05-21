package sudoku

// Node represents a node in the Dancing Links (DLX) matrix.
type Node struct {
	Left, Right, Up, Down *Node
	Col                   *Column
	RowVal, ColIdx, Digit int // Metadata representing the candidate (row, column, digit)
}

// Column represents a column header in the DLX matrix.
type Column struct {
	Head Node
	Size int
	ID   int
}

// cover removes a column from the DLX matrix, along with all rows that have a node in this column.
func cover(c *Column) {
	c.Head.Right.Left = c.Head.Left
	c.Head.Left.Right = c.Head.Right
	for i := c.Head.Down; i != &c.Head; i = i.Down {
		for j := i.Right; j != i; j = j.Right {
			j.Down.Up = j.Up
			j.Up.Down = j.Down
			j.Col.Size--
		}
	}
}

// uncover restores a column to the DLX matrix, along with all rows that have a node in this column.
func uncover(c *Column) {
	for i := c.Head.Up; i != &c.Head; i = i.Up {
		for j := i.Left; j != i; j = j.Left {
			j.Col.Size++
			j.Down.Up = j
			j.Up.Down = j
		}
	}
	c.Head.Right.Left = &c.Head
	c.Head.Left.Right = &c.Head
}

// selectColumn selects the column with the minimum size (number of active rows) to minimize branching.
func selectColumn(root *Node) *Column {
	var minCol *Column
	minSize := 999999
	for curr := root.Right; curr != root; curr = curr.Right {
		col := curr.Col
		if col.Size < minSize {
			minSize = col.Size
			minCol = col
		}
	}
	return minCol
}

// SolveExactCover solves the Sudoku puzzle using Knuth's Algorithm X / Dancing Links (DLX).
func SolveExactCover(board *[9][9]int) bool {
	// 1. Initialize root node
	var root Node

	// 2. Initialize 324 constraint columns
	columns := make([]*Column, 324)
	for i := 0; i < 324; i++ {
		columns[i] = &Column{
			ID: i,
		}
		columns[i].Head.Col = columns[i]
		columns[i].Head.Up = &columns[i].Head
		columns[i].Head.Down = &columns[i].Head
	}

	// Link root and columns horizontally
	curr := &root
	for i := 0; i < 324; i++ {
		curr.Right = &columns[i].Head
		columns[i].Head.Left = curr
		curr = &columns[i].Head
	}
	curr.Right = &root
	root.Left = curr

	// 3. Build the DLX matrix with all 729 candidates (9 rows * 9 columns * 9 digits)
	// We keep a lookup table to easily locate the first node of a specific row candidate.
	rowNodes := make([][][]*Node, 9)
	for r := 0; r < 9; r++ {
		rowNodes[r] = make([][]*Node, 9)
		for c := 0; c < 9; c++ {
			rowNodes[r][c] = make([]*Node, 10) // 1-indexed digit values (1..9)
		}
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			for v := 1; v <= 9; v++ {
				// Constraint Column indices
				c1 := r*9 + c                                      // Cell constraint (0..80)
				c2 := 81 + r*9 + (v - 1)                           // Row constraint (81..161)
				c3 := 162 + c*9 + (v - 1)                          // Column constraint (162..242)
				c4 := 243 + ((r/3)*3+c/3)*9 + (v - 1)              // Box constraint (243..323)

				nodes := make([]*Node, 4)
				for i, colIdx := range []int{c1, c2, c3, c4} {
					col := columns[colIdx]
					node := &Node{
						Col:    col,
						RowVal: r,
						ColIdx: c,
						Digit:  v,
					}
					nodes[i] = node

					// Insert node at the bottom of the column's vertical circular list
					last := col.Head.Up
					node.Down = &col.Head
					col.Head.Up = node
					node.Up = last
					last.Down = node

					col.Size++
				}

				// Link the 4 nodes of this row horizontally in a circular list
				nodes[0].Right = nodes[1]
				nodes[1].Left = nodes[0]
				nodes[1].Right = nodes[2]
				nodes[2].Left = nodes[1]
				nodes[2].Right = nodes[3]
				nodes[3].Left = nodes[2]
				nodes[3].Right = nodes[0]
				nodes[0].Left = nodes[3]

				rowNodes[r][c][v] = nodes[0]
			}
		}
	}

	// 4. Pre-cover columns for cells that already contain values on the starting board
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if board[r][c] != 0 {
				v := board[r][c]
				node := rowNodes[r][c][v]
				// Cover the column containing the current node, and all other columns satisfied by this row
				cover(node.Col)
				for j := node.Right; j != node; j = j.Right {
					cover(j.Col)
				}
			}
		}
	}

	// 5. Perform the DLX backtracking search
	var solution []*Node
	var search func() bool
	search = func() bool {
		// If root.Right points back to root, all constraints (columns) are satisfied
		if root.Right == &root {
			return true
		}

		col := selectColumn(&root)
		cover(col)

		for r := col.Head.Down; r != &col.Head; r = r.Down {
			solution = append(solution, r)
			for j := r.Right; j != r; j = j.Right {
				cover(j.Col)
			}

			if search() {
				return true
			}

			// Backtrack
			solution = solution[:len(solution)-1]
			for j := r.Left; j != r; j = j.Left {
				uncover(j.Col)
			}
		}
		uncover(col)
		return false
	}

	if search() {
		// Copy the solution values back to the board
		for _, node := range solution {
			board[node.RowVal][node.ColIdx] = node.Digit
		}
		return true
	}

	return false
}
