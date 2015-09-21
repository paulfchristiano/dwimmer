package intern

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

	CachePack(int, interface{}, Packed) Packed
	GetCachedPack(int, interface{}) (Packed, bool)
	CacheUnpack(int, Packed, interface{})
	GetCachedUnpack(int, Packed) (interface{}, bool)
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

	packCaches   map[int]map[interface{}]ID
	unpackCaches map[int]map[ID]interface{}
}

func NewIDer() *IDer {
	return &IDer{
		intIDs:       make(map[int]ID),
		stringIDs:    make(map[string]ID),
		pairIDs:      make(map[idPair]ID),
		packCaches:   make(map[int](map[interface{}]ID)),
		unpackCaches: make(map[int](map[ID]interface{})),
	}
}

var _ Packer = (*IDer)(nil)

type idPair struct {
	a, b ID
}

func (ider *IDer) packCache(n int) map[interface{}]ID {
	result, ok := ider.packCaches[n]
	if !ok {
		result = make(map[interface{}]ID)
		ider.packCaches[n] = result
	}
	return result
}

func (ider *IDer) unpackCache(n int) map[ID]interface{} {
	result, ok := ider.unpackCaches[n]
	if !ok {
		result = make(map[ID]interface{})
		ider.unpackCaches[n] = result
	}
	return result
}

func (ider *IDer) CachePack(n int, key interface{}, value Packed) Packed {
	ider.packCache(n)[key] = value.(ID)
	ider.unpackCache(n)[value.(ID)] = key
	return value
}

func (ider *IDer) GetCachedPack(n int, key interface{}) (Packed, bool) {
	result, ok := ider.packCache(n)[key]
	return result, ok
}

func (ider *IDer) GetCachedUnpack(n int, key Packed) (interface{}, bool) {
	result, ok := ider.unpackCache(n)[key.(ID)]
	return result, ok
}

func (ider *IDer) CacheUnpack(n int, key Packed, value interface{}) {
	ider.unpackCache(n)[key.(ID)] = value
	ider.packCache(n)[value] = key.(ID)
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
	reverse(result)
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
