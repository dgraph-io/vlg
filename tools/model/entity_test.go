package model

import (
	"bytes"
	"testing"
	"time"
)

func TestEntity_ToRDF(t *testing.T) {

	date, _ := time.Parse("2006-01-02", "2020-01-01")
	type fields struct {
		Record            Record
		OriginalName      string
		FormerName        string
		Jurisdiction      string
		CompanyType       string
		Address           string
		IncorporationDate DateTime
		InactivationDate  DateTime
		StruckOffDate     DateTime
		DormDate          DateTime
		Status            string
		ServiceProvider   string
	}
	tests := []struct {
		name   string
		fields fields
		wantW  string
	}{
		{
			name: "valid entity",
			fields: fields{
				Record: Record{
					NodeID: "1",
					Name:   "Acme Inc.",
				},
				OriginalName:      "Acme Inc.",
				FormerName:        "Acme Inc.",
				Jurisdiction:      "US",
				CompanyType:       "LLC",
				Address:           "123 Main St, Anytown, USA",
				IncorporationDate: DateTime{date},
				Status:            "active",
				ServiceProvider:   "Acme SP Inc.",
			},
			wantW: `_:1 <dgraph.type> "Entity" .
_:1 <dgraph.type> "Record" .
_:1 <Record.nodeID> "1"^^<xsd:integer> .
_:1 <Record.sourceID> "0"^^<xsd:integer> .
_:1 <Record.name> "Acme Inc."^^<xsd:string> .
_:1 <Entity.originalName> "Acme Inc."^^<xsd:string> .
_:1 <Entity.formerName> "Acme Inc."^^<xsd:string> .
_:1 <Entity.jurisdiction> "US"^^<xsd:string> .
_:1 <Entity.companyType> "LLC"^^<xsd:string> .
_:1 <Entity.address> "123 Main St, Anytown, USA"^^<xsd:string> .
_:1 <Entity.incorporationDate> "2020-01-01 00:00:00 +0000 UTC"^^<xsd:string> .
_:1 <Entity.status> "active"^^<xsd:string> .
_:1 <Entity.serviceProvider> "Acme SP Inc."^^<xsd:string> .
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entity{
				Record:            tt.fields.Record,
				OriginalName:      tt.fields.OriginalName,
				FormerName:        tt.fields.FormerName,
				Jurisdiction:      tt.fields.Jurisdiction,
				CompanyType:       tt.fields.CompanyType,
				Address:           tt.fields.Address,
				IncorporationDate: tt.fields.IncorporationDate,
				InactivationDate:  tt.fields.InactivationDate,
				StruckOffDate:     tt.fields.StruckOffDate,
				DormDate:          tt.fields.DormDate,
				Status:            tt.fields.Status,
				ServiceProvider:   tt.fields.ServiceProvider,
			}
			w := &bytes.Buffer{}
			e.ToRDF(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Entity.ToRDF() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
