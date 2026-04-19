package hooks

import "michelprogram/photon-parser/internal/types"

type Hooks struct {
	types.AsyncHooks
	types.SyncHooks
}

func NewHooks() *Hooks {
	return &Hooks{
		AsyncHooks: types.AsyncHooks{
			OnSession:   nil,
			OnCommand:   nil,
			OnParameter: nil,
		},
		SyncHooks: types.SyncHooks{
			OnSession:   nil,
			OnCommand:   nil,
			OnParameter: nil,
		},
	}
}

func ensureChan[T types.Hookable](slot *chan T, size uint16) <-chan T {
	var minSize uint16 = 1

	if size != 0 {
		minSize = size
	}

	if *slot == nil {
		*slot = make(chan T, minSize)
	}
	return *slot
}
func (h *Hooks) OnSessionAsync(options types.HookOptions) <-chan types.Session {
	return ensureChan(&h.AsyncHooks.OnSession, options.Size)
}
func (h *Hooks) OnCommandAsync(options types.HookOptions) <-chan types.Command {
	return ensureChan(&h.AsyncHooks.OnCommand, options.Size)
}
func (h *Hooks) OnParameterAsync(options types.HookOptions) <-chan types.Parameter {
	return ensureChan(&h.AsyncHooks.OnParameter, options.Size)
}

func (h *Hooks) CloseAsyncHooks() {

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
