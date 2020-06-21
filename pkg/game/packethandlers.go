/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package game

import (
	"github.com/BurntSushi/toml"

	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/rand"
)

//HandlerFunc Represents a func that is to be called whenever a connected client receives
// a specific incoming handlers.
type HandlerFunc = func(*world.Player, *net.Packet)

//handlers A map with descriptive names for the keys, and functions to run for the value.
var Handlers = make(map[string]HandlerFunc)

//pDefinitions a collection of handlers pDefinitions.
var pDefinitions packetList

//packetDefinition Definition of a handlers handler.
type packetDefinition struct {
	Opcode int    `toml:"opcode"`
	Name   string `toml:"name"`
	//	Handler HandlerFunc
}

//packetList Represents a mapping of descriptive names to handlers opcodes.
type packetList struct {
	Set []packetDefinition `toml:"packets"`
}

func init() {
	// Just to prevent non-handled handlers message from spamming up the logs
	AddHandler("pingreq", func(*world.Player, *net.Packet) {})
	AddHandler("sessionreq", func(player *world.Player, p *net.Packet) {
		// TODO: Remove maybe...TLS deprecates the need for it
		player.SetConnected(true)
		p.ReadUint8() // UID, useful?
		player.SetServerSeed(rand.Rng.Uint64())
		player.SendPacket(net.NewReplyPacket(nil).AddUint64(player.ServerSeed()))
	})
}

//UnmarshalPackets Loads the handlers pDefinitions into memory from the configured TOML file
func UnmarshalPackets() {
	if _, err := toml.DecodeFile(config.PacketHandlers(), &pDefinitions); err != nil {
		log.Error.Fatalln("Could not open handlers handler pDefinitions data file:", err)
		return
	}
}

//Handler Returns the handlers handler function assigned to this opcode.  If it can't be found, returns nil.
func Handler(opcode byte) HandlerFunc {
	for _, h := range pDefinitions.Set {
		if byte(h.Opcode) == opcode {
			return Handlers[h.Name]
		}
	}
	return nil
}

//AddHandler Adds and assigns the packethandler to the handlers with the specified name.
func AddHandler(name string, h HandlerFunc) {
	if _, ok := Handlers[name]; ok {
		log.Warning.Printf("Attempted to bind a handler to handlers '%v' which is already handled elsewhere.  Ignoring bind.\n", name)
		return
	}
	Handlers[name] = h
}

//PacketCount returns the number of handlers pDefinitions
func PacketCount() int {
	return len(pDefinitions.Set)
}

//HandlerCount returns the number of pDefinitions that are handled
func HandlerCount() int {
	return len(Handlers)
}
