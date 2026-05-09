package context

import (
	"michelprogram/photon-parser/internal/assembler"
	"michelprogram/photon-parser/internal/hooks"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
)

// Context bundles parser state and protocol-specific decoders for one parse flow.
type Context[P types.ParameterView] struct {
	Reader    *reader.Reader
	Assembler *assembler.Assembler
	Hooks     *hooks.Hooks[P]
	Decoders  Decoders[P]
}

// NewContext creates a parsing context with shared reader, hooks, and decoders.
func NewContext[P types.ParameterView](reader *reader.Reader, assembler *assembler.Assembler, hooks *hooks.Hooks[P], decoders Decoders[P]) *Context[P] {
	return &Context[P]{
		Reader:    reader,
		Assembler: assembler,
		Hooks:     hooks,
		Decoders:  decoders,
	}
}
