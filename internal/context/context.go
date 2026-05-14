package context

import (
	"github.com/AutoDruid/photon-parser/internal/assembler"
	"github.com/AutoDruid/photon-parser/internal/hooks"
	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/types"
)

type Context[P types.ParameterView] struct {
	Reader        *reader.Reader
	Assembler     *assembler.Assembler
	Hooks         *hooks.Hooks[P]
	Decoders      Decoders[P]
	PoolParameter *Pool[P]
	PoolCommand   *Pool[types.Command[P]]
}

func NewContext[P types.ParameterView](reader *reader.Reader, assembler *assembler.Assembler, hooks *hooks.Hooks[P], decoders Decoders[P]) *Context[P] {
	return &Context[P]{
		Reader:        reader,
		Assembler:     assembler,
		Hooks:         hooks,
		Decoders:      decoders,
		PoolParameter: NewPool[P](500),
		PoolCommand:   NewPool[types.Command[P]](100),
	}
}
