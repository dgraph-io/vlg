package model

import (
	"fmt"
	"io"

	tstore "github.com/matthewmcneely/triplestore"
	"github.com/timshannon/badgerhold/v4"
)

type Officer struct {
	Record
}

func (officer *Officer) Normalize() {
	officer.Record.Normalize()
	if officer.Name == "" {
		officer.Name = fmt.Sprintf("Unknown Officer %s", officer.NodeID)
	}
}

func (officer *Officer) String() string {
	return fmt.Sprintf("Officer %s '%s'", officer.NodeID, officer.Name)
}

func (officer *Officer) ToRDF(w io.Writer) {
	id := officer.RDFID()

	fmt.Fprintf(w, "%s <dgraph.type> \"Officer\" .\n", id)
	officer.Record.ToRDF(w)
	RDFEncodeTriples(w, tstore.TriplesFromStruct(id, officer))
}

func (officer *Officer) ExportAll(w io.Writer, store *badgerhold.Store) error {
	q := &badgerhold.Query{}
	tx := store.Badger().NewTransaction(false)
	defer tx.Discard()
	err := store.TxForEach(tx, q, func(entry *Officer) error {
		entry.Normalize()
		entry.ToRDF(w)
		return nil
	})
	return err
}
