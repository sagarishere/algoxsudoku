package main

import (
	"fmt"
	"os"
)

// Establish Global Variables
var board [9][9]int
var validInput bool

func main() {
	inputBoard := os.Args[1:]
	var err bool
	board, err = createBoard(inputBoard)
	validInput = err

	canProceed := true

	if startValid(board) == false {
		canProceed = false
	} else if validInput == false {
		canProceed = false
	}

	if !canProceed {
		fmt.Println("Error")
		return
	}

	if solveExactCover(&board) {
		printBoard(board)
		fmt.Println()
	} else {
		fmt.Println("Error")
	}
}

// Node represents a node in the Dancing Links (DLX) matrix.
type Node struct {
	Left, Right, Up, Down *Node
	Col                   *Column
	RowVal, ColIdx, Digit int
}

// Column represents a column header in the DLX matrix.
type Column struct {
	Head Node
	Size int
	ID   int
}

// solveExactCover solves the Sudoku puzzle using Knuth's Algorithm X / Dancing Links (DLX).
func solveExactCover(board *[9][9]int) bool {
	var root Node

	columns := make([]*Column, 324)
	for i := 0; i < 324; i++ {
		columns[i] = &Column{
			ID: i,
		}
		columns[i].Head.Col = columns[i]
		columns[i].Head.Up = &columns[i].Head
		columns[i].Head.Down = &columns[i].Head
	}

	curr := &root
	for i := 0; i < 324; i++ {
		curr.Right = &columns[i].Head
		columns[i].Head.Left = curr
		curr = &columns[i].Head
	}
	curr.Right = &root
	root.Left = curr

	rowNodes := make([][][]*Node, 9)
	for r := 0; r < 9; r++ {
		rowNodes[r] = make([][]*Node, 9)
		for c := 0; c < 9; c++ {
			rowNodes[r][c] = make([]*Node, 10)
		}
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			for v := 1; v <= 9; v++ {
				c1 := r*9 + c
				c2 := 81 + r*9 + (v - 1)
				c3 := 162 + c*9 + (v - 1)
				c4 := 243 + ((r/3)*3+c/3)*9 + (v - 1)

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

					last := col.Head.Up
					node.Down = &col.Head
					col.Head.Up = node
					node.Up = last
					last.Down = node

					col.Size++
				}

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

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if board[r][c] != 0 {
				v := board[r][c]
				node := rowNodes[r][c][v]
				cover(node.Col)
				for j := node.Right; j != node; j = j.Right {
					cover(j.Col)
				}
			}
		}
	}

	var solution []*Node
	var search func() bool
	search = func() bool {
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

			solution = solution[:len(solution)-1]
			for j := r.Left; j != r; j = j.Left {
				uncover(j.Col)
			}
		}
		uncover(col)
		return false
	}

	if search() {
		for _, node := range solution {
			board[node.RowVal][node.ColIdx] = node.Digit
		}
		return true
	}

	return false
}

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

// createBoard parses nine row strings into a 9x9 board.
func createBoard(startCondition []string) ([9][9]int, bool) {
	sudokuSize := len(startCondition)
	validCreate := true

	var startBoard = [9][9]int{}
	if sudokuSize != 9 {
		validCreate = false
	}

	if validCreate == true {
		for i := 0; i < sudokuSize; i++ {
			if len(startCondition[i]) != 9 {
				validCreate = false
				break
			}
			for j := 0; j < sudokuSize; j++ {
				startBoard[i][j] = 0
				if startCondition[i][j] >= '1' && startCondition[i][j] <= '9' {
					startBoard[i][j] = int(startCondition[i][j] - 48)
				}
			}
		}
	}
	return startBoard, validCreate
}

// startValid checks minimum clue count and starting grid consistency.
func startValid(inputBoard [9][9]int) bool {
	uniqueNumSlice := make([]int, 9)
	uniqueNumCount := 0
	numberCount := 0
	canContinue := true

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if inputBoard[i][j] >= 1 && inputBoard[i][j] <= 9 {
				numberCount++
				if uniqueNumSlice[inputBoard[i][j]-1] < 1 {
					uniqueNumSlice[inputBoard[i][j]-1]++
					uniqueNumCount++
				}
			}
		}
	}
	if uniqueNumCount < 8 || numberCount < 17 || numberCount > 77 {
		canContinue = false
	}

	if canContinue == true {
		for k := 0; k < 9; k++ {
			for l := 0; l < 9; l++ {
				if inputBoard[k][l] != 0 {
					for m := 0; m < 9; m++ {
						if inputBoard[k][m] == inputBoard[k][l] && m != l {
							return false
						}
					}
					for n := 0; n < 9; n++ {
						if inputBoard[n][l] == inputBoard[k][l] && n != k {
							return false
						}
					}
					boxStartRow := (k / 3) * 3
					boxStartCol := (l / 3) * 3
					for p := boxStartRow; p < boxStartRow+3; p++ {
						for q := boxStartCol; q < boxStartCol+3; q++ {
							if inputBoard[p][q] == inputBoard[k][l] && p != k && q != l {
								return false
							}
						}
					}
				}
			}
		}
	}
	return canContinue
}

// printBoard prints the solved board with spaces between values.
func printBoard(board [9][9]int) {
	sudokuSize := len(board)

	for i := 0; i < sudokuSize; i++ {
		for j := 0; j < sudokuSize; j++ {
			if board[i][j] < 1 && j != sudokuSize-1 {
				fmt.Print(board[i][j])
				fmt.Print(" ")
			} else if board[i][j] < 1 && j == sudokuSize-1 {
				fmt.Print(board[i][j])
				fmt.Print("\n")
			} else if board[i][j] >= 1 && j != sudokuSize-1 {
				fmt.Print(board[i][j])
				fmt.Print(" ")
			} else {
				fmt.Print(board[i][j])
				fmt.Print("\n")
			}
		}
	}
}

// ---------------------------------------------------------------------------
// DLX helpers
// ---------------------------------------------------------------------------

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
