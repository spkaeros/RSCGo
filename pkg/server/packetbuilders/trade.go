package packetbuilders

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
)

//TradeClose Closes a trade window
var TradeClose = NewOutgoingPacket(128)

//TradeUpdate Builds a packet to update a trade offer
func TradeUpdate(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(97)
	player.TradeOffer.Lock.RLock()
	p.AddByte(uint8(len(player.TradeOffer.List)))
	for _, item := range player.TradeOffer.List {
		p.AddShort(uint16(item.ID))
		p.AddInt(uint32(item.Amount))
	}
	player.TradeOffer.Lock.RUnlock()
	return
}

//TradeOpen Builds a packet to open a trade window
func TradeOpen(player *world.Player) *Packet {
	return NewOutgoingPacket(92).AddShort(uint16(player.TradeTarget()))
}

//TradeTargetAccept Builds a packet to change trade targets accepted status
func TradeTargetAccept(accepted bool) *Packet {
	if accepted {
		return NewOutgoingPacket(162).AddByte(1)
	}
	return NewOutgoingPacket(162).AddByte(0)
}

//TradeAccept Builds a packet to change trade targets accepted status
func TradeAccept(accepted bool) *Packet {
	if accepted {
		return NewOutgoingPacket(15).AddByte(1)
	}
	return NewOutgoingPacket(15).AddByte(0)
}

//TradeConfirmationOpen Builds a packet to open the trade confirmation page
func TradeConfirmationOpen(player, other *world.Player) *Packet {
	p := NewOutgoingPacket(20)
	p.AddLong(other.UserBase37)

	other.TradeOffer.Lock.RLock()
	p.AddByte(uint8(len(other.TradeOffer.List)))
	for _, item := range other.TradeOffer.List {
		p.AddShort(uint16(item.ID))
		p.AddInt(uint32(item.Amount))
	}
	other.TradeOffer.Lock.RUnlock()

	player.TradeOffer.Lock.RLock()
	p.AddByte(uint8(len(player.TradeOffer.List)))
	for _, item := range player.TradeOffer.List {
		p.AddShort(uint16(item.ID))
		p.AddInt(uint32(item.Amount))
	}
	player.TradeOffer.Lock.RUnlock()
	return p
}
