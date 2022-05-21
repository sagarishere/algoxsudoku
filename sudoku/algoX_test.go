package sudoku

import "testing"

func TestSolveExactCover(t *testing.T) {
	// A standard valid Sudoku puzzle
	input := []string{
		"096040001",
		"100060004",
		"504810390",
		"007950043",
		"030080000",
		"405023018",
		"010630059",
		"059070830",
		"003590007",
	}

	board, valid := CreateBoard(input)
	if !valid {
		t.Fatalf("Expected board to be validly created")
	}

	solved := SolveExactCover(&board)
	if !solved {
		t.Fatalf("Expected exact cover solver to find a solution")
	}

	expectedSolution := [9][9]int{
		{3, 9, 6, 2, 4, 5, 7, 8, 1},
		{1, 7, 8, 3, 6, 9, 5, 2, 4},
		{5, 2, 4, 8, 1, 7, 3, 9, 6},
		{2, 8, 7, 9, 5, 1, 6, 4, 3},
		{9, 3, 1, 4, 8, 6, 2, 7, 5},
		{4, 6, 5, 7, 2, 3, 9, 1, 8},
		{7, 1, 2, 6, 3, 8, 4, 5, 9},
		{6, 5, 9, 1, 7, 4, 8, 3, 2},
		{8, 4, 3, 5, 9, 2, 1, 6, 7},
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if board[r][c] != expectedSolution[r][c] {
				t.Errorf("At row %d, col %d: expected %d, got %d", r, c, expectedSolution[r][c], board[r][c])
			}
		}
	}
}
