package storage

import (
	"github.com/paulfchristiano/dwimmer/storage/database"
	"github.com/paulfchristiano/dwimmer/term"
)

type StorageImplementer interface {
	CloseStorage()
	GetStorage() term.T
	SetStorage(term.T)
}

type DBStorage struct {
	collection database.C
	name       string
	current    term.T
}

func (s *DBStorage) SetStorage(t term.T) {
	s.current = t
}

func (s *DBStorage) GetStorage() term.T {
	return s.current
}

func (s *DBStorage) CloseStorage() {
	//TODO this mechanism could clearly be nicer
	saved := term.SaveT(s.current)
	s.collection.Set(s.name, saved)
	s.collection.Set(s.collection.Count(), saved)
}

func NewStorage(name string) *DBStorage {
	collection := database.Collection("newterms")
	stateRecord := collection.Get(name)
	state := term.Make("an uninitialized state").T()
	if stateRecord != nil {
		state = term.LoadT(collection.Get(name))
	}
	return &DBStorage{
		collection: collection,
		name:       name,
		current:    state,
	}
}
