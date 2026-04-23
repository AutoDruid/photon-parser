package v18

type ParameterType uint8

const (
	UnknownType         ParameterType = 0 //Tested
	BooleanType         ParameterType = 2
	Int8Type            ParameterType = 3
	Int16Type           ParameterType = 4 //Tested
	Float32Type         ParameterType = 5
	Float64Type         ParameterType = 6
	StringType          ParameterType = 7 //Tested
	NilType             ParameterType = 8
	CompressedInt32Type ParameterType = 9 //Tested
	CompressedInt64Type ParameterType = 10 //Tested

	Int8Positive  ParameterType = 11 // 1 byte unsigned, cast to +int32
	Int8Negative  ParameterType = 12 // 1 byte unsigned, cast to -int32
	Int16Positive ParameterType = 13 // 2 bytes unsigned, cast to +int32 Tested
	Int16Negative ParameterType = 14 // 2 bytes unsigned, cast to -int32

	Long8Positive  ParameterType = 15 // 1 byte unsigned, cast to +int64
	Long8Negative  ParameterType = 16 // 1 byte unsigned, cast to -int64
	Long16Positive ParameterType = 17 // 2 bytes unsigned, cast to +int64
	Long16Negative ParameterType = 18 // 2 bytes unsigned, cast to -int64

	CustomType     ParameterType = 19
	CustomTypeSlim ParameterType = 0x80

	// Complex types
	DictionaryType        ParameterType = 20
	HashtableType         ParameterType = 21
	ObjectArrayType       ParameterType = 23
	OperationRequestType  ParameterType = 24
	OperationResponseType ParameterType = 25
	EventDataType         ParameterType = 26

	// Zero shorthands — no payload bytes, type code is the entire value
	BooleanFalseType ParameterType = 27
	BooleanTrueType  ParameterType = 28
	ShortZeroType    ParameterType = 29
	IntZeroType      ParameterType = 30
	LongZeroType     ParameterType = 31
	FloatZeroType    ParameterType = 32
	DoubleZeroType   ParameterType = 33
	ByteZeroType     ParameterType = 34 //Tested

	// Array container — element count + element type follows
	ArrayType ParameterType = 0x40

	// Typed arrays — element type is baked into the type code (elemType | 0x40)
	BooleanArrayType        ParameterType = BooleanType | ArrayType         // 0x42
	ByteArrayType           ParameterType = Int8Type | ArrayType            // 0x43
	ShortArrayType          ParameterType = Int16Type | ArrayType           // 0x44 Tested
	Float32ArrayType        ParameterType = Float32Type | ArrayType         // 0x45 Tested
	Float64ArrayType        ParameterType = Float64Type | ArrayType         // 0x46
	StringArrayType         ParameterType = StringType | ArrayType          // 0x47
	CompressedIntArrayType  ParameterType = CompressedInt32Type | ArrayType // 0x49
	CompressedLongArrayType ParameterType = CompressedInt64Type | ArrayType // 0x4A
	CustomTypeArrayType     ParameterType = CustomType | ArrayType          // 0x53
	DictionaryArrayType     ParameterType = DictionaryType | ArrayType      // 0x54
	HashtableArrayType      ParameterType = HashtableType | ArrayType       // 0x55
)
