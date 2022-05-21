# Traditional Grid Backtracking

The Traditional Grid Backtracking algorithm is the foundational method for solving Sudoku. It uses a **Depth-First Search (DFS)** recursion to explore potential board configurations step-by-step.

---

## The Analogy: Solving a Maze

Imagine you are trying to solve a physical maze. 
1. You walk forward until you reach an intersection.
2. At the intersection, you decide to turn left.
3. You continue walking. If you hit a dead end, you do not panic. You walk **backwards** to the last intersection you encountered and try the next direction (e.g., turning right).
4. You repeat this until you reach the exit of the maze.

In Sudoku, each empty cell is an "intersection," and the digits $1$ through $9$ are the "paths" you can take. If a digit violates Sudoku rules later down the line (a dead end), you erase it (backtrack) and try the next number.

---

## How It Works

1. **Find an Empty Cell**: Search the board cell-by-cell (from top-left to bottom-right) to find an empty cell (represented by `0`).
2. **Try Placement**: Try placing numbers from $1$ to $9$ in that empty cell.
3. **Validate Placement**: Check if the number is valid in the current row, column, and $3 \times 3$ box.
4. **Recurse**: If it's valid, place the number and recursively try to solve the rest of the board.
5. **Backtrack**: If no numbers from $1$ to $9$ lead to a solution, reset the cell to `0` and return `false` to let the previous cell try its next number.

---

## Core Code Snippet

Here is the central backtracking solver implemented in [backtracking.go](file:///Users/sagar/Downloads/sudoku/sudoku/backtracking.go):

```go
func SolveBacktracking(board *[9][9]int) bool {
	var solve func(row, col int) bool
	solve = func(row, col int) bool {
		// If we reach row 9, we have successfully filled the grid
		if row == 9 {
			return true
		}
		
		// Determine the next cell to solve (moves left-to-right, row-by-row)
		nextR, nextC := NextCell(row, col)
		
		// If the cell is already filled, skip to the next cell
		if board[row][col] != 0 {
			return solve(nextR, nextC)
		}
		
		// Try placing numbers from 1 to 9
		for i := 1; i <= 9; i++ {
			if CheckValid(*board, row, col, i) {
				board[row][col] = i // Place value
				
				// Recurse to see if this placement leads to a solution
				if solve(nextR, nextC) {
					return true
				}
				
				// Undo placement (backtrack)
				board[row][col] = 0
			}
		}
		return false // Trigger backtracking in the parent caller
	}
	return solve(0, 0)
}
```

---

## Key Characteristics

*   **Memory Efficiency**: Requires minimal memory since the recursion stack uses negligible space and the board is solved in-place.
*   **Time Complexity**: In the worst case, it can be slow ($\mathcal{O}(9^{81})$ in theory, though constraint checking makes it much faster in practice).
*   **Suitability**: Excellent baseline, but runs into performance bottlenecks on extremely difficult or adversarial Sudoku puzzles.
