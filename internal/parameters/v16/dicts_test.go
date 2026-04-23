package v16_test

import (
	. "michelprogram/photon-parser/internal/parameters/v16"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
	"reflect"
	"testing"
)

func TestReadDictionnary(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    map[any]any
		wantErr bool
	}{
		{
			name: "int32 to int32 dictionary",
			input: []byte{
				0x00, 0x44,
				byte(types.Int32Type), byte(types.Int32Type), // keyType, valueType
				0x00, 0x02, // size = 2
				0x00, 0x00, 0x00, 0x01, // key = 1
				0x00, 0x00, 0x28, 0xC1, // value = 10433
				0x00, 0x00, 0x00, 0x02, // key = 2
				0x00, 0x00, 0x02, 0x21, // value = 545
			},
			want: map[any]any{
				int32(1): int32(10433),
				int32(2): int32(545),
			},
		},
		{
			name: "empty dictionary",
			input: []byte{
				0x00, 0x44,
				byte(types.Int32Type), byte(types.Int32Type),
				0x00, 0x00, // size = 0
			},
			want: map[any]any{},
		},
		{
			name: "string to boolean dictionary",
			input: []byte{
				0x00, 0x44,
				byte(types.StringType), byte(types.BooleanType),
				0x00, 0x02, // size = 2
				0x00, 0x01, 'a', // key = "a"
				0x01,                 // value = true
				0x00, 0x02, 'b', 'b', // key = "bb"
				0x00, // value = false
			},
			want: map[any]any{
				"a":  true,
				"bb": false,
			},
		},
		{
			name:    "truncated key type",
			input:   []byte{0x00, 0x44},
			wantErr: true,
		},
		{
			name:    "missing value type",
			input:   []byte{0x00, 0x44, byte(types.Int32Type)},
			wantErr: true,
		},
		{
			name: "truncated size",
			input: []byte{
				0x00, 0x44,
				byte(types.Int32Type), byte(types.Int32Type),
				0x00, // size needs 2 bytes
			},
			wantErr: true,
		},
		{
			name: "truncated data",
			input: []byte{
				0x00, 0x44,
				byte(types.Int32Type), byte(types.Int32Type),
				0x00, 0x01, // size = 1
				0x00, 0x00, 0x00, 0x01, // key = 1 (value missing)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &Parameter{},
				ReliableHeaderParameterCount: &ReliableHeaderParameterCountV16{},
			})
			p := &Parameter{}
			out := &types.Parameter{}
			err := p.Parse(reader, out, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDictionary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(out.Value, tt.want) {
				t.Errorf("ReadDictionary() = %v, want %v", out.Value, tt.want)
			}
		})
	}
}

func TestReadHashtable(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    map[any]any
		wantErr bool
	}{
		{
			name: "mixed key/value types",
			input: []byte{
				0x00, 0x68,
				0x00, 0x02, // size = 2
				byte(types.Int32Type),  // key type
				0x00, 0x00, 0x00, 0x01, // key = 1
				byte(types.StringType), // value type
				0x00, 0x02, 'h', 'i',   // value = "hi"
				byte(types.StringType), // key type
				0x00, 0x01, 'k',        // key = "k"
				byte(types.BooleanType), // value type
				0x01,                    // value = true
			},
			want: map[any]any{
				int32(1): "hi",
				"k":      true,
			},
		},
		{
			name:    "truncated size",
			input:   []byte{0x00, 0x68, 0x00},
			wantErr: true,
		},
		{
			name: "truncated data",
			input: []byte{
				0x00, 0x68,
				0x00, 0x01, // size = 1
				byte(types.Int32Type),
				0x00, 0x00, 0x00, 0x01, // key = 1 (value type/value missing)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &Parameter{},
				ReliableHeaderParameterCount: &ReliableHeaderParameterCountV16{},
			})
			p := &Parameter{}
			out := &types.Parameter{}
			err := p.Parse(reader, out, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadHashTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(out.Value, tt.want) {
				t.Errorf("ReadHashTable() = %v, want %v", out.Value, tt.want)
			}
		})
	}
}
