package kenken

import (
	"testing"
)

func TestBoxValue(t *testing.T) {
	b := NewBox(Index{0, 0}, 3)
	if b.GetValue() != 0 || b.IsValueSet() {
		t.Error("Box had value on init")
	}
	b.SetValue(2)
	if b.GetValue() != 2 || !b.IsValueSet() {
		t.Error("Box value did not update after set")
	}
	b.UnsetValue()
	if b.GetValue() != 0 || b.IsValueSet() {
		t.Error("Box value did not unset on UnsetValue")
	}
}

func TestBoxPossibles(t *testing.T) {
	b := NewBox(Index{0, 0}, 3)
	if b.NumPossible() != 0 {
		t.Error("Box had non-zero possibles on init")
	}
	b.AddPossible(2)
	b.AddPossible(2)
	b.AddPossible(1)
	if b.NumPossible() != 2 {
		t.Errorf("Box had %v num possible, expected %v", b.NumPossible(), 2)
	}
	if !b.HasPossible(2) || !b.HasPossible(1) || b.HasPossible(3) {
		t.Errorf("Box had incorrect possible values after adding 2 and 1: %v", b)
	}
	b.DeletePossible(2)
	if b.NumPossible() != 1 || b.HasPossible(2) || !b.HasPossible(1) {
		t.Errorf("Box did not properly delete 2 from possibles: %v", b)
	}
	b.DeletePossible(3)
	if b.NumPossible() != 1 || !b.HasPossible(1) {
		t.Errorf("Box should ignore delete of unknown possible")
	}
}

func TestBoxGetPossibles(t *testing.T) {
	b := NewBox(Index{0, 0}, 3)
	p := b.GetPossibles()
	if len(p) != 0 {
		t.Error("Box had non-empty possibles list on init")
	}
	b.AddPossible(2)
	b.AddPossible(2)
	b.AddPossible(3)
	b.AddPossible(1)
	b.DeletePossible(3)
	p = b.GetPossibles()
	b.AddPossible(3)
	b.AddPossible(4)
	b.AddPossible(5)
	b.AddPossible(6)
	b.DeletePossible(2)
	b.DeletePossible(1)
	if len(p) != 2 {
		t.Errorf("Box had wrong length of possibles: %v, expected %v", len(p), 2)
	}
	expected := []byte{1, 2}
	for _, exp := range expected {
		isFound := false
		for _, pos := range p {
			if exp == pos {
				isFound = true
			}
		}
		if !isFound {
			t.Errorf("Found wrong possibles: %v, expected %v", p, exp)
		}
	}
}
