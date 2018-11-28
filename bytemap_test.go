package kenken

import "testing"

func TestByteMapAddGet(t *testing.T) {
	m := NewByteMap()
	m.Add(1)
	m.Add(2)
	m.Add(1)
	m.Add(4)
	result := m.GetSortedList()
	expected := []byte{1, 1, 2, 4}
	if len(result) != len(expected) {
		t.Fatalf("Result was the wrong length! Should have been %v, but was %v", len(expected), len(result))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Fatalf("Unexpected bytemap result: %v", result)
		}
	}
}

func TestByteMapLen(t *testing.T) {
	m := NewByteMap()
	if m.Len() != 0 {
		t.Errorf("Initial array should have Len() == 0, not %v", m.Len())
	}
	m.Add(1)
	m.Add(1)
	m.Add(4)
	if m.Len() != 3 {
		t.Errorf("ByteMap length was %v, should have been 3", m.Len())
	}
}

func TestByteMapCopy(t *testing.T) {
	m := NewByteMap()
	m.Add(3)
	n := m.Copy()
	n.Add(3)
	if m.Len() != 1 || n.Len() != 2 {
		t.Fatal("Copy modified the original ByteMap")
	}
	listM := m.GetSortedList()
	listN := n.GetSortedList()
	if len(listN) == len(listM) {
		t.Fatal("Copy modified the original map in ByteMap")
	}
}

func TestByteMapEquals(t *testing.T) {
	m := *NewByteMap()
	m.Add(3)
	m.Add(2)
	n := m.Copy()
	if !n.Equals(&m) || !m.Equals(&m) {
		t.Fatal("Maps should be equal")
	}
	n.Add(2)
	if m.Equals(&n) || n.Equals(&m) {
		t.Fatal("Arrays should not be equal")
	}
}
