package storage

import "github.com/paulfchristiano/dwimmer/term"

func Dummy() *DummyStorage {
	return &DummyStorage{}
}

var (
	DummyState = term.Make("a placeholder returned when accessing the state of a test harness")
)

type DummyStorage struct{}

func (_ *DummyStorage) CloseStorage()      {}
func (_ *DummyStorage) GetStorage() term.T { return DummyState.T() }
func (_ *DummyStorage) SetStorage(term.T)  {}

var _ StorageImplementer = &DummyStorage{}
