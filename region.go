package kenken

import (
	"fmt"
)

type Operation uint8

const (
	Sum     Operation = 1
	Sub     Operation = 2
	Mul     Operation = 3
	Div     Operation = 4
	Nothing Operation = 5
)

func (o Operation) String() string {
	switch o {
	case Sum:
		return "Sum"
	case Sub:
		return "Sub"
	case Mul:
		return "Mul"
	case Div:
		return "Div"
	case Nothing:
		return "Nothing"
	default:
		return "Unknown"
	}
}

type Region struct {
	result  uint
	op      Operation
	indices IndexSet
}

func (r Region) GetResult() uint {
	return r.result
}

func (r Region) GetOp() Operation {
	return r.op
}

func (r Region) GetIndices() []Index {
	idxs := make([]Index, len(r.indices))
	i := 0
	for idx := range r.indices {
		idxs[i] = idx
		i++
	}
	return idxs
}

func (r Region) String() string {
	return fmt.Sprintf("Result: %v, Operation: %v, Indices: %v", r.result, r.op, r.indices)
}

func (r *Region) GetPossibleMaps(size uint8) ByteMapList {
	switch (*r).op {
	case Sum:
		return r.getSumMaps(size)
	case Sub:
		return r.getSubMaps(size)
	case Mul:
		return r.getMulMaps(size)
	case Div:
		return r.getDivMaps(size)
	case Nothing:
		return r.getNothingMaps()
	}
	return nil
}

func (r *Region) getNothingMaps() ByteMapList {
	maps := make(ByteMapList, 1)
	maps[0] = *NewByteMap()
	maps[0].Add(byte((*r).result))
	return maps
}

func (r *Region) getSumMaps(size uint8) ByteMapList {
	numArgs := r.indices.Len()
	result := r.result
	maps := getSumMapsForResult(size, uint(numArgs), result)
	return maps
}

func (r *Region) getSubMaps(size uint8) ByteMapList {
	numArgs := uint(r.indices.Len())
	result := uint(r.result)
	key := opMapKey{Sub, size, numArgs, result}
	maps, present := opMaps[key]
	if present {
		// return maps
	}
	maps = make(ByteMapList, 0)
	for i := uint8(result + 1); i <= size; i++ {
		innerMaps := getSumMapsForResult(size, uint(numArgs-1), uint(i)-result)
		maps.appendValueAndAdd(&innerMaps, i)
	}
	return maps
}

func (r *Region) getMulMaps(size uint8) ByteMapList {
	numArgs := r.indices.Len()
	result := r.result
	maps := getMulMapsForResult(size, uint(numArgs), result)
	return maps
}

func (r *Region) getDivMaps(size uint8) ByteMapList {
	numArgs := uint(r.indices.Len())
	result := uint(r.result)
	key := opMapKey{Div, size, numArgs, result}
	maps, present := opMaps[key]
	if present {
		// return maps
	}
	maps = make(ByteMapList, 0)
	for i := uint8(1); i <= size; i++ {
		if i%uint8(result) != 0 {
			continue
		}
		innerMaps := getMulMapsForResult(size, numArgs-1, uint(i)/result)
		maps.appendValueAndAdd(&innerMaps, i)
	}
	return maps
}

type opMapKey struct {
	op      Operation
	size    uint8
	numArgs uint
	result  uint
}

type ByteMapList []ByteMap

func (l *ByteMapList) Contains(bm *ByteMap) bool {
	for _, m := range *l {
		if m.Equals(bm) {
			return true
		}
	}
	return false
}

func (l *ByteMapList) appendValueAndAdd(il *ByteMapList, v byte) {
	for _, innerMap := range *il {
		innerMap = innerMap.Copy()
		innerMap.Add(v)
		if l.Contains(&innerMap) {
			continue
		}
		*l = append(*l, innerMap)
	}
}

var opMaps = make(map[opMapKey]ByteMapList)

func getSumMapsForResult(size uint8, numArgs uint, result uint) ByteMapList {
	key := opMapKey{Sum, size, numArgs, result}
	maps, present := opMaps[key]
	if present {
		// return maps
	}
	maps = make(ByteMapList, 0)
	for i := uint8(1); i <= size && uint(i) <= result; i++ {
		if numArgs == 1 {
			if result != uint(i) {
				continue
			}
			m := *NewByteMap()
			m.Add(i)
			maps = append(maps, m)
		} else {
			if result <= uint(i) {
				break
			}
			innerMaps := getSumMapsForResult(size, numArgs-1, result-uint(i))
			maps.appendValueAndAdd(&innerMaps, i)
		}
	}
	opMaps[key] = maps
	return maps
}

func getMulMapsForResult(size uint8, numArgs uint, result uint) ByteMapList {
	key := opMapKey{Mul, size, numArgs, result}
	maps, present := opMaps[key]
	if present {
		// return maps
	}
	maps = make(ByteMapList, 0)
	if numArgs == 1 {
		if result <= uint(size) {
			m := *NewByteMap()
			m.Add(byte(result))
			maps = append(maps, m)
		}
	} else {
		for i := uint8(1); i <= size; i++ {
			if result%uint(i) != 0 {
				continue
			}
			innerMaps := getMulMapsForResult(size, numArgs-1, result/uint(i))
			maps.appendValueAndAdd(&innerMaps, i)
		}
	}
	opMaps[key] = maps
	return maps
}
