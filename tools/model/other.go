package model

import "fmt"

type Other struct {
	Record

	FormerName        string   `csv:"former_name"`
	Type              string   `csv:"type"`
	IncorporationDate DateTime `csv:"incorporation_date"`
	StruckOffDate     DateTime `csv:"struck_off_date"`
	ClosedDate        DateTime `csv:"closed_date"`
	Jurisdiction      string   `csv:"jurisdiction"`
}

func (other *Other) String() string {
	return fmt.Sprintf("Other %d '%s'", other.NodeID, other.Name)
}
