package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"vlg/tools/model"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
	"github.com/storozhukBM/verifier"
	"github.com/timshannon/badgerhold/v4"
)

const dataDirectory = "data"

func main() {
	verify := verifier.New()
	verify.That(len(os.Args) == 2, "Usage: preload <directory>")
	verify.That(directoryExists(os.Args[1]), "Directory does not exist: "+os.Args[1])
	verify.PanicOnError()

	files := []struct {
		name string
		f    func(*os.File, *badgerhold.Store) error
	}{
		{
			"nodes-addresses.csv",
			loadAddresses,
		},
		{
			"nodes-entities.csv",
			loadEntities,
		},
		{
			"nodes-intermediaries.csv",
			loadIntermediaries,
		},
		{
			"nodes-officers.csv",
			loadOfficers,
		},
		{
			"nodes-others.csv",
			loadOthers,
		},
		{
			"relationships.csv",
			loadRelationships,
		},
	}
	for _, file := range files {
		verify.That(fileExists(os.Args[1]+"/"+file.name), "File does not exist: "+file.name)
	}
	verify.PanicOnError()

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		return gocsv.LazyCSVReader(in) // Allows use of quotes in CSV
	})

	options := badgerhold.DefaultOptions
	options.Dir = dataDirectory
	options.ValueDir = dataDirectory
	store, err := badgerhold.Open(options)
	defer store.Close()
	if err != nil {
		panic(err)
	}

	if true {
		err = updateRelationships(store)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		f, err := os.OpenFile(os.Args[1]+"/"+file.name, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		log.Println("Processing", file.name)
		go func(file struct {
			name string
			f    func(*os.File, *badgerhold.Store) error
		}, store *badgerhold.Store) {
			err = file.f(f, store)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(file, store)
	}
	wg.Wait()

	err = updateRelationships(store)
	if err != nil {
		panic(err)
	}
}

func loadAddresses(f *os.File, store *badgerhold.Store) error {
	addresses := []*model.Address{}
	if err := gocsv.UnmarshalFile(f, &addresses); err != nil {
		return err
	}
	log.Printf("Upserting %d addresses...", len(addresses))
	for i := range addresses {
		err := store.Upsert(addresses[i].NodeID, addresses[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func loadEntities(f *os.File, store *badgerhold.Store) error {
	entities := []*model.Entity{}
	if err := gocsv.UnmarshalFile(f, &entities); err != nil {
		return err
	}
	log.Printf("Upserting %d entities...", len(entities))
	for i := range entities {
		err := store.Upsert(entities[i].NodeID, entities[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func loadIntermediaries(f *os.File, store *badgerhold.Store) error {
	intermediaries := []*model.Intermediary{}
	if err := gocsv.UnmarshalFile(f, &intermediaries); err != nil {
		return err
	}
	log.Printf("Upserting %d intermediaries...", len(intermediaries))
	for i := range intermediaries {
		err := store.Upsert(intermediaries[i].NodeID, intermediaries[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func loadOfficers(f *os.File, store *badgerhold.Store) error {
	officers := []*model.Officer{}
	if err := gocsv.UnmarshalFile(f, &officers); err != nil {
		return err
	}
	log.Printf("Upserting %d officers...", len(officers))
	for i := range officers {
		err := store.Upsert(officers[i].NodeID, officers[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func loadRelationships(f *os.File, store *badgerhold.Store) error {
	relationships := []*model.Relationship{}
	if err := gocsv.UnmarshalFile(f, &relationships); err != nil {
		return err
	}
	log.Printf("Upserting %d relationships...", len(relationships))
	for i := range relationships {
		key := fmt.Sprintf("%d-%d", relationships[i].FromID, relationships[i].ToID)
		err := store.Upsert(key, relationships[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func loadOthers(f *os.File, store *badgerhold.Store) error {
	others := []*model.Other{}
	if err := gocsv.UnmarshalFile(f, &others); err != nil {
		return err
	}
	log.Printf("Upserting %d others...", len(others))
	for i := range others {
		err := store.Upsert(others[i].NodeID, others[i])
		if err != nil {
			return err
		}
	}
	return nil
}

const xactLimit = 20000

func updateRelationships(store *badgerhold.Store) error {
	q := &badgerhold.Query{}
	n := 0
	tx := store.Badger().NewTransaction(true)
	//defer tx.Discard()
	err := store.ForEach(q, func(record *model.Relationship) error {
		n++
		if n%xactLimit == 0 {
			log.Printf("Updating relationships... progress %d", n)
			tx = store.Badger().NewTransaction(true)
		}
		var err error
		_, record.FromType, err = model.RecordByID(store, record.FromID)
		if err != nil {
			return err
		}
		_, record.ToType, err = model.RecordByID(store, record.ToID)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%d-%d", record.FromID, record.ToID)
		err = store.TxUpdate(tx, key, record)
		if n%xactLimit == 0 {
			err := tx.Commit()
			if err != nil {
				return err
			}
			tx = store.Badger().NewTransaction(true)
		}
		return err
	})
	if err != nil {
		return err
	}
	return tx.Commit()
}

// DirectoryExists return true if the path exists.
func directoryExists(path string) bool {
	if fileinfo, err := os.Stat(path); err == nil {
		return fileinfo.IsDir()
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		panic(err)
	}
}

// FileExists return true if the file exists.
func fileExists(path string) bool {
	if fileinfo, err := os.Stat(path); err == nil {
		return !fileinfo.IsDir()
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		panic(err)
	}
}
