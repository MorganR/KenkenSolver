package kenken

import "sort"

type ByteMap struct {
	m    map[byte]int
	size int
}

func NewByteMap() *ByteMap {
	return &ByteMap{make(map[byte]int), 0}
}

func (m *ByteMap) Add(i byte) {
	n, _ := (*m).m[i]
	(*m).m[i] = n + 1
	(*m).size++
}

func (m *ByteMap) Map() map[byte]int {
	return m.m
}

func (m *ByteMap) Len() int {
	return (*m).size
}

func (m *ByteMap) Copy() ByteMap {
	c := *NewByteMap()
	for k, v := range (*m).m {
		c.m[k] = v
	}
	c.size = (*m).size
	return c
}

func (m *ByteMap) GetSortedList() []byte {
	list := make([]byte, m.Len())
	i := 0
	for k, v := range (*m).m {
		for ; v > 0; v-- {
			list[i] = k
			i++
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
	return list
}

func (m *ByteMap) Equals(o *ByteMap) bool {
	if m.Len() != o.Len() {
		return false
	}
	list := m.GetSortedList()
	otherList := o.GetSortedList()
	for i, v := range list {
		if v != otherList[i] {
			return false
		}
	}
	return true
}
