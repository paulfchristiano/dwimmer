package intern

import "fmt"

type Packer interface {
	PackString(string) Packed
	PackInt(int) Packed
	PackList([]Packed) Packed
	PackPair(Packed, Packed) Packed
	AppendToPacked(Packed, Packed) Packed

	UnpackString(Packed) string
	UnpackInt(Packed) int
	UnpackList(Packed) []Packed
	UnpackPair(Packed) (Packed, Packed)
	UnpackLast(Packed) Packed
	UnpackInit(Packed) Packed
}

type Packed interface {
	packable()
}

type ID int

func (id ID) packable() {}

type IDer struct {
	ints      []int
	intIDs    map[int]ID
	strings   []string
	stringIDs map[string]ID
	pairs     []idPair
	pairIDs   map[idPair]ID
}

func NewIDer() *IDer {
	return &IDer{
		intIDs:    make(map[int]ID),
		stringIDs: make(map[string]ID),
		pairIDs:   make(map[idPair]ID),
	}
}

var DefaultIDer Packer = NewIDer()

type idPair struct {
	a, b ID
}

func (ider *IDer) PackString(s string) Packed {
	result, ok := ider.stringIDs[s]
	if !ok {
		result = ID(len(ider.strings))
		ider.strings = append(ider.strings, s)
		ider.stringIDs[s] = result
	}
	return result
}

func (ider *IDer) PackInt(n int) Packed {
	return ID(n)
}

func (ider *IDer) PackPair(x, y Packed) Packed {
	pair := idPair{x.(ID), y.(ID)}
	result, ok := ider.pairIDs[pair]
	if !ok {
		result = ID(len(ider.pairs))
		ider.pairs = append(ider.pairs, pair)
		ider.pairIDs[pair] = result
	}
	return result
}

func (ider *IDer) PackList(xs []Packed) Packed {
	var result Packed
	result = ID(-1)
	for _, x := range xs {
		result = ider.PackPair(x, result)
	}
	return result
}

func (ider *IDer) UnpackInt(id Packed) int {
	return int(id.(ID))
}

func (ider *IDer) UnpackString(id Packed) string {
	return ider.strings[id.(ID)]
}

func (ider *IDer) UnpackPair(id Packed) (Packed, Packed) {
	pair := ider.pairs[id.(ID)]
	return pair.a, pair.b
}

func (ider *IDer) UnpackList(x Packed) []Packed {
	var result []Packed
	var next Packed
	for id := x.(ID); id >= 0; id = x.(ID) {
		next, x = ider.UnpackPair(id)
		result = append(result, next)
	}
	fmt.Println(result)
	reverse(result)
	fmt.Println(result)
	return result
}

func (ider *IDer) UnpackInit(list Packed) Packed {
	_, init := ider.UnpackPair(list)
	return init
}

func (ider *IDer) UnpackLast(list Packed) Packed {
	last, _ := ider.UnpackPair(list)
	return last
}

func reverse(l []Packed) {
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
}

func (ider *IDer) AppendToPacked(packedList, other Packed) Packed {
	return ider.PackPair(other, packedList)
}
