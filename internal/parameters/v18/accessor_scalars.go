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
	case Int8Type:
		return int64(int8(uint8(p.Value.Num)))
	case Int8Positive:
		return int64(uint8(p.Value.Num))
	case Int8Negative:
		return int64(int32(p.Value.Num))
	case Int16Type:
		return int64(int16(uint16(p.Value.Num)))
	case Int16Positive:
		return int64(uint16(p.Value.Num))
	case Int16Negative, CompressedInt32Type:
		return int64(int32(uint32(p.Value.Num)))
	case Long8Positive:
		return int64(uint8(p.Value.Num))
	case Long8Negative:
		return int64(p.Value.Num)
	case Long16Positive:
		return int64(uint16(p.Value.Num))
	case Long16Negative, CompressedInt64Type:
		return int64(p.Value.Num)
	case IntZeroType, ShortZeroType, LongZeroType, ByteZeroType:
		return 0
	default:
		return 0
	}
}
