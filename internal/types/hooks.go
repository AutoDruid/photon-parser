package types

// SyncHooks contains synchronous callbacks fired during parsing.
type SyncHooks[P ParameterView] struct {
	OnSession   func(Session)
	OnCommand   func(Command)
	OnParameter func(P)
}

// AsyncHooks contains asynchronous channels that receive parsed entities.
type AsyncHooks[P ParameterView] struct {
	OnSession   chan Session
	OnCommand   chan Command
	OnParameter chan P
}

// HookOptions configures asynchronous hook channels.
type HookOptions struct {
	Size uint16
}
