package model

import (
	"fmt"
)

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Address struct {
	Record

	Address   string    `csv:"address" predicate:"Address.address"`
	GeoSource string    `predicate:"Address.geoSource"`
	Location  *Location `json:"location"`
}

func (a *Address) Normalize() {
	a.Record.Normalize()
	if a.Name == "" {
		a.Name = a.Address
	}
}

func (a *Address) String() string {
	return fmt.Sprintf("Address %s '%s'", a.NodeID, a.Address)
}
