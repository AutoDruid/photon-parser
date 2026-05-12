// Package photon decodes Photon session envelopes and reliable payloads from raw bytes.
// It supports protocol versions 16 and 18, which differ in parameter layout and how
// reliable headers declare parameter counts.
//
// Register hooks on a Parser before calling ParsePacket; synchronous and asynchronous
// callbacks run as the decoder walks the buffer.
package photon

import (
	"github.com/AutoDruid/photon-parser/internal/assembler"
	"github.com/AutoDruid/photon-parser/internal/context"
	"github.com/AutoDruid/photon-parser/internal/hooks"
	v16 "github.com/AutoDruid/photon-parser/internal/parameters/v16"
	v18 "github.com/AutoDruid/photon-parser/internal/parameters/v18"
	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/session"
	"github.com/AutoDruid/photon-parser/internal/types"
)

// Parser decodes Photon UDP payloads for a fixed protocol version (16 or 18).
// Create one with NewV16 or NewV18; it is safe to reuse for multiple ParsePacket calls.
type Parser[P types.ParameterView] struct {
	ctx *context.Context[P]
}

// NewV16 returns a Parser that interprets parameters and reliable headers using protocol 16 rules.
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

// ParseV16 parses data using a newly allocated protocol 16 Parser and returns the resulting Session.
func ParseV16(data []byte) (*Session, error) {
	p := NewV16()
		var sess Session
	err := p.ParsePacketInto(data, &sess)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// NewV18 returns a Parser that interprets parameters and reliable headers using protocol 18 rules.
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

// ParseV18 parses data using a newly allocated protocol 18 Parser and returns the resulting Session.
func ParseV18(data []byte) (*Session, error) {
	p := NewV18()
	var sess Session
	err := p.ParsePacketInto(data, &sess)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// ParsePacket resets the internal reader to data and decodes one Photon session.
// Hooks registered on the Parser are invoked during this call.
func (p *Parser[P]) ParsePacketInto(data []byte, sess *Session) error {

	p.ctx.Reader.Reset(data)

	err := session.Parse(p.ctx, sess)
	if err != nil {
		return err
	}

	return nil
}

// OnEventData registers a callback for reliable messages of type event data (server-raised events).
// The callback runs during ParsePacket when such a message is decoded.
func (p *Parser[P]) OnEventData(fn func(Reliable[P])) {
	p.ctx.Hooks.OnEvents[types.EventDataType] = fn
}

// OnOperationResponse registers a callback for reliable operation response messages.
func (p *Parser[P]) OnOperationResponse(fn func(Reliable[P])) {
	p.ctx.Hooks.OnEvents[types.OperationResponse] = fn
}

// OnOperationRequest registers a callback for reliable operation request messages.
func (p *Parser[P]) OnOperationRequest(fn func(Reliable[P])) {
	p.ctx.Hooks.OnEvents[types.OperationRequest] = fn
}

// OnOtherOperationResponse registers a callback for reliable other-operation-response messages.
func (p *Parser[P]) OnOtherOperationResponse(fn func(Reliable[P])) {
	p.ctx.Hooks.OnEvents[types.OtherOperationResponse] = fn
}

// OnSessionSync registers a function called once per ParsePacket when the session has been parsed.
func (p *Parser[P]) OnSessionSync(fn func(Session)) {
	p.ctx.Hooks.SyncHooks.OnSession = fn
}

// OnCommandSync registers a function called for each top-level command during ParsePacket.
func (p *Parser[P]) OnCommandSync(fn func(Command)) {
	p.ctx.Hooks.SyncHooks.OnCommand = fn
}

// OnParameterSync registers a function called for each decoded parameter during ParsePacket.
func (p *Parser[P]) OnParameterSync(fn func(P)) {
	p.ctx.Hooks.SyncHooks.OnParameter = fn
}

// OnSessionAsync returns a receive-only channel of parsed sessions.
// HookOptions.Size sets the channel buffer capacity; see Close when finished with async hooks.
func (p *Parser[P]) OnSessionAsync(options types.HookOptions) <-chan Session {
	return p.ctx.Hooks.OnSessionAsync(options)
}

// OnCommandAsync returns a receive-only channel that receives each command as it is parsed.
func (p *Parser[P]) OnCommandAsync(options types.HookOptions) <-chan Command {
	return p.ctx.Hooks.OnCommandAsync(options)
}

// OnParameterAsync returns a receive-only channel that receives each parameter as it is parsed.
func (p *Parser[P]) OnParameterAsync(options types.HookOptions) <-chan P {
	return p.ctx.Hooks.OnParameterAsync(options)
}

// Close shuts down asynchronous hook channels created by OnSessionAsync, OnCommandAsync, and OnParameterAsync.
func (p *Parser[P]) Close() {
	p.ctx.Hooks.CloseAsyncHooks()
}
