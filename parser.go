package photon

import (
	"AutoDruid/photon-parser/internal/assembler"
	"AutoDruid/photon-parser/internal/context"
	"AutoDruid/photon-parser/internal/hooks"
	v16 "AutoDruid/photon-parser/internal/parameters/v16"
	v18 "AutoDruid/photon-parser/internal/parameters/v18"
	"AutoDruid/photon-parser/internal/reader"
	"AutoDruid/photon-parser/internal/session"
	"AutoDruid/photon-parser/internal/types"
)

type Parser[P types.ParameterView] struct {
	ctx *context.Context[P]
}

func NewV16() *Parser[v16.Parameter] {
	return &Parser[v16.Parameter]{
		ctx: context.NewContext(
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

func ParseV16(data []byte) (*Session, error) {
	p := NewV16()
	return p.ParsePacket(data)
}

func NewV18() *Parser[v18.Parameter] {
	return &Parser[v18.Parameter]{
		ctx: context.NewContext(
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

func ParseV18(data []byte) (*Session, error) {
	p := NewV18()
	return p.ParsePacket(data)
}

func (p *Parser[P]) ParsePacket(data []byte) (*Session, error) {

	p.ctx.Reader.Reset(data)

	var sess Session

	err := session.Parse(p.ctx, &sess)
	if err != nil {
		return nil, err
	}

	return &sess, nil
}

func (p *Parser[P]) OnSessionSync(fn func(Session)) {
	p.ctx.Hooks.SyncHooks.OnSession = fn
}

func (p *Parser[P]) OnCommandSync(fn func(Command)) {
	p.ctx.Hooks.SyncHooks.OnCommand = fn
}

func (p *Parser[P]) OnParameterSync(fn func(P)) {
	p.ctx.Hooks.SyncHooks.OnParameter = fn
}

func (p *Parser[P]) OnSessionAsync(options types.HookOptions) <-chan Session {
	return p.ctx.Hooks.OnSessionAsync(options)
}

func (p *Parser[P]) OnCommandAsync(options types.HookOptions) <-chan Command {
	return p.ctx.Hooks.OnCommandAsync(options)
}

func (p *Parser[P]) OnParameterAsync(options types.HookOptions) <-chan P {
	return p.ctx.Hooks.OnParameterAsync(options)
}

func (p *Parser[P]) Close() {
	p.ctx.Hooks.CloseAsyncHooks()
}
