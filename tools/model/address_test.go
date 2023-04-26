package model

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestAddress_ToRDF(t *testing.T) {
	type fields struct {
		Record    Record
		Address   string
		GeoSource string
		Location  *Location
	}
	tests := []struct {
		name   string
		fields fields
		wantW  string
	}{
		{
			name: "valid address",
			fields: fields{
				Record: Record{
					NodeID:   "2",
					Name:     "123 Main St, Anytown, USA",
					SourceID: "Paradise Papers - Malta corporate registry",
				},
				Address:   "123 Main St, Anytown, USA",
				GeoSource: "google",
				Location: &Location{
					Lat:  37.772318,
					Long: -122.4220186,
				},
			},
			wantW: `_:2 <dgraph.type> "Address" .
_:2 <dgraph.type> "Record" .
_:2 <Record.nodeID> "2"^^<xs:string> .
_:2 <Record.name> "123 Main St, Anytown, USA"^^<xs:string> .
_:2 <Record.sourceID> "ParadisePapers"^^<xs:string> .
_:2 <Address.geoSource> "google"^^<xs:string> .
_:2 <Address.location> "{'type':'Point','coordinates':[-122.4220186,37.7723180]}"^^<geo:geojson> .
`,
		},
		{
			name: "empty name",
			fields: fields{
				Record: Record{
					NodeID:   "2",
					SourceID: "Bahamas Leaks",
				},
				Address:   "123 Main St, Anytown, USA",
				GeoSource: "google",
				Location: &Location{
					Lat:  37.772318,
					Long: -122.4220186,
				},
			},
			wantW: `_:2 <dgraph.type> "Address" .
_:2 <dgraph.type> "Record" .
_:2 <Record.nodeID> "2"^^<xs:string> .
_:2 <Record.name> "123 Main St, Anytown, USA"^^<xs:string> .
_:2 <Record.sourceID> "BahamasLeaks"^^<xs:string> .
_:2 <Address.geoSource> "google"^^<xs:string> .
_:2 <Address.location> "{'type':'Point','coordinates':[-122.4220186,37.7723180]}"^^<geo:geojson> .
`,
		},
		{
			name: "empty name and address",
			fields: fields{
				Record: Record{
					NodeID:   "2",
					SourceID: "",
				},
				GeoSource: "google",
				Location: &Location{
					Lat:  37.772318,
					Long: -122.4220186,
				},
			},
			wantW: `_:2 <dgraph.type> "Address" .
_:2 <dgraph.type> "Record" .
_:2 <Record.nodeID> "2"^^<xs:string> .
_:2 <Record.name> "Unknown Address 2"^^<xs:string> .
_:2 <Record.sourceID> "OffshoreLeaks"^^<xs:string> .
_:2 <Address.geoSource> "google"^^<xs:string> .
_:2 <Address.location> "{'type':'Point','coordinates':[-122.4220186,37.7723180]}"^^<geo:geojson> .
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Address{
				Record:    tt.fields.Record,
				Address:   tt.fields.Address,
				GeoSource: tt.fields.GeoSource,
				Location:  tt.fields.Location,
			}
			w := &bytes.Buffer{}
			a.ToRDF(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Address.ToRDF() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_LoadGeoEncodedAddress(t *testing.T) {
	os.Setenv("VLG_US_CENSUS_ADDRESSES", "../addresses-geoencoded.csv")

	address := &Address{
		Record: Record{
			NodeID: "67299",
			Name:   "Test",
		},
	}

	stringWriter := &bytes.Buffer{}
	address.ToRDF(stringWriter)

	if !strings.Contains(stringWriter.String(), "-122.1402772,37.4267879") {
		t.Errorf("Address.ToRDF() = %v, want %v", stringWriter.String(), "-122.1402772,37.4267879")
	}
}
