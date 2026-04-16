package types

// Type represents a Photon Protocol16 type code.
// Each type code indicates how the following bytes should be interpreted.
type ParameterType uint8

// Photon Protocol16 type codes.
// These constants define the binary type codes used in Photon's serialization format.
const (
	UnknownType           ParameterType = 0x00 // Unknown or unsupported type
	NilType               ParameterType = 0x2a // Null/nil value
	DictionaryType        ParameterType = 0x44 // Dictionary with fixed key/value types
	StringArrayType       ParameterType = 0x61 // Array of strings
	Int8Type              ParameterType = 0x62 // 8-bit signed integer
	Custom                ParameterType = 0x63 // Custom serialized object
	DoubleType            ParameterType = 0x64 // Alias for Float64Type
	EventDateType         ParameterType = 0x65 // Event date/time
	Float32Type           ParameterType = 0x66 // 32-bit floating point
	Float64Type           ParameterType = 0x67 // 64-bit floating point
	HashTableType         ParameterType = 0x68 // Hashtable with mixed key/value types
	Int32Type             ParameterType = 0x69 // 32-bit signed integer
	Int16Type             ParameterType = 0x6b // 16-bit signed integer
	Int64Type             ParameterType = 0x6c // 64-bit signed integer
	Int32ArrayType        ParameterType = 0x6e // Array of 32-bit integers
	BooleanType           ParameterType = 0x6f // Boolean (0x00=false, 0x01=true)
	OperationResponseType ParameterType = 0x70 // Operation response message
	OperationRequestType  ParameterType = 0x71 // Operation request message
	StringType            ParameterType = 0x73 // UTF-8 string with uint16 length prefix
	Int8ArrayType         ParameterType = 0x78 // Array of 8-bit integers
	ArrayType             ParameterType = 0x79 // Generic typed array
	ObjectArrayType       ParameterType = 0x7a // Array of serialized objects
)


// Header represents the parameter header containing the parameter ID and type.
// This appears at the beginning of each serialized parameter.
type ParameterHeader struct {
	ID   uint8 // Parameter identifier (application-specific)
	Type ParameterType  // Protocol16 type code indicating how to decode the value
}

// Parameters represents a complete Photon Protocol parameter with its header and decoded value.
// The Value field contains the decoded data according to the Type specified in the Header.
type Parameter struct {
	ParameterHeader

	Value interface{} // Decoded value, type depends on Header.Type
}
