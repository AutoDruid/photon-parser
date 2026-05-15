package photon

import "github.com/AutoDruid/photon-parser/internal/types"

type Option func(*types.Config)

func DefaultConfig() types.Config {

	skipCommands := map[types.CommandType]bool{
		types.SendReliableCommand:         false,
		types.SendUnreliableCommand:       false,
		types.SendReliableFragmentCommand: false,
		types.AcknowledgeCommand:          false,
		types.ConnectCommand:              false,
		types.VerifyConnectCommand:        false,
		types.PingCommand:                 false,
		types.DisconnectCommand:           false,
	}

	skipTargetEventCodes := map[types.MessageType]bool{
		types.OperationRequest:       false,
		types.OperationResponse:      false,
		types.OtherOperationResponse: false,
		types.EventDataType:          false,
		types.ExchangeKeys:           false,
	}

	return types.Config{
		SkipUnknownPayloads:  false,
		SkipParameterParsing: false,
		SkipCommands:         skipCommands,
		SkipTargetEventCodes: skipTargetEventCodes,
	}
}

// SkipUnknownPayloads allows to skip the parsing of the unknown payloads.
// This option allows to skip these allocations and continue parsing the packet.
func SkipUnknownPayloads(skip bool) Option {
	return func(c *types.Config) {
		c.SkipUnknownPayloads = skip
	}
}

// SkipParameterParsing allows to skip the parsing of the parameters.
// This option is useful when you are only interested in the commands and the events and explore the requests.
func SkipParameterParsing(skip bool) Option {
	return func(c *types.Config) {
		c.SkipParameterParsing = skip
	}
}

// SkipCommands allows to skip the parsing of recognized commands.
// This option is useful when you are only interested in some specific commands.
func SkipCommands(commands ...types.CommandType) Option {
	return func(c *types.Config) {
		if c.SkipCommands == nil {
			c.SkipCommands = make(map[types.CommandType]bool)
		}
		for _, t := range commands {
			c.SkipCommands[t] = true
		}
	}
}

// SkipTargetEventCodes allows to skip the parsing of the events that are not in the list.
// This option is useful when you are only interested in some specific events.
func SkipTargetEventCodes(codes ...types.MessageType) Option {
	return func(c *types.Config) {
		if c.SkipTargetEventCodes == nil {
			c.SkipTargetEventCodes = make(map[types.MessageType]bool)
		}
		for _, code := range codes {
			c.SkipTargetEventCodes[code] = true
		}
	}
}
