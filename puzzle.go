package kenken

import (
	"container/heap"
	"fmt"
	"strings"

	tm "github.com/buger/goterm"
)

type Char uint

const (
	Esc               Char = 27
	LeftSquareBracket Char = 91
	Up                Char = 65
	Down              Char = 66
	Right             Char = 67
	Left              Char = 68
	Enter             Char = 10
)

type Puzzle struct {
	size           uint8
	puzzle         [][]Box
	regions        []Region
	regionsByIndex map[Index]*Region
	heap           BoxHeap
}

func NewPuzzle(size uint8) *Puzzle {
	p := make([][]Box, size)
	for i := range p {
		p[i] = make([]Box, size)
	}
	return &Puzzle{size, p, nil, make(map[Index]*Region), nil}
}

func RequestPuzzle(size uint8) *Puzzle {
	p := make([][]Box, size)
	selected := make([][]bool, size)
	for i := range p {
		p[i] = make([]Box, size)
		selected[i] = make([]bool, size)
	}
	pzl := &Puzzle{size, p, nil, make(map[Index]*Region), nil}
	cursor := Index{0, size - 1}

	region := *NewIndexSet()
	var input uint
	for numUnset := size * size; numUnset > 0; fmt.Scanf("%c", &input) {
		switch input {
		case uint(Esc):
			fmt.Scanf("%c", &input) // LeftSquareBracket
			fmt.Scanf("%c", &input) // Direction
			cursor = pzl.moveCursor(cursor, input)
		case 'c':
			for ; Char(input) != Enter; fmt.Scanf("%c", &input) {
				// Clear the enter char
			}
			pzl.confirmRegion(region)
			numUnset -= uint8(len(region))
			region = *NewIndexSet()
		case 'u':
			for ; Char(input) != Enter; fmt.Scanf("%c", &input) {
				// Clear the enter char
			}
			if region.Len() > 1 {
				selected[cursor.Y][cursor.X] = false
				region.Drop(cursor)
				cursor = region.Slice()[0]
			}
		default:
			// Do nothing
		}
		if !selected[cursor.Y][cursor.X] {
			selected[cursor.Y][cursor.X] = true
			region.Add(cursor)
		}
		tm.Clear()
		tm.MoveCursor(1, 1)
		tm.Printf("Input the formula regions in the puzzle. Start by selecting a region, then enter the result and necessary operation.\n\n")
		tm.Println("Use the arrow keys to move the 'X' and select a region in the grid.")
		tm.Println("Press 'c' when complete, or press 'u' to remove the current square from the selection.")
		tm.Println("All inputs must be followed by enter.")
		pzl.printWithCursor(cursor, selected, region)
		tm.Flush()
	}
	pzl.prepareRegionsByIndex()
	pzl.prepareBoxesFromRegions()
	pzl.buildHeap()
	return pzl
}

func (p *Puzzle) moveCursor(cursor Index, dir uint) Index {
	switch Char(dir) {
	case Up:
		if cursor.Y < p.size-1 {
			cursor.Y += 1
		}
	case Down:
		if cursor.Y > 0 {
			cursor.Y -= 1
		}
	case Left:
		if cursor.X > 0 {
			cursor.X -= 1
		}
	case Right:
		if cursor.X < p.size-1 {
			cursor.X += 1
		}
	}
	return cursor
}

func (p *Puzzle) confirmRegion(region IndexSet) {
	tm.Println("What is the result of this region?")
	tm.Flush()
	var result uint
	fmt.Scan(&result)
	tm.Println("What is the op for this region? (eg: *, /, +, -, or nothing)")
	tm.Flush()
	var op Operation
	var opChar uint
	fmt.Scanf("%c", &opChar)
	switch opChar {
	case '*':
		op = Mul
	case '/':
		op = Div
	case '+':
		op = Sum
	case '-':
		op = Sub
	case uint(Enter):
		op = Nothing
	}
	p.regions = append(p.regions, Region{result, op, region})
}

// Fill the p.regionsByIndex container. Must be done once no more modifications will be made to p.regions.
func (p *Puzzle) prepareRegionsByIndex() {
	for i := range p.regions {
		for _, idx := range p.regions[i].GetIndices() {
			p.regionsByIndex[idx] = &p.regions[i]
		}
	}
}

