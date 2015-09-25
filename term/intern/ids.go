package intern

type ID int

func (id ID) packable() {}

type IDer struct {
	nextID    ID
	ints      map[ID]int
	intIDs    map[int]ID
	strings   map[ID]string
	stringIDs map[string]ID
	pairs     map[ID]idPair
	pairIDs   map[idPair]ID

	packCache   map[interface{}]Packed
	unpackCache map[Packed]Pickler
}

func NewIDer() *IDer {
	return &IDer{
		nextID:      ID(0),
		ints:        make(map[ID]int),
		intIDs:      make(map[int]ID),
		strings:     make(map[ID]string),
		stringIDs:   make(map[string]ID),
		pairIDs:     make(map[idPair]ID),
		pairs:       make(map[ID]idPair),
		packCache:   make(map[interface{}]Packed),
		unpackCache: make(map[Packed]Pickler),
	}
}

var _ Packer = (*IDer)(nil)

type idPair struct {
	a, b ID
}

func (ider *IDer) NextID() ID {
	ider.nextID++
	return ider.nextID
}

func (ider *IDer) PackString(s string) Packed {
	id, ok := ider.stringIDs[s]
	if !ok {
		id = ider.NextID()
		ider.strings[id] = s
		ider.stringIDs[s] = id
	}
	return id
}

func (ider *IDer) PackInt(n int) Packed {
	id, ok := ider.intIDs[n]
	if !ok {
		id = ider.NextID()
		ider.ints[id] = n
		ider.intIDs[n] = id
	}
	return id
}

func (ider *IDer) PackPair(x, y Packed) Packed {
	pair := idPair{x.(ID), y.(ID)}
	id, ok := ider.pairIDs[pair]
	if !ok {
		id = ider.NextID()
		ider.pairs[id] = pair
		ider.pairIDs[pair] = id
	}
	return id
}

func (ider *IDer) PackList(xs []Packed) Packed {
	var result Packed
	result = ID(-1)
	for _, x := range xs {
		result = ider.PackPair(x, result)
	}
	return result
}

func (ider *IDer) PackPickler(pickler Pickler) Packed {
	key := pickler.Key()
	result, ok := ider.packCache[key]
	if !ok || key == nil {
		result = pickler.Pickle(ider)
		ider.packCache[key] = result
		ider.unpackCache[result] = pickler
	}
	return result
}

func (ider *IDer) UnpackPickler(id Packed, pickler Pickler) (result Pickler, ok bool) {
	defer func() { ok = ok && pickler.Test(result) }()
	result, ok = ider.unpackCache[id.(ID)]
	if !ok {
		result, ok := pickler.Unpickle(ider, id)
		if ok {
			ider.unpackCache[id] = result
		}
	}
	return result, ok
}

func (ider *IDer) UnpackInt(id Packed) (int, bool) {
	result, ok := ider.ints[id.(ID)]
	return result, ok
}

func (ider *IDer) UnpackString(id Packed) (string, bool) {
	result, ok := ider.strings[id.(ID)]
	return result, ok
}

func (ider *IDer) UnpackPair(id Packed) (Packed, Packed, bool) {
	pair, ok := ider.pairs[id.(ID)]
	return pair.a, pair.b, ok
}

func (ider *IDer) UnpackList(x Packed) ([]Packed, bool) {
	var result []Packed
	var next Packed
	var ok bool
	for id := x.(ID); id >= 0; id = x.(ID) {
		next, x, ok = ider.UnpackPair(id)
		if !ok {
			return nil, false
		}
		result = append(result, next)
	}
	reverse(result)
	return result, true
}

func (ider *IDer) UnpackInit(list Packed) (Packed, bool) {
	_, init, ok := ider.UnpackPair(list)
	return init, ok
}

func (ider *IDer) UnpackLast(list Packed) (Packed, bool) {
	last, _, ok := ider.UnpackPair(list)
	return last, ok
}

func reverse(l []Packed) {
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
}

func (ider *IDer) AppendToPacked(packedList, other Packed) Packed {
	return ider.PackPair(other, packedList)
}
