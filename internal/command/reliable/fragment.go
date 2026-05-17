package reliable

import (
	"github.com/AutoDruid/photon-parser/internal/context"
	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/types"
)

// ParseIntoFragment parses a Photon reliable fragment from a parser.Reader.
// It reads the fragment header, then feeds the fragment into the assembler.
// If the fragment is complete, it parses the reliable message into the destination.
// Returns an error if any part of parsing fails.
func ParseIntoFragment[P types.ParameterView](ctx *context.Context[P], length uint32, destFragment *types.Fragment, destReliable *types.Reliable[P]) error {

	err := readFragmentHeader(ctx.Reader, destFragment)
	if err != nil {
		return err
	}

	data, completed := ctx.Assembler.Feed(*destFragment)

	if completed {
		ctx.Reader.Reset(data)
		return Parse(ctx, destReliable, uint32(len(data)))
	}

	return nil
}

func readFragmentHeader(reader *reader.Reader, dest *types.Fragment) error {
	var err error

	dest.ID, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	dest.Count, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	dest.Index, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	dest.Size, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	dest.Offset, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	dest.Data = reader.ReadRemaining()

	return nil
}
