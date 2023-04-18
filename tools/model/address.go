package model

import (
	"fmt"
	"io"
	"os"

	"github.com/gocarina/gocsv"
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
	Location  *Location `json:"Location"`
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
	if a.Location == nil {
		if censusAddress, ok := censusAddresses[a.NodeID]; ok {
			a.Location = &Location{
				Lat:  censusAddress.Lat,
				Long: censusAddress.Long,
			}
		}
	}
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

type USCensusAddress struct {
	NodeID string  `csv:"node_id"`
	Name   string  `csv:"name"`
	Lat    float64 `csv:"lat"`
	Long   float64 `csv:"long"`
}

var censusAddresses map[string]*USCensusAddress

// Load the census addresses CSV file and create a map of node IDs to addresses.
// Fails silently if the file is not found or address nodeID is not found.
// See [/notes/3. Data Importing.md](/notes/3.%20Data%20Importing.md) for more information.
func init() {
	censusAddresses = make(map[string]*USCensusAddress)

	f, err := os.OpenFile("addresses-geoencoded.csv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		wd, _ := os.Getwd()
		fmt.Println("No census addresses CSV found, pwd:", wd)
		return
	}

	addresses := []*USCensusAddress{}
	if err := gocsv.UnmarshalFile(f, &addresses); err != nil {
		fmt.Println("Error unmarshalling census addresses CSV", err)
		return
	}
	for _, address := range addresses {
		censusAddresses[address.NodeID] = address
	}
	f.Close()
}
