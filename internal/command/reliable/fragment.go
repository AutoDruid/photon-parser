package reliable

import (
	"michelprogram/photon-parser/internal/context"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
)

type Fragment struct {
	types.Fragment
}

func ParseFragment(ctx *context.Context, length uint32) (*Reliable, error) {

	fragment, err := parseMetadata(ctx.Reader)
	if err != nil {
		return nil, err
	}

	data, completed := ctx.Assembler.Feed(fragment.Fragment)

	if completed {
		ctx.Reader.Reset(data)
		return Parse(ctx, length)
	}

	return nil, nil
}

func parseMetadata(reader *reader.Reader) (*Fragment, error) {
	var fragment Fragment
	var err error

	fragment.ID, err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}

	fragment.Count, err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}

	fragment.Index, err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}

	fragment.Size, err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}

	fragment.Offset, err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}

	fragment.Data = reader.ReadRest()

	return &fragment, nil

}
