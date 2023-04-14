package model

import (
	"fmt"
)

type Officer struct {
	Record
}

func (officer *Officer) Normalize() {
	officer.Record.Normalize()
	if officer.Name == "" {
		officer.Name = fmt.Sprintf("Unknown Officer %s", officer.NodeID)
	}
}

func (officer *Officer) String() string {
	return fmt.Sprintf("Officer %s '%s'", officer.NodeID, officer.Name)
}
