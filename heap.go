package kenken

type BoxHeap []*Box

func (h BoxHeap) Len() int { return len(h) }

func (h BoxHeap) Less(i, j int) bool {
	return h[i].NumPossible() < h[j].NumPossible()
}

func (h *BoxHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
	(*(*h)[i]).heapIndex = i
	(*(*h)[j]).heapIndex = j
}

func (h *BoxHeap) Push(b interface{}) {
	box := b.(*Box)
	(*box).heapIndex = h.Len()
	*h = append(*h, box)
}

func (h *BoxHeap) Pop() interface{} {
	b := (*h)[h.Len()-1]
	(*b).heapIndex = -1
	*h = (*h)[0 : h.Len()-1]
	return b
}
