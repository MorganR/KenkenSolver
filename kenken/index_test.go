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
