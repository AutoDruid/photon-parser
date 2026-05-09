package v18

import (
	"AutoDruid/photon-parser/internal/reader"
	"iter"
)

func (p Parameter) DictionaryValue() iter.Seq2[any, any] {
	if p.Kind != DictionaryType || p.Num == 0 || len(p.Blob) == 0 {
		return nil
	}
	return func(yield func(any, any) bool) {
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
