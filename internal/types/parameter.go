package types

import "iter"

type VersionedParameter interface {
	ID() uint8
	Float32s() iter.Seq2[int, float32]
	Float32() float32
	MarshalJSON() ([]byte, error)
}
