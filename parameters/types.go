package parameters

import (
	"fmt"
	"michelprogram/photon-parser/parameters/readers"
)

type Header struct {
	ID   uint8
	Type readers.Type
}

type Parameters struct {
	Header

	Value interface{}
}

func (p Parameters) String() string {
	return fmt.Sprintf("ID: %d\nType: %d\nValue: %v\n", p.ID, p.Type, p.Value)
}
