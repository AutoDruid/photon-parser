package photon

import (
	v16 "github.com/AutoDruid/photon-parser/internal/parameters/v16"
	v18 "github.com/AutoDruid/photon-parser/internal/parameters/v18"
	"github.com/AutoDruid/photon-parser/internal/types"
)

// Photon Protocol session types.
// These define the different kinds of session that can be exchanged.
type Session[P types.ParameterView] = types.Session[P]

// Photon Protocol command types.
// These define the various operations that can be performed in a Photon session.
type Command[P types.ParameterView] = types.Command[P]

// Photon Protocol reliable types.
// These define the different kinds of reliable messages that can be exchanged.
type Reliable[P types.ParameterView] = types.Reliable[P]

// Photon Protocol hook options types.
// These define the different kinds of hook options that can be exchanged.
type HookOptions = types.HookOptions

// Photon Protocol parameter types.
// These define the different kinds of parameters that can be exchanged.
type ParameterV16 = v16.Parameter

// Photon Protocol parameter types.
// These define the different kinds of parameters that can be exchanged.
type ParameterV18 = v18.Parameter

// Photon Protocol parameter types.
// These define the different kinds of parameters that can be exchanged.
type ParameterV18Type = v18.ParameterType

// Photon Protocol parameter types.
// These define the different kinds of parameters that can be exchanged.
type ParameterV16Type = v16.ParameterType

// Photon Protocol reliable types.
// These define the different kinds of reliable messages that can be exchanged.
type ReliableV18 = Reliable[ParameterV18]

// Photon Protocol reliable types.
// These define the different kinds of reliable messages that can be exchanged.
type ReliableV16 = Reliable[ParameterV16]

// Photon Protocol command types.
// These define the various operations that can be performed in a Photon session.
const (
	AcknowledgeCommand          types.CommandType = 0x01 // Acknowledges receipt of reliable commands
	ConnectCommand              types.CommandType = 0x02 // Initiates a connection
	VerifyConnectCommand        types.CommandType = 0x03 // Verifies connection establishment
	DisconnectCommand           types.CommandType = 0x04 // Gracefully closes a connection
	PingCommand                 types.CommandType = 0x05 // Keep-alive ping message
	SendReliableCommand         types.CommandType = 0x06 // Sends reliable data (guaranteed delivery)
	SendUnreliableCommand       types.CommandType = 0x07 // Sends unreliable data (best effort)
	SendReliableFragmentCommand types.CommandType = 0x08 // Sends a fragment of a large reliable message
)

// Photon Protocol reliable message types.
// These define the different kinds of reliable messages that can be exchanged.
const (
	OperationRequest       types.Type = 0x02 // Client requests an operation
	OperationResponse      types.Type = 0x07 // Server responds to an operation
	OtherOperationResponse types.Type = 0x03 // Alternative response format
	EventDataType          types.Type = 0x04 // Server sends an event to client
	ExchangeKeys           types.Type = 0x06 // Key exchange for encryption
)
