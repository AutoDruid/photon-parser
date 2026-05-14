package photon

import (
	v16 "github.com/AutoDruid/photon-parser/internal/parameters/v16"
	v18 "github.com/AutoDruid/photon-parser/internal/parameters/v18"
	"github.com/AutoDruid/photon-parser/internal/types"
)

type Session[P types.ParameterView] = types.Session[P]
type Command[P types.ParameterView] = types.Command[P]
type HookOptions = types.HookOptions
type ParameterV16 = v16.Parameter
type ParameterV18 = v18.Parameter
type ParameterV18Type = v18.ParameterType
type Reliable[P types.ParameterView] = types.Reliable[P]
type ReliableV18 = Reliable[ParameterV18]
