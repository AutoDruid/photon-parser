// Package readers provides Protocol16 type readers for Photon Protocol parameters.
// It handles all Photon data types including primitives, arrays, and dictionaries.
package parameters

// Type represents a Photon Protocol16 type code.
// Each type code indicates how the following bytes should be interpreted.
type Type uint8

// Photon Protocol16 type codes.
// These constants define the binary type codes used in Photon's serialization format.
const (
	UnknownType           Type = 0x00 // Unknown or unsupported type
	NilType               Type = 0x2a // Null/nil value
	DictionaryType        Type = 0x44 // Dictionary with fixed key/value types
	StringArrayType       Type = 0x61 // Array of strings
	Int8Type              Type = 0x62 // 8-bit signed integer
	Custom                Type = 0x63 // Custom serialized object
	DoubleType            Type = 0x64 // Alias for Float64Type
	EventDateType         Type = 0x65 // Event date/time
	Float32Type           Type = 0x66 // 32-bit floating point
	Float64Type           Type = 0x67 // 64-bit floating point
	HashTableType         Type = 0x68 // Hashtable with mixed key/value types
	Int32Type             Type = 0x69 // 32-bit signed integer
	Int16Type             Type = 0x6b // 16-bit signed integer
	Int64Type             Type = 0x6c // 64-bit signed integer
	Int32ArrayType        Type = 0x6e // Array of 32-bit integers
	BooleanType           Type = 0x6f // Boolean (0x00=false, 0x01=true)
	OperationResponseType Type = 0x70 // Operation response message
	OperationRequestType  Type = 0x71 // Operation request message
	StringType            Type = 0x73 // UTF-8 string with uint16 length prefix
	Int8ArrayType         Type = 0x78 // Array of 8-bit integers
	ArrayType             Type = 0x79 // Generic typed array
	ObjectArrayType       Type = 0x7a // Array of serialized objects
)
