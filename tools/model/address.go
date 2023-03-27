package model

import (
	"fmt"
	"io"

	tstore "github.com/matthewmcneely/triplestore"
	"github.com/timshannon/badgerhold/v4"
)

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Address struct {
	Record

	Address   string    `csv:"address" predicate:"Address.address"`
	GeoSource string    `predicate:"Address.geoSource"`
	Location  *Location `json:"location"`
}

func (a *Address) Normalize() {
	a.Record.Normalize()
	if a.Name == "" {
		a.Name = a.Address
	}
}

func (a *Address) String() string {
	return fmt.Sprintf("Address %s '%s'", a.NodeID, a.Address)
}

func (a *Address) ToRDF(w io.Writer) {
	id := a.RDFID()

	fmt.Fprintf(w, "%s <dgraph.type> \"Address\" .\n", id)
	a.Record.ToRDF(w)
	RDFEncodeTriples(w, tstore.TriplesFromStruct(id, a))
	if a.Location != nil {
		geoJSON := fmt.Sprintf("\"{'type':'Point','coordinates':[%0.7f,%0.7f]}\"^^<geo:geojson>", a.Location.Long, a.Location.Lat)
		fmt.Fprintf(w, "%s <Address.location> %s .\n", id, geoJSON)
	}
}

func (address *Address) ExportAll(w io.Writer, store *badgerhold.Store) error {
	q := &badgerhold.Query{}
	tx := store.Badger().NewTransaction(false)
	defer tx.Discard()
	err := store.TxForEach(tx, q, func(entry *Address) error {
		entry.Normalize()
		entry.ToRDF(w)
		return nil
	})
	return err
}
