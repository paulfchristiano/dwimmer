package storage

import (
	"fmt"

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
	fmt.Println("Closing storage...")
	saved := term.SaveT(s.current)
	fmt.Println("Saved state:", saved)
	s.collection.Set(s.name, saved)
	fmt.Println("Stored in database!")
	s.collection.Set(s.collection.Count(), saved)
	fmt.Println("Backed up in database!")
}

func NewStorage(name string) *DBStorage {
	collection := database.Collection("newterms")
	stateRecord, ok := collection.Get(name)
	var state term.T
	if ok {
		state, ok = term.LoadT(stateRecord)
		if !ok {
			fmt.Printf("failed to load state %v\n", stateRecord)
		}
	}
	if state == nil {
		state = term.Make("an uninitialized state").T()
	}
	return &DBStorage{
		collection: collection,
		name:       name,
		current:    state,
	}
}
