package disconnect

// Disconnect represents the payload of a disconnect command.
type Disconnect struct{}

// Parse decodes a disconnect payload.
func Parse() *Disconnect {
	return &Disconnect{}
}
