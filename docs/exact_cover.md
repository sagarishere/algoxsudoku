# Knuth's Algorithm X (Exact Cover via Dancing Links)

Knuth's Algorithm X solves the **Exact Cover Problem**. It treats Sudoku not as a grid-filling game, but as a constraint-satisfaction puzzle where we must select a set of choices that satisfy all rules exactly once.

---

## The Analogy: A Jigsaw Puzzle

Imagine you are putting together a complex jigsaw puzzle.
1. There are **324 slots** on the board that must be filled.
2. You have a bucket of **729 puzzle pieces** (representing every possible digit placement in every cell).
3. Each puzzle piece is designed to cover exactly **4 slots** on the board.
4. Your goal is to select exactly **81 puzzle pieces** from the bucket such that all **324 slots** are filled, with **no overlaps** (no slot is covered by more than one piece).

This is an **Exact Cover**. If you try a piece and it blocks other required slots, you remove it and try another piece.

---

## How It Works

### 1. The Constraints (Columns)
A valid Sudoku board must satisfy $4$ sets of constraints, each having $81$ unique slots (totaling $324$ columns):
1. **Cell Constraint**: Every cell must contain exactly one number.
2. **Row Constraint**: Every row must contain digits $1$ through $9$ exactly once.
3. **Column Constraint**: Every column must contain digits $1$ through $9$ exactly once.
4. **Box Constraint**: Every $3 \times 3$ box must contain digits $1$ through $9$ exactly once.

### 2. The Candidates (Rows)
There are $9 \times 9 \times 9 = 729$ possible choices. Each choice represents placing a specific number in a specific cell. Each choice satisfies exactly $4$ constraints.

### 3. Dancing Links (DLX)
Donald Knuth's **Dancing Links** technique represents this sparse matrix of constraints as a toroidal, circularly doubly-linked list. Each node has pointers to its neighbors (`Left`, `Right`, `Up`, `Down`). 

Covering a column removes it and all intersecting rows from the matrix. Uncovering them restores them. The pointer arithmetic is simple and elegant:
```go
// Covering a node
node.Right.Left = node.Left
node.Left.Right = node.Right

// Restoring (Uncovering) a node
node.Right.Left = node
node.Left.Right = node
```
This is why it's called "Dancing Links"—nodes seem to dance out of the list and back in.

---

## Core Code Snippet

Here is the recursive constraint solver implemented in [algoX.go](file:///Users/sagar/Downloads/sudoku/sudoku/algoX.go):

```go
func solveDLX(root *Node, solution *[]*Node, board *[9][9]int) bool {
	// If the root points to itself, all constraints are satisfied!
	if root.Right == root {
		return true
	}

	// Choose the column with the fewest active rows (Minimum Size Heuristic)
	col := selectColumn(root)
	cover(col)

	// Try each row candidate in the column
	for rowNode := col.Down; rowNode != col; rowNode = rowNode.Down {
		*solution = append(*solution, rowNode)

		// Cover all columns satisfied by this candidate
		for sibling := rowNode.Right; sibling != rowNode; sibling = sibling.Right {
			cover(sibling.Col)
		}

		// Recurse to solve remaining constraints
		if solveDLX(root, solution, board) {
			// Decode solution back into 9x9 board values
			for _, node := range *solution {
				r := node.RowIdx / 81
				c := (node.RowIdx % 81) / 9
				val := (node.RowIdx % 9) + 1
				board[r][c] = val
			}
			return true
		}

		// Backtrack: uncover siblings
		for sibling := rowNode.Left; sibling != rowNode; sibling = sibling.Left {
			uncover(sibling.Col)
		}
		*solution = (*solution)[:len(*solution)-1]
	}

	uncover(col)
	return false
}
```

---

## Key Characteristics

*   **Heuristics**: Choosing the column with the fewest candidates (MRV heuristic) minimizes the branching factor immediately.
*   **Speed**: Extremely fast for hard or adversarial puzzles, as it cuts down search space size rapidly.
*   **Memory Overhead**: High memory overhead relative to backtracking because of the dense network of pointers and column headers.
