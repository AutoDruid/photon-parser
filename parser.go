package photonparser

import (
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/session"
	"michelprogram/photon-parser/internal/types"
)

type Parser struct {
	reader *reader.Reader
}

func NewParser() *Parser {
	return &Parser{
		reader: reader.NewReader(nil),
	}
}

func (p *Parser) ParsePacket(data []byte) (*Session, error) {

	p.reader.Reset(data)

	sess := session.Session{}
	err := sess.Parse(p.reader)
	if err != nil {
		return nil, err
	}

	return &sess.Session, nil
}

func (p *Parser) OnSessionSync(fn func(Session)) {
	p.reader.SyncHooks.OnSession = fn
}

func (p *Parser) OnCommandSync(fn func(Command)) {
	p.reader.SyncHooks.OnCommand = fn
}

func (p *Parser) OnParameterSync(fn func(Parameter)) {
	p.reader.SyncHooks.OnParameter = fn
}

func (p *Parser) OnSessionAsync(options types.HookOptions) <-chan Session {
	return p.reader.OnSessionAsync(options)
}

func (p *Parser) OnCommandAsync(options types.HookOptions) <-chan Command {
	return p.reader.OnCommandAsync(options)
}

func (p *Parser) OnParameterAsync(options types.HookOptions) <-chan Parameter {
	return p.reader.OnParameterAsync(options)
}

func (p *Parser) Close() {
	p.reader.CloseAsyncHooks()
}