func (p *Puzzle) prepareBoxesFromRegions() {
	for _, r := range p.regions {
		valueMaps := r.GetPossibleMaps(p.Size())
		for _, idx := range r.GetIndices() {
			box := p.getBox(idx)
			*box = *NewBox(idx, p.Size())
			for _, m := range valueMaps {
				for _, possible := range m.GetSortedList() {
					if !box.HasPossible(possible) {
						box.AddPossible(possible)
					}
				}
			}
		}
	}
}

func (p *Puzzle) buildHeap() {
	(*p).heap = make(BoxHeap, 0, p.Size()*p.Size())
	for y := range p.puzzle {
		for x := range p.puzzle[y] {
			(*p).heap.Push(&p.puzzle[y][x])
		}
	}
	heap.Init(&(*p).heap)
}

func (p *Puzzle) getBox(i Index) *Box {
	return &(*p).puzzle[i.Y][i.X]
}

func (p *Puzzle) Size() uint8 {
	return p.size
}

type UnsolveableError struct {
	failedPaths uint
}

func (e UnsolveableError) Error() string {
	return fmt.Sprintf("Failed solving puzzle after trying %v paths", e.failedPaths)
}

func (p *Puzzle) Solve() error {
	return p.trySolve()
}

func (p *Puzzle) trySolve() error {
	if p.heap.Len() == 0 {
		return nil
	}
	numFailedPaths := uint(0)
	topBox := heap.Pop(&p.heap).(*Box)
	possibles := topBox.GetPossibles()
	for _, v := range possibles {
		if !p.isRegionValidIfSet(*topBox, v) {
			numFailedPaths++
			continue
		}
		topBox.SetValue(v)
		modifications := make([]Index, 0)
		p.deletePossibilityFromRow(v, topBox.idx.Y, &modifications)
		p.deletePossibilityFromCol(v, topBox.idx.X, &modifications)
		err := p.trySolve()
		if err == nil {
			return nil
		}
		switch err.(type) {
		case UnsolveableError:
			numFailedPaths += err.(UnsolveableError).failedPaths
		default:
			numFailedPaths++
		}
		p.resetPossibilities(v, modifications)
		topBox.UnsetValue()
	}
	heap.Push(&p.heap, topBox)
	return UnsolveableError{numFailedPaths}
}

func (p *Puzzle) isRegionValidIfSet(b Box, v byte) bool {
	setValues := *NewByteMap()
	setValues.Add(v)
	r := *p.regionsByIndex[b.idx]
	for _, idx := range r.GetIndices() {
		box := p.puzzle[idx.Y][idx.X]
		if box.GetValue() != 0 {
			setValues.Add(box.GetValue())
		}
	}
	possibleMaps := r.GetPossibleMaps(p.Size())
	for _, posMap := range possibleMaps {
		isPos := true
		for setPos, setNum := range setValues.Map() {
			posNum, present := posMap.Map()[setPos]
			if !present || posNum < setNum {
				isPos = false
			}
		}
		// If at least one map is still possible, return no error
		if isPos {
			return true
		}
	}
	return false
}

func (p *Puzzle) deletePossibilityFromRow(v, y byte, m *[]Index) {
	for x := byte(0); x < p.Size(); x++ {
		if !p.puzzle[y][x].HasPossible(v) {
			continue
		}
		idx := Index{x, y}
		*m = append(*m, idx)
		p.puzzle[y][x].DeletePossible(v)
		if p.puzzle[y][x].heapIndex >= 0 {
			heap.Fix(&p.heap, p.puzzle[y][x].heapIndex)
		}
	}
}

func (p *Puzzle) deletePossibilityFromCol(v, x byte, m *[]Index) {
	for y := byte(0); y < p.Size(); y++ {
		if !p.puzzle[y][x].HasPossible(v) {
			continue
		}
		idx := Index{x, y}
		*m = append(*m, idx)
		p.puzzle[y][x].DeletePossible(v)
		if p.puzzle[y][x].heapIndex >= 0 {
			heap.Fix(&p.heap, p.puzzle[y][x].heapIndex)
		}
	}
}

