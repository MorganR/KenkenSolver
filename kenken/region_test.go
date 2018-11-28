package kenken

import "testing"

func TestGetResult(t *testing.T) {
	indices := make(IndexSet)
	indices.Add(Index{0, 0})
	expected := uint(3)
	r := Region{expected, Nothing, indices}

	result := r.GetResult()
	if result != expected {
		t.Fatalf("Returned the wrong result: %v, expected %v", result, expected)
	}

	r.result = 2
	if result != expected {
		t.Fatalf("Modifying the original Region changed the result")
	}
}

func TestGetOp(t *testing.T) {
	indices := make(IndexSet)
	indices.Add(Index{0, 0})
	expected := Nothing
	r := Region{3, expected, indices}

	result := r.GetOp()
	if result != expected {
		t.Fatalf("Returned the wrong result: %v, expected %v", result, expected)
	}

	r.op = Sum
	if result != expected {
		t.Fatalf("Modifying the original Region changed the result")
	}
}

func TestGetIndices(t *testing.T) {
	indices := make(IndexSet)
	indices.Add(Index{0, 1})
	indices.Add(Index{0, 2})
	indices.Add(Index{1, 2})
	expected := make([]Index, 3)
	expected[0] = Index{0, 1}
	expected[1] = Index{0, 2}
	expected[2] = Index{1, 2}
	r := Region{3, Nothing, indices}

	result := r.GetIndices()
	if len(result) != len(expected) {
		t.Errorf("Returned unexpected number of indices: %v, expected %v", len(result), len(expected))
	}
	for _, expIdx := range expected {
		isFound := false
		for _, idx := range result {
			if idx == expIdx {
				isFound = true
				break
			}
		}
		if !isFound {
			t.Fatalf("Returned the wrong result: %v, expected %v", result, expected)
		}
	}

	r.indices.Add(Index{1, 1})
	r.indices.Add(Index{2, 2})
	if len(result) != len(expected) {
		t.Errorf("Number of indices changed: %v, expected %v", len(result), len(expected))
	}
	for _, expIdx := range expected {
		isFound := false
		for _, idx := range result {
			if idx == expIdx {
				isFound = true
				break
			}
		}
		if !isFound {
			t.Fatalf("Result changed: %v, expected %v", result, expected)
		}
	}
}

func TestGetNothingMaps(t *testing.T) {
	indices := make(IndexSet)
	indices.Add(Index{0, 0})
	r := Region{3, Nothing, indices}
	size := uint8(6)
	expected := ByteMapList{
		*NewByteMap(),
	}
	expected[0].Add(3)

	results := r.GetPossibleMaps(size)
	compareByteMapLists(t, &results, &expected)
}

func TestGetDivMaps(t *testing.T) {
	indices := make(IndexSet)
	indices.Add(Index{0, 0})
	indices.Add(Index{0, 1})
	indices.Add(Index{0, 2})
	r := Region{2, Div, indices}
	size := uint8(6)
	expected := ByteMapList{
		*NewByteMap(),
		*NewByteMap(),
		*NewByteMap(),
	}
	expected[0].Add(4)
	expected[0].Add(2)
	expected[0].Add(1)
	expected[1].Add(6)
	expected[1].Add(3)
	expected[1].Add(1)
	expected[2].Add(2)
	expected[2].Add(1)
	expected[2].Add(1)

	results := r.GetPossibleMaps(size)
	compareByteMapLists(t, &results, &expected)
}

func TestGetMulMaps(t *testing.T) {
	indices := make(IndexSet)
	indices.Add(Index{0, 0})
	indices.Add(Index{0, 1})
	indices.Add(Index{0, 2})
	r := Region{4, Mul, indices}
	size := uint8(5)
	expected := ByteMapList{
		*NewByteMap(),
		*NewByteMap(),
	}
	expected[0].Add(4)
	expected[0].Add(1)
	expected[0].Add(1)
	expected[1].Add(2)
	expected[1].Add(2)
	expected[1].Add(1)

	results := r.GetPossibleMaps(size)
	compareByteMapLists(t, &results, &expected)
}

func TestGetSubMapsForThreeIndices(t *testing.T) {
	indices := make(IndexSet)
	indices.Add(Index{0, 0})
	indices.Add(Index{0, 1})
	indices.Add(Index{0, 2})
	r := Region{2, Sub, indices}
	size := uint8(5)
	expected := ByteMapList{
		*NewByteMap(),
		*NewByteMap(),
	}
	expected[0].Add(4)
	expected[0].Add(1)
	expected[0].Add(1)
	expected[1].Add(5)
	expected[1].Add(2)
	expected[1].Add(1)

	results := r.GetPossibleMaps(size)
	compareByteMapLists(t, &results, &expected)
}

func TestGetSubMapsForTwoIndices(t *testing.T) {
	indices := make(IndexSet)
	indices.Add(Index{0, 0})
	indices.Add(Index{0, 1})
	r := Region{3, Sub, indices}
	size := uint8(5)
	expected := ByteMapList{
		*NewByteMap(),
		*NewByteMap(),
	}
	expected[0].Add(4)
	expected[0].Add(1)
	expected[1].Add(5)
	expected[1].Add(2)

	results := r.GetPossibleMaps(size)
	compareByteMapLists(t, &results, &expected)
}

func TestGetSumSetsForThreeResults(t *testing.T) {
	indices := make(IndexSet)
	indices.Add(Index{0, 0})
	indices.Add(Index{0, 1})
	indices.Add(Index{0, 2})
	r := Region{6, Sum, indices}
	size := uint8(3)
	expected := ByteMapList{
		*NewByteMap(),
		*NewByteMap(),
	}
	expected[0].Add(1)
	expected[0].Add(2)
	expected[0].Add(3)
	expected[1].Add(2)
	expected[1].Add(2)
	expected[1].Add(2)

	results := r.GetPossibleMaps(size)
	compareByteMapLists(t, &results, &expected)
}

func TestOpMapCache(t *testing.T) {
	result := getSumMapsForResult(3, 3, 6)
	resultb := getSumMapsForResult(3, 3, 6)
	if len(result) != len(resultb) {
		t.Errorf("Results were inconsistent")
	}

	var opMaps = make(map[opMapKey]int)
	keya := opMapKey{Sum, 2, 2, 2}
	keyb := opMapKey{Sum, 2, 2, 2}
	opMaps[keya] = 2
	v, present := opMaps[keyb]
	if v != 2 || !present {
		t.Errorf("Key not recognized. Map: %v", opMaps)
	}
}

func compareByteMapLists(t *testing.T, r, exp *ByteMapList) {
	if len(*r) != len(*exp) {
		t.Fatalf("Returned wrong length %v, expected %v", len(*r), len(*exp))
	}
	for _, v := range *r {
		if !exp.Contains(&v) {
			t.Errorf("Returned unexpected result %v", v)
		}
	}
	for _, v := range *exp {
		if !r.Contains(&v) {
			t.Errorf("Missing result %v", v)
		}
	}
}
