# Go Sudoku Solver

A Sudoku solver implemented in Go using **Knuth's Algorithm X (Exact Cover)** with **Dancing Links (DLX)**. All code lives in a single `main.go`, read top-to-bottom with helpers at the bottom.

## Features

- **Algorithm X (Exact Cover)**: Formulates Sudoku as an exact cover problem and solves it with Dancing Links.
- **Grading compliant**: Prints only the final solution or `Error` for invalid boards.
- **Dependency-free**: Uses only `os` and `fmt`.
- **Robust validation**: Checks dimensions, characters, and minimum clues (at least 17 numbers, 8 unique).

---

## How It Works

Knuth's Algorithm X solves the **Exact Cover Problem**. Sudoku is treated as constraint satisfaction: pick a set of placements that satisfy every rule exactly once.

### Jigsaw analogy

1. There are **324 slots** that must be filled (constraints).
2. There are **729 pieces** (every digit in every cell).
3. Each piece covers exactly **4 slots** (cell, row, column, box).
4. Pick **81 pieces** so all slots are filled with no overlap.

### Constraints (columns)

Four constraint types, 81 each (**324 columns** total):

1. **Cell**: each cell has exactly one digit.
2. **Row**: each row contains 1–9 exactly once.
3. **Column**: each column contains 1–9 exactly once.
4. **Box**: each 3×3 box contains 1–9 exactly once.

### Candidates (rows)

There are **729** choices (9×9×9): place digit `v` in cell `(r, c)`. Each choice satisfies four constraints.

### Dancing Links (DLX)

The sparse matrix is a circular doubly-linked list. Each node has `Left`, `Right`, `Up`, `Down`. Covering a column removes it and intersecting rows; uncovering restores them:

```go
// Cover
node.Right.Left = node.Left
node.Left.Right = node.Right

// Uncover
node.Right.Left = node
node.Left.Right = node
```

The recursive search in `main.go` picks the column with fewest candidates (MRV heuristic), tries each row, covers related columns, and backtracks on failure:

```go
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
```

---

## Usage

Run with exactly 9 arguments (one row each). Use `.` or `0` for empty cells.

### Valid example

```bash
go run . ".96.4...1" "1...6...4" "5.481.39." "..795..43" ".3..8...." "4.5.23.18" ".1.63..59" ".59.7.83." "..359...7"
```

**Output:**

```
3 9 6 2 4 5 7 8 1
1 7 8 3 6 9 5 2 4
5 2 4 8 1 7 3 9 6
2 8 7 9 5 1 6 4 3
9 3 1 4 8 6 2 7 5
4 6 5 7 2 3 9 1 8
7 1 2 6 3 8 4 5 9
6 5 9 1 7 4 8 3 2
8 4 3 5 9 2 1 6 7

```

### Invalid example

```bash
go run . "invalid" "args"
```

**Output:**

```
Error
```

---

## Project layout

```
├── main.go        # Entry point, solver, and helpers
├── main_test.go   # Integration and unit tests
└── go.mod
```

---

## Testing

```bash
go test -v .
```

`TestAllScenarios` runs 18 integration cases via the built binary. `TestSolveExactCover` checks the DLX solver directly.
