package kenken

import (
	"container/heap"
	"fmt"
	"testing"
)

func TestUnsolveablePuzzle(t *testing.T) {
	p := *NewPuzzle(2)
	idxs := NewIndexSet()
	idxs.Add(Index{0, 0})
	p.regions = append(p.regions, Region{1, Nothing, *idxs})
	idxs = NewIndexSet()
	idxs.Add(Index{0, 1})
	p.regions = append(p.regions, Region{1, Nothing, *idxs})
	idxs = NewIndexSet()
	idxs.Add(Index{1, 0})
	idxs.Add(Index{1, 1})
	p.regions = append(p.regions, Region{2, Mul, *idxs})

	p.prepareRegionsByIndex()
	p.prepareBoxesFromRegions()
	p.buildHeap()

	err := p.Solve()
	if err == nil {
		t.Fatalf("Puzzle returned success for unsolveable puzzle")
	}
}

func TestPuzzleSolve(t *testing.T) {
	pzls := make([]Puzzle, 0, 2)
	pzls = append(pzls, examplePuzzle())
	pzls = append(pzls, examplePuzzle2())
	sols := make([][][]byte, 0, 2)
	sols = append(sols, exampleSolution())
	sols = append(sols, exampleSolution2())
	for i := range pzls {
		p := pzls[i]
		p.prepareRegionsByIndex()
		p.prepareBoxesFromRegions()
		p.buildHeap()
		s := sols[i]

		err := p.Solve()
		fmt.Println("Solved puzzle")
		if err != nil {
			t.Error("Solve failed with error: ", err.Error())
		}
		hasError := false
		for y := range p.puzzle {
			for x := range p.puzzle[y] {
				if p.puzzle[y][x].GetValue() != s[y][x] {
					hasError = true
					t.Errorf("Mistake at box:\n%v\nexpected: %v", p.puzzle[y][x], s[y][x])
				}
			}
		}
		if hasError {
			t.Errorf("Solution was wrong:\n%v\n%v", p.String(), s)
		}
	}
}

func TestPuzzleBuildHeap(t *testing.T) {
	p := examplePuzzle()
	p.prepareBoxesFromRegions()

	p.buildHeap()

	previous := uint8(0)
	if p.heap.Len() < int(p.Size())*int(p.Size()) {
		t.Fatalf("Heap was size %v instead of %v", p.heap.Len(), p.Size()*p.Size())
	}
	for p.heap.Len() > 0 {
		top := heap.Pop(&p.heap).(*Box)
		if top.NumPossible() < previous {
			t.Fatalf("Heap returned values out of order")
		}
		previous = top.NumPossible()
	}
}

func TestPuzzlePrepBoxesFromRegions(t *testing.T) {
	pzls := make([]Puzzle, 0, 2)
	pzls = append(pzls, examplePuzzle())
	pzls = append(pzls, examplePuzzle2())
	sols := make([][][]byte, 0, 2)
	sols = append(sols, exampleSolution())
	sols = append(sols, exampleSolution2())
	for i := range pzls {
		p := pzls[i]
		s := sols[i]
		p.prepareBoxesFromRegions()
		for y, boxes := range p.puzzle {
			for x, box := range boxes {
				if !box.HasPossible(s[y][x]) {
					t.Errorf("Box did not include solution as possibility at (%v,%v)", x, y)
				}
			}
		}

		if i == 0 {
			// Test select possible values:
			box := p.puzzle[3][2]
			checkBoxPossibles(t, &box, []byte{3})
			box = p.puzzle[1][4]
			checkBoxPossibles(t, &box, []byte{1, 2, 4, 5})
		}
	}
}

func checkBoxPossibles(t *testing.T, b *Box, possibles []byte) {
	if b.NumPossible() > uint8(len(possibles)) {
		t.Errorf("Box %v had wrong num possibles %v, expected %v", b, b.NumPossible(), len(possibles))
	}
	for _, v := range possibles {
		if !b.HasPossible(v) {
			t.Errorf("Box %v was missing possible: %v", b, v)
		}
	}
}

func TestPuzzleMoveCursorDisallowed(t *testing.T) {
	p := Puzzle{2, nil, nil, nil, nil}
	tl := Index{0, 1}
	tr := Index{1, 1}
	bl := Index{0, 0}
	r := p.moveCursor(bl, uint(Left))
	if r != bl {
		t.Error("moveCursor allowed move left at boundary")
	}
	r = p.moveCursor(bl, uint(Down))
	if r != bl {
		t.Error("moveCursor allowed move down at boundary")
	}
	r = p.moveCursor(tl, uint(Left))
	if r != tl {
		t.Error("moveCursor allowed move left at boundary")
	}
	r = p.moveCursor(tl, uint(Up))
	if r != tl {
		t.Error("moveCursor allowed move up at boundary")
	}
	r = p.moveCursor(tr, uint(Right))
	if r != tr {
		t.Error("moveCursor allowed  move right at boundary")
	}
}

