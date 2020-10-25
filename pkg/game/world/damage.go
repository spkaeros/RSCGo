package world

import `github.com/spkaeros/rscgo/pkg/game/entity`
import `math`

type HitSplat struct {
	Owner entity.MobileEntity
	Damage int
}

func NewHitsplat(target entity.MobileEntity, damage int) interface{} {
	return &HitSplat{target, damage}
}

func (p *Player) DamageFrom(m entity.MobileEntity, damage int, kind int) {
	damage = int(math.Min(float64(p.Skills().Current(entity.StatHits)), float64(damage)))
	splat := NewHitsplat(p, damage)
	list := NewMobList()
	for _, r := range Region(p.X(), p.Y()).neighbors() {
		r.Players.RangePlayers(func(p1 *Player) bool {
			if p.Near(p1, 16) && !list.Contains(p1) {
				list.Add(p1)
				p1.enqueue(playerEvents, splat)
			}
			return false
		})
	}
	// p.Skills().SetCur(entity.StatHits, p.Skills().Current(entity.StatHits) - damage)
	// if p.Skills().Current(entity.StatHits) <= 0 {
		// if attacker := AsPlayer(m); attacker != nil {
			// attacker.PlaySound("victory")
		// }
		// p.Killed(m)
	// }
}

func (n *NPC) DamageFrom(m entity.MobileEntity, damage int, kind int) {
	damage = int(math.Min(float64(n.Skills().Current(entity.StatHits)), float64(damage)))
	splat := NewHitsplat(n, damage)
	list := NewMobList()
	for _, r := range Region(n.X(), n.Y()).neighbors() {
		r.Players.RangePlayers(func(p1 *Player) bool {
			if n.Near(p1, 16) && !list.Contains(p1) {
				list.Add(p1)
				p1.enqueue(npcEvents, splat)
			}
			return false
		})
	}
	n.Skills().SetCur(entity.StatHits, n.Skills().Current(entity.StatHits) - damage)
	if attacker := AsPlayer(m); attacker != nil {
		n.meleeRangeDamage.Put(attacker.UsernameHash(), damage)
	}
	// if n.Skills().Current(entity.StatHits) <= 0 {
		// if attacker := AsPlayer(m); attacker != nil {
			// attacker.PlaySound("victory")
		// }
		// n.Killed(m)
	// }
}
