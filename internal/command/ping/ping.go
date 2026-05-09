package ping

// Ping represents the payload of a ping command.
type Ping struct{}

// Parse decodes a ping payload.
func Parse() *Ping {
	return &Ping{}
}