func (p *Puzzle) resetPossibilities(v byte, modifications []Index) {
	for _, idx := range modifications {
		p.puzzle[idx.Y][idx.X].AddPossible(v)
		if p.puzzle[idx.Y][idx.X].heapIndex >= 0 {
			heap.Fix(&p.heap, p.puzzle[idx.Y][idx.X].heapIndex)
		}
	}
}

func (p *Puzzle) Print() {
	getValue := func(i Index) string {
		v := p.puzzle[i.Y][i.X]
		return v.ValueString()
	}
	p.printValueAs(getValue)
	tm.Flush()
}

func (p *Puzzle) printWithCursor(cursor Index, selected [][]bool, region IndexSet) {
	getValue := func(i Index) string {
		if cursor == i {
			return "X"
		} else if region.Contains(i) {
			return "\u25a1"
		} else if selected[i.Y][i.X] {
			return "\u25a6"
		}
		return " "
	}
	p.printValueAs(getValue)
}

func (p Puzzle) String() string {
	getValue := func(i Index) string {
		v := p.puzzle[i.Y][i.X]
		return v.ValueString()
	}
	return p.stringValueAs(getValue)
}

func (p *Puzzle) stringValueAs(getValue func(Index) string) string {
	var sb strings.Builder
	sb.WriteString("\t ")
	for x := uint8(0); x < p.size; x++ {
		sb.WriteString(fmt.Sprintf(" %v", x))
	}
	sb.WriteString("\n")
	sb.WriteString("\t \u250f")
	for x := uint8(0); x < p.size; x++ {
		sb.WriteString("\u2501")
		if x < p.size-1 {
			sb.WriteString("\u252f")
		}
	}
	sb.WriteString("\u2513\n")
	for y := int16(p.size - 1); y >= 0; y-- {
		sb.WriteString(fmt.Sprintf("\t%v\u2503", y))
		for x := uint8(0); x < p.size; x++ {
			sb.WriteString(fmt.Sprint(getValue(Index{x, uint8(y)})))
			if x < p.size-1 {
				sb.WriteString("\u2502")
			}
		}
		sb.WriteString("\u2503\n")
		if y > 0 {
			sb.WriteString("\t \u2520")
			for x := uint8(0); x < p.size; x++ {
				sb.WriteString("\u2500")
				if x < p.size-1 {
					sb.WriteString("\u253c")
				}
			}
			sb.WriteString("\u2528\n")
		}
	}
	sb.WriteString("\t \u2517")
	for x := uint8(0); x < p.size; x++ {
		sb.WriteString("\u2501")
		if x < p.size-1 {
			sb.WriteString("\u2537")
		}
	}
	sb.WriteString("\u251b\n")
	sb.WriteString("Regions:\n")
	for _, region := range p.regions {
		sb.WriteString(fmt.Sprintf("%v\n", region))
	}
	return sb.String()
}

func (p *Puzzle) printValueAs(getValue func(Index) string) {
	tm.Printf("\t ")
	for x := uint8(0); x < p.size; x++ {
		tm.Printf(" %v", x)
	}
	tm.Printf("\n")
	tm.Printf("\t \u250f")
	for x := uint8(0); x < p.size; x++ {
		tm.Printf("\u2501")
		if x < p.size-1 {
			tm.Printf("\u252f")
		}
	}
	tm.Printf("\u2513\n")
	for y := int16(p.size - 1); y >= 0; y-- {
		tm.Printf("\t%v\u2503", y)
		for x := uint8(0); x < p.size; x++ {
			tm.Printf(getValue(Index{x, uint8(y)}))
			if x < p.size-1 {
				tm.Printf("\u2502")
			}
		}
		tm.Printf("\u2503\n")
		if y > 0 {
			tm.Printf("\t \u2520")
			for x := uint8(0); x < p.size; x++ {
				tm.Printf("\u2500")
				if x < p.size-1 {
					tm.Printf("\u253c")
				}
			}
			tm.Printf("\u2528\n")
		}
	}
	tm.Printf("\t \u2517")
	for x := uint8(0); x < p.size; x++ {
		tm.Printf("\u2501")
		if x < p.size-1 {
			tm.Printf("\u2537")
		}
	}
	tm.Printf("\u251b\n")
	tm.Println("Regions:")
	for _, region := range p.regions {
		tm.Printf("%v\n", region)
	}
}
