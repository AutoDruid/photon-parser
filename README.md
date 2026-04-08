# Photon Parser

A Go library for parsing Photon Protocol network packets.

## Overview

This library provides a complete implementation of the Photon Protocol binary format parser, enabling developers to decode and analyze network traffic captured from Photon-based applications. The parser handles multiple protocol layers including session headers, commands, reliable messages, and Protocol16 parameter encoding.

## Features

- **Session Layer Parsing**: Decode Photon session headers with peer ID, timestamps, and CRC validation
- **Command Layer Parsing**: Handle multiple command types (Acknowledge, Connect, SendReliable, SendUnreliable, etc.)
- **Protocol16 Parameter Decoding**: Full support for Photon's Protocol16 binary format including:
  - Primitive types (int8, int16, int32, int64, float32, float64, string, boolean)
  - Array types (int8[], int32[], string[], generic arrays)
  - Dictionary types (nested key-value pairs)
- **Type-Safe Generic Readers**: Leverages Go 1.18+ generics for efficient, type-safe binary parsing
- **Comprehensive Test Coverage**: Extensive test suites for all parsing components

## Installation

```bash
go get michelprogram/photon-parser
```

## Quick Start

### Parsing a Complete Session

```go
package main

import (
    "encoding/hex"
    "fmt"
    "michelprogram/photon-parser/session"
    "strings"
)

func main() {
    // Example: Wireshark hex dump of Photon packet
    payload := "00:00:00:07:4e:71:2d:5f:40:41:ad:15:06:00:01:00..."
    
    // Convert hex string to bytes
    cleared, err := hex.DecodeString(strings.ReplaceAll(payload, ":", ""))
    if err != nil {
        panic(err)
    }
    
    // Parse the session
    session, err := session.Parse(cleared)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Peer ID: %d\n", session.PeerID)
    fmt.Printf("Commands: %d\n", len(session.Commands))
}
```

### Parsing Commands

```go
import (
    "michelprogram/photon-parser/command"
    "michelprogram/photon-parser/parser"
)

func parseCommand(data []byte) {
    reader := parser.NewReader(data)
    
    cmd, err := command.Parse(reader)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Command Type: %d\n", cmd.Type)
    fmt.Printf("Channel ID: %d\n", cmd.ChannelID)
    fmt.Printf("Reliable Seq: %d\n", cmd.ReliableSequenceNumber)
}
```

### Parsing Reliable Messages

```go
import (
    "michelprogram/photon-parser/reliable"
    "michelprogram/photon-parser/parser"
)

func parseReliable(data []byte) {
    reader := parser.NewReader(data)
    
    rel, err := reliable.Parse(reader)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Message Type: %d\n", rel.Type)
    fmt.Printf("Event Code: %d\n", rel.EventCode)
    fmt.Printf("Parameters: %d\n", len(rel.Parameters))
    
    // Access parameters
    for _, param := range rel.Parameters {
        fmt.Printf("  Param ID: %d, Value: %v\n", param.ID, param.Value)
    }
}
```

### Using Protocol16 Readers Directly

```go
import (
    "bytes"
    "michelprogram/photon-parser/parameters/readers"
    "michelprogram/photon-parser/parser"
)

func readParameters(data []byte) {
    reader := parser.NewReader(data)
    
    // Read an int32
    val, err := readers.ReadInt32(reader)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Int32 value: %d\n", val)
    
    // Read a string
    str, err := readers.ReadString(reader)
    if err != nil {
        panic(err)
    }
    fmt.Printf("String value: %s\n", str)
    
    // Read an array
    arr, err := readers.ReadInt32Array(reader)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Array values: %v\n", arr)
    
    // Read with automatic type detection
    value, err := readers.Decode(reader, readers.StringType)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Decoded value: %v\n", value)
}
```

## Protocol16 Type System

The library supports all Photon Protocol16 types:

| Type Code | Go Type | Description |
|-----------|---------|-------------|
| `0x62` | `int8` | 8-bit signed integer |
| `0x6b` | `int16` | 16-bit signed integer |
| `0x69` | `int32` | 32-bit signed integer |
| `0x6c` | `int64` | 64-bit signed integer |
| `0x66` | `float32` | 32-bit floating point |
| `0x67` | `float64` | 64-bit floating point |
| `0x73` | `string` | UTF-8 string (uint16 length prefix) |
| `0x6f` | `bool` | Boolean (0x00 or 0x01) |
| `0x62` | `[]int8` | Array of int8 (uint32 size prefix) |
| `0x6e` | `[]int32` | Array of int32 (uint32 size prefix) |
| `0x61` | `[]string` | Array of strings (uint32 size prefix) |
| `0x79` | `[]any` | Generic array (type byte + uint32 size) |
| `0x44` | `map[any]any` | Dictionary (nested parameters) |


## Testing

The project includes comprehensive test coverage:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./parameters/readers -v

# Run benchmarks
go test ./parameters/readers -bench=. -benchmem
```

## Development

### Requirements

- Go 1.25.0 or higher (uses generics)
- `golang.org/x/exp` for generic constraints

## Use Cases

- **Network Traffic Analysis**: Analyze video game sessions
- **Packet Inspection**: Debug multiplayer game communication
- **Protocol Research**: Understand Photon Protocol implementation details
- **Bot Development**: Parse game events for automation (educational purposes)
- **Security Research**: Analyze network protocol security

## License

See LICENSE file for details.

## Acknowledgments

This parser was developed for educational purposes to understand network protocols and game communication patterns. It is not affiliated with Exit Games (Photon Engine).
