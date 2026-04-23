package context

import (
	"michelprogram/photon-parser/internal/assembler"
	"michelprogram/photon-parser/internal/hooks"
	"michelprogram/photon-parser/internal/reader"
)

type Context struct {
	Reader    *reader.Reader
	Assembler *assembler.Assembler
	Hooks     *hooks.Hooks
}
