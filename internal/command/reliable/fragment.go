package reliable

import (
	"github.com/AutoDruid/photon-parser/internal/context"
	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/types"
)

func ParseFragment[P types.ParameterView](ctx *context.Context[P], out *types.Fragment, outt *types.Reliable[P], length uint32) error {

	err := parseMetadata(ctx.Reader, out)
	if err != nil {
		return err
	}

	data, completed := ctx.Assembler.Feed(*out)

	if completed {
		ctx.Reader.Reset(data)
		return Parse(ctx, outt, uint32(len(data)))
	}

	return nil
}

func parseMetadata(reader *reader.Reader, out *types.Fragment) error {
	var err error

	out.ID, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	out.Count, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	out.Index, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	out.Size, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	out.Offset, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	out.Data = reader.ReadRemaining()

	return nil
}
