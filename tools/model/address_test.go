package model

import (
	"bytes"
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
					NodeID: "2",
					Name:   "123 Main St, Anytown, USA",
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
_:2 <Record.nodeID> "2"^^<xsd:integer> .
_:2 <Record.sourceID> "0"^^<xsd:integer> .
_:2 <Record.name> "123 Main St, Anytown, USA"^^<xsd:string> .
_:2 <Address.address> "123 Main St, Anytown, USA"^^<xsd:string> .
_:2 <Address.geoSource> "google"^^<xsd:string> .
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
