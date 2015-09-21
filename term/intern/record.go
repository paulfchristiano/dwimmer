package intern

import (
	"errors"
	"math/rand"
	"time"

	"github.com/paulfchristiano/dwimmer/storage/database"
)

var WrongType = errors.New("wrong type of record provided")

type Recorder struct {
	packCaches   map[int](map[interface{}]int64)
	unpackCaches map[int](map[int64]interface{})
	database.C
	Accesses int
}

func NewRecorder(c database.C) *Recorder {
	return &Recorder{
		packCaches:   make(map[int]map[interface{}]int64),
		unpackCaches: make(map[int]map[int64]interface{}),
		C:            c,
		Accesses:     0,
	}
}

var _ Packer = &Recorder{}

type Record struct {
	Value interface{}
}

func (r Record) packable() {}

func (rec *Recorder) packCache(n int) map[interface{}]int64 {
	result, ok := rec.packCaches[n]
	if !ok {
		result = make(map[interface{}]int64)
		rec.packCaches[n] = result
	}
	return result
}

func (rec *Recorder) unpackCache(n int) map[int64]interface{} {
	result, ok := rec.unpackCaches[n]
	if !ok {
		result = make(map[int64]interface{})
		rec.unpackCaches[n] = result
	}
	return result
}

var source = rand.NewSource(time.Now().UnixNano())

func (rec *Recorder) CachePack(n int, key interface{}, value Packed) Packed {
	var taken interface{}
	var cacheKey int64
	for taken = true; taken != nil; taken = rec.Get(cacheKey) {
		if taken == key {
			return Record{cacheKey}
		}
		cacheKey = source.Int63()
		rec.Accesses++
	}
	rec.Set(cacheKey, value.(Record).Value)
	rec.Accesses++
	rec.packCache(n)[key] = cacheKey
	rec.unpackCache(n)[cacheKey] = key
	return Record{cacheKey}
}

func indirect(r Record) (int64, bool) {
	val, ok := r.Value.(int64)
	return val, ok
}

func (rec *Recorder) GetCachedPack(n int, key interface{}) (Packed, bool) {
	result, ok := rec.packCache(n)[key]
	return Record{result}, ok
}

func (rec *Recorder) CacheUnpack(n int, key Packed, value interface{}) {
	if k, isIndirect := indirect(key.(Record)); isIndirect {
		rec.unpackCache(n)[k] = value
		rec.packCache(n)[value] = k
	}
}

func (rec *Recorder) GetCachedUnpack(n int, key Packed) (interface{}, bool) {
	r := key.(Record)
	if k, isIndirect := indirect(r); isIndirect {
		result, ok := rec.unpackCache(n)[k]
		if ok {
			return result, true
		}
		result = rec.Get(k)
		rec.Accesses++
		rec.unpackCache(n)[k] = result
		return result, true
	}
	return nil, false
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

func (_ *Recorder) UnpackString(record Packed) string {
	result, ok := record.(Record).Value.(string)
	if !ok {
		panic(WrongType)
	}
	return result
}

func (_ *Recorder) UnpackInt(record Packed) int {
	result, ok := record.(Record).Value.(int)
	if !ok {
		panic(WrongType)
	}
	return result
}

func (_ *Recorder) UnpackPair(record Packed) (Packed, Packed) {
	elems, ok := record.(Record).Value.([]interface{})
	if !ok {
		panic(WrongType)
	}
	return Record{elems[0]}, Record{elems[1]}
}

func (_ *Recorder) UnpackList(record Packed) []Packed {
	elems, ok := record.(Record).Value.([]interface{})
	if !ok {
		panic(WrongType)
	}
	result := make([]Packed, len(elems))
	for i, elem := range elems {
		result[i] = Record{elem}
	}
	return result
}

func (_ *Recorder) UnpackLast(record Packed) Packed {
	result, ok := record.(Record).Value.([]interface{})
	if !ok {
		panic(WrongType)
	}
	return Record{result[len(result)-1]}
}

func (_ *Recorder) UnpackInit(record Packed) Packed {
	result, ok := record.(Record).Value.([]interface{})
	if !ok {
		panic(WrongType)
	}
	return Record{result[:len(result)-1]}
}

func MakeRecord(b interface{}) Record {
	return Record{b}
	/*
		switch b := b.(type) {
		case []interface{}:
			result := make([]Packed, len(b))
			for i, x := range b {
				result[i] = MakeRecord(x)
			}
			return Record{result}
		}
		return Record{b}
	*/
}

func FromRecord(r Record) interface{} {
	return r.Value
	/*
		switch v := r.Value.(type) {
		case []Packed:
			result := make([]interface{}, len(v))
			for i, x := range v {
				result[i] = FromRecord(x.(Record))
			}
			return result
		}
		return r.Value
	*/
}
