# photon-parser

## Download

**As a Go dependency** (recommended when you only need the library in your module):

```bash
go get github.com/AutoDruid/photon-parser@latest
```

**Clone the repository** (for contributing, running examples, or local development):

```bash
git clone https://github.com/AutoDruid/photon-parser.git
cd photon-parser
```

## Introduction

**photon-parser** is a Go library that decodes Photon session envelopes and command payloads from raw bytes—useful for inspecting captures, learning how the wire format is structured, and building tools around Photon traffic. This project is **explanatory and educational**: it documents and parses protocol data for understanding; it does **not** endorse any particular use, and **maintainers are not responsible** for how you apply it (including compliance, game or service terms, or any risk you take when handling real traffic). The library supports **protocol version 16** and **protocol version 18** (distinct parameter layouts and reliable-header parameter counting).

## How it works

Parsing is **sequential**: an internal `Reader` holds the input buffer and a **cursor** that advances as fields are read—there is no upfront full decode into a giant tree. **Lazy access** shows up in how you read values after parse: scalars are read when you call accessor methods; **arrays, dictionaries, and other composite parameters** are exposed as Go `iter` sequences (or similar pull-based access) so underlying bytes are processed as you iterate, instead of allocating everything up front.

## Examples

### `ParseV18`

```go
import "github.com/AutoDruid/photon-parser"

func exampleParseV18(payload []byte) error {
	sess, err := photon.ParseV18(payload)
	if err != nil {
		return err
	}
	_ = sess // use decoded session (commands, headers, …)
	return nil
}
```

### `ParsePacket`

Use a **reusable parser** when you want to attach hooks (below) or parse many packets with the same protocol version.

```go
import "github.com/AutoDruid/photon-parser"

func exampleParsePacket(payload []byte) error {
	p := photon.NewV18()
	sess, err := p.ParsePacket(payload)
	if err != nil {
		return err
	}
	_ = sess
	return nil
}
```

### `OnEventData`

Register a callback for **event data** reliable messages (the common path for server-raised events). It runs during `ParsePacket` when a matching reliable payload is decoded.

```go
import "github.com/AutoDruid/photon-parser"

func exampleOnEventData(payload []byte) error {
	p := photon.NewV18()
	p.OnEventData(func(msg photon.ReliableV18) {
		_ = msg.EventCode
		_ = msg.Parameters
	})
	_, err := p.ParsePacket(payload)
	return err
}
```

`OnOperationRequest`, `OnOperationResponse`, and `OnOtherOperationResponse` are the same pattern: register the hook, then call `ParsePacket`.

### `OnSessionSync`

```go
import "github.com/AutoDruid/photon-parser"

func exampleOnSessionSync(payload []byte) error {
	p := photon.NewV18()
	p.OnSessionSync(func(s photon.Session) {
		_ = s.CommandCount
	})
	_, err := p.ParsePacket(payload)
	return err
}
```

`OnCommandSync` and `OnParameterSync` follow the same pattern: register a `func(Command)` or `func(ParameterV16|ParameterV18)` before calling `ParsePacket`.

### `OnSessionAsync`

```go
import (
	"time"

	"github.com/AutoDruid/photon-parser"
)

func exampleOnSessionAsync(payload []byte) error {
	p := photon.NewV18()
	ch := p.OnSessionAsync(photon.HookOptions{Size: 8})
	defer p.Close()

	go func() {
		_, _ = p.ParsePacket(payload)
	}()

	select {
	case s := <-ch:
		_ = s.Timestamp
	case <-time.After(time.Second):
	}
	return nil
}
```

`OnCommandAsync` and `OnParameterAsync` are analogous: they return `<-chan Command` or `<-chan P` with the same `HookOptions` channel sizing.

## Local development

This repo pins Go and tooling in `.mise/config.toml`. With [mise](https://mise.jdx.dev/) installed:

1. Clone the repository and `cd` into it.
2. Run `mise install` in the project root so the pinned **Go** and **golangci-lint** versions are available.
3. Use the defined tasks (run `mise tasks` to list them)

From any subdirectory, mise resolves config upward; staying under the repo root ensures these tasks use the pinned toolchain.

## Contributing

Contributions are welcome. **Open an issue** to describe a bug, a protocol gap, or a feature you have in mind so maintainers can discuss scope and approach. When you are ready to implement a fix or enhancement, **resolve it with a pull request** that references the issue and keeps the change focused and easy to review.

