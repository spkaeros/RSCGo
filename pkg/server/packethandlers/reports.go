package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

var reasons = []string{
	"Offensive Language", // 1
	"Item scamming", // 2
	"Password scamming", // 3
	"Bug abuse", // 4
	"Staff impersonation", // 5
	"Account sharing", // 6
	"Macroing", // 7
	"Multiple logging in", // 8
	"Encouraging others to break rules", // 9
	"Misuse of customer support", // 10
	"Advertising/website", // 11
	"Real world item trading", // 12
}

var actions = []string {
	"reported",
	"muted",
	"banned",
}

func init() {
	PacketHandlers["reportabuse"] = func(player *world.Player, p *packet.Packet) {
		userHash := p.ReadLong()
		reasonIndex := int(p.ReadByte()-1)
		actionIndex := int(p.ReadByte())

		if reasonIndex < 0 || reasonIndex > len(reasons)-1 {
			log.Suspicious.Printf("Report had invalid reason:\n[\n\taction:%d ('%s'),\n\tsender:'%s',\n\ttarget:'%s',\n\treason:%d\n];\n", actionIndex, actions[actionIndex], player.Username(), strutil.Base37.Decode(userHash), reasonIndex+1)
			log.Info.Printf("Report had invalid reason:\n[\n\taction:%d ('%s'),\n\tsender:'%s',\n\ttarget:'%s',\n\treason:%d\n];\n", actionIndex, actions[actionIndex], player.Username(), strutil.Base37.Decode(userHash), reasonIndex+1)
			return
		}
		if actionIndex < 0 || actionIndex > len(actions)-1 {
			log.Suspicious.Printf("Report had invalid action:\n[\n\taction:%d,\n\tsender:'%s',\n\ttarget:'%s',\n\treason:%d ('%s')\n];\n", actionIndex, player.Username(), strutil.Base37.Decode(userHash), reasonIndex+1, reasons[reasonIndex])
			log.Info.Printf("Report had invalid action:\n[\n\taction:%d,\n\tsender:'%s',\n\ttarget:'%s',\n\treason:%d ('%s')\n];\n", actionIndex, player.Username(), strutil.Base37.Decode(userHash), reasonIndex+1, reasons[reasonIndex])
			return
		}
		log.Info.Printf("Report:\n[\n\taction:%d ('%s'),\n\tsender:'%s',\n\ttarget:'%s',\n\treason:%d ('%s')\n];\n", actionIndex, actions[actionIndex], player.Username(), strutil.Base37.Decode(userHash), reasonIndex+1, reasons[reasonIndex])

		log.Info.Println(player.Username(), actions[actionIndex], strutil.Base37.Decode(userHash), "for breaking rule", reasonIndex+1, "('" + reasons[reasonIndex] + "')")
		log.Suspicious.Println(player.Username(), actions[actionIndex], strutil.Base37.Decode(userHash), "for breaking rule", reasonIndex+1, "('" + reasons[reasonIndex] + "')")
	}
}
