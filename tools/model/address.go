package model

import "fmt"

type Location struct {
	Lat    float64 `json:"lat"`
	Long   float64 `json:"long"`
	Source string  `json:"source"`
}

type Address struct {
	Record

	Address string `csv:"address"`

	Location *Location `json:"location"`
}

func (a *Address) String() string {
	return fmt.Sprintf("Address %d '%s'", a.NodeID, a.Address)
}
