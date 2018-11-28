package kenken

import "fmt"

type PossibleSet map[uint8]struct{}

type Box struct {
	idx       Index
	possibles PossibleSet
	value     uint8
	heapIndex int
}

func (b Box) GetValue() uint8 { return b.value }

func (b *Box) SetValue(v uint8) { b.value = v }

func (b *Box) UnsetValue() { b.value = 0 }

func (b Box) IsValueSet() bool { return b.value != 0 }

func (b Box) NumPossible() uint8 { return uint8(len(b.possibles)) }

func (b *Box) DeletePossible(p uint8) {
	delete(b.possibles, p)
}

func (b *Box) AddPossible(p uint8) {
	b.possibles.Add(p)
}

func (b Box) HasPossible(p uint8) bool {
	return b.possibles.Contains(p)
}

func (b Box) GetPossibles() []byte {
	p := make([]byte, b.NumPossible())
	i := 0
	for k := range b.possibles {
		p[i] = k
		i++
	}
	return p
}

func NewBox(idx Index, size uint8) *Box {
	box := new(Box)
	box.idx = idx
	box.possibles = make(PossibleSet)
	box.value = 0
	box.heapIndex = -1
	return box
}

func (b Box) ValueString() string {
	if b.value != 0 {
		return fmt.Sprintf("%v", b.value)
	}
	return " "
}

func (ps *PossibleSet) Add(x uint8) {
	var empty struct{}
	(*ps)[x] = empty
}

func (ps *PossibleSet) Contains(x uint8) bool {
	_, present := (*ps)[x]
	return present
}
