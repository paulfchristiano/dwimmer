package intern

type Packed interface {
	Str() string
	Int() int
	Pair() (Packed, Packed)
	List() []Packed
	New() Packed
	StoreInt(int) Packed
	StoreStr(string) Packed
	StoreList([]Packed) Packed
	StorePair(Packed, Packed) Packed
	Empty() Packed
	Append(Packed) Packed
}

type Id int

var (
	strings   = make([]string, 0)
	stringIds = make(map[string]Id)
)

func (id *Id) StoreStr(s string) Packed {
	result, ok := stringIds[s]
	if !ok {
		result = Id(len(strings))
		strings = append(strings, s)
		stringIds[s] = result
	}
	*id = result
	return id
}

func (id *Id) StoreInt(n int) Packed {
	*id = Id(n)
	return id
}

type idPair struct {
	a, b Id
}

var (
	pairs   = make([]idPair, 0)
	pairIds = make(map[idPair]Id)
)

func (id *Id) StorePair(x, y Packed) Packed {
	pair := idPair{*x.(*Id), *y.(*Id)}
	result, ok := pairIds[pair]
	if !ok {
		result = Id(len(pairs))
		pairs = append(pairs, pair)
		pairIds[pair] = result
	}
	*id = result
	return id
}

func (id *Id) New() Packed {
	return new(Id)
}

func (id *Id) StoreList(xs []Packed) Packed {
	//FIXME this is an awkward way to handle the end of a list...
	*id = Id(-1)
	for _, x := range xs {
		id.StorePair(x, id)
	}
	return id
}

func (id *Id) Empty() Packed {
	*id = -1
	return id
}

func (id Id) Int() int {
	return int(id)
}

func (id Id) Str() string {
	return strings[id]
}

func (id Id) Pair() (Packed, Packed) {
	pair := pairs[id]
	return &pair.a, &pair.b
}

func (id Id) List() []Packed {
	result := make([]Packed, 0)
	var next, p Packed
	for id >= 0 {
		next, p = id.Pair()
		id = *p.(*Id)
		result = append(result, next)
	}
	reverse(result)
	return result
}

func (id Id) Init() Packed {
	_, init := id.Pair()
	return init
}

func (id Id) Last() Packed {
	last, _ := id.Pair()
	return last
}

func reverse(l []Packed) {
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
}

func (id *Id) Append(other Packed) Packed {
	id.StorePair(other, id)
	return id
}
