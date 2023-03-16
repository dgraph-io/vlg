package model

import "fmt"

type Intermediary struct {
	Record

	Status  string `csv:"status"`
	Address string `csv:"address"`
}

func (intermediary *Intermediary) String() string {
	return fmt.Sprintf("Intermediary %d '%s'", intermediary.NodeID, intermediary.Name)
}
