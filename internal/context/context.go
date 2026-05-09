package context

import (
	"AutoDruid/photon-parser/internal/assembler"
	"AutoDruid/photon-parser/internal/hooks"
	"AutoDruid/photon-parser/internal/reader"
	"AutoDruid/photon-parser/internal/types"
)

type Context[P types.ParameterView] struct {
	Reader    *reader.Reader
	Assembler *assembler.Assembler
	Hooks     *hooks.Hooks[P]
	Decoders  Decoders[P]
}

func NewContext[P types.ParameterView](reader *reader.Reader, assembler *assembler.Assembler, hooks *hooks.Hooks[P], decoders Decoders[P]) *Context[P] {
	return &Context[P]{
		Reader:    reader,
		Assembler: assembler,
		Hooks:     hooks,
		Decoders:  decoders,
	}
}
