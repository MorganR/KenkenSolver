package kenken

import (
	"fmt"
	"strings"
)

type Index struct {
	X uint8
	Y uint8
}

type IndexSet map[Index]struct{}

func (i Index) Equals(x, y uint8) bool {
	return x == i.X && y == i.Y
}

func (i Index) String() string {
	return fmt.Sprintf("(%v,%v)", i.X, i.Y)
}

func NewIndexSet() *IndexSet {
	is := make(IndexSet)
	return &is
}

func (is *IndexSet) Add(i Index) {
	var empty struct{}
	(*is)[i] = empty
}

func (is *IndexSet) Drop(i Index) {
	delete(*is, i)
}

func (is IndexSet) Contains(i Index) bool {
	_, present := is[i]
	return present
}

func (is IndexSet) Len() int {
	return len(is)
}

func (is IndexSet) Slice() []Index {
	s := make([]Index, is.Len())
	i := 0
	for idx := range is {
		s[i] = idx
		i++
	}
	return s
}

func (is IndexSet) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	counter := 0
	for i := range is {
		counter += 1
		if counter < len(is) {
			sb.WriteString(fmt.Sprintf("%v,", i))
		} else {
			sb.WriteString(fmt.Sprintf("%v", i))
		}
	}
	sb.WriteString("]")
	return sb.String()
}
