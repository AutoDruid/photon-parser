package reliable

import "michelprogram/photon-parser/parameters"

const HEADER_SIZE = 5

type Type uint8

const (
	OperationRequest       Type = 0x02
	OperationResponse      Type = 0x07
	OtherOperationResponse Type = 0x03
	EventDataType          Type = 0x04
	ExchangeKeys           Type = 0x06
)

type Header struct {
	Signature      uint8
	Type           Type
	EventCode      uint8
	ParameterCount int16
}

type Reliable struct {
	Header
	Parameters []parameters.Parameters
}
