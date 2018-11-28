package kenken

import (
	"fmt"
	"testing"
)

func TestIndexEquals(t *testing.T) {
	i := Index{1, 1}

	if !i.Equals(1, 1) {
		t.Errorf("Equals returned false incorrectly")
	}
	if i.Equals(0, 1) || i.Equals(1, 0) {
		t.Errorf("Equals returned true incorrectly")
	}
}

func ExampleIndexString() {
	i := Index{1, 2}
	fmt.Println(i)
	// Output: (1,2)
}

func TestIndexSetLen(t *testing.T) {
	is := NewIndexSet()
	if is.Len() != 0 {
		t.Errorf("IndexSet's length not initialized to zero.")
	}

	is.Add(Index{0, 1})
	is.Add(Index{0, 1})
	if is.Len() != 1 {
		t.Errorf("IndexSet's length was %v, expected %v", is.Len(), 1)
	}
}

func TestIndexSetContains(t *testing.T) {
	is := NewIndexSet()
	is.Add(Index{0, 1})
	is.Add(Index{0, 1})
	is.Add(Index{1, 2})
	if !is.Contains(Index{0, 1}) || !is.Contains(Index{1, 2}) {
		t.Errorf("IndexSet did not contain expected contents, instead: %v", is)
	}
	if is.Contains(Index{1, 1}) {
		t.Errorf("IndexSet contained unexpected contents: %v", Index{1, 1})
	}
}

func TestIndexSetDrop(t *testing.T) {
	is := NewIndexSet()
	i := Index{0, 1}
	is.Add(i)
	is.Add(Index{0, 2})
	is.Drop(i)
	if is.Contains(i) || is.Len() != 1 {
		t.Errorf("Did not properly drop index %v from set %v", i, is)
	}
}

func TestIndexSetSlice(t *testing.T) {
	is := NewIndexSet()
	i1, i2 := Index{0,1}, Index{0,2}
	is.Add(i1)
	is.Add(i1)
	is.Add(i2)
	s := is.Slice()
	if len(s) != 2 {
		t.Errorf("Slice had wrong length %v, expected %v", len(s), 2)
	}
	isSliceCorrect := false
	if s[0] == i1 {
		if s[1] == i2 {
			isSliceCorrect = true
		}
	} else if s[0] == i2 {
		if s[1] == i1 {
			isSliceCorrect = true
		}
	}
	if !isSliceCorrect {
		t.Errorf("Slice contained wrong indices: %v, expected %v and %v", s, i1, i2)
	}
}

func ExampleIndexSetDropMissingIndex() {
	is := NewIndexSet()
	is.Drop(Index{0, 1})
	fmt.Println(is)
	// Output: []
}

func ExampleIndexSetString() {
	is := NewIndexSet()
	is.Add(Index{0, 1})
	fmt.Println(is)
	// Output: [(0,1)]
}
