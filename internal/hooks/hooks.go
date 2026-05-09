package hooks

import "michelprogram/photon-parser/internal/types"

// Hooks contains synchronous and asynchronous parser hooks.
type Hooks[P types.ParameterView] struct {
	types.AsyncHooks[P]
	types.SyncHooks[P]
}

// NewHooks returns a Hooks value with all hook slots initialized to nil.
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

// OnSessionAsync returns the asynchronous session hook channel.
func (h *Hooks[P]) OnSessionAsync(options types.HookOptions) <-chan types.Session {
	return ensureChan(&h.AsyncHooks.OnSession, options.Size)
}

// OnCommandAsync returns the asynchronous command hook channel.
func (h *Hooks[P]) OnCommandAsync(options types.HookOptions) <-chan types.Command {
	return ensureChan(&h.AsyncHooks.OnCommand, options.Size)
}

// OnParameterAsync returns the asynchronous parameter hook channel.
func (h *Hooks[P]) OnParameterAsync(options types.HookOptions) <-chan P {
	return ensureChan(&h.AsyncHooks.OnParameter, options.Size)
}

// CloseAsyncHooks closes and resets all asynchronous hook channels.
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
