/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-20-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-22-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

func init() {
	Handlers[32] = sessionRequest
	Handlers[0] = loginRequest
	Handlers[145] = func(c *Client, p *packets.Packet) {
		c.outgoingPackets <- packets.Logout
		c.kill <- struct{}{}
	}
}

func sessionRequest(c *Client, p *packets.Packet) {
	c.uID, _ = p.ReadByte()
	seed := GenerateSessionID()
	c.isaacSeed[1] = seed >> 32
	c.outgoingPackets <- packets.NewBarePacket(nil).AddLong(seed)
}

func loginRequest(c *Client, p *packets.Packet) {
	// TODO: Handle reconnect slightly different
	recon, _ := p.ReadByte()
	version, _ := p.ReadInt()
	if version != uint32(Version) {
		if len(Flags.Verbose) >= 1 {
			LogWarning.Printf("Player tried logging in with invalid client version. Got %d, expected %d\n", version, Version)
		}
		c.sendLoginResponse(5)
		return
	}
	seed := make([]uint64, 2)
	for i := 0; i < 2; i++ {
		seed[i], _ = p.ReadLong()
	}
	cipher := c.SeedISAAC(seed)
	if cipher == nil {
		c.sendLoginResponse(8)
		return
	}
	c.isaacStream = cipher
	c.player.Username, _ = p.ReadString()
	c.player.Username = strutil.DecodeBase37(strutil.Base37(c.player.Username))
	c.player.Password, _ = p.ReadString()
	c.player.Index = c.index
	//	entity.GetRegion(c.player.X(), c.player.Y()).AddPlayer(c.player)
	LogInfo.Printf("Registered Player{idx:%v,ip:'%v'username:'%v',password:'%v',reconnecting:%v,version:%v}\n", c.index, c.ip, c.player.Username, c.player.Password, recon, version)
	c.sendLoginResponse(0)
}
