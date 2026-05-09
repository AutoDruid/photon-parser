package v16

import "math"

// StringValue returns the decoded string when the parameter kind is StringType.
func (p Parameter) StringValue() (string, bool) {
	if p.Kind != StringType {
		return "", false
	}
	return string(p.Blob), true
}

// Float32Value returns the decoded float32 when the parameter kind is Float32Type.
func (p Parameter) Float32Value() (float32, bool) {
	if p.Kind != Float32Type {
		return 0, false
	}
	return math.Float32frombits(uint32(p.Num)), true
}

// Float64Value returns the decoded float64 when the parameter kind is Float64Type.
func (p Parameter) Float64Value() (float64, bool) {
	if p.Kind != Float64Type {
		return 0, false
	}
	return math.Float64frombits(uint64(p.Num)), true
}

// IntValue returns the decoded integer for supported integer kinds.
func (p Parameter) IntValue() (int64, bool) {
	switch p.Kind {
	case Int8Type:
		return int64(int8(uint8(p.Num))), true
	case Int16Type:
		return int64(int16(uint16(p.Num))), true
	case Int32Type:
		return int64(int32(uint32(p.Num))), true
	case Int64Type:
		return int64(p.Num), true
	default:
		return 0, false
	}
}

// BooleanValue returns the decoded boolean when the parameter kind is BooleanType.
func (p Parameter) BooleanValue() (bool, bool) {
	if p.Kind != BooleanType {
		return false, false
	}
	return p.Num == 1, true
}
