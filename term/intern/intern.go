package intern

type Packer interface {
	PackString(string) Packed
	PackInt(int) Packed
	PackList([]Packed) Packed
	PackPair(Packed, Packed) Packed
	PackPickler(Pickler) Packed
	AppendToPacked(Packed, Packed) Packed

	UnpackString(Packed) (string, bool)
	UnpackInt(Packed) (int, bool)
	UnpackList(Packed) ([]Packed, bool)
	UnpackPair(Packed) (Packed, Packed, bool)
	UnpackPickler(Packed, Pickler) (Pickler, bool)
	UnpackLast(Packed) (Packed, bool)
	UnpackInit(Packed) (Packed, bool)
}

type Packed interface {
	packable()
}
