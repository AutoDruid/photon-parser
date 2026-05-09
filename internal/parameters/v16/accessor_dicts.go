package v16

import (
	"iter"
	"michelprogram/photon-parser/internal/reader"
)

// DictionaryValue iterates over decoded key-value pairs for DictionaryType parameters.
func (p Parameter) DictionaryValue() iter.Seq2[any, any] {
	return func(yield func(any, any) bool) {
		if p.Kind != DictionaryType || p.Num == 0 || len(p.Blob) == 0 {
			return
		}
		r := reader.NewReader(p.Blob)
		for i := uint64(0); i < p.Num; i++ {
			k, err := scanPayload(r, p.KeyType)
			if err != nil {
				return
			}
			v, err := scanPayload(r, p.ValType)
			if err != nil {
				return
			}
			if !yield(decodeValue(k), decodeValue(v)) {
				return
			}
		}
	}
}
