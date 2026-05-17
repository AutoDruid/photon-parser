// Package photon decodes Photon session envelopes and reliable payloads from raw bytes.
// It supports protocol versions 16 and 18, which differ in parameter layout and how
// reliable headers declare parameter counts.
//
// # Basic usage
//
//	parser := photon.NewParserV18()
//	var sess photon.SessionV18
//	if err := parser.ParsePacketInto(payload, &sess); err != nil {
//	    log.Fatal(err)
//	}
//	for _, cmd := range sess.Commands {
//	    fmt.Println(cmd.Type)
//	}
//
// # Performance
//
// ParsePacketInto is the zero-allocation hot path: it reuses internal pools and
// writes into the caller-supplied Session. The Session and its Commands slice are
// valid only until the next call to ParsePacketInto on the same Parser.
//
// Hooks (sync and async) are a convenience layer and carry additional cost:
//   - Sync hooks run inline during parsing; a slow callback blocks the parser.
//   - Async hooks copy the Commands slice before sending to avoid use-after-pool
//     corruption; prefer the direct Session() approach for high-throughput paths.
//
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

// Parser decodes Photon UDP payloads for a fixed protocol version.
// Create one with [NewParserV16] or [NewParserV18].
//
// A Parser is not safe for concurrent use. For parallel workloads, create one
// Parser per goroutine.
type Parser[P types.ParameterView] struct {
	ctx *context.Context[P]
}

// NewParserV16 returns a Parser that decodes parameters and reliable headers
// using Photon protocol version 16 rules.
//
// Pass functional [Option] values to override the default configuration, for
// example to skip unknown payloads or filter specific command types.
func NewParserV16(options ...Option) *Parser[v16.Parameter] {
	config := defaultConfig()
	for _, option := range options {
		option(&config)
	}

	return &Parser[v16.Parameter]{
		ctx: context.NewContext(
			reader.NewReader(nil),
			assembler.NewAssembler(),
			hooks.NewHooks[v16.Parameter](),
			context.Decoders[v16.Parameter]{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
			},
			config,
		),
	}
}

