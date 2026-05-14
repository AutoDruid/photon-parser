package v18

import (
	"iter"

	"github.com/AutoDruid/photon-parser/internal/reader"
)

func (p Parameter) DictionaryValue() iter.Seq2[any, any] {
	if p.Kind != DictionaryType || p.Num == 0 || len(p.Blob) == 0 {
		return nil
	}
	return func(yield func(any, any) bool) {
		r := reader.NewReader(p.Blob)
		var key, value Value
		for i := uint64(0); i < p.Num; i++ {
			key.Kind = p.KeyType
			err := scanPayload(r, &key)
			if err != nil {
				return
			}

			value.Kind = p.ValType
			err = scanPayload(r, &value)
			if err != nil {
				return
			}
			if !yield(decodeValue(key), decodeValue(value)) {
				return
			}
		}
	}
}
