package server

// TODO: Maybe load this from some sort of persistent storage medium, e.g YAML/TOML/JSON file

var handlers = make(map[byte]func(*Client, *Packet))

func (c *Client) HandlePacket(p *Packet) {
	handler, ok := handlers[p.Opcode]
	if !ok {
		LogDebug(0, "Unhandled Packet: {opcode:%d; length:%d};\n", p.Opcode, p.Length)
		return
	}
	handler(c, p)
}