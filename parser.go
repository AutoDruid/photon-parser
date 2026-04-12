package photonparser

import "michelprogram/photon-parser/internal/session"

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParsePacket(data []byte) (*session.Session, error) {

	reader := NewReader(data)

	sess := session.Session{}
	err := sess.Parse(reader)
	if err != nil {
		return nil, err
	}

	return &sess, nil
}

/* type EventHandler func(Event)
type RequestHandler func(Request)
type ResponseHandler func(Response)
func (p *Parser) OnEvent(handler EventHandler)
func (p *Parser) OnRequest(handler RequestHandler)
func (p *Parser) OnResponse(handler ResponseHandler)
func (p *Parser) RegisterCustomType(code byte, decoder CustomDecoder) */
