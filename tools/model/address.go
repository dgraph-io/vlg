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

	Address string `csv:"address" predicate:"Address.address"`

	// If true, the address is not resolved to a location
	Unresolved bool
	// Level of confidence in the geo-encoded location
	Confidence float64

	// The source of the geo-encoding
	GeoSource string    `predicate:"Address.geoSource"`
	Location  *Location `json:"location"`
}

func (a *Address) Normalize() {
	a.Record.Normalize()
	// if the address is empty or "None", then the name is the address (and the address is blank)
	if a.Name == "" || a.Name == "None" {
		a.Name = a.Address
		a.Address = ""
	}
	if a.Name == "None" || a.Name == "" {
		a.Unresolved = true
		a.Confidence = 0.0
		a.Name = fmt.Sprintf("Unknown Address %s", a.NodeID)
	}
	if a.Name == a.Address {
		a.Address = ""
	}
}

func (a *Address) String() string {
	return fmt.Sprintf("Address %s '%s'", a.NodeID, a.Address)
}

func (a *Address) ToRDF(w io.Writer) {
	id := a.RDFID()
	a.Normalize()

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
		entry.ToRDF(w)
		return nil
	})
	return err
}
