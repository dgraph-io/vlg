package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"vlg/tools/model"

	"github.com/davecgh/go-spew/spew"
	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
	"github.com/storozhukBM/verifier"
	"github.com/timshannon/badgerhold/v4"
)

const dataDirectory = "data"

var typeMap sync.Map

func main() {
	verify := verifier.New()
	verify.That(len(os.Args) == 2, "Usage: preload <directory>")
	verify.That(directoryExists(os.Args[1]), "Directory does not exist: "+os.Args[1])
	verify.PanicOnError()

	typeMap = sync.Map{}
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
	}
	for _, file := range files {
		verify.That(fileExists(os.Args[1]+"/"+file.name), "File does not exist: "+file.name)
	}
	verify.PanicOnError()

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		return gocsv.LazyCSVReader(in) // Allows use of quotes in CSV
	})

	store, err := model.GetStore(false)
	defer func() {
		_ = store.Close()
	}()
	if err != nil {
		panic(err)
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

	// Load relationships
	log.Println("Processing relationships.csv")
	f, err := os.OpenFile(os.Args[1]+"/relationships.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = loadRelationships(f, store)
	if err != nil {
		panic(err)
	}

	/*
		err = updateRelationships(store)
		if err != nil {
			panic(err)
		}
	*/
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
		typeMap.Store(addresses[i].NodeID, "address")
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
		typeMap.Store(entities[i].NodeID, "entity")
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
		typeMap.Store(intermediaries[i].NodeID, "intermediary")
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
		typeMap.Store(officers[i].NodeID, "officer")
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
		typeMap.Store(others[i].NodeID, "other")
	}
	return nil
}

func loadRelationships(f *os.File, store *badgerhold.Store) error {
	relationships := []*model.Relationship{}
	if err := gocsv.UnmarshalFile(f, &relationships); err != nil {
		return err
	}
	relMap := make(map[string]int)
	log.Printf("Upserting %d relationships...", len(relationships))
	for i := range relationships {
		var err error
		var ok bool
		key := fmt.Sprintf("%d-%d", relationships[i].FromID, relationships[i].ToID)
		record := relationships[i]
		t, ok := typeMap.Load(record.FromID)
		if !ok {
			return errors.Errorf("Finding FromID, From: %d To: %d Type %s", record.FromID, record.ToID, record.RelationshipType)
		}
		record.FromType = t.(string)
		t, ok = typeMap.Load(record.ToID)
		if !ok {
			return errors.Errorf("Finding ToID, From: %d To: %d Type %s", record.FromID, record.ToID, record.RelationshipType)
		}
		record.ToType = t.(string)
		if record.ToType == "" || record.FromType == "" {
			spew.Dump(record)
			return errors.Errorf("Missing type %d %d %s %s", record.FromID, record.ToID, record.FromType, record.ToType)
		}
		err = store.Upsert(key, record)
		if err != nil {
			return err
		}
		key = record.FromType + " -> " + record.ToType + " -> " + record.RelationshipType
		_, ok = relMap[key]
		if !ok {
			relMap[key] = 1
		} else {
			relMap[key]++
		}
	}
	foundRelationships := 0
	for k, v := range relMap {
		fmt.Printf("%s: %d\n", k, v)
		foundRelationships += v
	}
	if len(relationships) != foundRelationships {
		return errors.Errorf("Relationships mismatch, processed: %d paired: %d", len(relationships), foundRelationships)
	}
	return nil
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
