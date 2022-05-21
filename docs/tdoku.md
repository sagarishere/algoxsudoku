# SIMD-Optimized (Tdoku-Inspired) Solver

The SIMD-Optimized solver is a state-of-the-art engine inspired by **Tdoku** (a famous, lightning-fast C++ solver). It simulates parallel vector registers and employs advanced constraint propagation rules before fallback backtracking.

---

## The Analogy: Sherlock Holmes

Imagine you are Sherlock Holmes solving a mystery:
1.  **Deduction (Propagation)**: Instead of randomly guessing who the culprit is, you gather clues.
    *   *Clue 1 (Naked Single)*: "This room only has one door key that fits. It must be Key A." (If a cell has only one possible candidate remaining, it must be that candidate).
    *   *Clue 2 (Hidden Single)*: "Across all rooms in the hall, only Room 3 has a closet big enough to hide the diamond." (If a number can only fit in one cell within a row, column, or box, it must go there, even if that cell has other candidate options).
2.  **Making a Guess (MRV Search)**: When you have deduced everything you can and are still stuck, you must guess. However, you choose the question with the **fewest possible answers** (e.g., a choice between 2 suspects rather than 9) to minimize the chance of making a wrong turn.

This solver runs deductions in parallel using bitboards, and fallback searches on the cell with the Minimum Remaining Values (MRV).

---

## How It Works

### 1. Bitboard Candidates
We model the grid as an array of 81 `uint16` bitboards. Each bit represents whether a number is still possible for that cell:
```go
Candidates [81]uint16 // Bits 1..9 represent digits 1..9
```

### 2. Constraint Propagation (Naked & Hidden Singles)
*   **Naked Singles**: If a cell's candidate bitboard has only a single bit set (e.g., `0b000010000`), we instantly assign that digit and propagate the constraint to clear that bit from all other cells in the same row, column, and box.
*   **Hidden Singles**: For each row, column, and box, we count how many times each digit appears as a candidate. If a digit is only possible in a single cell, we place it there immediately.

### 3. Hardware-Accelerated Bit counting
To count candidates or find the index of a set bit, we use the Go compiler's built-in `math/bits` package:
*   `bits.OnesCount16(candidates)`: Instantly returns the number of possibilities for a cell.
*   `bits.TrailingZeros16(candidates)`: Instantly returns the index of the only set bit.
On modern CPUs (x86_64 and ARM64), these compile down to a single instruction like `POPCNT` or `TZCNT`.

---

## Core Code Snippet

Here is the propagation loop implemented in [tdoku.go](file:///Users/sagar/Downloads/sudoku/sudoku/tdoku.go):

```go
func (s *BoardState) Propagate() bool {
	changed := true
	for changed {
		changed = false
		for i := 0; i < 81; i++ {
			if s.Board[i] != 0 {
				continue
			}
			cands := s.Candidates[i]
			count := bits.OnesCount16(cands)
			if count == 0 {
				return false // Dead end: cell has 0 options
			}
			if count == 1 {
				// Naked Single found!
				val := bits.TrailingZeros16(cands)
				if !s.Assign(i, val) {
					return false
				}
				changed = true
			}
		}
		
		// Propagate Hidden Singles
		hiddenChanged, ok := s.PropagateHidden()
		if !ok {
			return false
		}
		if hiddenChanged {
			changed = true
		}
	}
	return true
}
```

---

## Key Characteristics

*   **Extremely Optimized**: Bypasses full DFS recursion for easy/medium boards, solving them instantly via propagation alone.
*   **Heuristic Fallback**: Uses the MRV (Minimum Remaining Values) heuristic to select the next cell with the fewest candidates for backtracking if propagation gets stuck.
*   **Hardware Acceleration**: Direct mapping to processor bitwise instructions provides incredible speedups on ARM64 and modern Intel/AMD CPUs.
