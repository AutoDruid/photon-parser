package photon

import (
	"michelprogram/photon-parser/internal/assembler"
	"michelprogram/photon-parser/internal/command/reliable"
	"michelprogram/photon-parser/internal/context"
	"michelprogram/photon-parser/internal/hooks"
	v16 "michelprogram/photon-parser/internal/parameters/v16"
	v18 "michelprogram/photon-parser/internal/parameters/v18"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/session"
	"michelprogram/photon-parser/internal/types"
)

type Parser[P types.VersionedParameter] struct {
	Ctx *context.Context[P]
}

func NewParserV16() *Parser[v16.Parameter] {
	return &Parser[v16.Parameter]{
		Ctx: context.NewContext[v16.Parameter](
			reader.NewReader(nil),
			assembler.NewAssembler(),
			hooks.NewHooks[v16.Parameter](),
			context.Decoders[v16.Parameter]{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
			},
		),
	}
}

func NewParserV18() *Parser[v18.Parameter] {
	return &Parser[v18.Parameter]{
		Ctx: context.NewContext(
			reader.NewReader(nil),
			assembler.NewAssembler(),
			hooks.NewHooks[v18.Parameter](),
			context.Decoders[v18.Parameter]{
				ParameterParser:              &v18.Parameter{},
				ReliableHeaderParameterCount: &v18.ReliableHeaderParameterCountV18{},
			},
		),
	}
}

func (p *Parser[P]) ParsePacket(data []byte) (*Session, error) {

	p.Ctx.Reader.Reset(data)

	var sess Session

	err := session.Parse(p.Ctx, &sess)
	if err != nil {
		return nil, err
	}

	return &sess, nil
}

func (p *Parser[P]) OnSessionSync(fn func(Session)) {
	p.Ctx.Hooks.SyncHooks.OnSession = fn
}

func (p *Parser[P]) OnCommandSync(fn func(Command)) {
	p.Ctx.Hooks.SyncHooks.OnCommand = fn
}

func (p *Parser[P]) OnParameterSync(fn func(P)) {
	p.Ctx.Hooks.SyncHooks.OnParameter = fn
}

func (p *Parser[P]) OnSessionAsync(options types.HookOptions) <-chan Session {
	return p.Ctx.Hooks.OnSessionAsync(options)
}

func (p *Parser[P]) OnCommandAsync(options types.HookOptions) <-chan Command {
	return p.Ctx.Hooks.OnCommandAsync(options)
}

func (p *Parser[P]) OnParameterAsync(options types.HookOptions) <-chan P {
	return p.Ctx.Hooks.OnParameterAsync(options)
}

func (p *Parser[P]) Close() {
	p.Ctx.Hooks.CloseAsyncHooks()
}

func (p *Parser[P]) Release(s *Session) {
	if s == nil {
		return
	}
	for _, cmd := range s.Commands {
		if rel, ok := cmd.Payload.(*reliable.Reliable[P]); ok {
			rel.Release()
		}
	}
}
