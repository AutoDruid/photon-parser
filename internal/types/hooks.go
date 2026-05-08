package types

type SyncHooks[P ParameterView] struct {
	OnSession   func(Session)
	OnCommand   func(Command)
	OnParameter func(P)
}

type AsyncHooks[P ParameterView] struct {
	OnSession   chan Session
	OnCommand   chan Command
	OnParameter chan P
}

type HookOptions struct {
	Size uint16
}

