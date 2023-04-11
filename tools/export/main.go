package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/timshannon/badgerhold/v4"

	"vlg/tools/model"
)

const exportDir = "../rdf"

// exports all records to RDF files.
func main() {

	now := time.Now()
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
			panic(errors.Wrapf(err, "creating file %s", k))
		}
		defer func() {
			if err = f.Close(); err != nil {
				log.Println(errors.Wrapf(err, "closing file %s", f.Name()))
			}
		}()
		err = v(f, store)
		if err != nil {
			panic(errors.Wrapf(err, "exporting %s", k))
		}
	}

	// load all relationships
	var relationships []*model.Relationship
	err = store.Find(&relationships, &badgerhold.Query{})
	if err != nil {
		panic(err)
	}
	log.Println("Exporting", len(relationships), "relationships")
	f, err := os.Create(exportDir + "/relationships.rdf")
	if err != nil {
		panic(errors.Wrap(err, "creating file relationships.rdf"))
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Println(errors.Wrapf(err, "closing file %s", f.Name()))
		}
	}()
	for i := range relationships {
		relationships[i].ToRDF(f)
	}

	log.Println("Exported all records in", time.Since(now))
}