// ParsePacketV16 is a convenience wrapper that allocates a one-shot protocol 16
// Parser, parses data, and returns the resulting Session.
//
// For repeated parsing of multiple packets, prefer [NewParserV16] and reuse the
// Parser to avoid per-call allocation overhead.
func ParsePacketV16(data []byte, options ...Option) (*SessionV16, error) {
	p := NewParserV16(options...)
	var sess SessionV16
	err := p.ParsePacketInto(data, &sess)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// NewParserV18 returns a Parser that decodes parameters and reliable headers
// using Photon protocol version 18 rules.
//
// Pass functional [Option] values to override the default configuration, for
// example to skip unknown payloads or filter specific command types.
func NewParserV18(options ...Option) *Parser[v18.Parameter] {
	config := defaultConfig()
	for _, option := range options {
		option(&config)
	}
	return &Parser[v18.Parameter]{
		ctx: context.NewContext(
			reader.NewReader(nil),
			assembler.NewAssembler(),
			hooks.NewHooks[v18.Parameter](),
			context.Decoders[v18.Parameter]{
				ParameterParser:              &v18.Parameter{},
				ReliableHeaderParameterCount: &v18.ReliableHeaderParameterCountV18{},
			},
			config,
		),
	}
}

// ParsePacketV18 is a convenience wrapper that allocates a one-shot protocol 18
// Parser, parses data, and returns the resulting Session.
//
// For repeated parsing of multiple packets, prefer [NewParserV18] and reuse the
// Parser to avoid per-call allocation overhead.
func ParsePacketV18(data []byte, options ...Option) (*SessionV18, error) {
	p := NewParserV18(options...)
	var sess SessionV18
	err := p.ParsePacketInto(data, &sess)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// ParsePacketInto resets the internal reader to data and decodes one Photon
// session into sess.
//
// This is the zero-allocation hot path. The parsed Session and its Commands
// slice are valid only until the next call to ParsePacketInto on the same
// Parser; do not retain references across calls unless you copy the data you
// need.
//
// Any sync or async hooks registered on the Parser are invoked during this
// call.
func (p *Parser[P]) ParsePacketInto(data []byte, sess *Session[P]) error {

	p.ctx.Reader.Reset(data)

	err := session.ParseInto(p.ctx, sess)
	if err != nil {
		return err
	}

	return nil
}

// OnEventData registers a sync callback invoked for each reliable EventData
// message (server-raised events) decoded during [ParsePacketInto].
//
// The callback runs on the parser goroutine; a slow implementation will reduce
// parsing throughput. For non-blocking processing use [OnSessionAsync] or
// decouple work with your own channel.
func (p *Parser[P]) OnEventData(fn func(Reliable[P])) {
	p.ctx.Hooks.OnEvents[types.EventDataType] = fn
}

// OnOperationResponse registers a sync callback invoked for each reliable
// OperationResponse message decoded during [ParsePacketInto].
func (p *Parser[P]) OnOperationResponse(fn func(Reliable[P])) {
	p.ctx.Hooks.OnEvents[types.OperationResponse] = fn
}

// OnOperationRequest registers a sync callback invoked for each reliable
// OperationRequest message decoded during [ParsePacketInto].
func (p *Parser[P]) OnOperationRequest(fn func(Reliable[P])) {
	p.ctx.Hooks.OnEvents[types.OperationRequest] = fn
}

// OnOtherOperationResponse registers a sync callback invoked for each reliable
// OtherOperationResponse message decoded during [ParsePacketInto].
func (p *Parser[P]) OnOtherOperationResponse(fn func(Reliable[P])) {
	p.ctx.Hooks.OnEvents[types.OtherOperationResponse] = fn
}

// OnSessionSync registers a sync callback invoked once per [ParsePacketInto]
// call, after the full session (header + all commands) has been decoded.
//
// The Session passed to fn is valid only for the duration of the call; copy
// any fields you need to retain.
func (p *Parser[P]) OnSessionSync(fn func(Session[P])) {
	p.ctx.Hooks.SyncHooks.OnSession = fn
}

// OnCommandSync registers a sync callback invoked once per decoded command
// during [ParsePacketInto].
//
// The Command passed to fn is valid only for the duration of the call; copy
// any fields you need to retain.
func (p *Parser[P]) OnCommandSync(fn func(Command[P])) {
	p.ctx.Hooks.SyncHooks.OnCommand = fn
}

// OnParameterSync registers a sync callback invoked once per decoded parameter
// during [ParsePacketInto].
func (p *Parser[P]) OnParameterSync(fn func(P)) {
	p.ctx.Hooks.SyncHooks.OnParameter = fn
}

// OnSessionAsync returns a receive-only channel that receives a copy of each
// decoded Session after [ParsePacketInto] completes.
//
// The Commands slice is copied before sending to guarantee safe access after
// the parser has moved on; this allocation is the trade-off for non-blocking
// delivery. If the channel buffer is full the session is dropped silently —
// size the buffer (via HookOptions.Size) to match your consumer throughput.
//
// Call [Close] when you are done to drain and close the underlying channel.
func (p *Parser[P]) OnSessionAsync(options types.HookOptions) <-chan Session[P] {
	return p.ctx.Hooks.OnSessionAsync(options)
}

// OnCommandAsync returns a receive-only channel that receives each decoded
// Command during [ParsePacketInto].
//
// If the channel buffer is full the command is dropped silently.
// Call [Close] when you are done.
func (p *Parser[P]) OnCommandAsync(options types.HookOptions) <-chan Command[P] {
	return p.ctx.Hooks.OnCommandAsync(options)
}

// OnParameterAsync returns a receive-only channel that receives each decoded
// parameter during [ParsePacketInto].
//
// If the channel buffer is full the parameter is dropped silently.
// Call [Close] when you are done.
func (p *Parser[P]) OnParameterAsync(options types.HookOptions) <-chan P {
	return p.ctx.Hooks.OnParameterAsync(options)
}

// Close closes all async hook channels registered via [OnSessionAsync],
// [OnCommandAsync], and [OnParameterAsync], signalling consumers that no
// further values will be sent.
//
// After Close, new async hooks can be registered by calling the On*Async
// methods again.
func (p *Parser[P]) Close() {
	p.ctx.Hooks.CloseAsyncHooks()
}
