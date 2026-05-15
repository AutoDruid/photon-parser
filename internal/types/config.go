package types

type Config struct {
	SkipUnknownPayloads  bool                 `json:"skip_unknown_payloads"`
	SkipParameterParsing bool                 `json:"skip_parameter_parsing"`
	SkipCommands         map[CommandType]bool `json:"skip_commands"`
	SkipTargetEventCodes map[MessageType]bool `json:"skip_target_event_codes"`
}
