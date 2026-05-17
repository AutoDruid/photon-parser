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

**photon-parser** is a Go library that decodes Photon session envelopes and command payloads from raw bytes—useful for inspecting captures, learning how the wire format is structured, and building tools around Photon traffic. This project is **explanatory and educational**: it documents and parses protocol data for understanding; it does **not** endorse any particular use, and **maintainers are not responsible** for how you apply it (including compliance, game or service terms, or any risk you take when handling real traffic). The library supports **protocol version 16** and **protocol version 18** (distinct parameter layouts and reliable-header parameter counting). For more information look at the [documentation](https://autodruid.github.io/documentation/photon-parser/).

## How it works

Parsing is **sequential**: an internal `Reader` holds the input buffer and a **cursor** that advances as fields are read—there is no upfront full decode into a giant tree. **Lazy access** shows up in how you read values after parse: scalars are read when you call accessor methods; **arrays, dictionaries, and other composite parameters** are exposed as Go `iter` sequences (or similar pull-based access) so underlying bytes are processed as you iterate, instead of allocating everything up front.

## Local development

This repo pins Go and tooling in `.mise/config.toml`. With [mise](https://mise.jdx.dev/) installed:

1. Clone the repository and `cd` into it.
2. Run `mise install` in the project root so the pinned **Go** and **golangci-lint** versions are available.
3. Use the defined tasks (run `mise tasks` to list them)

From any subdirectory, mise resolves config upward; staying under the repo root ensures these tasks use the pinned toolchain.

## Contributing

Contributions are welcome. **Open an issue** to describe a bug, a protocol gap, or a feature you have in mind so maintainers can discuss scope and approach. When you are ready to implement a fix or enhancement, **resolve it with a pull request** that references the issue and keeps the change focused and easy to review.

