package photonparser

import (
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/session"
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

func (p *Parser) OnSessionAsync() <-chan Session {
	if p.reader.AsyncHooks.OnSession == nil {
		ch := make(chan Session, 1024)
		p.reader.AsyncHooks.OnSession = ch
	}
	return p.reader.AsyncHooks.OnSession
}

func (p *Parser) OnCommandSync(fn func(Command)) {
	p.reader.SyncHooks.OnCommand = fn
}

func (p *Parser) OnCommandAsync() <-chan Command {
	if p.reader.AsyncHooks.OnCommand == nil {
		ch := make(chan Command, 1024)
		p.reader.AsyncHooks.OnCommand = ch
	}
	return p.reader.AsyncHooks.OnCommand
}

func (p *Parser) Close() {
	p.reader.CloseAsyncHooks()
}
