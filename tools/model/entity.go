package model

import (
	"fmt"
)

type Entity struct {
	Record

	OriginalName      string   `csv:"original_name" predicate:"Entity.originalName"`
	FormerName        string   `csv:"former_name" predicate:"Entity.formerName"`
	Jurisdiction      string   `csv:"jurisdiction" predicate:"Entity.jurisdiction"`
	CompanyType       string   `csv:"company_type" predicate:"Entity.companyType"`
	Address           string   `csv:"address" predicate:"Entity.address"`
	IncorporationDate DateTime `csv:"incorporation_date" predicate:"Entity.incorporationDate"`
	InactivationDate  DateTime `csv:"inactivation_date" predicate:"Entity.inactivationDate"`
	StruckOffDate     DateTime `csv:"struck_off_date" predicate:"Entity.struckOffDate"`
	DormDate          DateTime `csv:"dorm_date" predicate:"Entity.dormDate"`
	Status            string   `csv:"status" predicate:"Entity.status"`
	ServiceProvider   string   `csv:"service_provider" predicate:"Entity.serviceProvider"`
}

func (e *Entity) Normalize() {
	e.Record.Normalize()
	if e.Name == "" {
		e.Name = fmt.Sprintf("Unknown Entity %s", e.NodeID)
	}
	// Triple 'X' is the recognized value for an unknown jurisdiction (Wikipedia)
	if e.Jurisdiction == "XX" {
		e.Jurisdiction = "XXX"
	}
}

func (e *Entity) String() string {
	return fmt.Sprintf("Entity %s '%s'", e.NodeID, e.Name)
}
