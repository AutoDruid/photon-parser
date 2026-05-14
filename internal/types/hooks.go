package types

type SyncHooks[P ParameterView] struct {
	OnSession   func(Session[P])
	OnCommand   func(Command[P])
	OnEvents    map[Type]func(Reliable[P])
	OnParameter func(P)
}

type AsyncHooks[P ParameterView] struct {
	OnSession   chan Session[P]
	OnCommand   chan Command[P]
	OnParameter chan P
}

type HookOptions struct {
	Size uint16
}
