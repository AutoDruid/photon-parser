package types

type SyncHooks struct {
	OnSession   func(Session)
	OnCommand   func(Command)
	OnParameter func(Parameter)
}

type AsyncHooks struct {
	OnSession   chan Session
	OnCommand   chan Command
	OnParameter chan Parameter
}
