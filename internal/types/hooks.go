package types

type SyncHooks[P VersionedParameter] struct {
	OnSession   func(Session)
	OnCommand   func(Command)
	OnParameter func(P)
}

type AsyncHooks[P VersionedParameter] struct {
	OnSession   chan Session
	OnCommand   chan Command
	OnParameter chan P
}

type HookOptions struct {
	Size uint16
}

type Hookable interface {
	Session | Command
}
