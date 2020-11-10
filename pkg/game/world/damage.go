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

func (p *Player) DamageFrom(m entity.MobileEntity, damage int, kind int) bool {
	damage = int(math.Min(float64(p.Skills().Current(entity.StatHits)), float64(damage)))
	splat := NewHitsplat(p, damage)
	p.Enqueue(playerEvents, splat)
	p.Skills().SetCur(entity.StatHits, p.Skills().Current(entity.StatHits) - damage)
	if p.Skills().Current(entity.StatHits) <= 0 {
		if attacker := AsPlayer(m); attacker != nil {
			attacker.PlaySound("victory")
		}
		p.Killed(m)
		return true
	}
	return false
}

func (n *NPC) DamageFrom(m entity.MobileEntity, damage int, kind int) bool {
	damage = int(math.Min(float64(n.Skills().Current(entity.StatHits)), float64(damage)))
	splat := NewHitsplat(n, damage)
	n.enqueueArea(npcEvents, splat)
	n.Skills().SetCur(entity.StatHits, n.Skills().Current(entity.StatHits) - damage)
	if damage > 0 {
		n.CacheDamage(m.SessionCache().VarLong("username", 0), damage)
	}
	if n.Skills().Current(entity.StatHits) <= 0 {
		if attacker := AsPlayer(m); attacker != nil {
			attacker.PlaySound("victory")
		}
		n.Killed(m)
		return true
	}
	return false
}
