package intern

type Record struct {
	Value interface{}
}

func (record *Record) New() Packed {
	return &Record{}
}

func (record *Record) Empty() Packed {
	record.Value = []interface{}{}
	return record
}

func (record *Record) StoreStr(s string) Packed {
	record.Value = s
	return record
}

func (record *Record) StoreInt(n int) Packed {
	record.Value = n
	return record
}

func (record *Record) StorePair(x, y Packed) Packed {
	record.Value = []interface{}{x.(*Record).Value, y.(*Record).Value}
	return record
}

func (record *Record) StoreList(xs []Packed) Packed {
	result := make([]interface{}, len(xs))
	for i, x := range xs {
		result[i] = x.(*Record).Value
	}
	record.Value = result
	return record
}

func (record *Record) Append(other Packed) Packed {
	record.Value = append((record.Value).([]interface{}), other.(*Record).Value)
	return record
}

func (record *Record) Str() string {
	return record.Value.(string)
}

func (record *Record) Int() int {
	return record.Value.(int)
}

func (record *Record) Pair() (Packed, Packed) {
	l := record.Value.([]interface{})
	return &Record{l[0]}, &Record{l[1]}
}

func (record *Record) List() []Packed {
	values := record.Value.([]interface{})
	result := make([]Packed, len(values))
	for i, value := range values {
		result[i] = &Record{value}
	}
	return result
}
