package intern

type Recorder struct{}

var _ Packer = Recorder{}

type Record struct {
	Value interface{}
}

func (r Record) packable() {}

func (_ Recorder) PackString(s string) Packed {
	return Record{s}
}

func (_ Recorder) PackInt(n int) Packed {
	return Record{n}
}

func (_ Recorder) PackPair(x, y Packed) Packed {
	return Record{[]Packed{x, y}}
}

func (_ Recorder) PackList(xs []Packed) Packed {
	return Record{xs}
}

func (_ Recorder) AppendToPacked(init, last Packed) Packed {
	result := init.(Record).Value.([]Packed)
	result = append(result, last)
	return Record{result}
}

func (_ Recorder) UnpackString(record Packed) string {
	return record.(Record).Value.(string)
}

func (_ Recorder) UnpackInt(record Packed) int {
	return record.(Record).Value.(int)
}

func (_ Recorder) UnpackPair(record Packed) (Packed, Packed) {
	elems := record.(Record).Value.([]Packed)
	return elems[0], elems[1]
}

func (_ Recorder) UnpackList(record Packed) []Packed {
	return record.(Record).Value.([]Packed)
}

func (_ Recorder) UnpackLast(record Packed) Packed {
	result := record.(Record).Value.([]Packed)
	return result[len(result)-1]
}

func (_ Recorder) UnpackInit(record Packed) Packed {
	result := record.(Record).Value.([]Packed)
	return Record{result[:len(result)-1]}
}
