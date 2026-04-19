package v18

import (
	"michelprogram/photon-parser/internal/hooks"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
)

type Parameter struct{}

var _ reader.ParameterParser = (*Parameter)(nil)

func (p *Parameter) Parse(r *reader.Reader, out *types.Parameter, hooks *hooks.Hooks) error {
	return nil
}
