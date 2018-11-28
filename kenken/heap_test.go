package kenken

import (
	"container/heap"
	"testing"
)

func TestHeapInit(t *testing.T) {
	h := make(BoxHeap, 0, 4)
	if h.Len() != 0 {
		t.Errorf("Heap length did not start at 0")
	}
	size := uint8(2)
	h.Push(NewBox(Index{0, 0}, size))
	h[0].AddPossible(1)
	h.Push(NewBox(Index{0, 1}, size))
	h.Push(NewBox(Index{1, 0}, size))
	h[2].AddPossible(2)
	h[2].AddPossible(1)
	h.Push(NewBox(Index{1, 1}, size))
	h[3].AddPossible(2)

	heap.Init(&h)

	previous := uint8(0)
	for h.Len() > 0 {
		top := heap.Pop(&h).(*Box)
		if top.NumPossible() < previous {
			t.Fatal("Heap returned values out of order!")
		}
		previous = top.NumPossible()
	}
}

func TestHeapSwap(t *testing.T) {
	h := make(BoxHeap, 0, 2)
	size := uint8(2)
	b0 := NewBox(Index{0, 0}, size)
	b0.AddPossible(1)
	h.Push(b0)
	b1 := NewBox(Index{0, 1}, size)
	h.Push(b1)

	h.Swap(0, 1)
	if h[0] != b1 || h[1] != b0 {
		t.Fatalf("Did not swap values properly. Result: %v", h)
	}
	if b0.heapIndex != 1 || b1.heapIndex != 0 {
		t.Fatalf("Did not update heapIndex of swapped values")
	}
}

func TestHeapPopPush(t *testing.T) {
	h := make(BoxHeap, 0, 2)
	size := uint8(2)
	smallBox := NewBox(Index{0, 1}, size)
	heap.Push(&h, smallBox)
	bigBox := NewBox(Index{0, 0}, size)
	bigBox.AddPossible(1)
	heap.Push(&h, bigBox)
	if bigBox.heapIndex != 1 || smallBox.heapIndex != 0 {
		t.Errorf("Did not update heapIndex of pushed values.\nValues:\n%v\n%v", smallBox, bigBox)
	}
	if h.Len() != 2 {
		t.Fatalf("Heap length did not update after push")
	}
	top := heap.Pop(&h).(*Box)
	if top != smallBox {
		t.Errorf("Popped the wrong box first: %v", top)
	}
	if top.heapIndex != -1 {
		t.Errorf("Did not update heapIndex of popped value: %v", top)
	}
	if h.Len() != 1 {
		t.Errorf("Did not decrement length on pop")
	}
	top = heap.Pop(&h).(*Box)
	if top != bigBox {
		t.Errorf("Popped the wrong box second: %v", top)
	}
	if top.heapIndex != -1 {
		t.Errorf("Did not update heapIndex of popped value: %v", top)
	}
	if h.Len() != 0 {
		t.Errorf("Did not decrement length on pop")
	}
}

func TestHeapFix(t *testing.T) {
	h := make(BoxHeap, 0, 2)
	size := uint8(2)
	smallBox := *NewBox(Index{0, 1}, size)
	heap.Push(&h, &smallBox)
	bigBox := *NewBox(Index{0, 0}, size)
	bigBox.AddPossible(1)
	heap.Push(&h, &bigBox)
	smallBox.AddPossible(1)
	smallBox.AddPossible(0)
	heap.Fix(&h, smallBox.heapIndex)
	top := heap.Pop(&h).(*Box)
	if top != &bigBox {
		t.Error("Smaller box was not returned first after fix")
	}
	top = heap.Pop(&h).(*Box)
	if top != &smallBox {
		t.Error("Bigger box was not returned last")
	}
}
