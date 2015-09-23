package intern

import (
	"errors"
	"math/rand"
	"time"

	"github.com/paulfchristiano/dwimmer/storage/database"
)

var WrongType = errors.New("wrong type of record provided")

type Recorder struct {
	packCache   map[interface{}]int64
	unpackCache map[int64]Pickler
	database.C
}

func NewRecorder(c database.C) *Recorder {
	return &Recorder{
		packCache:   make(map[interface{}]int64),
		unpackCache: make(map[int64]Pickler),
		C:           c,
	}
}

var _ Packer = &Recorder{}

type Record struct {
	Value interface{}
}

func (r Record) packable() {}

var source = rand.NewSource(time.Now().UnixNano())

func (rec *Recorder) UnpackPickler(record Packed, pickler Pickler) (Pickler, bool) {
	if k, isIndirect := indirect(record.(Record)); isIndirect {
		result, ok := rec.unpackCache[k]
		if ok {
			return result, true
		}
		pickled, found := rec.Get(k)
		if !found {
			return nil, false
		}
		result, ok = pickler.Unpickle(rec, pickled)
		if !ok {
			return nil, false
		}
		rec.unpackCache[k] = result
		rec.packCache[result.Key()] = k
		return result, true
	}
	pickled, ok := UnpackPickled(rec, record)
	if !ok {
		return nil, false
	}
	return pickler.Unpickle(rec, pickled)
}

func (rec *Recorder) PackPickler(pickler Pickler) Packed {
	key := pickler.Key()
	if result, ok := rec.packCache[key]; ok && key != nil {
		return Record{result}
	}
	result := pickler.Pickle(rec)
	var cacheKey int64
	for taken := true; taken; _, taken = rec.Get(cacheKey) {
		cacheKey = source.Int63()
	}
	rec.Set(cacheKey, result)
	rec.packCache[key] = cacheKey
	rec.unpackCache[cacheKey] = pickler
	return Record{cacheKey}
}

func indirect(r Record) (int64, bool) {
	val, ok := r.Value.(int64)
	return val, ok
}

func (_ *Recorder) PackString(s string) Packed {
	return Record{s}
}

func (_ *Recorder) PackInt(n int) Packed {
	return Record{n}
}

func (_ *Recorder) PackPair(x, y Packed) Packed {
	return Record{[]interface{}{x.(Record).Value, y.(Record).Value}}
}

func (_ *Recorder) PackList(xs []Packed) Packed {
	result := make([]interface{}, len(xs))
	for i, x := range xs {
		result[i] = x.(Record).Value
	}
	return Record{result}
}

func (_ *Recorder) AppendToPacked(init, last Packed) Packed {
	result, ok := init.(Record).Value.([]interface{})
	if !ok {
		panic(WrongType)
	}
	result = append(result, last.(Record).Value)
	return Record{result}
}

func (_ *Recorder) UnpackString(record Packed) (string, bool) {
	result, ok := record.(Record).Value.(string)
	return result, ok
}

func (_ *Recorder) UnpackInt(record Packed) (int, bool) {
	result, ok := record.(Record).Value.(int)
	return result, ok
}

func (_ *Recorder) UnpackPair(record Packed) (Packed, Packed, bool) {
	elems, ok := record.(Record).Value.([]interface{})
	if !ok || len(elems) != 2 {
		return nil, nil, false
	}
	return Record{elems[0]}, Record{elems[1]}, true
}

func (_ *Recorder) UnpackList(record Packed) ([]Packed, bool) {
	elems, ok := record.(Record).Value.([]interface{})
	if !ok {
		return nil, false
	}
	result := make([]Packed, len(elems))
	for i, elem := range elems {
		result[i] = Record{elem}
	}
	return result, true
}

func (_ *Recorder) UnpackLast(record Packed) (Packed, bool) {
	result, ok := record.(Record).Value.([]interface{})
	if !ok || len(result) == 0 {
		return nil, false
	}
	return Record{result[len(result)-1]}, true
}

func (_ *Recorder) UnpackInit(record Packed) (Packed, bool) {
	result, ok := record.(Record).Value.([]interface{})
	if !ok || len(result) == 0 {
		return nil, false
	}
	return Record{result[:len(result)-1]}, true
}
