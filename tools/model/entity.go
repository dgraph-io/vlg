package model

import "fmt"

type Entity struct {
	Record

	OriginalName      string   `csv:"original_name"`
	FormerName        string   `csv:"former_name"`
	Jurisdiction      string   `csv:"jurisdiction"`
	CompanyType       string   `csv:"company_type"`
	Address           string   `csv:"address"`
	IncorporationDate DateTime `csv:"incorporation_date"`
	InactivationDate  DateTime `csv:"inactivation_date"`
	StruckOffDate     DateTime `csv:"struck_off_date"`
	DormDate          DateTime `csv:"dorm_date"`
	Status            string   `csv:"status"`
	ServiceProvider   string   `csv:"service_provider"`
}

func (e *Entity) String() string {
	return fmt.Sprintf("Entity %d '%s'", e.NodeID, e.Name)
}
