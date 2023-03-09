package main

import (
	"bytes"
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

	// retrieve the data
	record, err := model.RecordByID(store, int64(240491946))
	if err != nil {
		panic(err)
	}
	switch typed := record.(type) {
	case model.Entity:
		spew.Dump(typed)
	case model.Other:
		spew.Dump(typed)
	case model.Address:
		spew.Dump(typed)
	default:
		panic(errors.Errorf("Unknown type %T", typed))
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

func mainOLD() {
	b := bytes.NewBufferString(`node_id,name,type,incorporation_date,struck_off_date,closed_date,jurisdiction,jurisdiction_description,countries,country_codes,sourceID,valid_until,note
85004929,ANTAM ENTERPRISES N.V.,LIMITED LIABILITY COMPANY,18-MAY-1983,,28-NOV-2012,AW,Aruba,,,Paradise Papers - Aruba corporate registry,Aruba corporate registry data is current through 2016,Closed date stands for Cancelled date.
80004686,AAK Company Ltd.,,,,,,,Bermuda;Isle of Man,BMU;IMN,Paradise Papers - Appleby,Appleby data is current through 2014,
`)

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		//return csv.NewReader(in)
		return gocsv.LazyCSVReader(in) // Allows use of quotes in CSV
	})

	others := []*model.Other{}
	if err := gocsv.Unmarshal(b, &others); err != nil {
		panic(err)
	}

	for i := 0; i < len(others); i++ {
		spew.Dump(others[i])
	}
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
