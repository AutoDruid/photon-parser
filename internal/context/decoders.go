package context

import (
	"michelprogram/photon-parser/internal/hooks"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
)

// ParameterParser is implemented by each protocol-version parameters package
// (v16, v18).
// It is wired once at Parser construction so the hot path
type ParameterParser[P types.ParameterView] interface {
	Parse(*reader.Reader, *P, *hooks.Hooks[P]) error
}

// ReliableHeaderParameterCount is implemented by each protocol-version reliable header parameter count package
// (v16, v18).
// It is used to count the number of parameters in a reliable header.
type ReliableHeaderParameterCount interface {
	Count(*reader.Reader) (int, error)
}

type Decoders[P types.ParameterView] struct {
	ParameterParser              ParameterParser[P]
	ReliableHeaderParameterCount ReliableHeaderParameterCount
}
