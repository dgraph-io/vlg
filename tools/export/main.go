package main

import (
	"io"
	"log"
	"os"

	"github.com/timshannon/badgerhold/v4"

	"vlg/tools/model"
)

const exportDir = "../rdf"

// exports all records to RDF files.
func main() {

	store, err := model.GetStore(true)
	defer func() {
		_ = store.Close()
	}()
	if err != nil {
		panic(err)
	}

	m := map[string]func(io.Writer, *badgerhold.Store) error{
		"others.rdf":         new(model.Other).ExportAll,
		"entities.rdf":       new(model.Entity).ExportAll,
		"officers.rdf":       new(model.Officer).ExportAll,
		"intermediaries.rdf": new(model.Intermediary).ExportAll,
		"addresses.rdf":      new(model.Address).ExportAll,
	}

	for k, v := range m {
		log.Println("Exporting", k)
		f, err := os.Create(exportDir + "/" + k)
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = f.Close()
		}()
		err = v(f, store)
		if err != nil {
			panic(err)
		}
	}
}
