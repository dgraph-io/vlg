package model

import (
	"fmt"
	"io"
	"strings"

	tstore "github.com/matthewmcneely/triplestore"
	"github.com/timshannon/badgerhold/v4"
)

type Intermediary struct {
	Record

	Status  string `csv:"status" predicate:"Intermediary.status"`
	Address string `csv:"address"`
}

var intermediaryTypeLookup = map[string]string{
	"":                                   "",
	"none":                               "",
	"active":                             "Active",
	"active legal":                       "ActiveLegal",
	"client in representative territory": "ClientInRepresentativeTerritory",
	"delinquent":                         "Delinquent",
	"inactive":                           "Inactive",
	"prospect":                           "Prospect",
	"suspended":                          "Suspended",
	"suspended legal":                    "SuspendedLegal",
	"unrecoverable accounts":             "UnrecoverableAccounts",
}

func (intermediary *Intermediary) Normalize() {
	intermediary.Record.Normalize()
	if intermediary.Name == "" {
		intermediary.Name = fmt.Sprintf("Unknown Intermediary %s", intermediary.NodeID)
	}
	intermediary.Status = intermediaryTypeLookup[strings.ToLower(intermediary.Status)]
}

func (intermediary *Intermediary) String() string {
	return fmt.Sprintf("Intermediary %s '%s'", intermediary.NodeID, intermediary.Name)
}

func (intermediary *Intermediary) ToRDF(w io.Writer) {
	id := intermediary.RDFID()

	fmt.Fprintf(w, "%s <dgraph.type> \"Intermediary\" .\n", id)
	intermediary.Record.ToRDF(w)
	RDFEncodeTriples(w, tstore.TriplesFromStruct(id, intermediary))
}

func (intermediary *Intermediary) ExportAll(w io.Writer, store *badgerhold.Store) error {
	q := &badgerhold.Query{}
	tx := store.Badger().NewTransaction(false)
	defer tx.Discard()
	err := store.TxForEach(tx, q, func(entry *Intermediary) error {
		entry.Normalize()
		entry.ToRDF(w)
		return nil
	})
	return err
}
