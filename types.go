package photon

import (
	v16 "michelprogram/photon-parser/internal/parameters/v16"
	v18 "michelprogram/photon-parser/internal/parameters/v18"
	"michelprogram/photon-parser/internal/types"
)

// Session is the parsed Photon session model exposed by this package.
type Session = types.Session

// Command is the parsed Photon command model exposed by this package.
type Command = types.Command

// HookOptions configures asynchronous parser hook channels.
type HookOptions = types.HookOptions

// ParameterV16 is the Photon protocol v16 parameter model.
type ParameterV16 = v16.Parameter

// ParameterV18 is the Photon protocol v18 parameter model.
type ParameterV18 = v18.Parameter

// ParameterV18Type is the Photon protocol v18 parameter type enum.
type ParameterV18Type = v18.ParameterType
