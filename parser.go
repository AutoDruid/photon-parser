package photon

import (
	"michelprogram/photon-parser/internal/assembler"
	"michelprogram/photon-parser/internal/context"
	"michelprogram/photon-parser/internal/hooks"
	v16 "michelprogram/photon-parser/internal/parameters/v16"
	v18 "michelprogram/photon-parser/internal/parameters/v18"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/session"
	"michelprogram/photon-parser/internal/types"
)

// Parser decodes Photon packets into sessions and emits parsing hooks.
type Parser[P types.ParameterView] struct {
	ctx *context.Context[P]
}

// NewV16 returns a parser configured for Photon protocol v16 parameters.
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

// ParseV16 parses a single Photon packet using the v16 parameter format.
func ParseV16(data []byte) (*Session, error) {
	p := NewV16()
	return p.ParsePacket(data)
}

// NewV18 returns a parser configured for Photon protocol v18 parameters.
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

// ParseV18 parses a single Photon packet using the v18 parameter format.
func ParseV18(data []byte) (*Session, error) {
	p := NewV18()
	return p.ParsePacket(data)
}

// ParsePacket parses one Photon packet and returns its decoded session.
func (p *Parser[P]) ParsePacket(data []byte) (*Session, error) {

	p.ctx.Reader.Reset(data)

	var sess Session

	err := session.Parse(p.ctx, &sess)
	if err != nil {
		return nil, err
	}

	return &sess, nil
}

// OnSessionSync registers a synchronous callback invoked for each parsed session.
func (p *Parser[P]) OnSessionSync(fn func(Session)) {
	p.ctx.Hooks.SyncHooks.OnSession = fn
}

// OnCommandSync registers a synchronous callback invoked for each parsed command.
func (p *Parser[P]) OnCommandSync(fn func(Command)) {
	p.ctx.Hooks.SyncHooks.OnCommand = fn
}

// OnParameterSync registers a synchronous callback invoked for each parsed parameter.
func (p *Parser[P]) OnParameterSync(fn func(P)) {
	p.ctx.Hooks.SyncHooks.OnParameter = fn
}

// OnSessionAsync returns a channel that receives parsed sessions asynchronously.
func (p *Parser[P]) OnSessionAsync(options types.HookOptions) <-chan Session {
	return p.ctx.Hooks.OnSessionAsync(options)
}

// OnCommandAsync returns a channel that receives parsed commands asynchronously.
func (p *Parser[P]) OnCommandAsync(options types.HookOptions) <-chan Command {
	return p.ctx.Hooks.OnCommandAsync(options)
}

// OnParameterAsync returns a channel that receives parsed parameters asynchronously.
func (p *Parser[P]) OnParameterAsync(options types.HookOptions) <-chan P {
	return p.ctx.Hooks.OnParameterAsync(options)
}

// Close closes all asynchronous hook channels owned by this parser.
func (p *Parser[P]) Close() {
	p.ctx.Hooks.CloseAsyncHooks()
}
