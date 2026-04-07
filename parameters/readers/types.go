package readers

type Type uint8

const (
	UnknownType           Type = 0x00 // Done
	NilType               Type = 0x2a // Done
	DictionaryType        Type = 0x44 // Done
	StringArrayType       Type = 0x61 // Done
	Int8Type              Type = 0x62 //Done
	Custom                Type = 0x63
	DoubleType            Type = 0x64 //Done
	EventDateType         Type = 0x65
	Float32Type           Type = 0x66 //Done
	Float64Type           Type = 0x67 //Done
	HashTableType         Type = 0x68 //Done
	Int32Type             Type = 0x69 //Done
	Int16Type             Type = 0x6b //Done
	Int64Type             Type = 0x6c //Done
	Int32ArrayType        Type = 0x6e //Done
	BooleanType           Type = 0x6f //Done
	OperationResponseType Type = 0x70
	OperationRequestType  Type = 0x71
	StringType            Type = 0x73 //Done
	Int8ArrayType         Type = 0x78 //Done
	ArrayType             Type = 0x79 //Done
	ObjectArrayType       Type = 0x7a
)
