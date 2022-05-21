package main

import (
	"flag"
	"fmt"
	"sudoku/sudoku"
)

// Establish Global Variables
var board [9][9]int
var validInput bool

// An algorithm which solves a given sudoku puzzle using backtracking
func recursiveSolve(rowPosition, columnPosition int) bool {

	size := len(board)

	// End condition which should be recursively reached if solution found.
	// i.e. Finishes 9th row, moves to 10th row (non-existent)
	if rowPosition == 9 {
		return true
	}
	// Move to next cell if current cell already filled in
	if board[rowPosition][columnPosition] != 0 {
		return recursiveSolve(sudoku.NextCell(rowPosition, columnPosition))
	} else {
		for i := 1; i <= size; i++ {
			if sudoku.CheckValid(board, rowPosition, columnPosition, i) == true {
				board[rowPosition][columnPosition] = i
				if recursiveSolve(sudoku.NextCell(rowPosition, columnPosition)) {
					return true
				}
				board[rowPosition][columnPosition] = 0
			}
		}
		return false
	}
}

// See below for inspiration
// INSPIRATION: https://charltonaustin.com/posts/sudoku-using-go-lang/
// INSPIRATION: https://www.geeksforgeeks.org/sudoku-backtracking-7/
// INSPIRATION: https://www.5minsofcode.com/sodoku_solver.html
func main() {
	algoFlag := flag.String("algo", "backtracking", "Solver algorithm to use: 'backtracking' or 'exact-cover' (or 'algo-x')")
	flag.Parse()

	if *algoFlag != "backtracking" && *algoFlag != "exact-cover" && *algoFlag != "algo-x" {
		fmt.Printf("Error: Unknown algorithm '%s'. Supported values are 'backtracking' and 'exact-cover'\n", *algoFlag)
		return
	}

	inputBoard := flag.Args()
	var err bool
	board, err = sudoku.CreateBoard(inputBoard)
	validInput = err

	canProceed := true

	// Check starting board validity according to minimum number requirements
	if sudoku.StartValid(board) == false {
		canProceed = false
		fmt.Printf("Error: Input configuration is not valid.")
	} else if validInput == false {
		fmt.Printf("Error: Incorrect input - string cannot be read according to standard 9 x 9 dimensions")
	}
	// Recursively iterate through board and print results
	if canProceed {
		fmt.Println()
		fmt.Println("Initial sudoku board shown below:\n")
		sudoku.PrintBoard(board)

		var solved bool
		if *algoFlag == "exact-cover" || *algoFlag == "algo-x" {
			solved = sudoku.SolveExactCover(&board)
		} else {
			solved = recursiveSolve(0, 0)
		}

		if solved {
			fmt.Println()
			fmt.Println("The following solution was found:\n")
			sudoku.PrintBoard(board)
			fmt.Println()
		} else {
			fmt.Println("\nA solution for this start configuration does not exist.")
			fmt.Println()
		}
	}
}
