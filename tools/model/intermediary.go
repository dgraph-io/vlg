package model

import (
	"fmt"
	"strings"
)

type Intermediary struct {
	Record

	Status  string `csv:"status" predicate:"Intermediary.status"`
	Address string `csv:"address"`
}

var intermediaryTypeLookup = map[string]string{
	"":                                   "",
	"none":                               "",
	"active":                             "Active",
	"active legal":                       "ActiveLegal",
	"client in representative territory": "ClientInRepresentativeTerritory",
	"delinquent":                         "Delinquent",
	"inactive":                           "Inactive",
	"prospect":                           "Prospect",
	"suspended":                          "Suspended",
	"suspended legal":                    "SuspendedLegal",
	"unrecoverable accounts":             "UnrecoverableAccounts",
}

func (intermediary *Intermediary) Normalize() {
	intermediary.Record.Normalize()
	if intermediary.Name == "" {
		intermediary.Name = fmt.Sprintf("Unknown Intermediary %s", intermediary.NodeID)
	}
	intermediary.Status = intermediaryTypeLookup[strings.ToLower(intermediary.Status)]
}

func (intermediary *Intermediary) String() string {
	return fmt.Sprintf("Intermediary %s '%s'", intermediary.NodeID, intermediary.Name)
}
