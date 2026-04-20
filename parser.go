package photonparser

import (
	"encoding/binary"
	"michelprogram/photon-parser/internal/hooks"
	v16 "michelprogram/photon-parser/internal/parameters/v16"
	v18 "michelprogram/photon-parser/internal/parameters/v18"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/session"
	"michelprogram/photon-parser/internal/types"
)

type Parser struct {
	reader *reader.Reader
	hooks  *hooks.Hooks
}

func NewParserV16() *Parser {
	return &Parser{
		reader: reader.NewReader(nil, reader.Options{
			ParameterParser: &v16.Parameter{},
			ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
			BinaryOrder: binary.BigEndian,
		}),
		hooks: hooks.NewHooks(),
	}
}

func NewParserV18() *Parser {
	return &Parser{
		reader: reader.NewReader(nil, reader.Options{
			ParameterParser: &v18.Parameter{},
			ReliableHeaderParameterCount: &v18.ReliableHeaderParameterCountV18{},
			BinaryOrder: binary.LittleEndian,
		}),
		hooks: hooks.NewHooks(),
	}
}

func (p *Parser) ParsePacket(data []byte) (*Session, error) {

	p.reader.Reset(data)

	sess, err := session.Parse(p.reader, p.hooks)
	if err != nil {
		return nil, err
	}

	return &sess.Session, nil
}

func (p *Parser) OnSessionSync(fn func(Session)) {
	p.hooks.SyncHooks.OnSession = fn
}

func (p *Parser) OnCommandSync(fn func(Command)) {
	p.hooks.SyncHooks.OnCommand = fn
}

func (p *Parser) OnParameterSync(fn func(Parameter)) {
	p.hooks.SyncHooks.OnParameter = fn
}

func (p *Parser) OnSessionAsync(options types.HookOptions) <-chan Session {
	return p.hooks.OnSessionAsync(options)
}

func (p *Parser) OnCommandAsync(options types.HookOptions) <-chan Command {
	return p.hooks.OnCommandAsync(options)
}

func (p *Parser) OnParameterAsync(options types.HookOptions) <-chan Parameter {
	return p.hooks.OnParameterAsync(options)
}

func (p *Parser) Close() {
	p.hooks.CloseAsyncHooks()
}
