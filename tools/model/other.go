package model

import (
	"fmt"
	"strings"
)

type Other struct {
	Record

	FormerName        string   `csv:"former_name" predicate:"Other.formerName"`
	Type              string   `csv:"type" predicate:"Other.type"`
	IncorporationDate DateTime `csv:"incorporation_date" predicate:"Other.incorporationDate"`
	StruckOffDate     DateTime `csv:"struck_off_date" predicate:"Other.struckOffDate"`
	ClosedDate        DateTime `csv:"closed_date" predicate:"Other.closedDate"`
	Jurisdiction      string   `csv:"jurisdiction" predicate:"Other.jurisdiction"`
}

var otherTypeLookup = map[string]string{
	"":                          "",
	"limited liability company": "LLC",
	"sole ownership":            "SoleOwnership",
	"foreign formed":            "ForeignFormed",
}

func (other *Other) Normalize() {
	other.Record.Normalize()
	if other.Name == "" {
		other.Name = fmt.Sprintf("Unknown Other %s", other.NodeID)
	}
	other.Type = otherTypeLookup[strings.ToLower(other.Type)]
}

func (other *Other) String() string {
	return fmt.Sprintf("Other %s '%s'", other.NodeID, other.Name)
}
