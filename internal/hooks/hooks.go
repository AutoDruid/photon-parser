package hooks

import "AutoDruid/photon-parser/internal/types"

type Hooks[P types.ParameterView] struct {
	types.AsyncHooks[P]
	types.SyncHooks[P]
}

func NewHooks[P types.ParameterView]() *Hooks[P] {
	return &Hooks[P]{
		AsyncHooks: types.AsyncHooks[P]{
			OnSession:   nil,
			OnCommand:   nil,
			OnParameter: nil,
		},
		SyncHooks: types.SyncHooks[P]{
			OnSession:   nil,
			OnCommand:   nil,
			OnParameter: nil,
		},
	}
}

func ensureChan[T any](slot *chan T, size uint16) <-chan T {
	var minSize uint16 = 1

	if size != 0 {
		minSize = size
	}

	if *slot == nil {
		*slot = make(chan T, minSize)
	}
	return *slot
}
func (h *Hooks[P]) OnSessionAsync(options types.HookOptions) <-chan types.Session {
	return ensureChan(&h.AsyncHooks.OnSession, options.Size)
}
func (h *Hooks[P]) OnCommandAsync(options types.HookOptions) <-chan types.Command {
	return ensureChan(&h.AsyncHooks.OnCommand, options.Size)
}
func (h *Hooks[P]) OnParameterAsync(options types.HookOptions) <-chan P {
	return ensureChan(&h.AsyncHooks.OnParameter, options.Size)
}

func (h *Hooks[P]) CloseAsyncHooks() {

	if h.AsyncHooks.OnSession != nil {
		close(h.AsyncHooks.OnSession)
		h.AsyncHooks.OnSession = nil
	}

	if h.AsyncHooks.OnCommand != nil {
		close(h.AsyncHooks.OnCommand)
		h.AsyncHooks.OnCommand = nil
	}

	if h.AsyncHooks.OnParameter != nil {
		close(h.AsyncHooks.OnParameter)
		h.AsyncHooks.OnParameter = nil
	}
}
