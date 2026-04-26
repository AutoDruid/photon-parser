package v18

import "math"

func (p Parameter) StringValue() string {
	if p.Kind != StringType {
		return ""
	}
	return string(p.Blob)
}

func (p Parameter) Float32Value() float32 {
	if p.Kind != Float32Type {
		return 0
	}
	return math.Float32frombits(uint32(p.Num))
}

func (p Parameter) IntValue() int64 {
	switch p.Value.Kind {
	case Int8Type, Int8Positive, Int8Negative,
		Int16Type, Int16Positive, Int16Negative,
		CompressedInt32Type,
		Long8Positive, Long8Negative,
		Long16Positive, Long16Negative,
		CompressedInt64Type,
		IntZeroType, ShortZeroType, LongZeroType, ByteZeroType:
		return int64(p.Value.Num)
	default:
		return 0
	}
}