func examplePuzzle() Puzzle {
	p := *NewPuzzle(5)
	p.regions = make([]Region, 12)
	indices := *NewIndexSet()
	indices.Add(Index{0, 4})
	indices.Add(Index{1, 4})
	indices.Add(Index{2, 4})
	p.regions[0] = Region{12, Mul, indices}
	indices = *NewIndexSet()
	indices.Add(Index{3, 4})
	indices.Add(Index{3, 3})
	p.regions[1] = Region{2, Div, indices}
	indices = *NewIndexSet()
	indices.Add(Index{4, 4})
	indices.Add(Index{4, 3})
	indices.Add(Index{4, 2})
	p.regions[2] = Region{10, Sum, indices}
	indices = *NewIndexSet()
	indices.Add(Index{0, 3})
	indices.Add(Index{1, 3})
	p.regions[3] = Region{1, Sub, indices}
	indices = *NewIndexSet()
	indices.Add(Index{2, 3})
	p.regions[4] = Region{3, Nothing, indices}
	indices = *NewIndexSet()
	indices.Add(Index{0, 2})
	p.regions[5] = Region{2, Nothing, indices}
	indices = *NewIndexSet()
	indices.Add(Index{0, 1})
	indices.Add(Index{1, 1})
	indices.Add(Index{1, 2})
	p.regions[6] = Region{20, Mul, indices}
	indices = *NewIndexSet()
	indices.Add(Index{2, 2})
	indices.Add(Index{2, 1})
	p.regions[7] = Region{3, Sum, indices}
	indices = *NewIndexSet()
	indices.Add(Index{3, 2})
	indices.Add(Index{3, 1})
	p.regions[8] = Region{2, Sub, indices}
	indices = *NewIndexSet()
	indices.Add(Index{4, 1})
	indices.Add(Index{4, 0})
	p.regions[9] = Region{3, Sub, indices}
	indices = *NewIndexSet()
	indices.Add(Index{0, 0})
	indices.Add(Index{1, 0})
	p.regions[10] = Region{1, Sub, indices}
	indices = *NewIndexSet()
	indices.Add(Index{2, 0})
	indices.Add(Index{3, 0})
	p.regions[11] = Region{9, Sum, indices}
	return p
}

func exampleSolution() [][]uint8 {
	// Has to be written upside down
	return [][]uint8{
		[]uint8{3, 2, 5, 4, 1},
		[]uint8{5, 1, 2, 3, 4},
		[]uint8{2, 4, 1, 5, 3},
		[]uint8{4, 5, 3, 1, 2},
		[]uint8{1, 3, 4, 2, 5},
	}
}

func examplePuzzle2() Puzzle {
	p := *NewPuzzle(5)
	p.regions = make([]Region, 13)
	indices := *NewIndexSet()
	indices.Add(Index{0, 0})
	indices.Add(Index{1, 0})
	p.regions[0] = Region{1, Sub, indices}
	indices = *NewIndexSet()
	indices.Add(Index{2, 0})
	indices.Add(Index{2, 1})
	indices.Add(Index{2, 2})
	p.regions[1] = Region{9, Sum, indices}
	indices = *NewIndexSet()
	indices.Add(Index{3, 0})
	indices.Add(Index{4, 0})
	p.regions[2] = Region{1, Sub, indices}
	indices = *NewIndexSet()
	indices.Add(Index{0, 1})
	indices.Add(Index{0, 2})
	p.regions[3] = Region{3, Sum, indices}
	indices = *NewIndexSet()
	indices.Add(Index{1, 1})
	indices.Add(Index{1, 2})
	p.regions[4] = Region{12, Mul, indices}
	indices = *NewIndexSet()
	indices.Add(Index{3, 1})
	indices.Add(Index{4, 1})
	p.regions[5] = Region{2, Div, indices}
	indices = *NewIndexSet()
	indices.Add(Index{3, 2})
	p.regions[6] = Region{5, Nothing, indices}
	indices = *NewIndexSet()
	indices.Add(Index{4, 2})
	indices.Add(Index{4, 3})
	p.regions[7] = Region{6, Sum, indices}
	indices = *NewIndexSet()
	indices.Add(Index{0, 3})
	indices.Add(Index{0, 4})
	p.regions[8] = Region{2, Sub, indices}
	indices = *NewIndexSet()
	indices.Add(Index{1, 3})
	indices.Add(Index{1, 4})
	p.regions[9] = Region{2, Div, indices}
	indices = *NewIndexSet()
	indices.Add(Index{2, 3})
	indices.Add(Index{2, 4})
	p.regions[10] = Region{2, Sub, indices}
	indices = *NewIndexSet()
	indices.Add(Index{3, 3})
	indices.Add(Index{3, 4})
	p.regions[11] = Region{2, Sub, indices}
	indices = *NewIndexSet()
	indices.Add(Index{4, 4})
	p.regions[12] = Region{4, Nothing, indices}
	return p
}

func exampleSolution2() [][]uint8 {
	return [][]uint8{
		[]uint8{4, 5, 1, 2, 3},
		[]uint8{1, 3, 5, 4, 2},
		[]uint8{2, 4, 3, 5, 1},
		[]uint8{3, 2, 4, 1, 5},
		[]uint8{5, 1, 2, 3, 4},
	}
}
