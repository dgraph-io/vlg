package model

import "fmt"

type Officer struct {
	Record
}

func (officer *Officer) String() string {
	return fmt.Sprintf("Officer %d '%s'", officer.NodeID, officer.Name)
}
